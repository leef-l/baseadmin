package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gbaseadmin/app/system/internal/controller/auth"
	"gbaseadmin/app/system/internal/controller/dept"
	"gbaseadmin/app/system/internal/controller/health"
	"gbaseadmin/app/system/internal/controller/hello"
	"gbaseadmin/app/system/internal/controller/menu"
	"gbaseadmin/app/system/internal/controller/role"
	"gbaseadmin/app/system/internal/controller/users"
	"gbaseadmin/app/system/internal/middleware"
	"gbaseadmin/utility/httpmeta"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(httpmeta.RequestIDMiddleware, ghttp.MiddlewareHandlerResponse, httpmeta.AccessLogMiddleware)
				group.Bind(
					health.NewV1(),
					hello.NewV1(),
				)
				// 系统管理模块
				group.Group("/api/system", func(group *ghttp.RouterGroup) {
					// 公开接口（无需登录）
					group.Bind(
						auth.Auth.Login,
					)
					// 需要登录的接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.Bind(
							auth.Auth.Info,
							auth.Auth.ChangePassword,
							auth.Auth.Menus,
							dept.Dept,
							role.Role,
							menu.Menu,
							users.Users,
						)
					})
				})
			})
			s.Run()
			return nil
		},
	}
)
