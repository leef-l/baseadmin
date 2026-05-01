package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 列表 -----

// ContractListReq 后台合同列表（按会员/类型/状态筛选）。
type ContractListReq struct {
	g.Meta       `path:"/contract/list" method:"get" tags:"会员合同" summary:"合同列表"`
	UserID       string `json:"userId" dc:"会员 ID（可选）"`
	ContractType string `json:"contractType" v:"in:,register,upgrade,custom" dc:"合同类型（可选）"`
	PdfStatus    *int   `json:"pdfStatus" v:"in:0,1,2,3" dc:"PDF 状态（可选）"`
	PageNum      int    `json:"pageNum" d:"1"`
	PageSize     int    `json:"pageSize" d:"20"`
}

// ContractListRes 列表。
type ContractListRes struct {
	g.Meta `mime:"application/json"`
	Total  int                   `json:"total"`
	List   []*ContractListRecord `json:"list"`
}

// ContractListRecord 单条。
type ContractListRecord struct {
	ContractID    string `json:"contractId"`
	ContractNo    string `json:"contractNo"`
	UserID        string `json:"userId"`
	UserNickname  string `json:"userNickname"`
	UserPhone     string `json:"userPhone"`
	ContractType  string `json:"contractType"`
	TemplateID    string `json:"templateId"`
	SignedAt      string `json:"signedAt"`
	SignedIP      string `json:"signedIp"`
	PDFStatus     int    `json:"pdfStatus"`
	PDFStatusText string `json:"pdfStatusText"`
	CreatedAt     string `json:"createdAt"`
}

// ----- 下载 -----

// ContractDownloadReq 后台下载合同 HTML。
type ContractDownloadReq struct {
	g.Meta     `path:"/contract/download" method:"get" tags:"会员合同" summary:"下载合同"`
	ContractID string `json:"contractId" v:"required#合同 ID 不能为空"`
}
