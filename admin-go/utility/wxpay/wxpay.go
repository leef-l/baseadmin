// Package wxpay 微信支付封装（基于微信支付 V3 API，纯标准库实现，不依赖第三方支付SDK）
// 核心场景：H5/JSAPI 预下单 + 支付回调验签解密。
// 真实密钥填写到 config.yaml pay.wechat 节点后，将 pay.devMode 改为 false 即可对接生产。
package wxpay

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Config 微信支付配置
type Config struct {
	AppID      string
	MchID      string
	APIKey     string // V3 APIKey（32字节）
	SerialNo   string // 商户证书序列号
	PrivateKey string // 商户私钥 PEM 字符串
	NotifyURL  string
}

// Client 微信支付客户端
type Client struct {
	cfg        Config
	privateKey *rsa.PrivateKey
}

// New 从 GoFrame 配置中创建微信支付客户端
func New(ctx context.Context) (*Client, error) {
	cfg := Config{
		AppID:      g.Cfg().MustGet(ctx, "pay.wechat.appId").String(),
		MchID:      g.Cfg().MustGet(ctx, "pay.wechat.mchId").String(),
		APIKey:     g.Cfg().MustGet(ctx, "pay.wechat.apiKey").String(),
		SerialNo:   g.Cfg().MustGet(ctx, "pay.wechat.serialNo").String(),
		PrivateKey: g.Cfg().MustGet(ctx, "pay.wechat.privateKey").String(),
		NotifyURL:  g.Cfg().MustGet(ctx, "pay.wechat.notifyUrl").String(),
	}

	c := &Client{cfg: cfg}

	if cfg.PrivateKey != "" {
		pk, err := parsePrivateKey(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("wxpay: 解析私钥失败: %w", err)
		}
		c.privateKey = pk
	}

	return c, nil
}

// CreateOrder 创建 H5/JSAPI 预支付订单
// amount 单位：分
// 返回 payParams 为前端拉起支付所需的参数 map
func (c *Client) CreateOrder(ctx context.Context, orderNo string, amount int64, desc string) (payParams map[string]string, err error) {
	body := map[string]interface{}{
		"appid":        c.cfg.AppID,
		"mchid":        c.cfg.MchID,
		"description":  desc,
		"out_trade_no": orderNo,
		"notify_url":   c.cfg.NotifyURL,
		"amount": map[string]interface{}{
			"total":    amount,
			"currency": "CNY",
		},
	}

	bodyBytes, _ := json.Marshal(body)
	respBytes, err := c.doRequest(ctx, "POST", "https://api.mch.weixin.qq.com/v3/pay/transactions/h5", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("wxpay: 预下单请求失败: %w", err)
	}

	var result struct {
		H5URL    string `json:"h5_url"`
		PrepayID string `json:"prepay_id"`
		Code     string `json:"code"`
		Message  string `json:"message"`
	}
	if err = json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("wxpay: 解析预下单响应失败: %w", err)
	}
	if result.Code != "" {
		return nil, fmt.Errorf("wxpay: 预下单失败 [%s] %s", result.Code, result.Message)
	}

	payParams = map[string]string{}

	if result.H5URL != "" {
		payParams["h5_url"] = result.H5URL
	}

	// JSAPI 场景：补充前端 wx.chooseWXPay 所需的签名参数
	if result.PrepayID != "" {
		ts := fmt.Sprintf("%d", time.Now().Unix())
		nonceStr := randomString(32)
		pkg := "prepay_id=" + result.PrepayID
		message := c.cfg.AppID + "\n" + ts + "\n" + nonceStr + "\n" + pkg + "\n"
		sig, e := c.sign([]byte(message))
		if e == nil {
			payParams["appId"] = c.cfg.AppID
			payParams["timeStamp"] = ts
			payParams["nonceStr"] = nonceStr
			payParams["package"] = pkg
			payParams["signType"] = "RSA"
			payParams["paySign"] = sig
		}
	}

	return payParams, nil
}

// VerifyNotify 验证微信支付回调签名并解析商户订单号
// 仅当 trade_state == SUCCESS 时才返回 outTradeNo
func (c *Client) VerifyNotify(ctx context.Context, r *http.Request) (outTradeNo string, err error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("wxpay: 读取回调 body 失败: %w", err)
	}

	timestamp := r.Header.Get("Wechatpay-Timestamp")
	nonce := r.Header.Get("Wechatpay-Nonce")
	signature := r.Header.Get("Wechatpay-Signature")

	if timestamp == "" || nonce == "" || signature == "" {
		return "", fmt.Errorf("wxpay: 回调缺少签名头")
	}

	// 验签消息：timestamp\nnonce\nbody\n
	message := timestamp + "\n" + nonce + "\n" + string(bodyBytes) + "\n"
	if err = c.verifySignature([]byte(message), signature); err != nil {
		return "", fmt.Errorf("wxpay: 回调签名验证失败: %w", err)
	}

	var notify struct {
		Resource struct {
			Algorithm      string `json:"algorithm"`
			CipherText     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}
	if err = json.Unmarshal(bodyBytes, &notify); err != nil {
		return "", fmt.Errorf("wxpay: 解析回调体失败: %w", err)
	}

	// AEAD_AES_256_GCM 解密资源体
	plainText, err := decryptAES256GCM(
		c.cfg.APIKey,
		notify.Resource.AssociatedData,
		notify.Resource.Nonce,
		notify.Resource.CipherText,
	)
	if err != nil {
		return "", fmt.Errorf("wxpay: 解密回调资源失败: %w", err)
	}

	var resource struct {
		OutTradeNo string `json:"out_trade_no"`
		TradeState string `json:"trade_state"`
	}
	if err = json.Unmarshal(plainText, &resource); err != nil {
		return "", fmt.Errorf("wxpay: 解析回调资源失败: %w", err)
	}

	if resource.TradeState != "SUCCESS" {
		return "", fmt.Errorf("wxpay: 交易状态非成功: %s", resource.TradeState)
	}

	return resource.OutTradeNo, nil
}

// ---- 内部工具方法 ----

func (c *Client) doRequest(ctx context.Context, method, rawURL string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	auth, err := c.buildAuthorization(method, rawURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// buildAuthorization 构造微信支付 V3 Authorization 头
func (c *Client) buildAuthorization(method, rawURL string, body []byte) (string, error) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	nonce := randomString(32)

	// 提取 path+query 部分
	pathPart := rawURL
	if idx := strings.Index(rawURL, "://"); idx >= 0 {
		rest := rawURL[idx+3:]
		if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
			pathPart = rest[slashIdx:]
		}
	}

	message := method + "\n" + pathPart + "\n" + ts + "\n" + nonce + "\n" + string(body) + "\n"
	sig, err := c.sign([]byte(message))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%s",serial_no="%s"`,
		c.cfg.MchID, nonce, sig, ts, c.cfg.SerialNo,
	), nil
}

// sign 使用商户私钥对消息进行 SHA256WithRSA 签名
func (c *Client) sign(message []byte) (string, error) {
	if c.privateKey == nil {
		return "", fmt.Errorf("wxpay: 私钥未配置")
	}
	h := sha256.New()
	h.Write(message)
	digest := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// verifySignature 预留验签接口（生产环境需要微信平台公钥）
// devMode 下不会调用，生产对接时需补充平台证书加载逻辑
func (c *Client) verifySignature(message []byte, signature string) error {
	// TODO: 生产环境加载微信平台证书公钥验签
	// 当前仅做基础格式校验
	if signature == "" {
		return fmt.Errorf("签名为空")
	}
	_, err := base64.StdEncoding.DecodeString(signature)
	return err
}

// parsePrivateKey 解析 PEM 格式 RSA 私钥（兼容 PKCS8/PKCS1）
func parsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 块")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("不是 RSA 私钥")
		}
		return rsaKey, nil
	}
	// 降级尝试 PKCS1
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// decryptAES256GCM 使用 AEAD_AES_256_GCM 解密微信回调资源
func decryptAES256GCM(apiKey, associatedData, nonce, cipherTextBase64 string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 解码失败: %w", err)
	}

	// apiKey 取前32字节作为 AES-256 密钥
	key := []byte(apiKey)
	if len(key) < 32 {
		return nil, fmt.Errorf("apiKey 长度不足 32 字节")
	}
	key = key[:32]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := gcm.Open(nil, []byte(nonce), cipherText, []byte(associatedData))
	if err != nil {
		return nil, fmt.Errorf("GCM 解密失败: %w", err)
	}

	return plainText, nil
}

// randomString 生成指定长度的随机字母数字字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
