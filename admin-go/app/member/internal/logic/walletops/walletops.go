// Package walletops 是会员钱包的领域服务（基础设施层）。
//
// 设计要点：
//   - 钱包变动唯一入口：所有金额变化（注册赠送、商城下单、寄售确认、推广奖、后台调整）必须走 Debit/Credit/Freeze/Unfreeze 之一。
//   - 行级锁：每次操作都用 LockUpdate 锁住目标钱包行，避免并发扣减导致负余额或丢失流水。
//   - 流水一致：钱包余额更新与 wallet_log 写入在同一个事务内完成，业务方传入的 tx 可复用同一事务（避免嵌套事务）。
//   - 金额单位：全部按"分"传入（int64），调用方负责元↔分换算。
//   - 不做软删除：流水记录用 INSERT-only，DELETE 操作禁止。
package walletops

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/snowflake"
)

// Wallet 类型常量。
const (
	WalletTypeCoupon  = 1 // 优惠券余额（仅可消费商城商品）
	WalletTypeReward  = 2 // 奖金余额（寄售卖出收入进入此处）
	WalletTypePromote = 3 // 推广奖余额（拉新分成等）
)

// Change 类型常量（与 member_wallet_log.change_type 取值一致）。
const (
	ChangeTypeRecharge            = 1  // 充值
	ChangeTypeConsume             = 2  // 消费（商城下单）
	ChangeTypePromote             = 3  // 推广奖
	ChangeTypeWHIncome            = 4  // 仓库卖出收入（卖家奖金）
	ChangeTypeFee                 = 5  // 平台扣除
	ChangeTypeAdjust              = 6  // 后台调整
	ChangeTypeSelfRebateTier      = 11 // 自购阶梯返（如第 2/3/4 单 88/188/288）
	ChangeTypeSelfTurnoverReward  = 12 // 自购按金额比例返奖励钱包
	ChangeTypeDirectPromoteReward = 13 // 直推下级进货按金额比例返推广钱包
)

// MoveInput 描述一次钱包变动。
//
// 字段说明：
//   - UserID 必填，目标会员 ID。
//   - WalletType 必填，钱包类型。
//   - ChangeType 必填，变动业务类型。
//   - Amount 必填，正数（>0），由 Debit/Credit 决定方向。
//   - RelatedOrderNo 可选，关联订单号 / 交易号 / 业务编号，便于流水溯源。
//   - Remark 可选，业务备注。
//   - Operator 可选，操作人 ID（后台调整时填管理员 ID，C 端业务填 0）。
type MoveInput struct {
	UserID         int64
	WalletType     int
	ChangeType     int
	Amount         int64
	RelatedOrderNo string
	Remark         string
	Operator       int64
}

// Debit 扣减余额（用于消费场景）。
//
// 必须事务内调用：传入 tx 是上层业务事务，钱包余额更新和流水写入会复用此事务。
// 余额不足时返回明确错误，调用方不应吞掉。
//
// 返回值是变动后的钱包流水 ID（雪花），方便业务关联。
func Debit(ctx context.Context, tx gdb.TX, in *MoveInput) (snowflake.JsonInt64, error) {
	if err := validateMoveInput(in); err != nil {
		return 0, err
	}
	wallet, err := lockWallet(ctx, tx, in.UserID, in.WalletType)
	if err != nil {
		return 0, err
	}
	if wallet.Balance < in.Amount {
		return 0, gerror.New(walletShortageError(in.WalletType))
	}
	beforeBalance := wallet.Balance
	afterBalance := wallet.Balance - in.Amount

	if _, err := tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
		Where(dao.MemberWallet.Columns().Id, wallet.Id).
		Data(g.Map{
			dao.MemberWallet.Columns().Balance:      afterBalance,
			dao.MemberWallet.Columns().TotalExpense: gdb.Raw("total_expense + " + i64s(in.Amount)),
		}).Update(); err != nil {
		return 0, err
	}
	return appendLog(ctx, tx, wallet, in, -in.Amount, beforeBalance, afterBalance)
}

// Credit 增加余额（用于收入场景）。
func Credit(ctx context.Context, tx gdb.TX, in *MoveInput) (snowflake.JsonInt64, error) {
	if err := validateMoveInput(in); err != nil {
		return 0, err
	}
	wallet, err := lockWallet(ctx, tx, in.UserID, in.WalletType)
	if err != nil {
		return 0, err
	}
	beforeBalance := wallet.Balance
	afterBalance := wallet.Balance + in.Amount

	if _, err := tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
		Where(dao.MemberWallet.Columns().Id, wallet.Id).
		Data(g.Map{
			dao.MemberWallet.Columns().Balance:     afterBalance,
			dao.MemberWallet.Columns().TotalIncome: gdb.Raw("total_income + " + i64s(in.Amount)),
		}).Update(); err != nil {
		return 0, err
	}
	return appendLog(ctx, tx, wallet, in, in.Amount, beforeBalance, afterBalance)
}

// Freeze 冻结余额（金额从 balance 转到 frozen_amount，不写流水，业务后续 Unfreeze 或 Settle）。
//
// 暂未启用，留给后续"挂卖商品时锁定保证金"或"提现申请"扩展。
func Freeze(ctx context.Context, tx gdb.TX, in *MoveInput) error {
	if err := validateMoveInput(in); err != nil {
		return err
	}
	wallet, err := lockWallet(ctx, tx, in.UserID, in.WalletType)
	if err != nil {
		return err
	}
	if wallet.Balance < in.Amount {
		return gerror.New(walletShortageError(in.WalletType))
	}
	_, err = tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
		Where(dao.MemberWallet.Columns().Id, wallet.Id).
		Data(g.Map{
			dao.MemberWallet.Columns().Balance:      wallet.Balance - in.Amount,
			dao.MemberWallet.Columns().FrozenAmount: gdb.Raw("frozen_amount + " + i64s(in.Amount)),
		}).Update()
	return err
}

// Unfreeze 解冻（金额从 frozen_amount 转回 balance，不写流水）。
func Unfreeze(ctx context.Context, tx gdb.TX, in *MoveInput) error {
	if err := validateMoveInput(in); err != nil {
		return err
	}
	wallet, err := lockWallet(ctx, tx, in.UserID, in.WalletType)
	if err != nil {
		return err
	}
	if int64(wallet.FrozenAmount) < in.Amount {
		return gerror.New("冻结金额不足")
	}
	_, err = tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
		Where(dao.MemberWallet.Columns().Id, wallet.Id).
		Data(g.Map{
			dao.MemberWallet.Columns().Balance:      wallet.Balance + in.Amount,
			dao.MemberWallet.Columns().FrozenAmount: gdb.Raw("frozen_amount - " + i64s(in.Amount)),
		}).Update()
	return err
}

// GetBalance 查询会员某钱包余额（不加锁，纯读）。
//
// 返回 (balance, frozen, exists, error)。如果钱包不存在返回 exists=false 但不报错。
func GetBalance(ctx context.Context, userID int64, walletType int) (balance int64, frozen int64, exists bool, err error) {
	if userID <= 0 {
		return 0, 0, false, gerror.New("会员 ID 不能为空")
	}
	var w entity.MemberWallet
	err = dao.MemberWallet.Ctx(ctx).
		Where(dao.MemberWallet.Columns().UserId, userID).
		Where(dao.MemberWallet.Columns().WalletType, walletType).
		Where(dao.MemberWallet.Columns().DeletedAt, nil).
		Scan(&w)
	if err != nil {
		return 0, 0, false, err
	}
	if w.Id == 0 {
		return 0, 0, false, nil
	}
	return w.Balance, int64(w.FrozenAmount), true, nil
}

// ----- internal helpers -----

func validateMoveInput(in *MoveInput) error {
	if in == nil {
		return gerror.New("钱包变动参数不能为空")
	}
	if in.UserID <= 0 {
		return gerror.New("会员 ID 不能为空")
	}
	if in.WalletType != WalletTypeCoupon && in.WalletType != WalletTypeReward && in.WalletType != WalletTypePromote {
		return gerror.New("钱包类型不正确")
	}
	if in.Amount <= 0 {
		return gerror.New("金额必须为正数")
	}
	if in.ChangeType <= 0 {
		return gerror.New("变动类型不能为空")
	}
	return nil
}

// lockWallet 用 SELECT ... FOR UPDATE 锁住会员某钱包，确保后续 Update 不被并发覆盖。
func lockWallet(ctx context.Context, tx gdb.TX, userID int64, walletType int) (*entity.MemberWallet, error) {
	var w entity.MemberWallet
	err := tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
		Where(dao.MemberWallet.Columns().UserId, userID).
		Where(dao.MemberWallet.Columns().WalletType, walletType).
		Where(dao.MemberWallet.Columns().DeletedAt, nil).
		LockUpdate().
		Scan(&w)
	if err != nil {
		return nil, err
	}
	if w.Id == 0 {
		return nil, gerror.Newf("会员钱包(type=%d)不存在", walletType)
	}
	if w.Status != 1 {
		return nil, gerror.New("钱包已冻结，操作受限")
	}
	return &w, nil
}

// appendLog 写流水。changeAmount 已带正负号（Debit 传 -amount，Credit 传 +amount）。
func appendLog(
	ctx context.Context,
	tx gdb.TX,
	wallet *entity.MemberWallet,
	in *MoveInput,
	signedChangeAmount int64,
	beforeBalance int64,
	afterBalance int64,
) (snowflake.JsonInt64, error) {
	logID := snowflake.Generate()
	now := gtime.Now()
	if _, err := tx.Model(dao.MemberWalletLog.Table()).Ctx(ctx).Data(do.MemberWalletLog{
		Id:             logID,
		UserId:         in.UserID,
		WalletType:     in.WalletType,
		ChangeType:     in.ChangeType,
		ChangeAmount:   signedChangeAmount,
		BeforeBalance:  beforeBalance,
		AfterBalance:   afterBalance,
		RelatedOrderNo: strings.TrimSpace(in.RelatedOrderNo),
		Remark:         strings.TrimSpace(in.Remark),
		TenantId:       wallet.TenantId,
		MerchantId:     wallet.MerchantId,
		CreatedBy:      in.Operator,
		DeptId:         0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}).Insert(); err != nil {
		return 0, err
	}
	return logID, nil
}

// walletShortageError 返回针对不同钱包的友好错误文案。
func walletShortageError(walletType int) string {
	switch walletType {
	case WalletTypeCoupon:
		return "优惠券余额不足"
	case WalletTypeReward:
		return "奖金余额不足"
	case WalletTypePromote:
		return "推广奖余额不足"
	default:
		return "余额不足"
	}
}

// i64s 把 int64 转字符串，用于 gdb.Raw 拼 SQL（如 total_income + 100）。
func i64s(v int64) string {
	return strconv.FormatInt(v, 10)
}
