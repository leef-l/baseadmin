package appointment

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

func csvSafeAppointment(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Appointment = cAppointment{}

type cAppointment struct{}

// Create 创建体验预约
func (c *cAppointment) Create(ctx context.Context, req *v1.AppointmentCreateReq) (res *v1.AppointmentCreateRes, err error) {
	err = service.Appointment().Create(ctx, &model.AppointmentCreateInput{
		AppointmentNo: req.AppointmentNo,
		CustomerID: req.CustomerID,
		Subject: req.Subject,
		AppointmentAt: req.AppointmentAt,
		ContactPhone: req.ContactPhone,
		Address: req.Address,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验预约
func (c *cAppointment) Update(ctx context.Context, req *v1.AppointmentUpdateReq) (res *v1.AppointmentUpdateRes, err error) {
	err = service.Appointment().Update(ctx, &model.AppointmentUpdateInput{
		ID: req.ID,
		AppointmentNo: req.AppointmentNo,
		CustomerID: req.CustomerID,
		Subject: req.Subject,
		AppointmentAt: req.AppointmentAt,
		ContactPhone: req.ContactPhone,
		Address: req.Address,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验预约
func (c *cAppointment) Delete(ctx context.Context, req *v1.AppointmentDeleteReq) (res *v1.AppointmentDeleteRes, err error) {
	err = service.Appointment().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验预约
func (c *cAppointment) BatchDelete(ctx context.Context, req *v1.AppointmentBatchDeleteReq) (res *v1.AppointmentBatchDeleteRes, err error) {
	err = service.Appointment().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验预约
func (c *cAppointment) BatchUpdate(ctx context.Context, req *v1.AppointmentBatchUpdateReq) (res *v1.AppointmentBatchUpdateRes, err error) {
	err = service.Appointment().BatchUpdate(ctx, &model.AppointmentBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取体验预约详情
func (c *cAppointment) Detail(ctx context.Context, req *v1.AppointmentDetailReq) (res *v1.AppointmentDetailRes, err error) {
	res = &v1.AppointmentDetailRes{}
	res.AppointmentDetailOutput, err = service.Appointment().Detail(ctx, req.ID)
	return
}

// List 获取体验预约列表
func (c *cAppointment) List(ctx context.Context, req *v1.AppointmentListReq) (res *v1.AppointmentListRes, err error) {
	res = &v1.AppointmentListRes{}
	res.List, res.Total, err = service.Appointment().List(ctx, &model.AppointmentListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		AppointmentNo: req.AppointmentNo,
		Subject: req.Subject,
		ContactPhone: req.ContactPhone,
		CustomerID: req.CustomerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		AppointmentAtStart: req.AppointmentAtStart,
		AppointmentAtEnd: req.AppointmentAtEnd,
	})
	return
}
// Export 导出体验预约
func (c *cAppointment) Export(ctx context.Context, req *v1.AppointmentExportReq) (res *v1.AppointmentExportRes, err error) {
	list, err := service.Appointment().Export(ctx, &model.AppointmentListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		AppointmentNo: req.AppointmentNo,
		Subject: req.Subject,
		ContactPhone: req.ContactPhone,
		CustomerID: req.CustomerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		AppointmentAtStart: req.AppointmentAtStart,
		AppointmentAtEnd: req.AppointmentAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="appointment.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"预约编号", "客户", "预约主题", "联系电话", "预约地址", "备注", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeAppointment(item.AppointmentNo),
			csvSafeAppointment(item.CustomerName),
			csvSafeAppointment(item.Subject),
			csvSafeAppointment(item.ContactPhone),
			csvSafeAppointment(item.Address),
			csvSafeAppointment(item.Remark),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验预约
func (c *cAppointment) Import(ctx context.Context, req *v1.AppointmentImportReq) (res *v1.AppointmentImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Appointment().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.AppointmentImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验预约导入模板
func (c *cAppointment) ImportTemplate(ctx context.Context, req *v1.AppointmentImportTemplateReq) (res *v1.AppointmentImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="appointment_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"预约编号", "客户", "预约主题", "联系电话", "预约地址", "备注", "状态"})
	w.Flush()
	return
}
