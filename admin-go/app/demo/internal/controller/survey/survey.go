package survey

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

var Survey = cSurvey{}

type cSurvey struct{}

// Create 创建体验问卷
func (c *cSurvey) Create(ctx context.Context, req *v1.SurveyCreateReq) (res *v1.SurveyCreateRes, err error) {
	err = service.Survey().Create(ctx, &model.SurveyCreateInput{
		SurveyNo: req.SurveyNo,
		Title: req.Title,
		Poster: req.Poster,
		QuestionJSON: req.QuestionJSON,
		IntroContent: req.IntroContent,
		PublishAt: req.PublishAt,
		ExpireAt: req.ExpireAt,
		IsAnonymous: req.IsAnonymous,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验问卷
func (c *cSurvey) Update(ctx context.Context, req *v1.SurveyUpdateReq) (res *v1.SurveyUpdateRes, err error) {
	err = service.Survey().Update(ctx, &model.SurveyUpdateInput{
		ID: req.ID,
		SurveyNo: req.SurveyNo,
		Title: req.Title,
		Poster: req.Poster,
		QuestionJSON: req.QuestionJSON,
		IntroContent: req.IntroContent,
		PublishAt: req.PublishAt,
		ExpireAt: req.ExpireAt,
		IsAnonymous: req.IsAnonymous,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验问卷
func (c *cSurvey) Delete(ctx context.Context, req *v1.SurveyDeleteReq) (res *v1.SurveyDeleteRes, err error) {
	err = service.Survey().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验问卷
func (c *cSurvey) BatchDelete(ctx context.Context, req *v1.SurveyBatchDeleteReq) (res *v1.SurveyBatchDeleteRes, err error) {
	err = service.Survey().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验问卷
func (c *cSurvey) BatchUpdate(ctx context.Context, req *v1.SurveyBatchUpdateReq) (res *v1.SurveyBatchUpdateRes, err error) {
	err = service.Survey().BatchUpdate(ctx, &model.SurveyBatchUpdateInput{
		IDs: req.IDs,
		IsAnonymous: req.IsAnonymous,
		Status: req.Status,
	})
	return
}

// Detail 获取体验问卷详情
func (c *cSurvey) Detail(ctx context.Context, req *v1.SurveyDetailReq) (res *v1.SurveyDetailRes, err error) {
	res = &v1.SurveyDetailRes{}
	res.SurveyDetailOutput, err = service.Survey().Detail(ctx, req.ID)
	return
}

// List 获取体验问卷列表
func (c *cSurvey) List(ctx context.Context, req *v1.SurveyListReq) (res *v1.SurveyListRes, err error) {
	res = &v1.SurveyListRes{}
	res.List, res.Total, err = service.Survey().List(ctx, &model.SurveyListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		SurveyNo: req.SurveyNo,
		Title: req.Title,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsAnonymous: req.IsAnonymous,
		Status: req.Status,
		PublishAtStart: req.PublishAtStart,
		PublishAtEnd: req.PublishAtEnd,
		ExpireAtStart: req.ExpireAtStart,
		ExpireAtEnd: req.ExpireAtEnd,
	})
	return
}
// Export 导出体验问卷
func (c *cSurvey) Export(ctx context.Context, req *v1.SurveyExportReq) (res *v1.SurveyExportRes, err error) {
	list, err := service.Survey().Export(ctx, &model.SurveyListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		SurveyNo: req.SurveyNo,
		Title: req.Title,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsAnonymous: req.IsAnonymous,
		Status: req.Status,
		PublishAtStart: req.PublishAtStart,
		PublishAtEnd: req.PublishAtEnd,
		ExpireAtStart: req.ExpireAtStart,
		ExpireAtEnd: req.ExpireAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="survey.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头
	_ = w.Write([]string{"问卷编号", "问卷标题", "海报", "问题JSON", "问卷介绍", "发布时间", "过期时间", "是否匿名", "状态", "租户", "商户", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			item.SurveyNo,
			item.Title,
			item.Poster,
			item.QuestionJSON,
			item.IntroContent,
			func() string { if item.PublishAt != nil { return item.PublishAt.String() }; return "" }(),
			func() string { if item.ExpireAt != nil { return item.ExpireAt.String() }; return "" }(),
			fmt.Sprintf("%v", item.IsAnonymous),
			fmt.Sprintf("%v", item.Status),
			item.TenantName,
			item.MerchantName,
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验问卷
func (c *cSurvey) Import(ctx context.Context, req *v1.SurveyImportReq) (res *v1.SurveyImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Survey().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.SurveyImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验问卷导入模板
func (c *cSurvey) ImportTemplate(ctx context.Context, req *v1.SurveyImportTemplateReq) (res *v1.SurveyImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="survey_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"问卷编号", "问卷标题", "海报", "问题JSON", "问卷介绍", "是否匿名", "状态", "租户", "商户"})
	w.Flush()
	return
}
