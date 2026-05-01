package portal

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/service"
)

// Auth 控制器（注册 / 登录 / 找回密码 / 邀请码预览）。
var Auth = cAuth{}

type cAuth struct{}

// Register 注册并自动登录。
func (c *cAuth) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	out, err := service.PortalAuth().Register(ctx, &service.PortalRegisterInput{
		Phone:      req.Phone,
		SmsCode:    req.SmsCode,
		Password:   req.Password,
		InviteCode: req.InviteCode,
		Nickname:   req.Nickname,
		RegisterIP: clientIP(ctx),
	})
	if err != nil {
		return nil, err
	}
	return &v1.RegisterRes{LoginResult: convertLoginResult(out)}, nil
}

// Login 账号密码登录。
func (c *cAuth) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	out, err := service.PortalAuth().Login(ctx, &service.PortalLoginInput{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{LoginResult: convertLoginResult(out)}, nil
}

// ForgetPassword 找回密码。
func (c *cAuth) ForgetPassword(ctx context.Context, req *v1.ForgetPasswordReq) (res *v1.ForgetPasswordRes, err error) {
	if err = service.PortalAuth().ForgetPassword(ctx, &service.PortalForgetPasswordInput{
		Phone:       req.Phone,
		SmsCode:     req.SmsCode,
		NewPassword: req.NewPassword,
	}); err != nil {
		return nil, err
	}
	return &v1.ForgetPasswordRes{}, nil
}

// InvitePreview 邀请码反查上级，前端注册页用。
func (c *cAuth) InvitePreview(ctx context.Context, req *v1.InvitePreviewReq) (res *v1.InvitePreviewRes, err error) {
	out, err := service.PortalAuth().InvitePreview(ctx, req.InviteCode)
	if err != nil {
		return nil, err
	}
	return &v1.InvitePreviewRes{
		Found:    out.Found,
		Nickname: out.Nickname,
		Avatar:   out.Avatar,
	}, nil
}

// convertLoginResult 把 service 层结构转成 API 层结构，避免 service 依赖 api 包。
func convertLoginResult(out *service.PortalLoginOutput) *v1.LoginResult {
	if out == nil {
		return nil
	}
	return &v1.LoginResult{
		Token:       out.Token,
		MemberID:    out.MemberID,
		Phone:       out.Phone,
		Nickname:    out.Nickname,
		Avatar:      out.Avatar,
		InviteCode:  out.InviteCode,
		LevelID:     out.LevelID,
		IsQualified: out.IsQualified,
	}
}

// clientIP 取请求方 IP。优先 X-Forwarded-For，再 X-Real-IP，最后 GetClientIp。
func clientIP(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	if value := strings.TrimSpace(r.GetHeader("X-Forwarded-For")); value != "" {
		if idx := strings.Index(value, ","); idx > 0 {
			return strings.TrimSpace(value[:idx])
		}
		return value
	}
	if value := strings.TrimSpace(r.GetHeader("X-Real-IP")); value != "" {
		return value
	}
	return r.GetClientIp()
}
