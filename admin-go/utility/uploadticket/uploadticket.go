package uploadticket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Claims struct {
	MemberID            int64    `json:"memberId,omitempty"`
	Role                string   `json:"role,omitempty"`
	Scene               string   `json:"scene"`
	ConfigID            int64    `json:"configId,omitempty"`
	Dir                 string   `json:"dir,omitempty"`
	MaxSize             int64    `json:"maxSize,omitempty"`
	AllowedExts         []string `json:"allowedExts,omitempty"`
	AllowedMimePrefixes []string `json:"allowedMimePrefixes,omitempty"`
	ExpiresAt           int64    `json:"expiresAt"`
	Nonce               string   `json:"nonce,omitempty"`
}

func Sign(claims *Claims, secret string) (string, error) {
	if claims == nil {
		return "", fmt.Errorf("upload ticket claims 不能为空")
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", fmt.Errorf("upload ticket secret 不能为空")
	}

	body, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	payload := base64.RawURLEncoding.EncodeToString(body)
	signature := signPayload(payload, secret)
	return payload + "." + signature, nil
}

func Verify(token string, secret string) (*Claims, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, fmt.Errorf("upload ticket secret 不能为空")
	}

	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("upload ticket 格式无效")
	}

	payload := strings.TrimSpace(parts[0])
	signature := strings.TrimSpace(parts[1])
	if payload == "" || signature == "" {
		return nil, fmt.Errorf("upload ticket 内容无效")
	}

	expected := signPayload(payload, secret)
	if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(strings.ToLower(expected))) {
		return nil, fmt.Errorf("upload ticket 签名无效")
	}

	body, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("upload ticket 解码失败: %w", err)
	}

	var claims Claims
	if err = json.Unmarshal(body, &claims); err != nil {
		return nil, fmt.Errorf("upload ticket 解析失败: %w", err)
	}
	if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("upload ticket 已过期")
	}
	if strings.TrimSpace(claims.Scene) == "" {
		return nil, fmt.Errorf("upload ticket 场景不能为空")
	}
	return &claims, nil
}

func signPayload(payload string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
