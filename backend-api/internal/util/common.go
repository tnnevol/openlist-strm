package util

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/thanhpk/randstr"
)

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[\w.-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsValidUsername(username string) bool {
	// 用户名至少3个字符，最多20个字符，只能包含字母、数字、下划线
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return re.MatchString(username)
}

func IsStrongPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	hasUpper, hasLower, hasDigit := false, false, false
	for _, c := range pw {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

// GenerateVerificationCode 生成6位数字+字母验证码（不区分大小写）
func GenerateVerificationCode() string {
	code := strings.ToUpper(randstr.String(6))
	log.Println("[验证码生成] code:", code)
	return code
}
