package main

import (
	_ "gbaseadmin/app/system/internal/packed"

	_ "gbaseadmin/app/system/internal/logic/auth"
	_ "gbaseadmin/app/system/internal/logic/dept"
	_ "gbaseadmin/app/system/internal/logic/menu"
	_ "gbaseadmin/app/system/internal/logic/role"
	_ "gbaseadmin/app/system/internal/logic/users"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"gbaseadmin/app/system/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
