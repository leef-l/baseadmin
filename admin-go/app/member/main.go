package main

import (
	_ "gbaseadmin/app/member/internal/packed"

	_ "gbaseadmin/app/member/internal/logic/level"
	_ "gbaseadmin/app/member/internal/logic/level_log"
	_ "gbaseadmin/app/member/internal/logic/portal"
	_ "gbaseadmin/app/member/internal/logic/rebind_log"
	_ "gbaseadmin/app/member/internal/logic/shop_category"
	_ "gbaseadmin/app/member/internal/logic/shop_goods"
	_ "gbaseadmin/app/member/internal/logic/shop_order"
	_ "gbaseadmin/app/member/internal/logic/team_export"
	_ "gbaseadmin/app/member/internal/logic/user"
	_ "gbaseadmin/app/member/internal/logic/wallet"
	_ "gbaseadmin/app/member/internal/logic/wallet_log"
	_ "gbaseadmin/app/member/internal/logic/warehouse_goods"
	_ "gbaseadmin/app/member/internal/logic/warehouse_listing"
	_ "gbaseadmin/app/member/internal/logic/warehouse_trade"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"gbaseadmin/app/member/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
