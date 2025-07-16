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
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"go.uber.org/zap"
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
type Response[T any] struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    T `json:"data,omitempty"`
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
			logger.Info("[AuthMiddleware] 路径在白名单，直接放行", zap.String("path", path))
			c.Next()
			return
		}

		// 检查swagger路径（前缀匹配）
		if strings.HasPrefix(path, "/swagger/") {
			logger.Info("[AuthMiddleware] swagger路径放行", zap.String("path", path))
			c.Next()
			return
		}

		// 解析token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			logger.Info("[AuthMiddleware] 未获取到token，拒绝", zap.String("path", path))
			Unauthorized(c, "未登录或token缺失")
			c.Abort()
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// 检查token是否在黑名单中
		blacklist := service.GetTokenBlacklist()
		if blacklist.IsBlacklisted(tokenStr) {
			logger.Info("[AuthMiddleware] token在黑名单，拒绝", zap.String("path", path))
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
				logger.Info("[AuthMiddleware] token已过期（解析时），返回过期状态", zap.String("path", path))
				TokenExpired(c, "token已过期")
				c.Abort()
				return
			}
			// 兼容其他过期错误类型
			if err != nil && (strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "过期")) {
				logger.Info("[AuthMiddleware] token已过期（其他错误类型），返回过期状态", zap.String("path", path))
				TokenExpired(c, "token已过期")
				c.Abort()
				return
			}
			logger.Info("[AuthMiddleware] token解析失败或无效，拒绝", zap.String("path", path), zap.Error(err), zap.String("errType", fmt.Sprintf("%T", err)))
			Unauthorized(c, "token无效")
			c.Abort()
			return
		}

		// 添加token解析成功的详细日志
		logger.Info("[AuthMiddleware] token解析成功", zap.String("path", path), zap.Any("claims", claims))
		if exp, ok := claims["exp"].(float64); ok {
			logger.Info("[AuthMiddleware] token exp字段", zap.Int64("exp", int64(exp)), zap.Int64("now", time.Now().Unix()))
		}
		if iat, ok := claims["iat"].(float64); ok {
			logger.Info("[AuthMiddleware] token iat字段", zap.Int64("iat", int64(iat)))
		}

		// 检查token是否已过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				logger.Info("[AuthMiddleware] token已过期，返回过期状态", zap.String("path", path))
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
						logger.Info("[AuthMiddleware] 查找用户失败，拒绝", zap.String("path", path), zap.Error(err))
						Unauthorized(c, "token无效，请重新登录")
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
					logger.Info("[AuthMiddleware] 校验token_invalid_before username=", zap.String("username", username), zap.Int64("iat", int64(iat)), zap.String("token_invalid_before", tokenInvalidBeforeStr), zap.String("valid", validStr), zap.Error(err))
					
					if user == nil {
						logger.Info("[AuthMiddleware] 用户对象为空，拒绝", zap.String("path", path))
						Unauthorized(c, "token无效，请重新登录")
						c.Abort()
						return
					}
					
					if !user.TokenInvalidBefore.Valid {
						// token_invalid_before 无效，说明没有强制失效要求，直接通过
						c.Set("claims", claims)
						c.Next()
						return
					}
					if int64(iat) < user.TokenInvalidBefore.Time.Unix() {
						logger.Info("[AuthMiddleware] token签发时间早于token_invalid_before，拒绝", zap.String("path", path))
						logger.Info("[AuthMiddleware] 准备返回TokenExpired，code=40101")
						TokenExpired(c, "token已过期")
						logger.Info("[AuthMiddleware] TokenExpired已调用")
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
	SuccessWithMessage(c, "success", data)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response[any]{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response[any]{
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
	c.JSON(http.StatusOK, Response[any]{
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
