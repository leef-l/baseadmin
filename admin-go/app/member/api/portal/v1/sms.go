// Package v1 定义 C 端 portal HTTP API 的请求/响应结构。
//
// 与 admin-go/app/member/api/member/v1（后台 codegen 生成）完全分离，
// 命名空间为 /api/member-portal/*。
package v1

import "github.com/gogf/gf/v2/frame/g"

// SmsCodeReq 发送验证码请求。
//
// scene 取值与业务流程一一对应，前端发起时硬编码：
//   - register      注册（必须手机号未注册）
//   - forget_password 找回密码（必须手机号已注册）
//   - change_phone  修改手机号（已登录态校验旧号，新号校验未注册）
type SmsCodeReq struct {
	g.Meta `path:"/sms/code" method:"post" tags:"会员-认证" summary:"发送短信验证码"`
	Phone  string `json:"phone" v:"required|phone-loose#手机号不能为空|手机号格式不正确" dc:"手机号"`
	Scene  string `json:"scene" v:"required|in:register,forget_password,change_phone#场景不能为空|场景值非法" dc:"场景"`
}

// SmsCodeRes 发送验证码响应。
type SmsCodeRes struct {
	g.Meta    `mime:"application/json"`
	ExpiresIn int64 `json:"expiresIn" dc:"验证码有效秒数"`
}
