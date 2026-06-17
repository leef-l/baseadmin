package cmd

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gbaseadmin/app/system/internal/controller/auth"
	"gbaseadmin/app/system/internal/controller/cron"
	"gbaseadmin/app/system/internal/controller/daemon"
	"gbaseadmin/app/system/internal/controller/dept"
	"gbaseadmin/app/system/internal/controller/domain"
	"gbaseadmin/app/system/internal/controller/health"
	"gbaseadmin/app/system/internal/controller/hello"
	"gbaseadmin/app/system/internal/controller/menu"
	"gbaseadmin/app/system/internal/controller/merchant"
	"gbaseadmin/app/system/internal/controller/plan"
	"gbaseadmin/app/system/internal/controller/role"
	"gbaseadmin/app/system/internal/controller/tenant"
	"gbaseadmin/app/system/internal/controller/tenant_plan"
	"gbaseadmin/app/system/internal/controller/users"
	"gbaseadmin/app/system/internal/middleware"
	cronmgr "gbaseadmin/app/system/internal/logic/cron"
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
				group.Middleware(httpmeta.RequestIDMiddleware, middleware.DomainContext, ghttp.MiddlewareHandlerResponse, httpmeta.AccessLogMiddleware)
				group.Bind(
					health.NewV1(),
					hello.NewV1(),
				)
				group.Group("/api/system", func(group *ghttp.RouterGroup) {
					group.Bind(
						auth.Auth.Login,
						auth.Auth.TicketLogin,
					)
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.Bind(
							auth.Auth.Logout,
							auth.Auth.IssueTicket,
							auth.Auth.Info,
							auth.Auth.ChangePassword,
							auth.Auth.Menus,
							daemon.Daemon,
							cron.Cron,
							dept.Dept,
							domain.Domain,
							plan.Plan,
							tenant_plan.TenantPlan,
							tenant.Tenant,
							merchant.Merchant,
							role.Role,
							menu.Menu,
							users.Users,
						)
					})
				})
			})
			// 等服务完全启动后，再启动定时任务管理器
			go func() {
				time.Sleep(3 * time.Second)
				cronmgr.GetManager().Start(ctx)
			}()
			s.Run()
			return nil
		},
	}
)
