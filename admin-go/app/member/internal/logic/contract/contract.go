// Package contract 合同模板渲染 + 签署 + 异步生成 PDF。
//
// PDF 生成方案：用 headless Chromium 命令行（--headless --print-to-pdf）将签署后的 HTML
// 渲染为 PDF。Chromium 二进制路径默认 ungoogled-chromium / chromium / google-chrome 自动探测，
// 也可通过 member.contractChromiumPath 显式指定。
//
// 渲染失败时降级保存 HTML 文件并标记 PDF 状态为失败，前端可仍下载 HTML（浏览器打印 → PDF）。
package contract

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

// GetDownload 返回合同下载所需信息（优先 pdf_path 是 .pdf 时返回 PDF 流；否则返回 signed_html）。
//
// userID 传 0 表示后台调用（不校验归属）。
// 返回值：filePath（磁盘路径，可能 .pdf 或 .html）、contentBytes（filePath 为空时使用）、contractNo、isPDF、err。
func GetDownload(ctx context.Context, userID, contractID int64) (filePath string, contentBytes []byte, contractNo string, isPDF bool, err error) {
	cols := dao.MemberContract.Columns()
	q := dao.MemberContract.Ctx(ctx).
		Where(cols.Id, contractID).
		Where(cols.DeletedAt, nil)
	if userID > 0 {
		q = q.Where(cols.UserId, userID)
	}
	var c entity.MemberContract
	if err = q.Scan(&c); err != nil {
		return "", nil, "", false, err
	}
	if c.Id == 0 {
		return "", nil, "", false, gerror.New("合同不存在")
	}
	contractNo = c.ContractNo
	// PDF 已就绪
	if strings.HasSuffix(c.PdfPath, ".pdf") {
		if _, statErr := os.Stat(c.PdfPath); statErr == nil {
			return c.PdfPath, nil, contractNo, true, nil
		}
	}
	// 否则降级到 HTML 文件 / signed_html
	if strings.HasSuffix(c.PdfPath, ".html") {
		if _, statErr := os.Stat(c.PdfPath); statErr == nil {
			return c.PdfPath, nil, contractNo, false, nil
		}
	}
	return "", []byte(c.SignedHtml), contractNo, false, nil
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

// doGeneratePDFAsync 异步用 headless Chromium 把 signed_html 渲染成 PDF。
//
// 流程：
//  1. 把 signed_html 写到 storageDir/{contractNo}.html
//  2. 调 chromium --headless --print-to-pdf={contractNo}.pdf 渲染
//  3. 渲染成功 → pdf_path 指向 .pdf；失败 → pdf_path 指向 .html，标记 PDFStatusFailed
//
// Chromium 命令使用 --no-pdf-header-footer 去掉默认页眉页脚；--virtual-time-budget 等待资源加载；
// --no-sandbox 让 root 可运行（生产环境若不用 root 可移除）。
func doGeneratePDFAsync(contractID int64, signedHTML, contractNo string) {
	ctx := context.Background()
	dir := storageDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		markPDFFailed(ctx, contractID, err)
		return
	}
	htmlPath := filepath.Join(dir, contractNo+".html")
	if err := os.WriteFile(htmlPath, []byte(signedHTML), 0o644); err != nil {
		markPDFFailed(ctx, contractID, err)
		return
	}

	cols := dao.MemberContract.Columns()
	pdfPath := filepath.Join(dir, contractNo+".pdf")

	if err := renderHTMLToPDF(ctx, htmlPath, pdfPath); err != nil {
		// PDF 渲染失败，保留 HTML 作为兜底
		_, _ = dao.MemberContract.Ctx(ctx).
			Where(cols.Id, contractID).
			Data(g.Map{
				cols.PdfStatus: PDFStatusFailed,
				cols.PdfPath:   htmlPath,
				cols.PdfError:  "PDF 渲染失败：" + err.Error(),
			}).Update()
		g.Log().Errorf(ctx, "[contract] render pdf failed id=%d err=%v", contractID, err)
		return
	}

	// 渲染成功后可以删掉中间 html，只保留 pdf；这里保留 html 便于排查。
	if _, err := dao.MemberContract.Ctx(ctx).
		Where(cols.Id, contractID).
		Data(g.Map{
			cols.PdfStatus: PDFStatusReady,
			cols.PdfPath:   pdfPath,
			cols.PdfError:  "",
		}).Update(); err != nil {
		g.Log().Errorf(ctx, "[contract] update pdf path err id=%d err=%v", contractID, err)
	}
}

// renderHTMLToPDF 调 Chromium 命令行渲染 PDF。
func renderHTMLToPDF(ctx context.Context, htmlPath, pdfPath string) error {
	bin, err := resolveChromiumBinary(ctx)
	if err != nil {
		return err
	}
	args := []string{
		"--headless=new",
		"--no-sandbox",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--no-pdf-header-footer",
		"--virtual-time-budget=2000",
		"--print-to-pdf=" + pdfPath,
		"file://" + htmlPath,
	}
	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("chromium 渲染失败: %w (output=%s)", err, truncate(string(out), 400))
	}
	stat, err := os.Stat(pdfPath)
	if err != nil {
		return fmt.Errorf("PDF 输出文件不存在: %w (output=%s)", err, truncate(string(out), 400))
	}
	if stat.Size() == 0 {
		return fmt.Errorf("PDF 输出文件为空 (output=%s)", truncate(string(out), 400))
	}
	return nil
}

// resolveChromiumBinary 找到 chromium 可执行文件。
// 优先级：member.contractChromiumPath > $PATH 探测。
func resolveChromiumBinary(ctx context.Context) (string, error) {
	v := strings.TrimSpace(g.Cfg().MustGet(ctx, "member.contractChromiumPath").String())
	if v != "" {
		return v, nil
	}
	candidates := []string{
		"ungoogled-chromium",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			return path, nil
		}
	}
	return "", gerror.New("未找到 chromium 二进制（请安装 chromium 或在 member.contractChromiumPath 配置路径）")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "...(truncated)"
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

// storageDir 始终返回绝对路径（避免 cwd 漂移导致 .pdf/.html 找不到）。
func storageDir() string {
	dir := g.Cfg().MustGet(context.Background(), "member.contractStorageDir").String()
	if strings.TrimSpace(dir) == "" {
		dir = "./runtime/contracts"
	}
	if abs, err := filepath.Abs(dir); err == nil {
		return abs
	}
	return dir
}
