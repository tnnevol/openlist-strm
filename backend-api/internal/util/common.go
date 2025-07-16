package util

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// GetPageParams 获取分页参数，默认 page=1, pageSize=10
func GetPageParams(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
		if page < 1 {
			page = 1
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
		if pageSize < 1 {
			pageSize = 10
		}
	}
	return page, pageSize
}

// Paginate 对任意 slice 做分页，返回当前页数据和总数
func Paginate[T any](list []T, page, pageSize int) ([]T, int) {
	total := len(list)
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * pageSize
	if start >= total {
		return []T{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return list[start:end], total
}

// 从claims中提取user_id，兼容多种类型
func ExtractUserIDFromClaims(claims interface{}) int {
	if claims == nil {
		return 0
	}
	var m map[string]interface{}
	switch v := claims.(type) {
	case map[string]interface{}:
		m = v
	case gin.H:
		m = v
	case jwt.MapClaims:
		m = v
	default:
		if mv, ok := v.(map[string]interface{}); ok {
			m = mv
		}
	}
	if m != nil {
		if v, ok := m["user_id"]; ok {
			switch vv := v.(type) {
			case float64:
				return int(vv)
			case int:
				return vv
			}
		}
	}
	return 0
}

// parseEnabled 将字符串/数字/布尔值的 0/1 转为 bool
func ParseEnabled(val interface{}) bool {
    switch v := val.(type) {
    case string:
        return v == "1"
    case float64:
        return int(v) == 1
    case int:
        return v == 1
    case bool:
        return v
    default:
        return false
    }
}

func Bool2Int(val bool) int {
	var i int
	if val {
		i = 1
	} else {
		i = 0
	}
	return i
}

