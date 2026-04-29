package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	gojwt "github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT 载荷
type Claims struct {
	UserID     int64  `json:"userId"`
	Username   string `json:"username"`
	DeptID     int64  `json:"deptId"`
	TenantID   int64  `json:"tenantId"`
	MerchantID int64  `json:"merchantId"`
	gojwt.RegisteredClaims
}

var (
	secret       []byte
	memberSecret []byte
	expireTime   time.Duration
)

const defaultInsecureSecret = "gbaseadmin-secret-key"

func init() {
	ctx := gctx.New()
	key, _ := g.Cfg().Get(ctx, "jwt.secret", "")
	raw := strings.TrimSpace(key.String())
	if raw == "" || raw == defaultInsecureSecret || raw == "change_me" {
		panic("jwt.secret 未配置或仍为默认值，请在配置文件中设置一个安全的随机密钥（至少 32 字符）")
	}
	if len(raw) < 32 {
		panic("jwt.secret 长度不足 32 字符，请使用更长的随机密钥")
	}
	secret = []byte(raw)
	mKey, _ := g.Cfg().Get(ctx, "jwt.memberSecret", "")
	if memberKey := strings.TrimSpace(mKey.String()); memberKey != "" {
		memberSecret = []byte(memberKey)
	} else {
		memberSecret = secret
	}
	hours, _ := g.Cfg().Get(ctx, "jwt.expire", 24)
	expireHours := hours.Int()
	if expireHours <= 0 {
		expireHours = 24
	}
	expireTime = time.Duration(expireHours) * time.Hour
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, username string, deptID int64, tenantID int64, merchantID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:     userID,
		Username:   username,
		DeptID:     deptID,
		TenantID:   tenantID,
		MerchantID: merchantID,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(now.Add(expireTime)),
			IssuedAt:  gojwt.NewNumericDate(now),
			Issuer:    "gbaseadmin",
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析 JWT Token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := parseToken(tokenStr, &Claims{}, secret)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, gojwt.ErrTokenInvalidClaims
}

// MemberClaims C端会员 JWT 载荷
type MemberClaims struct {
	MemberID    int64  `json:"memberId"`
	Phone       string `json:"phone"`
	IsCoach     int    `json:"isCoach"`
	CoachID     int64  `json:"coachId"`
	CurrentRole string `json:"currentRole"` // "member" | "coach"
	gojwt.RegisteredClaims
}

// GenerateMemberToken 生成会员 JWT Token
func GenerateMemberToken(memberID int64, phone string, isCoach int, coachID int64, currentRole string) (string, error) {
	now := time.Now()
	claims := MemberClaims{
		MemberID:    memberID,
		Phone:       phone,
		IsCoach:     isCoach,
		CoachID:     coachID,
		CurrentRole: currentRole,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(now.Add(expireTime)),
			IssuedAt:  gojwt.NewNumericDate(now),
			Issuer:    "gbaseadmin-member",
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(memberSecret)
}

// VerifyAnyToken 只验证 token 签名合法且未过期，不关心是哪种身份
func VerifyAnyToken(tokenStr string) bool {
	_, err := ParseToken(tokenStr)
	if err == nil {
		return true
	}
	_, err = ParseMemberToken(tokenStr)
	return err == nil
}

// ParseMemberToken 解析会员 JWT Token
func ParseMemberToken(tokenStr string) (*MemberClaims, error) {
	token, err := parseToken(tokenStr, &MemberClaims{}, memberSecret)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MemberClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, gojwt.ErrTokenInvalidClaims
}

func parseToken(tokenStr string, claims gojwt.Claims, key []byte) (*gojwt.Token, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	return gojwt.ParseWithClaims(tokenStr, claims, func(t *gojwt.Token) (interface{}, error) {
		if t.Method == nil || t.Method.Alg() != gojwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
}

const blacklistKeyPrefix = "system:token:blacklist:"

func tokenHash(tokenStr string) string {
	h := sha256.Sum256([]byte(tokenStr))
	return hex.EncodeToString(h[:])
}

func BlacklistToken(ctx context.Context, tokenStr string, remainingTTL time.Duration) error {
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		return nil
	}
	if remainingTTL <= 0 {
		return nil
	}
	client := safeRedis(ctx)
	if client == nil {
		return nil
	}
	key := blacklistKeyPrefix + tokenHash(tokenStr)
	seconds := int64(remainingTTL/time.Second) + 1
	return client.SetEX(ctx, key, "1", seconds)
}

func IsBlacklisted(ctx context.Context, tokenStr string) bool {
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		return false
	}
	client := safeRedis(ctx)
	if client == nil {
		return false
	}
	key := blacklistKeyPrefix + tokenHash(tokenStr)
	val, err := client.Get(ctx, key)
	if err != nil || val == nil {
		return false
	}
	return val.String() != ""
}

func TokenRemainingTTL(claims *Claims) time.Duration {
	if claims == nil || claims.ExpiresAt == nil {
		return 0
	}
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return 0
	}
	return remaining
}

func safeRedis(ctx context.Context) (client *gredis.Redis) {
	defer func() {
		if recovered := recover(); recovered != nil {
			client = nil
		}
	}()
	return g.Redis()
}

