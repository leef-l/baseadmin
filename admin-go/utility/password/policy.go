package password

import (
	"errors"
	"strings"
	"unicode"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 64
)

func ValidatePolicy(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("密码不能为空")
	}
	if strings.ContainsAny(value, " \t\r\n") {
		return errors.New("密码不能包含空白字符")
	}
	length := len([]rune(value))
	if length < MinPasswordLength || length > MaxPasswordLength {
		return errors.New("密码长度需为8-64位")
	}

	hasLetter := false
	hasDigit := false
	for _, item := range value {
		switch {
		case unicode.IsLetter(item):
			hasLetter = true
		case unicode.IsDigit(item):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("密码必须同时包含字母和数字")
	}
	return nil
}
