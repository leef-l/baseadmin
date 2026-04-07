package main

import (
	_ "gbaseadmin/app/upload/internal/packed"

	_ "gbaseadmin/app/upload/internal/logic/config"
	_ "gbaseadmin/app/upload/internal/logic/dir"
	_ "gbaseadmin/app/upload/internal/logic/dir_rule"
	_ "gbaseadmin/app/upload/internal/logic/file"
	_ "gbaseadmin/app/upload/internal/logic/uploader"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"gbaseadmin/app/upload/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
