package portal

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
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
	OrderID    string // 订单 ID（雪花）
	OrderNo    string // 订单号（业务编号）
	TotalPrice int64  // 订单总价（分）
}

// PlaceShopOrder 处理商城下单的整个事务：
//
//  1. 锁库存：SELECT FOR UPDATE goods 行，校验上架 + 库存
//  2. 计算总价 = goods.price * quantity
//  3. 调 walletops.Debit 扣优惠券钱包（自动写流水 + 锁钱包行 + 余额校验）
//  4. 写订单（含商品快照）
//  5. 更新 goods.stock 和 sales
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

	orderID := snowflake.Generate()
	orderNo := fmt.Sprintf("M%d", int64(orderID))

	var totalPrice int64
	var goodsTitle, goodsCover string
	var tenantID, merchantID uint64

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 锁商品行
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

		// 2. 扣优惠券钱包（同事务）
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

		// 3. 写订单
		now := gtime.Now()
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
			OrderStatus: 1, // 已完成（线下交易，下单即完成）
			Remark:      strings.TrimSpace(in.Remark),
			Status:      1,
			TenantId:    tenantID,
			MerchantId:  merchantID,
			CreatedBy:   in.UserID,
			DeptId:      0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}).Insert(); err != nil {
			return err
		}

		// 4. 更新库存和销量
		if _, err := tx.Model(dao.MemberShopGoods.Table()).Ctx(ctx).
			Where(dao.MemberShopGoods.Columns().Id, in.GoodsID).
			Data(g.Map{
				dao.MemberShopGoods.Columns().Stock: gdb.Raw(fmt.Sprintf("stock - %d", qty)),
				dao.MemberShopGoods.Columns().Sales: gdb.Raw(fmt.Sprintf("sales + %d", qty)),
			}).Update(); err != nil {
			return err
		}

		// 5. 链式累加祖先 team_turnover，并触发各祖先 TryUpgrade
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
