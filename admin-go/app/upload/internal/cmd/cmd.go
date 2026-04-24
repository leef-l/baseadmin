package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gbaseadmin/app/upload/internal/controller/config"
	"gbaseadmin/app/upload/internal/controller/dir"
	"gbaseadmin/app/upload/internal/controller/dir_rule"
	"gbaseadmin/app/upload/internal/controller/file"
	"gbaseadmin/app/upload/internal/controller/health"
	"gbaseadmin/app/upload/internal/controller/uploader"

	"gbaseadmin/app/upload/internal/middleware"
	"gbaseadmin/utility/httpmeta"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start upload http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(httpmeta.RequestIDMiddleware, ghttp.MiddlewareHandlerResponse, httpmeta.AccessLogMiddleware)
				group.Bind(
					health.NewV1(),
				)
				group.Group("/api/upload", func(group *ghttp.RouterGroup) {
					group.Bind(
						uploader.PublicUploader,
					)
				})
				group.Group("/api/upload", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(
						config.Config,
						dir.Dir,
						dir_rule.DirRule,
						file.File,
						uploader.Uploader,
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
