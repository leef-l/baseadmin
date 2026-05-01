// Package contract 合同模板渲染 + 签署 + 异步生成。
//
// 当前 PDF 生成方案：把签署后的 HTML（含手写签名 PNG）落盘到 contractStorageDir。
// 客户端打开后可用浏览器"打印 → 另存为 PDF"或"分享"得到正式 PDF；
// 服务端只需保存 HTML，下载接口直接流式回吐文件。
//
// 后续如需服务端直接生成 PDF，可切到 chromedp 或 wkhtmltopdf，
// 仅需替换 generatePDF() 内部实现，对接调用方零变更。
package contract

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/snowflake"
)

// PDF 状态常量。
const (
	PDFStatusInit    = 0
	PDFStatusPending = 1
	PDFStatusReady   = 2
	PDFStatusFailed  = 3
)

// SignInput 提交签名入参。
type SignInput struct {
	UserID         int64
	ContractType   string // register/upgrade/custom
	TemplateID     int64  // 0 表示自动选择对应类型的默认模板
	SignatureImage string // data:image/png;base64,...
	UserNickname   string
	UserPhone      string
	IP             string
	UA             string
	RelatedID      int64
}

// SignResult 签署结果。
type SignResult struct {
	ContractID string
	ContractNo string
}

// Sign 提交手写签名 + 异步生成 PDF 文件。
func Sign(ctx context.Context, in *SignInput) (*SignResult, error) {
	if in == nil || in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	if !strings.HasPrefix(in.SignatureImage, "data:image/png;base64,") {
		return nil, gerror.New("签名图片格式不正确")
	}
	if len(in.SignatureImage) < 200 {
		return nil, gerror.New("签名为空，请书写后再提交")
	}
	contractType := strings.TrimSpace(in.ContractType)
	if contractType == "" {
		contractType = "register"
	}

	tpl, err := pickTemplate(ctx, in.TemplateID, contractType)
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return nil, gerror.New("未找到合同模板")
	}

	now := time.Now()
	contractID := snowflake.Generate()
	contractNo := fmt.Sprintf("C%d", int64(contractID))

	rendered := renderTemplate(tpl.Content, map[string]string{
		"nickname": fallback(in.UserNickname, "会员"),
		"phone":    in.UserPhone,
		"date":     now.Format("2006-01-02"),
	})
	signedHTML := wrapSignedHTML(tpl.TemplateName, rendered, in.SignatureImage, contractNo, in.UserNickname, in.UserPhone, now)

	row := do.MemberContract{
		Id:              int64(contractID),
		UserId:          in.UserID,
		ContractNo:      contractNo,
		ContractType:    contractType,
		TemplateId:      tpl.Id,
		RelatedId:       in.RelatedID,
		SignedHtml:      signedHTML,
		SignatureImage:  in.SignatureImage,
		SignedAt:        gtime.New(now),
		SignedIp:        in.IP,
		SignedUserAgent: in.UA,
		PdfStatus:       PDFStatusPending,
		Status:          1,
	}
	if _, err := dao.MemberContract.Ctx(ctx).Data(row).Insert(); err != nil {
		return nil, err
	}

	go doGeneratePDFAsync(int64(contractID), signedHTML, contractNo)

	return &SignResult{ContractID: fmt.Sprintf("%d", int64(contractID)), ContractNo: contractNo}, nil
}

// pickTemplate 选择模板：指定 ID > 类型默认 > 类型最新启用。
func pickTemplate(ctx context.Context, templateID int64, contractType string) (*entity.MemberContractTemplate, error) {
	cols := dao.MemberContractTemplate.Columns()
	if templateID > 0 {
		var tpl entity.MemberContractTemplate
		if err := dao.MemberContractTemplate.Ctx(ctx).
			Where(cols.Id, templateID).
			Where(cols.DeletedAt, nil).
			Where(cols.Status, 1).
			Scan(&tpl); err != nil {
			return nil, err
		}
		return &tpl, nil
	}
	var tpl entity.MemberContractTemplate
	if err := dao.MemberContractTemplate.Ctx(ctx).
		Where(cols.TemplateType, contractType).
		Where(cols.DeletedAt, nil).
		Where(cols.Status, 1).
		OrderDesc(cols.IsDefault).
		OrderDesc(cols.Sort).
		OrderDesc(cols.Id).
		Scan(&tpl); err != nil {
		return nil, err
	}
	return &tpl, nil
}

// GetSignedHTML 返回合同 signed_html（用于 download / 预览）。
func GetSignedHTML(ctx context.Context, userID int64, contractID int64) (htmlContent string, contractNo string, err error) {
	cols := dao.MemberContract.Columns()
	var c entity.MemberContract
	if err = dao.MemberContract.Ctx(ctx).
		Where(cols.Id, contractID).
		Where(cols.UserId, userID).
		Where(cols.DeletedAt, nil).
		Scan(&c); err != nil {
		return "", "", err
	}
	if c.Id == 0 {
		return "", "", gerror.New("合同不存在")
	}
	return c.SignedHtml, c.ContractNo, nil
}

// HasUserSigned 用户是否已签某类型合同（用于"未签"提示）。
func HasUserSigned(ctx context.Context, userID int64, contractType string) (bool, error) {
	cols := dao.MemberContract.Columns()
	count, err := dao.MemberContract.Ctx(ctx).
		Where(cols.UserId, userID).
		Where(cols.ContractType, contractType).
		Where(cols.Status, 1).
		Where(cols.DeletedAt, nil).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PickDefaultTemplate 公开接口：返回类型对应的默认模板（用于前端展示协议正文）。
func PickDefaultTemplate(ctx context.Context, contractType string) (*entity.MemberContractTemplate, error) {
	if contractType == "" {
		contractType = "register"
	}
	return pickTemplate(ctx, 0, contractType)
}

// renderTemplate 简易占位符替换 {{key}} → value。
func renderTemplate(tpl string, vars map[string]string) string {
	out := tpl
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
	}
	return out
}

func fallback(s, d string) string {
	if strings.TrimSpace(s) == "" {
		return d
	}
	return s
}

// wrapSignedHTML 把模板渲染后包成完整 HTML 文档（带签名章块）。
func wrapSignedHTML(title, body, signatureDataURL, contractNo, nickname, phone string, now time.Time) string {
	style := `
<style>
  body { font-family: "PingFang SC","Microsoft YaHei",sans-serif; max-width:720px;margin:24px auto;padding:0 16px;color:#333;line-height:1.7;background:#fff }
  h1 { text-align:center;font-size:24px;margin-bottom:8px }
  .meta { text-align:center;color:#888;font-size:12px;margin-bottom:24px }
  .sign-block { margin-top:48px;border-top:1px dashed #999;padding-top:16px;display:flex;justify-content:space-between;align-items:flex-end }
  .sign-block img { width:200px;height:auto;border-bottom:1px solid #333 }
  @media print { body{margin:0} }
</style>`
	return fmt.Sprintf(
		`<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="UTF-8"><title>%s - %s</title>%s</head>
<body>
%s
<div class="meta">合同编号：%s · 生成时间：%s</div>
<div class="sign-block">
  <div>
    <div style="font-size:14px;color:#555">乙方签字（手写）：</div>
    <img src="%s" alt="signature" />
  </div>
  <div style="text-align:right;font-size:13px;color:#555">
    <div>会员：%s</div>
    <div>手机号：%s</div>
    <div>签署日期：%s</div>
  </div>
</div>
</body>
</html>`,
		htmlEscape(title), contractNo, style,
		body,
		contractNo, now.Format("2006-01-02 15:04:05"),
		signatureDataURL,
		htmlEscape(nickname), htmlEscape(phone), now.Format("2006-01-02"),
	)
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}

// doGeneratePDFAsync 异步把 signed_html 写到磁盘文件。
// 当前实现：写为 .html 文件，前端下载后用浏览器查看 / 打印为 PDF。
// 占位符：未来可替换为 chromedp/wkhtmltopdf 真实 PDF 生成。
func doGeneratePDFAsync(contractID int64, signedHTML, contractNo string) {
	ctx := context.Background()
	dir := storageDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		markPDFFailed(ctx, contractID, err)
		return
	}
	filename := contractNo + ".html"
	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, []byte(signedHTML), 0o644); err != nil {
		markPDFFailed(ctx, contractID, err)
		return
	}
	cols := dao.MemberContract.Columns()
	if _, err := dao.MemberContract.Ctx(ctx).
		Where(cols.Id, contractID).
		Data(g.Map{
			cols.PdfStatus: PDFStatusReady,
			cols.PdfPath:   fullPath,
			cols.PdfError:  "",
		}).Update(); err != nil {
		g.Log().Errorf(ctx, "[contract] update pdf path err id=%d err=%v", contractID, err)
	}
}

func markPDFFailed(ctx context.Context, contractID int64, err error) {
	cols := dao.MemberContract.Columns()
	_, _ = dao.MemberContract.Ctx(ctx).
		Where(cols.Id, contractID).
		Data(g.Map{
			cols.PdfStatus: PDFStatusFailed,
			cols.PdfError:  err.Error(),
		}).Update()
	g.Log().Errorf(ctx, "[contract] generate pdf failed id=%d err=%v", contractID, err)
}

func storageDir() string {
	dir := g.Cfg().MustGet(context.Background(), "member.contractStorageDir").String()
	if strings.TrimSpace(dir) == "" {
		dir = "./runtime/contracts"
	}
	return dir
}
