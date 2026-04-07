package main

import (
	_ "gbaseadmin/app/{{.AppName}}/internal/packed"
{{range .Modules}}
	_ "gbaseadmin/app/{{$.AppName}}/internal/logic/{{.}}"{{end}}

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"gbaseadmin/app/{{.AppName}}/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
