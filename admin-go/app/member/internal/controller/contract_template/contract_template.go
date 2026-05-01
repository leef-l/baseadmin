package contract_template

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/service"
)

func csvSafeContractTemplate(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var ContractTemplate = cContractTemplate{}

type cContractTemplate struct{}

// Create 创建会员合同模板
func (c *cContractTemplate) Create(ctx context.Context, req *v1.ContractTemplateCreateReq) (res *v1.ContractTemplateCreateRes, err error) {
	err = service.ContractTemplate().Create(ctx, &model.ContractTemplateCreateInput{
		TemplateName: req.TemplateName,
		TemplateType: req.TemplateType,
		Content: req.Content,
		IsDefault: req.IsDefault,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新会员合同模板
func (c *cContractTemplate) Update(ctx context.Context, req *v1.ContractTemplateUpdateReq) (res *v1.ContractTemplateUpdateRes, err error) {
	err = service.ContractTemplate().Update(ctx, &model.ContractTemplateUpdateInput{
		ID: req.ID,
		TemplateName: req.TemplateName,
		TemplateType: req.TemplateType,
		Content: req.Content,
		IsDefault: req.IsDefault,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除会员合同模板
func (c *cContractTemplate) Delete(ctx context.Context, req *v1.ContractTemplateDeleteReq) (res *v1.ContractTemplateDeleteRes, err error) {
	err = service.ContractTemplate().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除会员合同模板
func (c *cContractTemplate) BatchDelete(ctx context.Context, req *v1.ContractTemplateBatchDeleteReq) (res *v1.ContractTemplateBatchDeleteRes, err error) {
	err = service.ContractTemplate().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑会员合同模板
func (c *cContractTemplate) BatchUpdate(ctx context.Context, req *v1.ContractTemplateBatchUpdateReq) (res *v1.ContractTemplateBatchUpdateRes, err error) {
	err = service.ContractTemplate().BatchUpdate(ctx, &model.ContractTemplateBatchUpdateInput{
		IDs: req.IDs,
		IsDefault: req.IsDefault,
		Status: req.Status,
	})
	return
}

// Detail 获取会员合同模板详情
func (c *cContractTemplate) Detail(ctx context.Context, req *v1.ContractTemplateDetailReq) (res *v1.ContractTemplateDetailRes, err error) {
	res = &v1.ContractTemplateDetailRes{}
	res.ContractTemplateDetailOutput, err = service.ContractTemplate().Detail(ctx, req.ID)
	return
}

// List 获取会员合同模板列表
func (c *cContractTemplate) List(ctx context.Context, req *v1.ContractTemplateListReq) (res *v1.ContractTemplateListRes, err error) {
	res = &v1.ContractTemplateListRes{}
	res.List, res.Total, err = service.ContractTemplate().List(ctx, &model.ContractTemplateListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TemplateName: req.TemplateName,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsDefault: req.IsDefault,
		Status: req.Status,
		TemplateType: req.TemplateType,
	})
	return
}
// Export 导出会员合同模板
func (c *cContractTemplate) Export(ctx context.Context, req *v1.ContractTemplateExportReq) (res *v1.ContractTemplateExportRes, err error) {
	list, err := service.ContractTemplate().Export(ctx, &model.ContractTemplateListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TemplateName: req.TemplateName,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsDefault: req.IsDefault,
		Status: req.Status,
		TemplateType: req.TemplateType,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="contract_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"模板名称", "模板类型", "模板正文", "是否默认模板", "备注", "排序", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeContractTemplate(item.TemplateName),
			csvSafeContractTemplate(item.TemplateType),
			csvSafeContractTemplate(item.Content),
			fmt.Sprintf("%v", item.IsDefault),
			csvSafeContractTemplate(item.Remark),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入会员合同模板
func (c *cContractTemplate) Import(ctx context.Context, req *v1.ContractTemplateImportReq) (res *v1.ContractTemplateImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.ContractTemplate().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.ContractTemplateImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载会员合同模板导入模板
func (c *cContractTemplate) ImportTemplate(ctx context.Context, req *v1.ContractTemplateImportTemplateReq) (res *v1.ContractTemplateImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="contract_template_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"模板名称", "模板类型", "模板正文", "是否默认模板", "备注", "排序", "状态"})
	w.Flush()
	return
}
