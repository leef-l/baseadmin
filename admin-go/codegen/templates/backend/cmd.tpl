package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
{{range .Modules}}
	"gbaseadmin/app/{{$.AppName}}/internal/controller/{{.}}"{{end}}

	"gbaseadmin/app/{{.AppName}}/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start {{.AppName}} http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Group("/api/{{.AppName}}", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind({{range .Modules}}
						{{.}}.{{ModuleCamel .}},{{end}}
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
