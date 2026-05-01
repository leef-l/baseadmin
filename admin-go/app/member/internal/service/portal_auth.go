package service

import (
	"context"
)

// PortalAuthService 是 C 端注册/登录/找回密码服务接口。
//
// 由 internal/logic/portal/auth.go 实现并通过 RegisterPortalAuth 注册到 localPortalAuth。
type IPortalAuthService interface {
	Register(ctx context.Context, in *PortalRegisterInput) (*PortalLoginOutput, error)
	Login(ctx context.Context, in *PortalLoginInput) (*PortalLoginOutput, error)
	ForgetPassword(ctx context.Context, in *PortalForgetPasswordInput) error
	InvitePreview(ctx context.Context, inviteCode string) (*PortalInvitePreviewOutput, error)
}

// PortalRegisterInput 注册入参。
type PortalRegisterInput struct {
	Phone      string
	SmsCode    string
	Password   string
	InviteCode string
	Nickname   string
	RegisterIP string
}

// PortalLoginInput 登录入参。
type PortalLoginInput struct {
	Account  string
	Password string
}

// PortalForgetPasswordInput 找回密码入参。
type PortalForgetPasswordInput struct {
	Phone       string
	SmsCode     string
	NewPassword string
}

// PortalLoginOutput 登录/注册成功后的统一返回。
type PortalLoginOutput struct {
	Token       string
	MemberID    string
	Phone       string
	Nickname    string
	Avatar      string
	InviteCode  string
	LevelID     string
	IsQualified int
}

// PortalInvitePreviewOutput 邀请码预览结果。
type PortalInvitePreviewOutput struct {
	Found    bool
	Nickname string
	Avatar   string
}

var localPortalAuth IPortalAuthService

// PortalAuth 暴露注册的 portal 认证服务。
func PortalAuth() IPortalAuthService { return localPortalAuth }

// RegisterPortalAuth 由 logic/portal 包在 init 时调用注入实现。
func RegisterPortalAuth(s IPortalAuthService) { localPortalAuth = s }
