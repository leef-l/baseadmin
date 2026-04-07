package sms

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gogf/gf/v2/frame/g"
)

// smsConfig 短信配置
type smsConfig struct {
	Provider        string
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
	DevMode         bool
}

// loadConfig 从 GoFrame 配置读取短信配置
func loadConfig(ctx context.Context) smsConfig {
	cfg := smsConfig{}

	if v, err := g.Cfg().Get(ctx, "sms.provider", "aliyun"); err == nil {
		cfg.Provider = v.String()
	}
	if v, err := g.Cfg().Get(ctx, "sms.accessKeyId", ""); err == nil {
		cfg.AccessKeyId = v.String()
	}
	if v, err := g.Cfg().Get(ctx, "sms.accessKeySecret", ""); err == nil {
		cfg.AccessKeySecret = v.String()
	}
	if v, err := g.Cfg().Get(ctx, "sms.signName", ""); err == nil {
		cfg.SignName = v.String()
	}
	if v, err := g.Cfg().Get(ctx, "sms.templateCode", ""); err == nil {
		cfg.TemplateCode = v.String()
	}
	if v, err := g.Cfg().Get(ctx, "sms.devMode", true); err == nil {
		cfg.DevMode = v.Bool()
	} else {
		cfg.DevMode = true
	}

	return normalizeConfig(cfg)
}

// generateCode 生成 6 位随机数字验证码
func generateCode() string {
	if n, err := crand.Int(crand.Reader, big.NewInt(1000000)); err == nil {
		return fmt.Sprintf("%06d", n.Int64())
	}
	return fmt.Sprintf("%06d", mrand.New(mrand.NewSource(time.Now().UnixNano())).Intn(1000000))
}

// SendCode 发送验证码短信
// 开发模式（devMode=true）：直接返回固定验证码 "123456"，不发送短信
// 生产模式（devMode=false）：生成 6 位随机验证码，调用阿里云短信 API 发送
func SendCode(ctx context.Context, phone string) (code string, err error) {
	cfg := loadConfig(ctx)
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "", fmt.Errorf("手机号不能为空")
	}

	if cfg.DevMode {
		g.Log().Infof(ctx, "[SMS DevMode] phone=%s, code=123456（未真实发送）", phone)
		return "123456", nil
	}
	if cfg.Provider != "aliyun" {
		return "", fmt.Errorf("暂不支持的短信服务商: %s", cfg.Provider)
	}
	if cfg.AccessKeyId == "" || cfg.AccessKeySecret == "" || cfg.SignName == "" || cfg.TemplateCode == "" {
		return "", fmt.Errorf("短信配置不完整")
	}

	// 生产模式：生成随机验证码
	code = generateCode()

	// 调用阿里云短信 API
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return "", fmt.Errorf("初始化短信客户端失败: %w", err)
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = cfg.SignName
	request.TemplateCode = cfg.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	response, err := client.SendSms(request)
	if err != nil {
		return "", fmt.Errorf("短信发送失败: %w", err)
	}
	if response.Code != "OK" {
		return "", fmt.Errorf("短信发送失败: %s", response.Message)
	}

	g.Log().Infof(ctx, "[SMS] 验证码已发送: phone=%s, requestId=%s", phone, response.RequestId)
	return code, nil
}

func normalizeConfig(cfg smsConfig) smsConfig {
	cfg.Provider = strings.ToLower(strings.TrimSpace(cfg.Provider))
	if cfg.Provider == "" {
		cfg.Provider = "aliyun"
	}
	cfg.AccessKeyId = strings.TrimSpace(cfg.AccessKeyId)
	cfg.AccessKeySecret = strings.TrimSpace(cfg.AccessKeySecret)
	cfg.SignName = strings.TrimSpace(cfg.SignName)
	cfg.TemplateCode = strings.TrimSpace(cfg.TemplateCode)
	return cfg
}
