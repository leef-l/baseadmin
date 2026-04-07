// Package alipay 支付宝支付封装（基于支付宝开放平台标准 API，纯标准库实现）
// 核心场景：手机网站支付（alipay.trade.wap.pay）+ 回调验签。
// 真实密钥填写到 config.yaml pay.alipay 节点后，将 pay.devMode 改为 false 即可对接生产。
package alipay

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Config 支付宝支付配置
type Config struct {
	AppID           string
	PrivateKey      string // 应用私钥 PEM 字符串
	AlipayPublicKey string // 支付宝公钥 PEM 字符串（用于验签）
	NotifyURL       string
	ReturnURL       string
	GatewayURL      string // 支付宝网关，默认 https://openapi.alipay.com/gateway.do
}

// Client 支付宝支付客户端
type Client struct {
	cfg        Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// New 从 GoFrame 配置中创建支付宝客户端
func New(ctx context.Context) (*Client, error) {
	cfg := Config{
		AppID:           g.Cfg().MustGet(ctx, "pay.alipay.appId").String(),
		PrivateKey:      g.Cfg().MustGet(ctx, "pay.alipay.privateKey").String(),
		AlipayPublicKey: g.Cfg().MustGet(ctx, "pay.alipay.alipayPublicKey").String(),
		NotifyURL:       g.Cfg().MustGet(ctx, "pay.alipay.notifyUrl").String(),
		ReturnURL:       g.Cfg().MustGet(ctx, "pay.alipay.returnUrl").String(),
		GatewayURL:      "https://openapi.alipay.com/gateway.do",
	}

	c := &Client{cfg: cfg}

	if cfg.PrivateKey != "" {
		pk, err := parsePrivateKey(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("alipay: 解析应用私钥失败: %w", err)
		}
		c.privateKey = pk
	}

	if cfg.AlipayPublicKey != "" {
		pub, err := parsePublicKey(cfg.AlipayPublicKey)
		if err != nil {
			return nil, fmt.Errorf("alipay: 解析支付宝公钥失败: %w", err)
		}
		c.publicKey = pub
	}

	return c, nil
}

// CreateOrder 创建手机网站支付订单
// amount 单位：分（内部转换为元）
// 返回支付跳转 URL（前端直接跳转该 URL 即可唤起支付宝）
func (c *Client) CreateOrder(ctx context.Context, orderNo string, amount int64, desc string) (payURL string, err error) {
	// 金额：分转元，保留两位小数
	amountYuan := fmt.Sprintf("%.2f", float64(amount)/100)

	bizContent := fmt.Sprintf(
		`{"out_trade_no":"%s","total_amount":"%s","subject":"%s","product_code":"QUICK_WAP_WAY"}`,
		orderNo, amountYuan, desc,
	)

	params := map[string]string{
		"app_id":      c.cfg.AppID,
		"method":      "alipay.trade.wap.pay",
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  c.cfg.NotifyURL,
		"biz_content": bizContent,
	}
	if returnURL := strings.TrimSpace(c.cfg.ReturnURL); returnURL != "" {
		params["return_url"] = returnURL
	}

	sign, err := c.sign(buildSignString(params))
	if err != nil {
		return "", fmt.Errorf("alipay: 签名失败: %w", err)
	}
	params["sign"] = sign

	// 构造 form 表单跳转 URL（GET 方式）
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}

	gatewayURL := strings.TrimSpace(c.cfg.GatewayURL)
	if gatewayURL == "" {
		gatewayURL = "https://openapi.alipay.com/gateway.do"
	}
	payURL = gatewayURL + "?" + vals.Encode()
	return payURL, nil
}

// VerifyNotify 验证支付宝异步通知签名并解析商户订单号
// params 为 POST form 表单参数（r.PostForm）
func (c *Client) VerifyNotify(ctx context.Context, params url.Values) (outTradeNo string, err error) {
	// 提取 sign 和 sign_type
	sign := params.Get("sign")
	if sign == "" {
		return "", fmt.Errorf("alipay: 回调缺少 sign 参数")
	}

	tradeStatus := params.Get("trade_status")
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		return "", fmt.Errorf("alipay: 交易状态非成功: %s", tradeStatus)
	}

	// 过滤 sign/sign_type 后重新排序构造待签名串
	filterParams := make(map[string]string)
	for k, v := range params {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if len(v) > 0 && v[0] != "" {
			filterParams[k] = v[0]
		}
	}

	signStr := buildSignString(filterParams)
	if err = c.verifySignature(signStr, sign); err != nil {
		return "", fmt.Errorf("alipay: 回调签名验证失败: %w", err)
	}

	outTradeNo = params.Get("out_trade_no")
	if strings.TrimSpace(outTradeNo) == "" {
		return "", fmt.Errorf("alipay: 回调缺少 out_trade_no 参数")
	}
	return outTradeNo, nil
}

// ---- 内部工具方法 ----

// buildSignString 按字母序排列参数并拼接为 k=v&k=v 格式
func buildSignString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v := params[k]
		if v != "" {
			parts = append(parts, k+"="+v)
		}
	}
	return strings.Join(parts, "&")
}

// sign 使用应用私钥对 content 进行 SHA256WithRSA 签名，返回 base64 字符串
func (c *Client) sign(content string) (string, error) {
	if c.privateKey == nil {
		return "", fmt.Errorf("alipay: 应用私钥未配置")
	}
	h := sha256.New()
	h.Write([]byte(content))
	digest := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// verifySignature 使用支付宝公钥验证 RSA2(SHA256WithRSA) 签名
func (c *Client) verifySignature(content, signBase64 string) error {
	if c.publicKey == nil {
		// 未配置支付宝公钥时，跳过验签（仅 devMode 会走到此处，正式环境必须配置）
		g.Log().Warning(context.Background(), "alipay: 支付宝公钥未配置，跳过验签（仅允许开发环境）")
		return nil
	}
	sig, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return fmt.Errorf("base64 解码签名失败: %w", err)
	}
	h := sha256.New()
	h.Write([]byte(content))
	digest := h.Sum(nil)
	return rsa.VerifyPKCS1v15(c.publicKey, crypto.SHA256, digest, sig)
}

// parsePrivateKey 解析 PEM 格式应用私钥（PKCS8/PKCS1 均兼容）
func parsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	// 支付宝私钥通常不含 -----BEGIN 头，需要手动包装
	pemStr = normalizePEM(pemStr, "RSA PRIVATE KEY")

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 块")
	}

	// 先尝试 PKCS8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("不是 RSA 私钥")
		}
		return rsaKey, nil
	}
	// 降级 PKCS1
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// parsePublicKey 解析 PEM 格式支付宝公钥
func parsePublicKey(pemStr string) (*rsa.PublicKey, error) {
	pemStr = normalizePEM(pemStr, "PUBLIC KEY")

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("无法解析公钥 PEM 块")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("不是 RSA 公钥")
		}
		return rsaPub, nil
	}
	if rsaPub, pkcs1Err := x509.ParsePKCS1PublicKey(block.Bytes); pkcs1Err == nil {
		return rsaPub, nil
	}
	return nil, err
}

// normalizePEM 如果 PEM 字符串没有 header/footer，则自动添加
func normalizePEM(key, keyType string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	header := "-----BEGIN " + keyType + "-----"
	footer := "-----END " + keyType + "-----"

	if strings.HasPrefix(key, "-----") {
		return key
	}

	// 按 64 字符换行
	var lines []string
	for len(key) > 64 {
		lines = append(lines, key[:64])
		key = key[64:]
	}
	if len(key) > 0 {
		lines = append(lines, key)
	}

	return header + "\n" + strings.Join(lines, "\n") + "\n" + footer
}
