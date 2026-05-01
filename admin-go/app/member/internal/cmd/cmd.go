package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gbaseadmin/app/member/internal/controller/biz_config"
	"gbaseadmin/app/member/internal/controller/level"
	"gbaseadmin/app/member/internal/controller/level_log"
	"gbaseadmin/app/member/internal/controller/portal"
	"gbaseadmin/app/member/internal/controller/rebind_log"
	"gbaseadmin/app/member/internal/cron"
	"gbaseadmin/app/member/internal/controller/shop_category"
	"gbaseadmin/app/member/internal/controller/shop_goods"
	"gbaseadmin/app/member/internal/controller/shop_order"
	"gbaseadmin/app/member/internal/controller/team_export"
	"gbaseadmin/app/member/internal/controller/user"
	"gbaseadmin/app/member/internal/controller/wallet"
	"gbaseadmin/app/member/internal/controller/wallet_log"
	"gbaseadmin/app/member/internal/controller/warehouse_goods"
	"gbaseadmin/app/member/internal/controller/warehouse_listing"
	"gbaseadmin/app/member/internal/controller/warehouse_trade"

	"gbaseadmin/app/member/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start member http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 注册定时任务（等级过期扫描等），server 启动前必须先注册
			if err := cron.Setup(ctx); err != nil {
				return err
			}
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// 后台 CRUD 路由：管理端账号 JWT，沿用 codegen 生成结构。
				group.Group("/api/member", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(
						biz_config.BizConfig,
						level.Level,
						level_log.LevelLog,
						rebind_log.RebindLog,
						shop_category.ShopCategory,
						shop_goods.ShopGoods,
						shop_order.ShopOrder,
						team_export.TeamExport,
						user.User,
						wallet.Wallet,
						wallet_log.WalletLog,
						warehouse_goods.WarehouseGoods,
						warehouse_listing.WarehouseListing,
						warehouse_trade.WarehouseTrade,
					)
				})

				// C 端 portal 路由：H5 会员账号 JWT。
				// 公开接口：发短信、注册、登录、找回密码、商城浏览、仓库市场浏览（无需登录）
				// 受保护接口：首页、个人中心、钱包、订单、团队、下单、挂卖、确认（需 PortalAuth）
				group.Group("/api/member-portal", func(group *ghttp.RouterGroup) {
					// 公开接口（无 PortalAuth）
					group.Bind(
						portal.Sms,
						portal.Auth,
					)
					// 公开浏览（不强制登录）
					group.Bind(portal.BizConfig)
					group.GET("/mall/categories", portal.Mall.Categories)
					group.GET("/mall/goods", portal.Mall.Goods)
					group.GET("/mall/goods/detail", portal.Mall.GoodsDetail)
					group.GET("/warehouse/market", portal.Warehouse.Market)

					// 受保护接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.PortalAuth)
						group.Bind(
							portal.Home,
							portal.Me,
						)
						// 商城受保护：下单 / 我的订单
						group.POST("/mall/order/place", portal.Mall.PlaceOrder)
						group.GET("/mall/orders", portal.Mall.MyOrders)
						// 仓库受保护：我的、挂卖、买家下单、卖家确认、我的交易
						group.GET("/warehouse/my", portal.Warehouse.MyHoldings)
						group.POST("/warehouse/list", portal.Warehouse.ListGoods)
						group.POST("/warehouse/trade/place", portal.Warehouse.PlaceTrade)
						group.POST("/warehouse/trade/confirm", portal.Warehouse.ConfirmTrade)
						group.GET("/warehouse/my-trades", portal.Warehouse.MyTrades)
					})
				})
			})
			s.Run()
			return nil
		},
	}
)
