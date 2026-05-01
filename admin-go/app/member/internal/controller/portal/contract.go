package portal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/logic/contract"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model/entity"
)

type cContract struct{}

// Contract 合同模块（部分接口公开模板，部分需要登录）。
var Contract = cContract{}

// Template 公开接口：返回当前类型默认模板。
func (c cContract) Template(ctx context.Context, req *v1.ContractTemplateReq) (res *v1.ContractTemplateRes, err error) {
	tpl, err := contract.PickDefaultTemplate(ctx, req.ContractType)
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return nil, gerror.New("未找到对应类型的合同模板")
	}
	return &v1.ContractTemplateRes{
		TemplateID:   fmt.Sprintf("%d", tpl.Id),
		TemplateName: tpl.TemplateName,
		Content:      tpl.Content,
	}, nil
}

// Sign 受保护接口：提交手写签名。
func (c cContract) Sign(ctx context.Context, req *v1.ContractSignReq) (res *v1.ContractSignRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	if memberID <= 0 {
		return nil, gerror.New("会员未登录")
	}

	var user entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, memberID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Scan(&user); err != nil {
		return nil, err
	}

	r := ghttp.RequestFromCtx(ctx)
	tplID, _ := strconv.ParseInt(req.TemplateID, 10, 64)
	relID, _ := strconv.ParseInt(req.RelatedID, 10, 64)

	out, err := contract.Sign(ctx, &contract.SignInput{
		UserID:         memberID,
		ContractType:   req.ContractType,
		TemplateID:     tplID,
		SignatureImage: req.SignatureImage,
		UserNickname:   user.Nickname,
		UserPhone:      user.Phone,
		IP:             r.GetClientIp(),
		UA:             r.UserAgent(),
		RelatedID:      relID,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ContractSignRes{ContractID: out.ContractID, ContractNo: out.ContractNo}, nil
}

// Status 受保护：查询是否已签。
func (c cContract) Status(ctx context.Context, req *v1.ContractStatusReq) (res *v1.ContractStatusRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	if memberID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	signed, err := contract.HasUserSigned(ctx, memberID, req.ContractType)
	if err != nil {
		return nil, err
	}
	return &v1.ContractStatusRes{HasSign: signed}, nil
}

// List 受保护：我的合同列表。
func (c cContract) List(ctx context.Context, req *v1.ContractListReq) (res *v1.ContractListRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	if memberID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	pageNum := req.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	cols := dao.MemberContract.Columns()
	m := dao.MemberContract.Ctx(ctx).
		Where(cols.UserId, memberID).
		Where(cols.DeletedAt, nil)
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberContract
	if err := m.OrderDesc(cols.Id).Page(pageNum, pageSize).Scan(&rows); err != nil {
		return nil, err
	}
	out := &v1.ContractListRes{Total: total, List: make([]*v1.ContractListItem, 0, len(rows))}
	for _, row := range rows {
		signedAt := ""
		if row.SignedAt != nil && !row.SignedAt.IsZero() {
			signedAt = row.SignedAt.String()
		}
		out.List = append(out.List, &v1.ContractListItem{
			ContractID:       fmt.Sprintf("%d", row.Id),
			ContractNo:       row.ContractNo,
			ContractType:     row.ContractType,
			ContractTypeText: contractTypeText(row.ContractType),
			SignedAt:         signedAt,
			PDFStatus:        row.PdfStatus,
			PDFStatusText:    pdfStatusText(row.PdfStatus),
		})
	}
	return out, nil
}

// Download 受保护：流式输出 signed_html（前端可在浏览器打印为 PDF）。
// 走 ghttp 直接写响应，不经 MiddlewareHandlerResponse 封装。
func DownloadContract(r *ghttp.Request) {
	ctx := r.Context()
	memberID := int64(middleware.CurrentMemberID(ctx))
	if memberID <= 0 {
		r.Response.WriteStatus(401, "未登录")
		return
	}
	contractID, _ := strconv.ParseInt(r.GetQuery("contractId").String(), 10, 64)
	if contractID <= 0 {
		r.Response.WriteStatus(400, "合同 ID 不能为空")
		return
	}
	html, contractNo, err := contract.GetSignedHTML(ctx, memberID, contractID)
	if err != nil {
		g.Log().Warningf(ctx, "DownloadContract err: %v", err)
		r.Response.WriteStatus(404, "合同不存在")
		return
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.html"`, contractNo))
	r.Response.Write(html)
}

func contractTypeText(t string) string {
	switch t {
	case "register":
		return "注册协议"
	case "upgrade":
		return "升级协议"
	case "custom":
		return "自定义"
	}
	return t
}

func pdfStatusText(s int) string {
	switch s {
	case 0:
		return "未生成"
	case 1:
		return "生成中"
	case 2:
		return "已就绪"
	case 3:
		return "失败"
	}
	return ""
}
