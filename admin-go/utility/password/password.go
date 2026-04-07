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
