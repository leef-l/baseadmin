package password

import (
	"strings"

	"github.com/gogf/gf/v2/crypto/gsha256"
	"golang.org/x/crypto/bcrypt"
)

// Hash 使用 bcrypt 生成密码摘要。
func Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify 校验密码，兼容历史 SHA-256 存量摘要。
func Verify(stored, plain string) bool {
	stored = normalizeStoredHash(stored)
	if stored == "" {
		return false
	}
	if isBCryptHash(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil
	}
	return stored == gsha256.Encrypt(plain)
}

// 预计算的 bcrypt hash，用于 DummyVerify 的恒定时间消耗。
// 即使 bcrypt.GenerateFromPassword 理论上失败，硬编码兜底也能保证非 nil。
var dummyHash = func() []byte {
	h, err := bcrypt.GenerateFromPassword([]byte("dummy-timing-pad"), bcrypt.DefaultCost)
	if err != nil {
		h = []byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
	}
	return h
}()

// DummyVerify 执行一次虚拟 bcrypt 验证，用于防止时序攻击探测用户是否存在。
func DummyVerify(plain string) {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plain))
}

// NeedsRehash 判断是否需要升级为 bcrypt 摘要。
func NeedsRehash(stored string) bool {
	stored = normalizeStoredHash(stored)
	if stored == "" {
		return false
	}
	if !isBCryptHash(stored) {
		return true
	}
	cost, err := bcrypt.Cost([]byte(stored))
	return err != nil || cost < bcrypt.DefaultCost
}

func isBCryptHash(stored string) bool {
	stored = normalizeStoredHash(stored)
	return strings.HasPrefix(stored, "$2a$") ||
		strings.HasPrefix(stored, "$2b$") ||
		strings.HasPrefix(stored, "$2y$")
}

func normalizeStoredHash(stored string) string {
	return strings.TrimSpace(stored)
}
