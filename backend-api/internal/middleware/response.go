package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
)

// 统一错误码定义
const (
	CodeSuccess         = 200
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeNotFound        = 404
	CodeTooManyRequests = 429
	CodeInternalError   = 500
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var WhiteList = map[string]bool{
	"/user/send-code": true,
	"/user/register": true,
	"/user/login": true,
	"/user/forgot-password/send-code": true,
	"/user/forgot-password/reset": true,
	"/swagger/": true,
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// 检查是否在白名单中（精确匹配）
		if WhiteList[path] {
			c.Next()
			return
		}
		
		// 检查swagger路径（前缀匹配）
		if strings.HasPrefix(path, "/swagger/") {
			c.Next()
			return
		}
		
		// 解析token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			Unauthorized(c, "未登录或token缺失")
			c.Abort()
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		
		// 检查token是否在黑名单中
		blacklist := service.GetTokenBlacklist()
		if blacklist.IsBlacklisted(tokenStr) {
			Unauthorized(c, "token已失效，请重新登录")
			c.Abort()
			return
		}
		
		jwtKey := []byte(os.Getenv("JWT_SECRET"))
		if len(jwtKey) == 0 {
			jwtKey = []byte("secret")
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			Unauthorized(c, "token无效或已过期")
			c.Abort()
			return
		}
		
		// 检查token是否已过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				Unauthorized(c, "token已过期")
				c.Abort()
				return
			}
		}
		
		c.Set("claims", claims)
		c.Next()
	}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// TooManyRequests 429错误
func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, message)
}

// ValidationError 参数验证错误
func ValidationError(c *gin.Context, field string) {
	BadRequest(c, "参数错误："+field)
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context) {
	InternalServerError(c, "数据库操作失败")
}
