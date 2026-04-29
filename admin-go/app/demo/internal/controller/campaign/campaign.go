package campaign

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

func csvSafeCampaign(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Campaign = cCampaign{}

type cCampaign struct{}

// Create 创建体验活动
func (c *cCampaign) Create(ctx context.Context, req *v1.CampaignCreateReq) (res *v1.CampaignCreateRes, err error) {
	err = service.Campaign().Create(ctx, &model.CampaignCreateInput{
		CampaignNo: req.CampaignNo,
		Title: req.Title,
		Banner: req.Banner,
		Type: req.Type,
		Channel: req.Channel,
		BudgetAmount: req.BudgetAmount,
		LandingURL: req.LandingURL,
		RuleJSON: req.RuleJSON,
		IntroContent: req.IntroContent,
		StartAt: req.StartAt,
		EndAt: req.EndAt,
		IsPublic: req.IsPublic,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验活动
func (c *cCampaign) Update(ctx context.Context, req *v1.CampaignUpdateReq) (res *v1.CampaignUpdateRes, err error) {
	err = service.Campaign().Update(ctx, &model.CampaignUpdateInput{
		ID: req.ID,
		CampaignNo: req.CampaignNo,
		Title: req.Title,
		Banner: req.Banner,
		Type: req.Type,
		Channel: req.Channel,
		BudgetAmount: req.BudgetAmount,
		LandingURL: req.LandingURL,
		RuleJSON: req.RuleJSON,
		IntroContent: req.IntroContent,
		StartAt: req.StartAt,
		EndAt: req.EndAt,
		IsPublic: req.IsPublic,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验活动
func (c *cCampaign) Delete(ctx context.Context, req *v1.CampaignDeleteReq) (res *v1.CampaignDeleteRes, err error) {
	err = service.Campaign().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验活动
func (c *cCampaign) BatchDelete(ctx context.Context, req *v1.CampaignBatchDeleteReq) (res *v1.CampaignBatchDeleteRes, err error) {
	err = service.Campaign().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验活动
func (c *cCampaign) BatchUpdate(ctx context.Context, req *v1.CampaignBatchUpdateReq) (res *v1.CampaignBatchUpdateRes, err error) {
	err = service.Campaign().BatchUpdate(ctx, &model.CampaignBatchUpdateInput{
		IDs: req.IDs,
		Type: req.Type,
		Channel: req.Channel,
		IsPublic: req.IsPublic,
		Status: req.Status,
	})
	return
}

// Detail 获取体验活动详情
func (c *cCampaign) Detail(ctx context.Context, req *v1.CampaignDetailReq) (res *v1.CampaignDetailRes, err error) {
	res = &v1.CampaignDetailRes{}
	res.CampaignDetailOutput, err = service.Campaign().Detail(ctx, req.ID)
	return
}

// List 获取体验活动列表
func (c *cCampaign) List(ctx context.Context, req *v1.CampaignListReq) (res *v1.CampaignListRes, err error) {
	res = &v1.CampaignListRes{}
	res.List, res.Total, err = service.Campaign().List(ctx, &model.CampaignListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		CampaignNo: req.CampaignNo,
		Title: req.Title,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Type: req.Type,
		Channel: req.Channel,
		IsPublic: req.IsPublic,
		Status: req.Status,
		StartAtStart: req.StartAtStart,
		StartAtEnd: req.StartAtEnd,
		EndAtStart: req.EndAtStart,
		EndAtEnd: req.EndAtEnd,
	})
	return
}
// Export 导出体验活动
func (c *cCampaign) Export(ctx context.Context, req *v1.CampaignExportReq) (res *v1.CampaignExportRes, err error) {
	list, err := service.Campaign().Export(ctx, &model.CampaignListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		CampaignNo: req.CampaignNo,
		Title: req.Title,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Type: req.Type,
		Channel: req.Channel,
		IsPublic: req.IsPublic,
		Status: req.Status,
		StartAtStart: req.StartAtStart,
		StartAtEnd: req.StartAtEnd,
		EndAtStart: req.EndAtStart,
		EndAtEnd: req.EndAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="campaign.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"活动编号", "活动标题", "横幅图", "活动类型", "投放渠道", "预算金额", "落地页URL", "规则JSON", "活动介绍", "是否公开", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeCampaign(item.CampaignNo),
			csvSafeCampaign(item.Title),
			csvSafeCampaign(item.Banner),
			fmt.Sprintf("%v", item.Type),
			fmt.Sprintf("%v", item.Channel),
			fmt.Sprintf("%.2f", float64(item.BudgetAmount)/100),
			csvSafeCampaign(item.LandingURL),
			csvSafeCampaign(item.RuleJSON),
			csvSafeCampaign(item.IntroContent),
			fmt.Sprintf("%v", item.IsPublic),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验活动
func (c *cCampaign) Import(ctx context.Context, req *v1.CampaignImportReq) (res *v1.CampaignImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Campaign().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.CampaignImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验活动导入模板
func (c *cCampaign) ImportTemplate(ctx context.Context, req *v1.CampaignImportTemplateReq) (res *v1.CampaignImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="campaign_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"活动编号", "活动标题", "横幅图", "活动类型", "投放渠道", "预算金额", "落地页URL", "规则JSON", "活动介绍", "是否公开", "状态"})
	w.Flush()
	return
}
