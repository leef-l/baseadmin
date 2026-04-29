package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IAppointment interface {
	Create(ctx context.Context, in *model.AppointmentCreateInput) error
	Update(ctx context.Context, in *model.AppointmentUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.AppointmentDetailOutput, err error)
	List(ctx context.Context, in *model.AppointmentListInput) (list []*model.AppointmentListOutput, total int, err error)
	Export(ctx context.Context, in *model.AppointmentListInput) (list []*model.AppointmentListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.AppointmentBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localAppointment IAppointment

func Appointment() IAppointment {
	return localAppointment
}

func RegisterAppointment(i IAppointment) {
	localAppointment = i
}
