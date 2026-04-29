package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gbaseadmin/app/demo/internal/controller/appointment"
	"gbaseadmin/app/demo/internal/controller/audit_log"
	"gbaseadmin/app/demo/internal/controller/campaign"
	"gbaseadmin/app/demo/internal/controller/category"
	"gbaseadmin/app/demo/internal/controller/contract"
	"gbaseadmin/app/demo/internal/controller/customer"
	"gbaseadmin/app/demo/internal/controller/order"
	"gbaseadmin/app/demo/internal/controller/product"
	"gbaseadmin/app/demo/internal/controller/survey"
	"gbaseadmin/app/demo/internal/controller/work_order"

	"gbaseadmin/app/demo/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start demo http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Group("/api/demo", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(
						appointment.Appointment,
						audit_log.AuditLog,
						campaign.Campaign,
						category.Category,
						contract.Contract,
						customer.Customer,
						order.Order,
						product.Product,
						survey.Survey,
						work_order.WorkOrder,
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
