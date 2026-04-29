package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gbaseadmin/app/member/internal/controller/level"
	"gbaseadmin/app/member/internal/controller/level_log"
	"gbaseadmin/app/member/internal/controller/rebind_log"
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
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Group("/api/member", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(
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
			})
			s.Run()
			return nil
		},
	}
)
