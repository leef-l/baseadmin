// Package contract 后台合同列表/下载接口（不走 codegen，手写）。
package contract

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/entity"
)

type cContract struct{}

// Contract 后台合同接口。
var Contract = cContract{}

// List 后台合同列表。
func (c cContract) List(ctx context.Context, req *v1.ContractListReq) (res *v1.ContractListRes, err error) {
	pageNum := req.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	cols := dao.MemberContract.Columns()
	m := dao.MemberContract.Ctx(ctx).Where(cols.DeletedAt, nil)
	if req.UserID != "" {
		uid, _ := strconv.ParseInt(req.UserID, 10, 64)
		if uid > 0 {
			m = m.Where(cols.UserId, uid)
		}
	}
	if req.ContractType != "" {
		m = m.Where(cols.ContractType, req.ContractType)
	}
	if req.PdfStatus != nil {
		m = m.Where(cols.PdfStatus, *req.PdfStatus)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberContract
	if err := m.OrderDesc(cols.Id).Page(pageNum, pageSize).Scan(&rows); err != nil {
		return nil, err
	}

	// 批量取关联会员 nickname/phone
	idSet := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		idSet[row.UserId] = struct{}{}
	}
	userMap := make(map[uint64]entity.MemberUser, len(idSet))
	if len(idSet) > 0 {
		ids := make([]uint64, 0, len(idSet))
		for id := range idSet {
			ids = append(ids, id)
		}
		var users []entity.MemberUser
		uc := dao.MemberUser.Columns()
		_ = dao.MemberUser.Ctx(ctx).Where(uc.Id, ids).Where(uc.DeletedAt, nil).Scan(&users)
		for _, u := range users {
			userMap[u.Id] = u
		}
	}

	out := &v1.ContractListRes{Total: total, List: make([]*v1.ContractListRecord, 0, len(rows))}
	for _, row := range rows {
		signedAt := ""
		if row.SignedAt != nil && !row.SignedAt.IsZero() {
			signedAt = row.SignedAt.String()
		}
		createdAt := ""
		if row.CreatedAt != nil && !row.CreatedAt.IsZero() {
			createdAt = row.CreatedAt.String()
		}
		u := userMap[row.UserId]
		out.List = append(out.List, &v1.ContractListRecord{
			ContractID:    fmt.Sprintf("%d", row.Id),
			ContractNo:    row.ContractNo,
			UserID:        fmt.Sprintf("%d", row.UserId),
			UserNickname:  u.Nickname,
			UserPhone:     u.Phone,
			ContractType:  row.ContractType,
			TemplateID:    fmt.Sprintf("%d", row.TemplateId),
			SignedAt:      signedAt,
			SignedIP:      row.SignedIp,
			PDFStatus:     row.PdfStatus,
			PDFStatusText: pdfStatusText(row.PdfStatus),
			CreatedAt:     createdAt,
		})
	}
	return out, nil
}

// Download 后台下载合同（流式）。
func Download(r *ghttp.Request) {
	ctx := r.Context()
	contractID, _ := strconv.ParseInt(r.GetQuery("contractId").String(), 10, 64)
	if contractID <= 0 {
		r.Response.WriteStatus(400, "合同 ID 不能为空")
		return
	}
	cols := dao.MemberContract.Columns()
	var row entity.MemberContract
	if err := dao.MemberContract.Ctx(ctx).
		Where(cols.Id, contractID).
		Where(cols.DeletedAt, nil).
		Scan(&row); err != nil {
		g.Log().Warningf(ctx, "admin Download err=%v", err)
		r.Response.WriteStatus(500, "查询失败")
		return
	}
	if row.Id == 0 {
		r.Response.WriteStatus(404, "合同不存在")
		return
	}
	if row.SignedHtml == "" {
		r.Response.WriteStatus(404, "合同尚未生成")
		return
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.html"`, row.ContractNo))
	r.Response.Write(row.SignedHtml)
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

// 兼容 unused-import 检查（gerror 在未来扩展会用到）
var _ = gerror.New
