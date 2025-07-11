package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
)

// 统一错误码定义
const (
	CodeSuccess         = 200
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeTokenExpired    = 40101 // token过期专用错误码
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
			// 新增日志
			println("[AuthMiddleware] 路径", path, "在白名单，直接放行")
			c.Next()
			return
		}

		// 检查swagger路径（前缀匹配）
		if strings.HasPrefix(path, "/swagger/") {
			println("[AuthMiddleware] swagger路径放行")
			c.Next()
			return
		}

		// 解析token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			println("[AuthMiddleware] 未获取到token，拒绝")
			Unauthorized(c, "未登录或token缺失")
			c.Abort()
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// 检查token是否在黑名单中
		blacklist := service.GetTokenBlacklist()
		if blacklist.IsBlacklisted(tokenStr) {
			println("[AuthMiddleware] token在黑名单，拒绝")
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
			// 自动区分token过期和其他无效
			if errors.Is(err, jwt.ErrTokenExpired) {
				println("[AuthMiddleware] token已过期（解析时），返回过期状态")
				TokenExpired(c, "token已过期")
				c.Abort()
				return
			}
			// 兼容其他过期错误类型
			if err != nil && (strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "过期")) {
				println("[AuthMiddleware] token已过期（其他错误类型），返回过期状态")
				TokenExpired(c, "token已过期")
				c.Abort()
				return
			}
			println("[AuthMiddleware] token解析失败或无效，拒绝", err, "类型:", fmt.Sprintf("%T", err))
			Unauthorized(c, "token无效或已过期")
			c.Abort()
			return
		}

		// 添加token解析成功的详细日志
		println("[AuthMiddleware] token解析成功，claims:", fmt.Sprintf("%+v", claims))
		if exp, ok := claims["exp"].(float64); ok {
			println("[AuthMiddleware] token exp字段:", int64(exp), "当前时间:", time.Now().Unix())
		}
		if iat, ok := claims["iat"].(float64); ok {
			println("[AuthMiddleware] token iat字段:", int64(iat))
		}

		// 检查token是否已过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				println("[AuthMiddleware] token已过期，返回过期状态")
				TokenExpired(c, "token已过期")
				c.Abort()
				return
			}
		}

		// 新增：校验token_invalid_before
		if iat, ok := claims["iat"].(float64); ok {
			// 获取用户名
			var username string
			if v, ok := claims["username"]; ok {
				switch vv := v.(type) {
				case string:
					username = vv
				case []byte:
					username = string(vv)
				default:
					username = ""
				}
			}
			if username != "" {
				db, err := service.GetDB()
				if err == nil {
					user, err := service.GetUserByUsername(db, username)
					if err != nil {
						println("[AuthMiddleware] 查找用户失败，拒绝")
						Unauthorized(c, "token已失效，请重新登录")
						c.Abort()
						return
					}
					
					// 安全地打印日志，避免空指针异常
					var tokenInvalidBeforeStr string
					var validStr string
					if user != nil && user.TokenInvalidBefore.Valid {
						tokenInvalidBeforeStr = fmt.Sprintf("%d", user.TokenInvalidBefore.Time.Unix())
						validStr = "true"
					} else {
						tokenInvalidBeforeStr = "nil"
						validStr = "false"
					}
					println("[AuthMiddleware] 校验token_invalid_before username=", username, "iat=", int64(iat), "token_invalid_before=", tokenInvalidBeforeStr, "valid=", validStr, "查找err=", err)
					
					if user == nil {
						println("[AuthMiddleware] 用户对象为空，拒绝")
						Unauthorized(c, "token已失效，请重新登录")
						c.Abort()
						return
					}
					
					if !user.TokenInvalidBefore.Valid {
						println("[AuthMiddleware] token_invalid_before无效，拒绝")
						TokenExpired(c, "token已过期")
						c.Abort()
						return
					}
					if int64(iat) < user.TokenInvalidBefore.Time.Unix() {
						println("[AuthMiddleware] token签发时间早于token_invalid_before，拒绝")
						println("[AuthMiddleware] 准备返回TokenExpired，code=40101")
						TokenExpired(c, "token已过期")
						println("[AuthMiddleware] TokenExpired已调用")
						c.Abort()
						return
					}
				}
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
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, CodeBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, CodeUnauthorized, message)
}

// TokenExpired token过期错误
func TokenExpired(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeTokenExpired,
		Message: message,
	})
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, CodeInternalError, message)
}

// TooManyRequests 429错误
func TooManyRequests(c *gin.Context, message string) {
	Error(c, CodeTooManyRequests, message)
}

// ValidationError 参数验证错误
func ValidationError(c *gin.Context, field string) {
	BadRequest(c, "参数错误："+field)
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context) {
	InternalServerError(c, "数据库操作失败")
}
