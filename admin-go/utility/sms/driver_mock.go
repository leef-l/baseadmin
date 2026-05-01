package sms

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// mockDriver 仅用于本地/测试。
// 不调用任何外部短信网关，把验证码写到日志，部署到生产前必须切到真实 provider。
type mockDriver struct{}

// Kind 实现 ProviderDriver。
func (d *mockDriver) Kind() string { return "mock" }

// SendCode 直接返回 cfg.FixedCode（默认 123456），方便联调。
func (d *mockDriver) SendCode(ctx context.Context, cfg *ProviderConfig, phone string, code string) (string, error) {
	actualCode := strings.TrimSpace(cfg.FixedCode)
	if actualCode == "" {
		actualCode = "123456"
	}
	g.Log().Infof(ctx, "[sms.mock] phone=%s provider=%s code=%s（未真实发送）", phone, cfg.Name, actualCode)
	return actualCode, nil
}
