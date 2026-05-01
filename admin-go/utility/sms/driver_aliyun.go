package sms

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// aliyunDriver 使用阿里云 dysmsapi SDK 发送验证码。
// 模板要求 ${code} 占位，例如 "您的验证码：${code}，5 分钟内有效"。
type aliyunDriver struct{}

// Kind 实现 ProviderDriver。
func (d *aliyunDriver) Kind() string { return "aliyun" }

// SendCode 真实下发到运营商短信网关。
func (d *aliyunDriver) SendCode(ctx context.Context, cfg *ProviderConfig, phone string, code string) (string, error) {
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.SignName == "" || cfg.TemplateCode == "" {
		return "", gerror.Newf("短信平台 %s 配置不完整", cfg.Name)
	}

	region := cfg.Region
	if region == "" {
		region = "cn-hangzhou"
	}
	client, err := dysmsapi.NewClientWithAccessKey(region, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return "", fmt.Errorf("初始化阿里云短信客户端失败: %w", err)
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = cfg.SignName
	request.TemplateCode = cfg.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	response, err := client.SendSms(request)
	if err != nil {
		return "", fmt.Errorf("阿里云短信发送失败: %w", err)
	}
	if response.Code != "OK" {
		return "", fmt.Errorf("阿里云短信发送失败: %s", response.Message)
	}

	g.Log().Infof(ctx, "[sms.aliyun] 验证码已发送 phone=%s provider=%s requestId=%s", phone, cfg.Name, response.RequestId)
	return code, nil
}
