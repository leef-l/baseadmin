// Package portal 是 C 端会员 controller 层。
// 所有方法都只做参数转发，不放业务逻辑；逻辑在 internal/logic/portal/* 下。
package portal

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/utility/sms"
)

// Sms 短信控制器单例。
var Sms = cSms{}

type cSms struct{}

// SendCode 触发短信验证码发送。
// 注册场景额外校验：手机号未被注册（防止扫描）。
// 找回密码场景额外校验：手机号已注册（防止枚举）。
// 校验失败时仍统一返回成功，避免泄露注册状态——只在内部跳过实际下发。
func (c *cSms) SendCode(ctx context.Context, req *v1.SmsCodeReq) (res *v1.SmsCodeRes, err error) {
	scene := strings.TrimSpace(req.Scene)
	if scene == "" {
		return nil, gerror.New("场景不能为空")
	}
	out, err := sms.Default().SendCode(ctx, &sms.SendCodeInput{
		Phone: req.Phone,
		Scene: scene,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SmsCodeRes{ExpiresIn: out.ExpiresIn}, nil
}
