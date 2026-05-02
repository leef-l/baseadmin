package portal

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/logic/bizconfig"
	"gbaseadmin/app/member/internal/logic/teamops"
	"gbaseadmin/app/member/internal/logic/walletops"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/snowflake"
)

// PlaceShopOrderInput 商城下单入参。
type PlaceShopOrderInput struct {
	UserID   int64
	GoodsID  int64
	Quantity int
	Remark   string
}

// PlaceShopOrderOutput 下单结果。
type PlaceShopOrderOutput struct {
	OrderID    string
	OrderNo    string
	TotalPrice int64
}

// PlaceShopOrder 商城下单事务（含时间窗、限购、阶梯返佣、自购返奖、直推返奖）。
//
// 流程：
//  1. 校验业务时间窗与工作日（来自 member_business_config）
//  2. 锁会员行：跨日重置 today_purchase_count；校验 today_purchase_count < daily_purchase_limit
//  3. 锁商品行：校验上架 + 库存
//  4. 扣优惠券钱包（同事务，自动写流水）
//  5. 写订单（含商品快照）
//  6. 更新库存 + 销量
//  7. 更新 user 限购计数 + total_purchase_count + last_purchase_date
//  8. 自购阶梯返：第 N 单匹配档位 → 优惠券钱包加余额
//  9. 自购按金额比例返：交易额 × selfTurnoverRewardRate% → 奖励钱包
// 10. 直推按金额比例返：上一级（仅 1 层）的推广钱包加金额
// 11. 链式累加祖先 team_turnover + 触发等级晋升
//
// 任一失败回滚整个事务。
func (s *sPortalAuth) PlaceShopOrder(ctx context.Context, in *PlaceShopOrderInput) (*PlaceShopOrderOutput, error) {
	if in == nil {
		return nil, gerror.New("下单参数不能为空")
	}
	if in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	if in.GoodsID <= 0 {
		return nil, gerror.New("商品不存在")
	}
	qty := in.Quantity
	if qty <= 0 {
		qty = 1
	}
	if qty > 999 {
		return nil, gerror.New("单次购买数量不能超过 999 件")
	}

	cfg, err := bizconfig.GetCachedConfig(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if ok, msg := cfg.IsPurchaseAllowed(now); !ok {
		return nil, gerror.New(msg)
	}

	orderID := snowflake.Generate()
	orderNo := fmt.Sprintf("M%d", int64(orderID))

	var (
		totalPrice int64
		goodsTitle string
		goodsCover string
		tenantID   uint64
		merchantID uint64
		nthOrder   int // 本次属于该会员历史第几单
	)

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// ---- 1. 锁会员行，校验限购、跨日重置 ----
		var user entity.MemberUser
		if err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, in.UserID).
			Where(dao.MemberUser.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&user); err != nil {
			return err
		}
		if user.Id == 0 {
			return gerror.New("会员不存在")
		}
		if user.Status != 1 {
			return gerror.New("会员账号已冻结")
		}
		todayCount := int(user.TodayPurchaseCount)
		todayStr := now.Format("2006-01-02")
		lastStr := ""
		if user.LastPurchaseDate != nil && !user.LastPurchaseDate.IsZero() {
			lastStr = user.LastPurchaseDate.Format("2006-01-02")
		}
		if lastStr != todayStr {
			todayCount = 0
		}
		dailyLimit := int(user.DailyPurchaseLimit)
		if dailyLimit <= 0 {
			dailyLimit = 1
		}
		if todayCount >= dailyLimit {
			return gerror.Newf("今日限购已用完（%d/%d 单）", todayCount, dailyLimit)
		}

		// ---- 2. 锁商品行 ----
		var goods entity.MemberShopGoods
		if err := tx.Model(dao.MemberShopGoods.Table()).Ctx(ctx).
			Where(dao.MemberShopGoods.Columns().Id, in.GoodsID).
			Where(dao.MemberShopGoods.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&goods); err != nil {
			return err
		}
		if goods.Id == 0 {
			return gerror.New("商品不存在或已下架")
		}
		if goods.Status != 1 {
			return gerror.New("商品已下架")
		}
		if int(goods.Stock) < qty {
			return gerror.New("商品库存不足")
		}
		totalPrice = int64(goods.Price) * int64(qty)
		if totalPrice <= 0 {
			return gerror.New("订单金额异常")
		}
		goodsTitle = goods.Title
		goodsCover = goods.Cover
		tenantID = goods.TenantId
		merchantID = goods.MerchantId

		// ---- 3. 扣优惠券钱包 ----
		if _, err := walletops.Debit(ctx, tx, &walletops.MoveInput{
			UserID:         in.UserID,
			WalletType:     walletops.WalletTypeCoupon,
			ChangeType:     walletops.ChangeTypeConsume,
			Amount:         totalPrice,
			RelatedOrderNo: orderNo,
			Remark:         "商城下单：" + goodsTitle,
			Operator:       in.UserID,
		}); err != nil {
			return err
		}

		// ---- 4. 写订单 ----
		nowT := gtime.Now()
		if _, err := tx.Model(dao.MemberShopOrder.Table()).Ctx(ctx).Data(do.MemberShopOrder{
			Id:          orderID,
			OrderNo:     orderNo,
			UserId:      in.UserID,
			GoodsId:     in.GoodsID,
			GoodsTitle:  goodsTitle,
			GoodsCover:  goodsCover,
			Quantity:    qty,
			TotalPrice:  totalPrice,
			PayWallet:   1,
			OrderStatus: 1,
			Remark:      strings.TrimSpace(in.Remark),
			Status:      1,
			TenantId:    tenantID,
			MerchantId:  merchantID,
			CreatedBy:   in.UserID,
			DeptId:      0,
			CreatedAt:   nowT,
			UpdatedAt:   nowT,
		}).Insert(); err != nil {
			return err
		}

		// ---- 5. 更新库存 + 销量 ----
		if _, err := tx.Model(dao.MemberShopGoods.Table()).Ctx(ctx).
			Where(dao.MemberShopGoods.Columns().Id, in.GoodsID).
			Data(g.Map{
				dao.MemberShopGoods.Columns().Stock: gdb.Raw(fmt.Sprintf("stock - %d", qty)),
				dao.MemberShopGoods.Columns().Sales: gdb.Raw(fmt.Sprintf("sales + %d", qty)),
			}).Update(); err != nil {
			return err
		}

		// ---- 6. 更新会员限购计数 / 历史购单数 ----
		// nthOrder = 历史累计第几单（一生维度，一旦达到某档位即奖一次，永不重复）。
		// today_purchase_count 是当日维度，仅用于限购判断，与阶梯返佣无关。
		// 例：用户一生第 2 单 → 奖 88；之后无论今天明天再下单都不再奖第 2 档；只有再到第 3、4 单才匹配新档位。
		nthOrder = int(user.TotalPurchaseCount) + 1
		if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, in.UserID).
			Data(g.Map{
				dao.MemberUser.Columns().TodayPurchaseCount: todayCount + 1,
				dao.MemberUser.Columns().LastPurchaseDate:   todayStr,
				dao.MemberUser.Columns().TotalPurchaseCount: nthOrder,
			}).Update(); err != nil {
			return err
		}

		// ---- 7. 自购阶梯返（优惠券钱包） ----
		if reward := cfg.FindRebateTier(nthOrder); reward > 0 {
			cents := yuanToCent(reward)
			if cents > 0 {
				if _, err := walletops.Credit(ctx, tx, &walletops.MoveInput{
					UserID:         in.UserID,
					WalletType:     walletops.WalletTypeCoupon,
					ChangeType:     walletops.ChangeTypeSelfRebateTier,
					Amount:         cents,
					RelatedOrderNo: orderNo,
					Remark:         fmt.Sprintf("第 %d 单阶梯奖励", nthOrder),
					Operator:       in.UserID,
				}); err != nil {
					return err
				}
			}
		}

		// ---- 8. 自购按金额比例返奖励钱包 ----
		if cfg.SelfTurnoverRewardRate > 0 {
			rewardCents := pctOfCent(totalPrice, cfg.SelfTurnoverRewardRate)
			if rewardCents > 0 {
				if _, err := walletops.Credit(ctx, tx, &walletops.MoveInput{
					UserID:         in.UserID,
					WalletType:     walletops.WalletTypeReward,
					ChangeType:     walletops.ChangeTypeSelfTurnoverReward,
					Amount:         rewardCents,
					RelatedOrderNo: orderNo,
					Remark:         fmt.Sprintf("自购返 %.2f%%", cfg.SelfTurnoverRewardRate),
					Operator:       in.UserID,
				}); err != nil {
					return err
				}
			}
		}

		// ---- 9. 直推（上一级，仅 1 层）按金额比例返推广钱包 ----
		if cfg.DirectPromoteRate > 0 && user.ParentId > 0 {
			promoteCents := pctOfCent(totalPrice, cfg.DirectPromoteRate)
			if promoteCents > 0 {
				if _, err := walletops.Credit(ctx, tx, &walletops.MoveInput{
					UserID:         int64(user.ParentId),
					WalletType:     walletops.WalletTypePromote,
					ChangeType:     walletops.ChangeTypeDirectPromoteReward,
					Amount:         promoteCents,
					RelatedOrderNo: orderNo,
					Remark:         fmt.Sprintf("直推下级进货返 %.2f%%", cfg.DirectPromoteRate),
					Operator:       in.UserID,
				}); err != nil {
					return err
				}
			}
		}

		// ---- 10. 链式累加祖先 team_turnover + 触发等级晋升 ----
		if err := teamops.AddAncestorTurnover(ctx, tx, in.UserID, totalPrice); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &PlaceShopOrderOutput{
		OrderID:    fmt.Sprintf("%d", int64(orderID)),
		OrderNo:    orderNo,
		TotalPrice: totalPrice,
	}, nil
}

// yuanToCent 元 → 分（四舍五入）。
func yuanToCent(yuan float64) int64 {
	if yuan <= 0 {
		return 0
	}
	return int64(math.Round(yuan * 100))
}

// pctOfCent 按百分比对金额计算分（math.Floor，不向上舍入避免多发奖）。
func pctOfCent(amountCent int64, percent float64) int64 {
	if amountCent <= 0 || percent <= 0 {
		return 0
	}
	return int64(math.Floor(float64(amountCent) * percent / 100.0))
}
