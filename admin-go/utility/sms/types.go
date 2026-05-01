// Package sms 提供短信验证码发送与校验能力。
//
// 设计取舍：
//   - 不引入独立 HTTP 微服务，所有调用方（portal/auth、portal/sms 等）直接 import 本包。
//   - 不引入消息队列（mqapi）和应用中心（appcenter），保持 funddisk 仓库精简。
//   - 验证码缓存复用 GoFrame 的 g.Redis()，限流键全局唯一（手机号+场景维度），不做租户隔离。
//   - Provider 抽象沿用 play-companion-v2 的 driver 模式，目前仅落地 mock 与 aliyun 两套。
package sms

import "context"

// SendCodeInput 发送验证码入参。
type SendCodeInput struct {
	Phone    string // 手机号，必填
	Scene    string // 场景标识，例如 register/login/forget_password
	Provider string // 指定短信平台（可选），为空走配置 sms.provider
}

// SendCodeOutput 发送验证码结果。
type SendCodeOutput struct {
	Provider  string // 实际使用的平台名
	ExpiresIn int64  // 验证码有效秒数
}

// VerifyCodeInput 校验验证码入参。
type VerifyCodeInput struct {
	Phone   string // 手机号
	Scene   string // 场景标识
	Code    string // 待校验验证码
	Consume bool   // 校验通过后是否立即销毁验证码（注册/找回密码场景为 true，其他视情况）
}

// VerifyCodeOutput 校验验证码结果。
type VerifyCodeOutput struct {
	Verified bool
}

// ProviderDriver 抽象短信发送通道。
// SendCode 返回的字符串是真正下发的验证码：mock 通道返回固定码或随机码；
// 真实通道把 code 透传给短信网关后返回原 code 即可。
type ProviderDriver interface {
	Kind() string
	SendCode(ctx context.Context, cfg *ProviderConfig, phone string, code string) (string, error)
}

// ProviderConfig 单个短信平台配置。
type ProviderConfig struct {
	Name            string
	Kind            string
	Region          string
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
	FixedCode       string // 仅 mock 通道使用，留空则随机生成 6 位数字
}

const (
	codeKeyPrefix          = "member:sms:code"
	limitKeyPrefix         = "member:sms:limit"
	verifyAttemptKeyPrefix = "member:sms:verify:attempt"

	defaultCodeExpireSeconds  int64 = 300
	defaultLimitSeconds       int64 = 60
	defaultVerifyMaxAttempts  int64 = 5
	defaultDriverFallbackKind       = "mock"
)
