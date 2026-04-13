package appticket

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	gojwt "github.com/golang-jwt/jwt/v5"

	"gbaseadmin/utility/snowflake"
)

const (
	defaultSharedKey     = "gbaseadmin-shared-sso-key"
	defaultAppID         = "gbaseadmin-admin"
	defaultTicketExpire  = 60 * time.Second
	ticketReplayTTLGrace = 5 * time.Second
)

// Claims 应用间票据载荷。
type Claims struct {
	Username  string `json:"username"`
	SourceApp string `json:"sourceApp,omitempty"`
	TargetApp string `json:"targetApp,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	gojwt.RegisteredClaims
}

// Generate 生成应用间登录票据。
func Generate(ctx context.Context, username, sourceApp, targetApp string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return "", fmt.Errorf("username is required")
	}

	var (
		now      = time.Now()
		appID    = CurrentAppID(ctx)
		source   = strings.TrimSpace(sourceApp)
		target   = strings.TrimSpace(targetApp)
		ticketID = fmt.Sprintf("%d", snowflake.Generate())
	)

	if source == "" {
		source = appID
	}
	if target == "" {
		target = appID
	}

	claims := Claims{
		Username:  username,
		SourceApp: source,
		TargetApp: target,
		Nonce:     fmt.Sprintf("%d", snowflake.Generate()),
		RegisteredClaims: gojwt.RegisteredClaims{
			Audience:  []string{target},
			ExpiresAt: gojwt.NewNumericDate(now.Add(ticketExpire(ctx))),
			ID:        ticketID,
			IssuedAt:  gojwt.NewNumericDate(now),
			Issuer:    source,
			Subject:   username,
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(sharedKey(ctx))
}

// Parse 解析应用间票据。
func Parse(ctx context.Context, ticket string) (*Claims, error) {
	ticket = strings.TrimSpace(ticket)
	token, err := gojwt.ParseWithClaims(ticket, &Claims{}, func(t *gojwt.Token) (interface{}, error) {
		if t.Method == nil || t.Method.Alg() != gojwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected ticket signing method: %v", t.Header["alg"])
		}
		return sharedKey(ctx), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.SourceApp == "" {
			claims.SourceApp = strings.TrimSpace(claims.Issuer)
		}
		if claims.TargetApp == "" && len(claims.Audience) > 0 {
			claims.TargetApp = strings.TrimSpace(claims.Audience[0])
		}
		return claims, nil
	}
	return nil, gojwt.ErrTokenInvalidClaims
}

// CurrentAppID 返回当前应用标识。
func CurrentAppID(ctx context.Context) string {
	value, _ := g.Cfg().Get(normalizeContext(ctx), "sso.appId", defaultAppID)
	return normalizeString(value.String(), defaultAppID)
}

// ValidateTarget 校验票据目标应用。
func ValidateTarget(ctx context.Context, claims *Claims) error {
	if claims == nil {
		return fmt.Errorf("ticket claims are required")
	}
	appID := CurrentAppID(ctx)
	if appID == "" {
		return nil
	}
	if len(claims.Audience) > 0 {
		for _, item := range claims.Audience {
			if strings.TrimSpace(item) == appID {
				return nil
			}
		}
		return fmt.Errorf("ticket audience mismatch")
	}
	if target := strings.TrimSpace(claims.TargetApp); target != "" && target != appID {
		return fmt.Errorf("ticket target mismatch")
	}
	return nil
}

// ReplayCacheKey 返回票据重放保护缓存键。
func ReplayCacheKey(claims *Claims, rawTicket string) string {
	var ticketID string
	if claims != nil {
		ticketID = strings.TrimSpace(claims.ID)
	}
	if ticketID == "" {
		sum := sha256.Sum256([]byte(strings.TrimSpace(rawTicket)))
		ticketID = hex.EncodeToString(sum[:16])
	}
	return "system:auth:ticket:used:" + ticketID
}

// ReplayTTL 返回票据重放保护缓存时长。
func ReplayTTL(claims *Claims) time.Duration {
	if claims != nil && claims.ExpiresAt != nil {
		ttl := time.Until(claims.ExpiresAt.Time) + ticketReplayTTLGrace
		if ttl < time.Second {
			return time.Second
		}
		return ttl
	}
	return defaultTicketExpire
}

func sharedKey(ctx context.Context) []byte {
	var (
		cfgValue, _ = g.Cfg().Get(normalizeContext(ctx), "sso.sharedKey", "")
		key         = strings.TrimSpace(cfgValue.String())
	)
	if key == "" {
		if jwtValue, err := g.Cfg().Get(normalizeContext(ctx), "jwt.secret", defaultSharedKey); err == nil {
			key = strings.TrimSpace(jwtValue.String())
		}
	}
	return []byte(normalizeString(key, defaultSharedKey))
}

func ticketExpire(ctx context.Context) time.Duration {
	value, _ := g.Cfg().Get(normalizeContext(ctx), "sso.ticketExpireSeconds", int(defaultTicketExpire/time.Second))
	seconds := value.Int()
	if seconds <= 0 {
		return defaultTicketExpire
	}
	return time.Duration(seconds) * time.Second
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func normalizeString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
