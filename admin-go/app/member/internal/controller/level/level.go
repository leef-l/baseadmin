package level

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

func csvSafeLevel(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Level = cLevel{}

type cLevel struct{}

// Create 创建会员等级配置
func (c *cLevel) Create(ctx context.Context, req *v1.LevelCreateReq) (res *v1.LevelCreateRes, err error) {
	err = service.Level().Create(ctx, &model.LevelCreateInput{
		Name: req.Name,
		LevelNo: req.LevelNo,
		Icon: req.Icon,
		DurationDays: req.DurationDays,
		NeedActiveCount: req.NeedActiveCount,
		NeedTeamTurnover: req.NeedTeamTurnover,
		IsTop: req.IsTop,
		AutoDeploy: req.AutoDeploy,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新会员等级配置
func (c *cLevel) Update(ctx context.Context, req *v1.LevelUpdateReq) (res *v1.LevelUpdateRes, err error) {
	err = service.Level().Update(ctx, &model.LevelUpdateInput{
		ID: req.ID,
		Name: req.Name,
		LevelNo: req.LevelNo,
		Icon: req.Icon,
		DurationDays: req.DurationDays,
		NeedActiveCount: req.NeedActiveCount,
		NeedTeamTurnover: req.NeedTeamTurnover,
		IsTop: req.IsTop,
		AutoDeploy: req.AutoDeploy,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除会员等级配置
func (c *cLevel) Delete(ctx context.Context, req *v1.LevelDeleteReq) (res *v1.LevelDeleteRes, err error) {
	err = service.Level().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除会员等级配置
func (c *cLevel) BatchDelete(ctx context.Context, req *v1.LevelBatchDeleteReq) (res *v1.LevelBatchDeleteRes, err error) {
	err = service.Level().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑会员等级配置
func (c *cLevel) BatchUpdate(ctx context.Context, req *v1.LevelBatchUpdateReq) (res *v1.LevelBatchUpdateRes, err error) {
	err = service.Level().BatchUpdate(ctx, &model.LevelBatchUpdateInput{
		IDs: req.IDs,
		IsTop: req.IsTop,
		AutoDeploy: req.AutoDeploy,
		Status: req.Status,
	})
	return
}

// Detail 获取会员等级配置详情
func (c *cLevel) Detail(ctx context.Context, req *v1.LevelDetailReq) (res *v1.LevelDetailRes, err error) {
	res = &v1.LevelDetailRes{}
	res.LevelDetailOutput, err = service.Level().Detail(ctx, req.ID)
	return
}

// List 获取会员等级配置列表
func (c *cLevel) List(ctx context.Context, req *v1.LevelListReq) (res *v1.LevelListRes, err error) {
	res = &v1.LevelListRes{}
	res.List, res.Total, err = service.Level().List(ctx, &model.LevelListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Name: req.Name,
		LevelNo: req.LevelNo,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsTop: req.IsTop,
		AutoDeploy: req.AutoDeploy,
		Status: req.Status,
	})
	return
}
// Export 导出会员等级配置
func (c *cLevel) Export(ctx context.Context, req *v1.LevelExportReq) (res *v1.LevelExportRes, err error) {
	list, err := service.Level().Export(ctx, &model.LevelListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Name: req.Name,
		LevelNo: req.LevelNo,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsTop: req.IsTop,
		AutoDeploy: req.AutoDeploy,
		Status: req.Status,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="level.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"等级名称", "等级编号", "等级图标", "有效天数", "升级所需有效用户数", "升级所需团队营业额", "是否最高等级", "到达后自动部署站点", "等级说明", "排序", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeLevel(item.Name),
			fmt.Sprintf("%v", item.LevelNo),
			csvSafeLevel(item.Icon),
			fmt.Sprintf("%v", item.DurationDays),
			fmt.Sprintf("%v", item.NeedActiveCount),
			fmt.Sprintf("%.2f", float64(item.NeedTeamTurnover)/100),
			fmt.Sprintf("%v", item.IsTop),
			fmt.Sprintf("%v", item.AutoDeploy),
			csvSafeLevel(item.Remark),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入会员等级配置
func (c *cLevel) Import(ctx context.Context, req *v1.LevelImportReq) (res *v1.LevelImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Level().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.LevelImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载会员等级配置导入模板
func (c *cLevel) ImportTemplate(ctx context.Context, req *v1.LevelImportTemplateReq) (res *v1.LevelImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="level_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"等级名称", "等级编号", "等级图标", "有效天数", "升级所需有效用户数", "升级所需团队营业额", "是否最高等级", "到达后自动部署站点", "等级说明", "排序", "状态"})
	w.Flush()
	return
}
