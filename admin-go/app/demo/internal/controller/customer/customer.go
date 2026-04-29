package customer

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

func csvSafeCustomer(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Customer = cCustomer{}

type cCustomer struct{}

// Create 创建体验客户
func (c *cCustomer) Create(ctx context.Context, req *v1.CustomerCreateReq) (res *v1.CustomerCreateRes, err error) {
	err = service.Customer().Create(ctx, &model.CustomerCreateInput{
		Avatar: req.Avatar,
		Name: req.Name,
		CustomerNo: req.CustomerNo,
		Phone: req.Phone,
		Email: req.Email,
		Gender: req.Gender,
		Level: req.Level,
		SourceType: req.SourceType,
		IsVip: req.IsVip,
		RegisteredAt: req.RegisteredAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验客户
func (c *cCustomer) Update(ctx context.Context, req *v1.CustomerUpdateReq) (res *v1.CustomerUpdateRes, err error) {
	err = service.Customer().Update(ctx, &model.CustomerUpdateInput{
		ID: req.ID,
		Avatar: req.Avatar,
		Name: req.Name,
		CustomerNo: req.CustomerNo,
		Phone: req.Phone,
		Email: req.Email,
		Gender: req.Gender,
		Level: req.Level,
		SourceType: req.SourceType,
		IsVip: req.IsVip,
		RegisteredAt: req.RegisteredAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验客户
func (c *cCustomer) Delete(ctx context.Context, req *v1.CustomerDeleteReq) (res *v1.CustomerDeleteRes, err error) {
	err = service.Customer().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验客户
func (c *cCustomer) BatchDelete(ctx context.Context, req *v1.CustomerBatchDeleteReq) (res *v1.CustomerBatchDeleteRes, err error) {
	err = service.Customer().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验客户
func (c *cCustomer) BatchUpdate(ctx context.Context, req *v1.CustomerBatchUpdateReq) (res *v1.CustomerBatchUpdateRes, err error) {
	err = service.Customer().BatchUpdate(ctx, &model.CustomerBatchUpdateInput{
		IDs: req.IDs,
		Gender: req.Gender,
		Level: req.Level,
		SourceType: req.SourceType,
		IsVip: req.IsVip,
		Status: req.Status,
	})
	return
}

// Detail 获取体验客户详情
func (c *cCustomer) Detail(ctx context.Context, req *v1.CustomerDetailReq) (res *v1.CustomerDetailRes, err error) {
	res = &v1.CustomerDetailRes{}
	res.CustomerDetailOutput, err = service.Customer().Detail(ctx, req.ID)
	return
}

// List 获取体验客户列表
func (c *cCustomer) List(ctx context.Context, req *v1.CustomerListReq) (res *v1.CustomerListRes, err error) {
	res = &v1.CustomerListRes{}
	res.List, res.Total, err = service.Customer().List(ctx, &model.CustomerListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		CustomerNo: req.CustomerNo,
		Name: req.Name,
		Phone: req.Phone,
		Email: req.Email,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Gender: req.Gender,
		Level: req.Level,
		SourceType: req.SourceType,
		IsVip: req.IsVip,
		Status: req.Status,
		RegisteredAtStart: req.RegisteredAtStart,
		RegisteredAtEnd: req.RegisteredAtEnd,
	})
	return
}
// Export 导出体验客户
func (c *cCustomer) Export(ctx context.Context, req *v1.CustomerExportReq) (res *v1.CustomerExportRes, err error) {
	list, err := service.Customer().Export(ctx, &model.CustomerListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		CustomerNo: req.CustomerNo,
		Name: req.Name,
		Phone: req.Phone,
		Email: req.Email,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Gender: req.Gender,
		Level: req.Level,
		SourceType: req.SourceType,
		IsVip: req.IsVip,
		Status: req.Status,
		RegisteredAtStart: req.RegisteredAtStart,
		RegisteredAtEnd: req.RegisteredAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="customer.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"头像", "客户名称", "客户编号", "联系电话", "邮箱", "性别", "等级", "来源", "是否VIP", "备注", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeCustomer(item.Avatar),
			csvSafeCustomer(item.Name),
			csvSafeCustomer(item.CustomerNo),
			csvSafeCustomer(item.Phone),
			csvSafeCustomer(item.Email),
			fmt.Sprintf("%v", item.Gender),
			fmt.Sprintf("%v", item.Level),
			fmt.Sprintf("%v", item.SourceType),
			fmt.Sprintf("%v", item.IsVip),
			csvSafeCustomer(item.Remark),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验客户
func (c *cCustomer) Import(ctx context.Context, req *v1.CustomerImportReq) (res *v1.CustomerImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Customer().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.CustomerImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验客户导入模板
func (c *cCustomer) ImportTemplate(ctx context.Context, req *v1.CustomerImportTemplateReq) (res *v1.CustomerImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="customer_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"头像", "客户名称", "客户编号", "联系电话", "邮箱", "性别", "等级", "来源", "是否VIP", "备注", "状态"})
	w.Flush()
	return
}
