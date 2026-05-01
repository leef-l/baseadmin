package portal

import (
	"context"
	"fmt"

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

// 挂卖商品状态。
const (
	GoodsStatusHolding = 1 // 持有中
	GoodsStatusListed  = 2 // 挂卖中
	GoodsStatusTrading = 3 // 交易中
)

// 挂卖记录状态。
const (
	ListingStatusActive   = 1 // 挂卖中
	ListingStatusSold     = 2 // 已售出
	ListingStatusCanceled = 3 // 已取消
)

// 交易记录状态。
const (
	TradeStatusPending   = 1 // 待卖家确认
	TradeStatusConfirmed = 2 // 已确认完成
	TradeStatusCanceled  = 3 // 已取消
)

// ListWarehouseGoodsInput 挂卖入参（卖家发起，但价格由系统按 price_rise_rate 自动算）。
type ListWarehouseGoodsInput struct {
	UserID  int64 // 当前操作会员（必须 = goods.owner_id）
	GoodsID int64
}

// ListWarehouseGoodsOutput 挂卖结果。
type ListWarehouseGoodsOutput struct {
	ListingID    string
	ListingPrice int64 // 自动加价后的挂卖价（分）
}

// ListWarehouseGoods 把持有中的仓库商品挂到市场。
//
// 挂卖价 = current_price × (1 + price_rise_rate%)，向下取整到分。
// 第一次挂卖时 current_price 还是 init_price，所以等价于 init_price × (1+r)。
func (s *sPortalAuth) ListWarehouseGoods(ctx context.Context, in *ListWarehouseGoodsInput) (*ListWarehouseGoodsOutput, error) {
	if in == nil {
		return nil, gerror.New("挂卖参数不能为空")
	}
	if in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	if in.GoodsID <= 0 {
		return nil, gerror.New("商品 ID 不能为空")
	}

	listingID := snowflake.Generate()
	var listingPrice int64

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 锁商品行
		var goods entity.MemberWarehouseGoods
		if err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, in.GoodsID).
			Where(dao.MemberWarehouseGoods.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&goods); err != nil {
			return err
		}
		if goods.Id == 0 {
			return gerror.New("商品不存在或已删除")
		}
		if int64(goods.OwnerId) != in.UserID {
			return gerror.New("不是该商品当前持有人，无法挂卖")
		}
		if goods.GoodsStatus != GoodsStatusHolding {
			return gerror.New("商品当前状态不可挂卖")
		}
		if goods.Status != 1 {
			return gerror.New("商品已停用")
		}

		// 2. 计算挂卖价（current_price * (1 + rate/100)，分级精度向下取整）
		listingPrice = computeListingPrice(int64(goods.CurrentPrice), int64(goods.PriceRiseRate))
		if listingPrice <= 0 {
			return gerror.New("挂卖价计算异常")
		}

		// 3. 写挂卖记录
		now := gtime.Now()
		if _, err := tx.Model(dao.MemberWarehouseListing.Table()).Ctx(ctx).Data(do.MemberWarehouseListing{
			Id:            listingID,
			GoodsId:       in.GoodsID,
			SellerId:      in.UserID,
			ListingPrice:  listingPrice,
			ListingStatus: ListingStatusActive,
			ListedAt:      now,
			Status:        1,
			TenantId:      goods.TenantId,
			MerchantId:    goods.MerchantId,
			CreatedBy:     in.UserID,
			DeptId:        0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}).Insert(); err != nil {
			return err
		}

		// 4. 商品状态切到挂卖中
		if _, err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, in.GoodsID).
			Data(g.Map{
				dao.MemberWarehouseGoods.Columns().GoodsStatus: GoodsStatusListed,
			}).Update(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ListWarehouseGoodsOutput{
		ListingID:    fmt.Sprintf("%d", int64(listingID)),
		ListingPrice: listingPrice,
	}, nil
}

// PlaceWarehouseTradeInput 买家下单入参。
type PlaceWarehouseTradeInput struct {
	BuyerID   int64
	ListingID int64
}

// PlaceWarehouseTradeOutput 下单结果（生成交易，待卖家确认）。
type PlaceWarehouseTradeOutput struct {
	TradeID string
	TradeNo string
}

// PlaceWarehouseTrade 买家下单：
//   - 校验买家 is_qualified=1
//   - 锁住挂卖记录 + 商品行
//   - 状态从 active → 不变（保留挂卖中），商品从 listed → trading
//   - 生成 trade(状态=待卖家确认)
//   - 买家钱包不变（线下交易）
func (s *sPortalAuth) PlaceWarehouseTrade(ctx context.Context, in *PlaceWarehouseTradeInput) (*PlaceWarehouseTradeOutput, error) {
	if in == nil {
		return nil, gerror.New("下单参数不能为空")
	}
	if in.BuyerID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	if in.ListingID <= 0 {
		return nil, gerror.New("挂卖记录不存在")
	}

	// 校验买家资格（不在事务内，资格短时间不变）
	var buyer entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, in.BuyerID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Scan(&buyer); err != nil {
		return nil, err
	}
	if buyer.Id == 0 {
		return nil, gerror.New("买家不存在")
	}
	if buyer.Status != 1 {
		return nil, gerror.New("账号已被禁用")
	}
	if buyer.IsQualified != 1 {
		return nil, gerror.New("您已失去仓库购买资格，请联系客服或等级续费")
	}

	tradeID := snowflake.Generate()
	tradeNo := fmt.Sprintf("W%d", int64(tradeID))

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 锁挂卖记录
		var listing entity.MemberWarehouseListing
		if err := tx.Model(dao.MemberWarehouseListing.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseListing.Columns().Id, in.ListingID).
			Where(dao.MemberWarehouseListing.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&listing); err != nil {
			return err
		}
		if listing.Id == 0 {
			return gerror.New("挂卖记录不存在")
		}
		if listing.ListingStatus != ListingStatusActive {
			return gerror.New("该商品已被其他买家拍下或已下架")
		}
		if int64(listing.SellerId) == in.BuyerID {
			return gerror.New("不能购买自己挂卖的商品")
		}

		// 2. 锁商品行
		var goods entity.MemberWarehouseGoods
		if err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, listing.GoodsId).
			Where(dao.MemberWarehouseGoods.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&goods); err != nil {
			return err
		}
		if goods.Id == 0 {
			return gerror.New("商品不存在")
		}
		if goods.GoodsStatus != GoodsStatusListed {
			return gerror.New("商品状态不可购买")
		}

		// 3. 写交易记录（待卖家确认）
		now := gtime.Now()
		if _, err := tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx).Data(do.MemberWarehouseTrade{
			Id:           tradeID,
			TradeNo:      tradeNo,
			GoodsId:      listing.GoodsId,
			ListingId:    listing.Id,
			SellerId:     listing.SellerId,
			BuyerId:      in.BuyerID,
			TradePrice:   listing.ListingPrice,
			PlatformFee:  0, // 待确认时还未结算
			SellerIncome: 0,
			TradeStatus:  TradeStatusPending,
			Status:       1,
			TenantId:     listing.TenantId,
			MerchantId:   listing.MerchantId,
			CreatedBy:    in.BuyerID,
			DeptId:       0,
			CreatedAt:    now,
			UpdatedAt:    now,
		}).Insert(); err != nil {
			return err
		}

		// 4. 商品切到交易中
		if _, err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, goods.Id).
			Data(g.Map{
				dao.MemberWarehouseGoods.Columns().GoodsStatus: GoodsStatusTrading,
			}).Update(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &PlaceWarehouseTradeOutput{
		TradeID: fmt.Sprintf("%d", int64(tradeID)),
		TradeNo: tradeNo,
	}, nil
}

// ConfirmWarehouseTradeInput 卖家确认入参。
type ConfirmWarehouseTradeInput struct {
	SellerID int64
	TradeID  int64
}

// ConfirmWarehouseTradeOutput 确认结果。
type ConfirmWarehouseTradeOutput struct {
	TradeID      string
	TradePrice   int64 // 成交价（分）
	PlatformFee  int2  // 平台抽成（仅从增值部分扣，分）
	SellerIncome int64 // 卖家奖金（增值-抽成，进 type=2 钱包，分）
}

// 简单别名，避免 IDE 误解为字符串。
type int2 = int64

// ConfirmWarehouseTrade 卖家确认交易（核心资金分配点）。
//
// 资金分配规则（用户已锁死）：
//   - 增值      = trade_price - current_price        ← 这次涨的差额
//   - 平台抽成   = 增值 × platform_fee_rate / 100    ← 仅从增值扣，向下取整
//   - 卖家奖金   = 增值 - 平台抽成                    ← 进卖家奖金钱包(type=2)
//   - 本金部分   = current_price                      ← 不进任何钱包，留在商品里继续滚下次
//
// 同事务执行：
//   1. 锁 trade、listing、goods 三行
//   2. 写卖家奖金钱包（Credit 到 type=2，自动写流水）
//   3. trade.status=2、写 platform_fee 和 seller_income、写 confirmed_at
//   4. listing.status=2 已售出、写 sold_at
//   5. goods.owner_id=buyer、current_price=trade_price、trade_count+=1、goods_status=持有中
func (s *sPortalAuth) ConfirmWarehouseTrade(ctx context.Context, in *ConfirmWarehouseTradeInput) (*ConfirmWarehouseTradeOutput, error) {
	if in == nil {
		return nil, gerror.New("确认参数不能为空")
	}
	if in.SellerID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	if in.TradeID <= 0 {
		return nil, gerror.New("交易 ID 不能为空")
	}

	var (
		tradePrice   int64
		platformFee  int64
		sellerIncome int64
	)

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 锁 trade
		var trade entity.MemberWarehouseTrade
		if err := tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseTrade.Columns().Id, in.TradeID).
			Where(dao.MemberWarehouseTrade.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&trade); err != nil {
			return err
		}
		if trade.Id == 0 {
			return gerror.New("交易记录不存在")
		}
		if int64(trade.SellerId) != in.SellerID {
			return gerror.New("您不是该交易的卖家")
		}
		if trade.TradeStatus != TradeStatusPending {
			return gerror.New("交易状态不可确认")
		}

		// 2. 锁 listing
		var listing entity.MemberWarehouseListing
		if err := tx.Model(dao.MemberWarehouseListing.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseListing.Columns().Id, trade.ListingId).
			Where(dao.MemberWarehouseListing.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&listing); err != nil {
			return err
		}
		if listing.Id == 0 || listing.ListingStatus != ListingStatusActive {
			return gerror.New("挂卖记录状态异常")
		}

		// 3. 锁 goods
		var goods entity.MemberWarehouseGoods
		if err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, trade.GoodsId).
			Where(dao.MemberWarehouseGoods.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&goods); err != nil {
			return err
		}
		if goods.Id == 0 || goods.GoodsStatus != GoodsStatusTrading {
			return gerror.New("商品状态不可确认")
		}

		// 4. 计算资金分配（基于 goods.current_price 当前快照 = 上次成交价）
		tradePrice = int64(trade.TradePrice)
		appreciation := tradePrice - int64(goods.CurrentPrice) // 增值部分
		if appreciation < 0 {
			// 理论不会发生：挂卖价 = current_price × (1+r) ≥ current_price
			return gerror.New("成交价低于商品当前价，交易异常")
		}
		platformFee = appreciation * int64(goods.PlatformFeeRate) / 100 // 向下取整
		sellerIncome = appreciation - platformFee
		if sellerIncome < 0 {
			sellerIncome = 0
		}

		// 5. 卖家奖金入账（同事务）
		if sellerIncome > 0 {
			if _, err := walletops.Credit(ctx, tx, &walletops.MoveInput{
				UserID:         in.SellerID,
				WalletType:     walletops.WalletTypeReward,
				ChangeType:     walletops.ChangeTypeWHIncome,
				Amount:         sellerIncome,
				RelatedOrderNo: trade.TradeNo,
				Remark:         fmt.Sprintf("仓库商品 %s 成交奖励", goods.GoodsNo),
				Operator:       in.SellerID,
			}); err != nil {
				return err
			}
		}

		now := gtime.Now()

		// 6. 更新 trade
		if _, err := tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseTrade.Columns().Id, trade.Id).
			Data(g.Map{
				dao.MemberWarehouseTrade.Columns().PlatformFee:  platformFee,
				dao.MemberWarehouseTrade.Columns().SellerIncome: sellerIncome,
				dao.MemberWarehouseTrade.Columns().TradeStatus:  TradeStatusConfirmed,
				dao.MemberWarehouseTrade.Columns().ConfirmedAt:  now,
			}).Update(); err != nil {
			return err
		}

		// 7. 更新 listing
		if _, err := tx.Model(dao.MemberWarehouseListing.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseListing.Columns().Id, listing.Id).
			Data(g.Map{
				dao.MemberWarehouseListing.Columns().ListingStatus: ListingStatusSold,
				dao.MemberWarehouseListing.Columns().SoldAt:        now,
			}).Update(); err != nil {
			return err
		}

		// 8. 更新 goods（owner 转移、current_price 更新、trade_count+1、状态回到持有）
		if _, err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, goods.Id).
			Data(g.Map{
				dao.MemberWarehouseGoods.Columns().OwnerId:      trade.BuyerId,
				dao.MemberWarehouseGoods.Columns().CurrentPrice: tradePrice,
				dao.MemberWarehouseGoods.Columns().TradeCount:   gdb.Raw("trade_count + 1"),
				dao.MemberWarehouseGoods.Columns().GoodsStatus:  GoodsStatusHolding,
			}).Update(); err != nil {
			return err
		}

		// 9. 累加买家祖先 team_turnover（成交价计入），并触发各祖先 TryUpgrade
		if err := teamops.AddAncestorTurnover(ctx, tx, int64(trade.BuyerId), tradePrice); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ConfirmWarehouseTradeOutput{
		TradeID:      fmt.Sprintf("%d", in.TradeID),
		TradePrice:   tradePrice,
		PlatformFee:  platformFee,
		SellerIncome: sellerIncome,
	}, nil
}

// computeListingPrice 计算挂卖价：current × (1 + rate/100)，向下取整到分。
//
// rate 是百分比整数：rate=10 表示 +10%。
// rate=0 时挂卖价 = current_price。
// 中间用 int64 防溢出（current_price 用 int64，rate 用 int64，乘积上限 = 10^10*10^4 不溢出）。
func computeListingPrice(currentPrice, rate int64) int64 {
	if rate <= 0 {
		return currentPrice
	}
	return currentPrice + currentPrice*rate/100
}
