package contract

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

func csvSafeContract(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Contract = cContract{}

type cContract struct{}

// Create 创建体验合同
func (c *cContract) Create(ctx context.Context, req *v1.ContractCreateReq) (res *v1.ContractCreateRes, err error) {
	err = service.Contract().Create(ctx, &model.ContractCreateInput{
		ContractNo: req.ContractNo,
		CustomerID: req.CustomerID,
		OrderID: req.OrderID,
		Title: req.Title,
		ContractFile: req.ContractFile,
		SignImage: req.SignImage,
		ContractAmount: req.ContractAmount,
		SignPassword: req.SignPassword,
		SignedAt: req.SignedAt,
		ExpiresAt: req.ExpiresAt,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验合同
func (c *cContract) Update(ctx context.Context, req *v1.ContractUpdateReq) (res *v1.ContractUpdateRes, err error) {
	err = service.Contract().Update(ctx, &model.ContractUpdateInput{
		ID: req.ID,
		ContractNo: req.ContractNo,
		CustomerID: req.CustomerID,
		OrderID: req.OrderID,
		Title: req.Title,
		ContractFile: req.ContractFile,
		SignImage: req.SignImage,
		ContractAmount: req.ContractAmount,
		SignPassword: req.SignPassword,
		SignedAt: req.SignedAt,
		ExpiresAt: req.ExpiresAt,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验合同
func (c *cContract) Delete(ctx context.Context, req *v1.ContractDeleteReq) (res *v1.ContractDeleteRes, err error) {
	err = service.Contract().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验合同
func (c *cContract) BatchDelete(ctx context.Context, req *v1.ContractBatchDeleteReq) (res *v1.ContractBatchDeleteRes, err error) {
	err = service.Contract().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验合同
func (c *cContract) BatchUpdate(ctx context.Context, req *v1.ContractBatchUpdateReq) (res *v1.ContractBatchUpdateRes, err error) {
	err = service.Contract().BatchUpdate(ctx, &model.ContractBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取体验合同详情
func (c *cContract) Detail(ctx context.Context, req *v1.ContractDetailReq) (res *v1.ContractDetailRes, err error) {
	res = &v1.ContractDetailRes{}
	res.ContractDetailOutput, err = service.Contract().Detail(ctx, req.ID)
	return
}

// List 获取体验合同列表
func (c *cContract) List(ctx context.Context, req *v1.ContractListReq) (res *v1.ContractListRes, err error) {
	res = &v1.ContractListRes{}
	res.List, res.Total, err = service.Contract().List(ctx, &model.ContractListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		ContractNo: req.ContractNo,
		Title: req.Title,
		CustomerID: req.CustomerID,
		OrderID: req.OrderID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		SignedAtStart: req.SignedAtStart,
		SignedAtEnd: req.SignedAtEnd,
		ExpiresAtStart: req.ExpiresAtStart,
		ExpiresAtEnd: req.ExpiresAtEnd,
	})
	return
}
// Export 导出体验合同
func (c *cContract) Export(ctx context.Context, req *v1.ContractExportReq) (res *v1.ContractExportRes, err error) {
	list, err := service.Contract().Export(ctx, &model.ContractListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		ContractNo: req.ContractNo,
		Title: req.Title,
		CustomerID: req.CustomerID,
		OrderID: req.OrderID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		SignedAtStart: req.SignedAtStart,
		SignedAtEnd: req.SignedAtEnd,
		ExpiresAtStart: req.ExpiresAtStart,
		ExpiresAtEnd: req.ExpiresAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="contract.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"合同编号", "客户", "订单", "合同标题", "合同文件", "签章图片", "合同金额", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeContract(item.ContractNo),
			csvSafeContract(item.CustomerName),
			csvSafeContract(item.OrderOrderNo),
			csvSafeContract(item.Title),
			csvSafeContract(item.ContractFile),
			csvSafeContract(item.SignImage),
			fmt.Sprintf("%.2f", float64(item.ContractAmount)/100),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验合同
func (c *cContract) Import(ctx context.Context, req *v1.ContractImportReq) (res *v1.ContractImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Contract().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.ContractImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验合同导入模板
func (c *cContract) ImportTemplate(ctx context.Context, req *v1.ContractImportTemplateReq) (res *v1.ContractImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="contract_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"合同编号", "客户", "订单", "合同标题", "合同文件", "签章图片", "合同金额", "状态"})
	w.Flush()
	return
}
