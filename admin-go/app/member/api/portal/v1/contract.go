package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 模板预览 -----

// ContractTemplateReq 获取指定类型的合同模板（默认）。
type ContractTemplateReq struct {
	g.Meta       `path:"/contract/template" method:"get" tags:"会员-合同" summary:"获取合同模板"`
	ContractType string `json:"contractType" v:"in:register,upgrade,custom" d:"register" dc:"合同类型"`
}

// ContractTemplateRes 模板正文。
type ContractTemplateRes struct {
	g.Meta       `mime:"application/json"`
	TemplateID   string `json:"templateId"`
	TemplateName string `json:"templateName"`
	Content      string `json:"content" dc:"HTML 模板（含占位符）"`
}

// ----- 提交签名 -----

// ContractSignReq 提交手写签名 + 合同信息。
type ContractSignReq struct {
	g.Meta         `path:"/contract/sign" method:"post" tags:"会员-合同" summary:"签署合同"`
	ContractType   string `json:"contractType" v:"in:register,upgrade,custom" d:"register"`
	TemplateID     string `json:"templateId" dc:"模板 ID（不传则用类型默认模板）"`
	SignatureImage string `json:"signatureImage" v:"required#签名不能为空" dc:"data:image/png;base64,..."`
	RelatedID      string `json:"relatedId" dc:"关联业务 ID（可选）"`
}

// ContractSignRes 签署结果。
type ContractSignRes struct {
	g.Meta     `mime:"application/json"`
	ContractID string `json:"contractId"`
	ContractNo string `json:"contractNo"`
}

// ----- 我的合同 -----

// ContractListReq 我的合同列表。
type ContractListReq struct {
	g.Meta   `path:"/contract/list" method:"get" tags:"会员-合同" summary:"我的合同"`
	PageNum  int `json:"pageNum" d:"1"`
	PageSize int `json:"pageSize" d:"20"`
}

// ContractListRes 合同列表。
type ContractListRes struct {
	g.Meta `mime:"application/json"`
	Total  int                 `json:"total"`
	List   []*ContractListItem `json:"list"`
}

// ContractListItem 单条合同。
type ContractListItem struct {
	ContractID       string `json:"contractId"`
	ContractNo       string `json:"contractNo"`
	ContractType     string `json:"contractType"`
	ContractTypeText string `json:"contractTypeText"`
	SignedAt         string `json:"signedAt"`
	PDFStatus        int    `json:"pdfStatus" dc:"0=未生成 1=生成中 2=已就绪 3=失败"`
	PDFStatusText    string `json:"pdfStatusText"`
}

// ----- 已签状态查询 -----

// ContractStatusReq 查询是否已签。
type ContractStatusReq struct {
	g.Meta       `path:"/contract/status" method:"get" tags:"会员-合同" summary:"是否已签"`
	ContractType string `json:"contractType" v:"in:register,upgrade,custom" d:"register"`
}

// ContractStatusRes 是否已签。
type ContractStatusRes struct {
	g.Meta   `mime:"application/json"`
	HasSign  bool `json:"hasSign"`
}

// ----- 下载 -----

// ContractDownloadReq 下载合同（HTML 流，浏览器可"打印为 PDF"）。
type ContractDownloadReq struct {
	g.Meta     `path:"/contract/download" method:"get" tags:"会员-合同" summary:"下载合同"`
	ContractID string `json:"contractId" v:"required#合同 ID 不能为空"`
}
