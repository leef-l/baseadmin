package main

import (
	_ "gbaseadmin/app/demo/internal/packed"

	_ "gbaseadmin/app/demo/internal/logic/appointment"
	_ "gbaseadmin/app/demo/internal/logic/audit_log"
	_ "gbaseadmin/app/demo/internal/logic/campaign"
	_ "gbaseadmin/app/demo/internal/logic/category"
	_ "gbaseadmin/app/demo/internal/logic/contract"
	_ "gbaseadmin/app/demo/internal/logic/customer"
	_ "gbaseadmin/app/demo/internal/logic/order"
	_ "gbaseadmin/app/demo/internal/logic/product"
	_ "gbaseadmin/app/demo/internal/logic/survey"
	_ "gbaseadmin/app/demo/internal/logic/work_order"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"gbaseadmin/app/demo/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
