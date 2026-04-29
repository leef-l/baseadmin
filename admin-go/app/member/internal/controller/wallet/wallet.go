package wallet

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

func csvSafeWallet(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Wallet = cWallet{}

type cWallet struct{}

// Create 创建会员钱包
func (c *cWallet) Create(ctx context.Context, req *v1.WalletCreateReq) (res *v1.WalletCreateRes, err error) {
	err = service.Wallet().Create(ctx, &model.WalletCreateInput{
		UserID: req.UserID,
		WalletType: req.WalletType,
		Balance: req.Balance,
		TotalIncome: req.TotalIncome,
		TotalExpense: req.TotalExpense,
		FrozenAmount: req.FrozenAmount,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新会员钱包
func (c *cWallet) Update(ctx context.Context, req *v1.WalletUpdateReq) (res *v1.WalletUpdateRes, err error) {
	err = service.Wallet().Update(ctx, &model.WalletUpdateInput{
		ID: req.ID,
		UserID: req.UserID,
		WalletType: req.WalletType,
		Balance: req.Balance,
		TotalIncome: req.TotalIncome,
		TotalExpense: req.TotalExpense,
		FrozenAmount: req.FrozenAmount,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除会员钱包
func (c *cWallet) Delete(ctx context.Context, req *v1.WalletDeleteReq) (res *v1.WalletDeleteRes, err error) {
	err = service.Wallet().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除会员钱包
func (c *cWallet) BatchDelete(ctx context.Context, req *v1.WalletBatchDeleteReq) (res *v1.WalletBatchDeleteRes, err error) {
	err = service.Wallet().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑会员钱包
func (c *cWallet) BatchUpdate(ctx context.Context, req *v1.WalletBatchUpdateReq) (res *v1.WalletBatchUpdateRes, err error) {
	err = service.Wallet().BatchUpdate(ctx, &model.WalletBatchUpdateInput{
		IDs: req.IDs,
		WalletType: req.WalletType,
		Status: req.Status,
	})
	return
}

// Detail 获取会员钱包详情
func (c *cWallet) Detail(ctx context.Context, req *v1.WalletDetailReq) (res *v1.WalletDetailRes, err error) {
	res = &v1.WalletDetailRes{}
	res.WalletDetailOutput, err = service.Wallet().Detail(ctx, req.ID)
	return
}

// List 获取会员钱包列表
func (c *cWallet) List(ctx context.Context, req *v1.WalletListReq) (res *v1.WalletListRes, err error) {
	res = &v1.WalletListRes{}
	res.List, res.Total, err = service.Wallet().List(ctx, &model.WalletListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		WalletType: req.WalletType,
		Status: req.Status,
	})
	return
}
// Export 导出会员钱包
func (c *cWallet) Export(ctx context.Context, req *v1.WalletExportReq) (res *v1.WalletExportRes, err error) {
	list, err := service.Wallet().Export(ctx, &model.WalletListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		WalletType: req.WalletType,
		Status: req.Status,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="wallet.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"会员", "钱包类型", "当前余额", "累计收入", "累计支出", "冻结金额", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeWallet(item.UserUsername),
			fmt.Sprintf("%v", item.WalletType),
			fmt.Sprintf("%.2f", float64(item.Balance)/100),
			fmt.Sprintf("%.2f", float64(item.TotalIncome)/100),
			fmt.Sprintf("%v", item.TotalExpense),
			fmt.Sprintf("%.2f", float64(item.FrozenAmount)/100),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入会员钱包
func (c *cWallet) Import(ctx context.Context, req *v1.WalletImportReq) (res *v1.WalletImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Wallet().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.WalletImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载会员钱包导入模板
func (c *cWallet) ImportTemplate(ctx context.Context, req *v1.WalletImportTemplateReq) (res *v1.WalletImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="wallet_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"会员", "钱包类型", "当前余额", "累计收入", "累计支出", "冻结金额", "状态"})
	w.Flush()
	return
}
