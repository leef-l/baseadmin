package wallet_log

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

func csvSafeWalletLog(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var WalletLog = cWalletLog{}

type cWalletLog struct{}

// Create 创建钱包流水记录
func (c *cWalletLog) Create(ctx context.Context, req *v1.WalletLogCreateReq) (res *v1.WalletLogCreateRes, err error) {
	err = service.WalletLog().Create(ctx, &model.WalletLogCreateInput{
		UserID: req.UserID,
		WalletType: req.WalletType,
		ChangeType: req.ChangeType,
		ChangeAmount: req.ChangeAmount,
		BeforeBalance: req.BeforeBalance,
		AfterBalance: req.AfterBalance,
		RelatedOrderNo: req.RelatedOrderNo,
		Remark: req.Remark,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新钱包流水记录
func (c *cWalletLog) Update(ctx context.Context, req *v1.WalletLogUpdateReq) (res *v1.WalletLogUpdateRes, err error) {
	err = service.WalletLog().Update(ctx, &model.WalletLogUpdateInput{
		ID: req.ID,
		UserID: req.UserID,
		WalletType: req.WalletType,
		ChangeType: req.ChangeType,
		ChangeAmount: req.ChangeAmount,
		BeforeBalance: req.BeforeBalance,
		AfterBalance: req.AfterBalance,
		RelatedOrderNo: req.RelatedOrderNo,
		Remark: req.Remark,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除钱包流水记录
func (c *cWalletLog) Delete(ctx context.Context, req *v1.WalletLogDeleteReq) (res *v1.WalletLogDeleteRes, err error) {
	err = service.WalletLog().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除钱包流水记录
func (c *cWalletLog) BatchDelete(ctx context.Context, req *v1.WalletLogBatchDeleteReq) (res *v1.WalletLogBatchDeleteRes, err error) {
	err = service.WalletLog().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取钱包流水记录详情
func (c *cWalletLog) Detail(ctx context.Context, req *v1.WalletLogDetailReq) (res *v1.WalletLogDetailRes, err error) {
	res = &v1.WalletLogDetailRes{}
	res.WalletLogDetailOutput, err = service.WalletLog().Detail(ctx, req.ID)
	return
}

// List 获取钱包流水记录列表
func (c *cWalletLog) List(ctx context.Context, req *v1.WalletLogListReq) (res *v1.WalletLogListRes, err error) {
	res = &v1.WalletLogListRes{}
	res.List, res.Total, err = service.WalletLog().List(ctx, &model.WalletLogListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		RelatedOrderNo: req.RelatedOrderNo,
		UserID: req.UserID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		WalletType: req.WalletType,
		ChangeType: req.ChangeType,
	})
	return
}
// Export 导出钱包流水记录
func (c *cWalletLog) Export(ctx context.Context, req *v1.WalletLogExportReq) (res *v1.WalletLogExportRes, err error) {
	list, err := service.WalletLog().Export(ctx, &model.WalletLogListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		RelatedOrderNo: req.RelatedOrderNo,
		UserID: req.UserID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		WalletType: req.WalletType,
		ChangeType: req.ChangeType,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="wallet_log.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"会员", "钱包类型", "变动类型", "变动金额", "变动前余额", "变动后余额", "关联单号", "备注说明", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeWalletLog(item.UserNickname),
			fmt.Sprintf("%v", item.WalletType),
			fmt.Sprintf("%v", item.ChangeType),
			fmt.Sprintf("%.2f", float64(item.ChangeAmount)/100),
			fmt.Sprintf("%.2f", float64(item.BeforeBalance)/100),
			fmt.Sprintf("%.2f", float64(item.AfterBalance)/100),
			csvSafeWalletLog(item.RelatedOrderNo),
			csvSafeWalletLog(item.Remark),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入钱包流水记录
func (c *cWalletLog) Import(ctx context.Context, req *v1.WalletLogImportReq) (res *v1.WalletLogImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.WalletLog().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.WalletLogImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载钱包流水记录导入模板
func (c *cWalletLog) ImportTemplate(ctx context.Context, req *v1.WalletLogImportTemplateReq) (res *v1.WalletLogImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="wallet_log_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"会员", "钱包类型", "变动类型", "变动金额", "变动前余额", "变动后余额", "关联单号", "备注说明"})
	w.Flush()
	return
}
