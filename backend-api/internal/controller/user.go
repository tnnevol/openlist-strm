package controller

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/middleware"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"github.com/tnnevol/openlist-strm/backend-api/internal/util"
	"go.uber.org/zap"
)

// SendCode godoc
// @Summary      发送验证码
// @Description  发送邮箱验证码
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        email  body  object{email=string}  true  "邮箱"
// @Success      200    {object}  model.Response
// @Router       /user/send-code [post]
func SendCode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/send-code called - 请求入口")
		var req struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /user/send-code 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "邮箱不能为空")
			return
		}
		logger.Info("[API] /user/send-code 参数校验", zap.String("email", req.Email))
		if !util.IsValidEmail(req.Email) {
			logger.Error("[API] /user/send-code 邮箱格式不正确", zap.String("email", req.Email))
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		logger.Info("[API] /user/send-code 调用service.SendCode", zap.String("email", req.Email))
		_, err := service.SendCode(db, req.Email)
		if err != nil {
			logger.Error("[API] /user/send-code service.SendCode失败", zap.Error(err))
			if err.Error() == "用户已注册" {
				middleware.BadRequest(c, "用户已注册")
			} else if err.Error() == "邮箱未注册，请先注册" {
				middleware.BadRequest(c, "邮箱未注册，请先注册")
			} else {
				middleware.InternalServerError(c, "发送验证码失败")
			}
			return
		}
		logger.Info("[API] /user/send-code 成功")
		middleware.SuccessWithMessage(c, "验证码已发送", nil)
	}
}

// Register godoc
// @Summary      注册并激活
// @Description  用户注册并激活账户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        register body object{email=string,username=string,password=string,confirmPassword=string,code=string} true "注册信息"
// @Success      200    {object}  model.Response
// @Router       /user/register [post]
func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/register called - 请求入口")
		
		// 记录请求体内容
		body, err := c.GetRawData()
		if err != nil {
			logger.Error("[API] /user/register 读取请求体失败", zap.Error(err))
			middleware.InternalServerError(c, "读取请求数据失败")
			return
		}
		logger.Info("[API] /user/register 请求体内容", zap.String("body", string(body)))
		
		// 重新设置请求体，因为GetRawData会消费掉请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		
		var req struct {
			Email           string `json:"email" binding:"required"`
			Username        string `json:"username" binding:"required"`
			Password        string `json:"password" binding:"required"`
			ConfirmPassword string `json:"confirmPassword" binding:"required"`
			Code            string `json:"code" binding:"required"`
		}
		
		logger.Info("[API] /user/register 开始参数绑定")
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /user/register 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "邮箱、用户名、密码、确认密码和验证码不能为空")
			return
		}
		
		logger.Info("[API] /user/register 参数绑定成功", 
			zap.String("email", req.Email),
			zap.String("username", req.Username),
			zap.String("password", "***"),
			zap.String("confirmPassword", "***"),
			zap.String("code", req.Code))
		
		// 参数校验
		logger.Info("[API] /user/register 开始参数校验")
		if !util.IsValidEmail(req.Email) {
			logger.Error("[API] /user/register 邮箱格式不正确", zap.String("email", req.Email))
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		if !util.IsValidUsername(req.Username) {
			logger.Error("[API] /user/register 用户名格式不正确", zap.String("username", req.Username))
			middleware.ValidationError(c, "用户名格式不正确，只能包含字母、数字、下划线，长度3-20位")
			return
		}
		if !util.IsStrongPassword(req.Password) {
			logger.Error("[API] /user/register 密码强度不够", zap.String("password", "***"))
			middleware.ValidationError(c, "密码需8位以上，含大小写字母和数字")
			return
		}
		if req.Password != req.ConfirmPassword {
			logger.Error("[API] /user/register 两次密码不一致")
			middleware.ValidationError(c, "两次输入的密码不一致")
			return
		}
		if len(req.Code) != 6 {
			logger.Error("[API] /user/register 验证码格式不正确", zap.String("code", req.Code))
			middleware.ValidationError(c, "验证码格式不正确")
			return
		}
		
		logger.Info("[API] /user/register 参数校验通过，调用service.ActivateUserWithPassword")
		err = service.ActivateUserWithPassword(db, req.Email, req.Username, req.Password, req.Code)
		if err != nil {
			logger.Error("[API] /user/register service.ActivateUserWithPassword失败", zap.Error(err))
			if err == sql.ErrNoRows {
				middleware.BadRequest(c, "验证码无效或已过期")
			} else {
				middleware.InternalServerError(c, "注册失败")
			}
			return
		}
		
		logger.Info("[API] /user/register 注册成功")
		middleware.SuccessWithMessage(c, "注册成功，请登录", nil)
	}
}

// Login godoc
// @Summary      用户登录
// @Description  用户登录接口
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        login body object{username=string,password=string} true "登录信息"
// @Success      200    {object}  model.Response
// @Router       /user/login [post]
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/login called - 请求入口")
		
		// 记录请求体内容
		body, err := c.GetRawData()
		if err != nil {
			logger.Error("[API] /user/login 读取请求体失败", zap.Error(err))
			middleware.InternalServerError(c, "读取请求数据失败")
			return
		}
		logger.Info("[API] /user/login 请求体内容", zap.String("body", string(body)))
		
		// 重新设置请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		
		logger.Info("[API] /user/login 开始参数绑定")
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /user/login 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "用户名和密码不能为空")
			return
		}
		
		logger.Info("[API] /user/login 参数绑定成功", 
			zap.String("username", req.Username),
			zap.String("password", "***"))
		
		if !util.IsValidUsername(req.Username) {
			logger.Error("[API] /user/login 用户名格式不正确", zap.String("username", req.Username))
			middleware.ValidationError(c, "用户名格式不正确，只能包含字母、数字、下划线，长度3-20位")
			return
		}
		
		logger.Info("[API] /user/login 调用service.LoginUser")
		token, err := service.LoginUser(db, req.Username, req.Password)
		if err != nil {
			logger.Error("[API] /user/login service.LoginUser失败", zap.Error(err))
			middleware.Unauthorized(c, "用户名或密码错误")
			return
		}
		if token == "not_activated" {
			logger.Error("[API] /user/login 账户未激活", zap.String("username", req.Username))
			middleware.Unauthorized(c, "账户未激活，请先验证邮箱")
			return
		}
		if token == "locked" {
			logger.Error("[API] /user/login 账户锁定", zap.String("username", req.Username))
			middleware.TooManyRequests(c, "账户暂时锁定，请稍后再试")
			return
		}
		
		logger.Info("[API] /user/login 登录成功", zap.String("username", req.Username))
		middleware.Success(c, gin.H{"accessToken": token})
	}
}

// ForgotPasswordSendCode godoc
// @Summary      忘记密码-发送验证码
// @Description  已注册用户发送重置密码验证码
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        email  body  object{email=string}  true  "邮箱"
// @Success      200    {object}  model.Response
// @Router       /user/forgot-password/send-code [post]
func ForgotPasswordSendCode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/forgot-password/send-code called - 请求入口")
		
		// 记录请求体内容
		body, err := c.GetRawData()
		if err != nil {
			logger.Error("[API] /user/forgot-password/send-code 读取请求体失败", zap.Error(err))
			middleware.InternalServerError(c, "读取请求数据失败")
			return
		}
		logger.Info("[API] /user/forgot-password/send-code 请求体内容", zap.String("body", string(body)))
		
		// 重新设置请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		
		var req struct {
			Email string `json:"email" binding:"required"`
		}
		
		logger.Info("[API] /user/forgot-password/send-code 开始参数绑定")
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /user/forgot-password/send-code 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "邮箱不能为空")
			return
		}
		
		logger.Info("[API] /user/forgot-password/send-code 参数绑定成功", zap.String("email", req.Email))
		
		if !util.IsValidEmail(req.Email) {
			logger.Error("[API] /user/forgot-password/send-code 邮箱格式不正确", zap.String("email", req.Email))
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		
		logger.Info("[API] /user/forgot-password/send-code 调用service.ForgotPasswordSendCode")
		_, err = service.ForgotPasswordSendCode(db, req.Email)
		if err != nil {
			logger.Error("[API] /user/forgot-password/send-code service.ForgotPasswordSendCode失败", zap.Error(err))
			if err.Error() == "用户未注册" {
				middleware.BadRequest(c, "用户未注册")
			} else if err.Error() == "账户未激活" {
				middleware.BadRequest(c, "账户未激活")
			} else {
				middleware.InternalServerError(c, "发送验证码失败")
			}
			return
		}
		
		logger.Info("[API] /user/forgot-password/send-code 发送成功")
		middleware.SuccessWithMessage(c, "验证码已发送", nil)
	}
}

// ForgotPasswordReset godoc
// @Summary      忘记密码-重置密码
// @Description  校验验证码并重置密码
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        reset body object{email=string,code=string,newPassword=string,confirmPassword=string} true "重置信息"
// @Success      200    {object}  model.Response
// @Router       /user/forgot-password/reset [post]
func ForgotPasswordReset(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/forgot-password/reset called - 请求入口")
		
		// 记录请求体内容
		body, err := c.GetRawData()
		if err != nil {
			logger.Error("[API] /user/forgot-password/reset 读取请求体失败", zap.Error(err))
			middleware.InternalServerError(c, "读取请求数据失败")
			return
		}
		logger.Info("[API] /user/forgot-password/reset 请求体内容", zap.String("body", string(body)))
		
		// 重新设置请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		
		var req struct {
			Email           string `json:"email" binding:"required"`
			Code            string `json:"code" binding:"required"`
			NewPassword     string `json:"newPassword" binding:"required"`
			ConfirmPassword string `json:"confirmPassword" binding:"required"`
		}
		
		logger.Info("[API] /user/forgot-password/reset 开始参数绑定")
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /user/forgot-password/reset 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "邮箱、验证码、新密码和确认密码不能为空")
			return
		}
		
		logger.Info("[API] /user/forgot-password/reset 参数绑定成功", 
			zap.String("email", req.Email),
			zap.String("code", req.Code),
			zap.String("newPassword", "***"),
			zap.String("confirmPassword", "***"))
		
		if !util.IsValidEmail(req.Email) {
			logger.Error("[API] /user/forgot-password/reset 邮箱格式不正确", zap.String("email", req.Email))
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		if !util.IsStrongPassword(req.NewPassword) {
			logger.Error("[API] /user/forgot-password/reset 新密码强度不够", zap.String("newPassword", "***"))
			middleware.ValidationError(c, "新密码需8位以上，含大小写字母和数字")
			return
		}
		if req.NewPassword != req.ConfirmPassword {
			logger.Error("[API] /user/forgot-password/reset 两次密码不一致")
			middleware.ValidationError(c, "两次输入的密码不一致")
			return
		}
		if len(req.Code) != 6 {
			logger.Error("[API] /user/forgot-password/reset 验证码格式不正确", zap.String("code", req.Code))
			middleware.ValidationError(c, "验证码格式不正确")
			return
		}
		
		logger.Info("[API] /user/forgot-password/reset 调用service.ForgotPasswordReset")
		err = service.ForgotPasswordReset(db, req.Email, req.Code, req.NewPassword)
		if err != nil {
			logger.Error("[API] /user/forgot-password/reset service.ForgotPasswordReset失败", zap.Error(err))
			if err.Error() == "验证码无效或已过期" {
				middleware.BadRequest(c, "验证码无效或已过期")
			} else {
				middleware.InternalServerError(c, "重置密码失败")
			}
			return
		}

		// 新增：将当前token加入黑名单
		tokenStr := c.GetHeader("Authorization")
		if tokenStr != "" {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			claims, ok := c.Get("claims")
			var expireTime time.Time
			if ok {
				switch m := claims.(type) {
				case map[string]interface{}:
					if v, ok := m["exp"].(float64); ok {
						expireTime = time.Unix(int64(v), 0)
					}
				case jwt.MapClaims:
					if v, ok := m["exp"].(float64); ok {
						expireTime = time.Unix(int64(v), 0)
					}
				}
			}
			if !expireTime.IsZero() {
				blacklist := service.GetTokenBlacklist()
				blacklist.AddToBlacklist(tokenStr, expireTime)
			}
		}

		logger.Info("[API] /user/forgot-password/reset 重置成功，token已加入黑名单")
		middleware.SuccessWithMessage(c, "密码重置成功，请登录", nil)
	}
}

// Logout godoc
// @Summary      用户登出
// @Description  用户登出接口，将token加入黑名单使其立即失效
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Success      200    {object}  model.Response
// @Router       /user/logout [post]
func Logout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/logout called - 请求入口")
		
		// 获取token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			logger.Error("[API] /user/logout 未获取到Authorization header")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		
		// 获取用户信息
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /user/logout 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		
		// 解析用户信息
		var username string
		var userID int
		var expireTime time.Time
		
		switch m := claims.(type) {
		case map[string]interface{}:
			if v, ok := m["username"]; ok {
				switch vv := v.(type) {
				case string:
					username = vv
				case []byte:
					username = string(vv)
				default:
					username = fmt.Sprintf("%v", vv)
				}
			}
			if v, ok := m["user_id"]; ok {
				switch vv := v.(type) {
				case float64:
					userID = int(vv)
				case int:
					userID = vv
				default:
					userID = 0
				}
			}
			if v, ok := m["exp"].(float64); ok {
				expireTime = time.Unix(int64(v), 0)
			}
		case jwt.MapClaims:
			if v, ok := m["username"]; ok {
				switch vv := v.(type) {
				case string:
					username = vv
				case []byte:
					username = string(vv)
				default:
					username = fmt.Sprintf("%v", vv)
				}
			}
			if v, ok := m["user_id"]; ok {
				switch vv := v.(type) {
				case float64:
					userID = int(vv)
				case int:
					userID = vv
				default:
					userID = 0
				}
			}
			if v, ok := m["exp"].(float64); ok {
				expireTime = time.Unix(int64(v), 0)
			}
		default:
			logger.Error("[API] /user/logout claims 类型异常", zap.Any("claims_type", fmt.Sprintf("%T", claims)))
			middleware.InternalServerError(c, "用户信息解析失败")
			return
		}
		
		logger.Info("[API] /user/logout 用户登出", 
			zap.String("username", username),
			zap.Int("user_id", userID),
			zap.Time("expireTime", expireTime))
		
		// 将token添加到黑名单
		blacklist := service.GetTokenBlacklist()
		blacklist.AddToBlacklist(tokenStr, expireTime)
		
		logger.Info("[API] /user/logout 登出成功，token已加入黑名单", zap.String("username", username))
		middleware.SuccessWithMessage(c, "登出成功", nil)
	}
}

// TokenBlacklistStatus godoc
// @Summary      获取Token黑名单状态
// @Description  获取当前token黑名单的状态信息（仅用于监控）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Success      200    {object}  model.Response{data=object{blacklistSize=int}}
// @Router       /user/token-blacklist-status [get]
func TokenBlacklistStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/token-blacklist-status called - 请求入口")
		
		// 获取黑名单状态
		blacklist := service.GetTokenBlacklist()
		blacklistSize := blacklist.GetBlacklistSize()
		
		logger.Info("[API] /user/token-blacklist-status 获取状态", zap.Int("blacklistSize", blacklistSize))
		
		middleware.Success(c, gin.H{
			"blacklistSize": blacklistSize,
			"timestamp":     time.Now().Unix(),
		})
	}
}

// UserInfo godoc
// @Summary 获取用户信息
// @Description 获取用户基本信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} model.Response{data=model.UserInfoResponse}
// @Router /user/info [get]
func UserInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/info called - 请求入口")
		
		// 记录请求头信息
		authHeader := c.GetHeader("Authorization")
		logger.Info("[API] /user/info Authorization header", zap.String("authorization", authHeader))
		
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /user/info 未获取到claims")
			c.JSON(401, gin.H{"code": 401, "message": "未登录", "data": nil})
			return
		}
		logger.Info("[API] /user/info claims原始内容", zap.Any("claims", claims))

		var username string
		switch m := claims.(type) {
		case map[string]interface{}:
			logger.Info("[API] /user/info claims类型为map[string]interface{}")
			if v, ok := m["username"]; ok {
				switch vv := v.(type) {
				case string:
					username = vv
				case []byte:
					username = string(vv)
				default:
					username = fmt.Sprintf("%v", vv)
				}
			} else {
				logger.Error("[API] /user/info claims中未找到username字段")
			}
		case jwt.MapClaims:
			logger.Info("[API] /user/info claims类型为jwt.MapClaims")
			if v, ok := m["username"]; ok {
				switch vv := v.(type) {
				case string:
					username = vv
				case []byte:
					username = string(vv)
				default:
					username = fmt.Sprintf("%v", vv)
				}
			} else {
				logger.Error("[API] /user/info claims中未找到username字段")
			}
		default:
			logger.Error("[API] /user/info claims 类型异常", zap.Any("claims_type", fmt.Sprintf("%T", claims)))
		}

		logger.Info("[API] /user/info 解析claims完成", zap.String("username", username))
		
		if username == "" {
			logger.Error("[API] /user/info username为空，无法获取用户信息")
			c.JSON(500, gin.H{"code": 500, "message": "用户信息解析失败", "data": nil})
			return
		}
		
		logger.Info("[API] /user/info 调用service.GetUserBaseInfo")
		user, err := service.GetUserBaseInfo(db, username)
		if err != nil {
			logger.Error("[API] /user/info 获取用户信息失败", zap.Error(err), zap.String("username", username))
			c.JSON(500, gin.H{"code": 500, "message": "获取用户信息失败", "data": nil})
			return
		}
		
		logger.Info("[API] /user/info 获取用户信息成功", zap.Any("user", user))
		
		resp := gin.H{
			"code":    200,
			"message": "success",
			"data": gin.H{
				"roles":    []string{},
				"realName": "Super",
			},
		}
		for k, v := range user {
			if k == "createdAt" || k == "isActive" {
				continue
			}
			resp["data"].(gin.H)[k] = v
		}
		
		logger.Info("[API] /user/info 返回响应", zap.Any("response", resp))
		c.JSON(200, resp)
	}
}
