package sms

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Service 是短信服务的对外门面，被 portal 接入层和后台调用。
//
// 调用方式：
//
//	out, err := sms.Default().SendCode(ctx, &sms.SendCodeInput{Phone: "13800138000", Scene: "register"})
//	ok, err  := sms.Default().VerifyCode(ctx, &sms.VerifyCodeInput{Phone, Scene, Code, Consume: true})
type Service struct {
	mu      sync.RWMutex
	drivers map[string]ProviderDriver
}

var defaultService = newService()

// Default 返回默认全局 sms 服务实例。
func Default() *Service { return defaultService }

func newService() *Service {
	s := &Service{drivers: make(map[string]ProviderDriver)}
	s.RegisterDriver(&mockDriver{})
	s.RegisterDriver(&aliyunDriver{})
	return s
}

// RegisterDriver 注册自定义 driver。
func (s *Service) RegisterDriver(d ProviderDriver) {
	if d == nil {
		return
	}
	kind := strings.ToLower(strings.TrimSpace(d.Kind()))
	if kind == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drivers[kind] = d
}

func (s *Service) resolveDriver(kind string) (ProviderDriver, error) {
	value := strings.ToLower(strings.TrimSpace(kind))
	if value == "" {
		return nil, gerror.New("短信平台类型不能为空")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	driver, ok := s.drivers[value]
	if !ok {
		return nil, gerror.Newf("未注册的短信平台类型: %s", kind)
	}
	return driver, nil
}

// SendCode 发送短信验证码，包含限流、生成随机码、写入 Redis、调用通道四步。
func (s *Service) SendCode(ctx context.Context, in *SendCodeInput) (*SendCodeOutput, error) {
	if in == nil {
		return nil, gerror.New("发送验证码请求不能为空")
	}
	phone := strings.TrimSpace(in.Phone)
	scene := strings.TrimSpace(in.Scene)
	if phone == "" {
		return nil, gerror.New("手机号不能为空")
	}
	if scene == "" {
		return nil, gerror.New("场景不能为空")
	}

	limitSeconds := codeLimitSeconds(ctx)
	limitKey := buildLimitKey(scene, phone)
	acquired, err := acquireLimit(ctx, limitKey, limitSeconds)
	if err != nil {
		return nil, gerror.New("短信频率校验失败")
	}
	if !acquired {
		return nil, gerror.New("发送过于频繁，请稍后再试")
	}
	success := false
	defer func() {
		if !success {
			releaseLimit(ctx, limitKey)
		}
	}()

	providerName := strings.TrimSpace(in.Provider)
	if providerName == "" {
		providerName = strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.provider").String())
	}
	if providerName == "" {
		providerName = defaultDriverFallbackKind
	}
	cfg := loadProviderConfig(ctx, providerName)
	driver, err := s.resolveDriver(cfg.Kind)
	if err != nil {
		return nil, err
	}

	code := generateCode()
	actualCode, err := driver.SendCode(ctx, cfg, phone, code)
	if err != nil {
		return nil, err
	}

	expireSeconds := codeExpireSeconds(ctx)
	codeKey := buildCodeKey(scene, phone)
	if err = g.Redis().SetEX(ctx, codeKey, actualCode, expireSeconds); err != nil {
		return nil, gerror.New("验证码存储失败")
	}
	releaseVerifyAttempts(ctx, buildVerifyAttemptKey(scene, phone))
	success = true

	g.Log().Infof(ctx, "[sms] code stored phone=%s scene=%s provider=%s", phone, scene, cfg.Name)
	return &SendCodeOutput{
		Provider:  cfg.Name,
		ExpiresIn: expireSeconds,
	}, nil
}

// VerifyCode 校验验证码。
//   - Consume=true：校验成功后立即销毁验证码（一次性）。
//   - Consume=false：校验后保留，便于多步流程复用。
//
// 错误次数超过 sms.verifyMaxAttempts 时直接报"重新获取"。
func (s *Service) VerifyCode(ctx context.Context, in *VerifyCodeInput) (*VerifyCodeOutput, error) {
	if in == nil {
		return nil, gerror.New("校验验证码请求不能为空")
	}
	phone := strings.TrimSpace(in.Phone)
	scene := strings.TrimSpace(in.Scene)
	code := strings.TrimSpace(in.Code)
	if phone == "" {
		return nil, gerror.New("手机号不能为空")
	}
	if scene == "" {
		return nil, gerror.New("场景不能为空")
	}
	if code == "" {
		return nil, gerror.New("验证码不能为空")
	}

	codeKey := buildCodeKey(scene, phone)
	attemptKey := buildVerifyAttemptKey(scene, phone)
	maxAttempts := codeVerifyMaxAttempts(ctx)
	if maxAttempts > 0 {
		current, err := currentVerifyAttempts(ctx, attemptKey)
		if err != nil {
			return nil, gerror.New("验证码校验失败")
		}
		if current >= maxAttempts {
			return nil, gerror.New("验证码错误次数过多，请重新获取")
		}
	}

	if in.Consume {
		consumed, err := consumeCode(ctx, codeKey, code)
		if err != nil {
			return nil, gerror.New("验证码校验失败")
		}
		if !consumed {
			return nil, verifyFailed(ctx, attemptKey, codeKey, maxAttempts)
		}
	} else {
		cached, err := g.Redis().Get(ctx, codeKey)
		if err != nil {
			return nil, gerror.New("验证码校验失败")
		}
		if cached.IsNil() || cached.IsEmpty() || strings.TrimSpace(cached.String()) != code {
			return nil, verifyFailed(ctx, attemptKey, codeKey, maxAttempts)
		}
	}
	releaseVerifyAttempts(ctx, attemptKey)
	return &VerifyCodeOutput{Verified: true}, nil
}

// loadProviderConfig 从 sms.providers.{name} 节点加载配置。
// 兼容老风格 sms.accessKeyId/.signName 平铺式配置（仅当 name == sms.provider 时生效）。
func loadProviderConfig(ctx context.Context, providerName string) *ProviderConfig {
	name := strings.ToLower(strings.TrimSpace(providerName))
	getStr := func(key string) string {
		return strings.TrimSpace(g.Cfg().MustGet(ctx, key).String())
	}
	cfg := &ProviderConfig{
		Name:            name,
		Kind:            strings.ToLower(getStr(fmt.Sprintf("sms.providers.%s.kind", name))),
		Region:          getStr(fmt.Sprintf("sms.providers.%s.region", name)),
		AccessKeyID:     getStr(fmt.Sprintf("sms.providers.%s.accessKeyId", name)),
		AccessKeySecret: getStr(fmt.Sprintf("sms.providers.%s.accessKeySecret", name)),
		SignName:        getStr(fmt.Sprintf("sms.providers.%s.signName", name)),
		TemplateCode:    getStr(fmt.Sprintf("sms.providers.%s.templateCode", name)),
		FixedCode:       getStr(fmt.Sprintf("sms.providers.%s.fixedCode", name)),
	}
	if cfg.Kind == "" {
		cfg.Kind = name
	}
	if cfg.Region == "" {
		cfg.Region = "cn-hangzhou"
	}

	// 平铺兜底：当请求的 provider 就是默认 sms.provider 时，再读 sms.{field}。
	if name == strings.ToLower(strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.provider").String())) {
		if cfg.AccessKeyID == "" {
			cfg.AccessKeyID = getStr("sms.accessKeyId")
		}
		if cfg.AccessKeySecret == "" {
			cfg.AccessKeySecret = getStr("sms.accessKeySecret")
		}
		if cfg.SignName == "" {
			cfg.SignName = getStr("sms.signName")
		}
		if cfg.TemplateCode == "" {
			cfg.TemplateCode = getStr("sms.templateCode")
		}
	}
	return cfg
}

// ----- Redis helpers -----

func acquireLimit(ctx context.Context, key string, ttlSeconds int64) (bool, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultLimitSeconds
	}
	result, err := g.Redis().Set(ctx, key, 1, gredis.SetOption{
		NX:        true,
		TTLOption: gredis.TTLOption{EX: &ttlSeconds},
	})
	if err != nil {
		return false, err
	}
	return result != nil && !result.IsNil() && !result.IsEmpty(), nil
}

func releaseLimit(ctx context.Context, key string) {
	if strings.TrimSpace(key) == "" {
		return
	}
	_, _ = g.Redis().Del(ctx, key)
}

func currentVerifyAttempts(ctx context.Context, key string) (int64, error) {
	value, err := g.Redis().Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if value == nil || value.IsNil() || value.IsEmpty() {
		return 0, nil
	}
	return value.Int64(), nil
}

func increaseVerifyAttempts(ctx context.Context, attemptKey, codeKey string, fallbackTTL int64) (int64, error) {
	attempts, err := g.Redis().Incr(ctx, attemptKey)
	if err != nil {
		return 0, err
	}
	if attempts == 1 {
		ttl, ttlErr := g.Redis().TTL(ctx, codeKey)
		if ttlErr != nil || ttl <= 0 {
			ttl = fallbackTTL
		}
		if ttl <= 0 {
			ttl = defaultCodeExpireSeconds
		}
		_, _ = g.Redis().Expire(ctx, attemptKey, ttl)
	}
	return attempts, nil
}

func releaseVerifyAttempts(ctx context.Context, key string) {
	if strings.TrimSpace(key) == "" {
		return
	}
	_, _ = g.Redis().Del(ctx, key)
}

// consumeCode 用 Lua 原子比较并删除验证码，避免并发误销毁。
func consumeCode(ctx context.Context, codeKey, code string) (bool, error) {
	result, err := g.Redis().Eval(ctx, `
local current = redis.call("GET", KEYS[1])
if not current then
  return 0
end
if tostring(current) ~= tostring(ARGV[1]) then
  return 0
end
redis.call("DEL", KEYS[1])
return 1
`, 1, []string{codeKey}, []any{code})
	if err != nil {
		return false, err
	}
	if result == nil || result.IsNil() {
		return false, nil
	}
	return result.Int() == 1, nil
}

func verifyFailed(ctx context.Context, attemptKey, codeKey string, maxAttempts int64) error {
	attempts, err := increaseVerifyAttempts(ctx, attemptKey, codeKey, codeExpireSeconds(ctx))
	if err != nil {
		return gerror.New("验证码校验失败")
	}
	if maxAttempts > 0 && attempts >= maxAttempts {
		return gerror.New("验证码错误次数过多，请重新获取")
	}
	return gerror.New("验证码错误或已过期")
}

// ----- key builder & config -----

func buildCodeKey(scene, phone string) string {
	return strings.Join([]string{codeKeyPrefix, strings.TrimSpace(scene), strings.TrimSpace(phone)}, ":")
}

func buildLimitKey(scene, phone string) string {
	return strings.Join([]string{limitKeyPrefix, strings.TrimSpace(scene), strings.TrimSpace(phone)}, ":")
}

func buildVerifyAttemptKey(scene, phone string) string {
	return strings.Join([]string{verifyAttemptKeyPrefix, strings.TrimSpace(scene), strings.TrimSpace(phone)}, ":")
}

func codeExpireSeconds(ctx context.Context) int64 {
	value := g.Cfg().MustGet(ctx, "sms.codeExpireSeconds").Int64()
	if value <= 0 {
		return defaultCodeExpireSeconds
	}
	return value
}

func codeLimitSeconds(ctx context.Context) int64 {
	value := g.Cfg().MustGet(ctx, "sms.limitSeconds").Int64()
	if value <= 0 {
		return defaultLimitSeconds
	}
	return value
}

func codeVerifyMaxAttempts(ctx context.Context) int64 {
	value := g.Cfg().MustGet(ctx, "sms.verifyMaxAttempts").Int64()
	if value <= 0 {
		return defaultVerifyMaxAttempts
	}
	return value
}

// generateCode 生成 6 位数字验证码。
func generateCode() string {
	const length = 6
	buf := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			buf = append(buf, '0')
			continue
		}
		buf = append(buf, byte('0'+v.Int64()))
	}
	return string(buf)
}
