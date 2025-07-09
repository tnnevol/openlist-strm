package controller

import (
	"database/sql"
	"fmt"

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
// @Param        register body object{email=string,username=string,password=string,confirm_password=string,code=string} true "注册信息"
// @Success      200    {object}  model.Response
// @Router       /user/register [post]
func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/register called")
		var req struct {
			Email           string `json:"email" binding:"required"`
			Username        string `json:"username" binding:"required"`
			Password        string `json:"password" binding:"required"`
			ConfirmPassword string `json:"confirm_password" binding:"required"`
			Code            string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ValidationError(c, "邮箱、用户名、密码、确认密码和验证码不能为空")
			return
		}
		if !util.IsValidEmail(req.Email) {
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		if !util.IsValidUsername(req.Username) {
			middleware.ValidationError(c, "用户名格式不正确，只能包含字母、数字、下划线，长度3-20位")
			return
		}
		if !util.IsStrongPassword(req.Password) {
			middleware.ValidationError(c, "密码需8位以上，含大小写字母和数字")
			return
		}
		if req.Password != req.ConfirmPassword {
			middleware.ValidationError(c, "两次输入的密码不一致")
			return
		}
		if len(req.Code) != 6 {
			middleware.ValidationError(c, "验证码格式不正确")
			return
		}
		err := service.ActivateUserWithPassword(db, req.Email, req.Username, req.Password, req.Code)
		if err != nil {
			if err == sql.ErrNoRows {
				middleware.BadRequest(c, "验证码无效或已过期")
			} else {
				middleware.InternalServerError(c, "注册失败")
			}
			return
		}
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
		logger.Info("[API] /user/login called")
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ValidationError(c, "用户名和密码不能为空")
			return
		}
		if !util.IsValidUsername(req.Username) {
			middleware.ValidationError(c, "用户名格式不正确，只能包含字母、数字、下划线，长度3-20位")
			return
		}
		token, err := service.LoginUser(db, req.Username, req.Password)
		if err != nil {
			middleware.Unauthorized(c, "用户名或密码错误")
			return
		}
		if token == "not_activated" {
			middleware.Unauthorized(c, "账户未激活，请先验证邮箱")
			return
		}
		if token == "locked" {
			middleware.TooManyRequests(c, "账户暂时锁定，请稍后再试")
			return
		}
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
		logger.Info("[API] /user/forgot-password/send-code called")
		var req struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ValidationError(c, "邮箱不能为空")
			return
		}
		if !util.IsValidEmail(req.Email) {
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		_, err := service.ForgotPasswordSendCode(db, req.Email)
		if err != nil {
			if err.Error() == "用户未注册" {
				middleware.BadRequest(c, "用户未注册")
			} else if err.Error() == "账户未激活" {
				middleware.BadRequest(c, "账户未激活")
			} else {
				middleware.InternalServerError(c, "发送验证码失败")
			}
			return
		}
		middleware.SuccessWithMessage(c, "验证码已发送", nil)
	}
}

// ForgotPasswordReset godoc
// @Summary      忘记密码-重置密码
// @Description  校验验证码并重置密码
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        reset body object{email=string,code=string,new_password=string,confirm_password=string} true "重置信息"
// @Success      200    {object}  model.Response
// @Router       /user/forgot-password/reset [post]
func ForgotPasswordReset(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /user/forgot-password/reset called")
		var req struct {
			Email           string `json:"email" binding:"required"`
			Code            string `json:"code" binding:"required"`
			NewPassword     string `json:"new_password" binding:"required"`
			ConfirmPassword string `json:"confirm_password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ValidationError(c, "邮箱、验证码、新密码和确认密码不能为空")
			return
		}
		if !util.IsValidEmail(req.Email) {
			middleware.ValidationError(c, "邮箱格式不正确")
			return
		}
		if !util.IsStrongPassword(req.NewPassword) {
			middleware.ValidationError(c, "新密码需8位以上，含大小写字母和数字")
			return
		}
		if req.NewPassword != req.ConfirmPassword {
			middleware.ValidationError(c, "两次输入的密码不一致")
			return
		}
		if len(req.Code) != 6 {
			middleware.ValidationError(c, "验证码格式不正确")
			return
		}
		err := service.ForgotPasswordReset(db, req.Email, req.Code, req.NewPassword)
		if err != nil {
			if err.Error() == "验证码无效或已过期" {
				middleware.BadRequest(c, "验证码无效或已过期")
			} else {
				middleware.InternalServerError(c, "重置密码失败")
			}
			return
		}
		middleware.SuccessWithMessage(c, "密码重置成功，请登录", nil)
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
		logger.Info("[API] /user/info called")
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
		default:
			logger.Error("claims 类型异常", zap.Any("claims_type", fmt.Sprintf("%T", claims)))
		}

		logger.Info("[API] /user/info 解析claims", zap.String("username", username))
		user, err := service.GetUserBaseInfo(db, username)
		if err != nil {
			logger.Error("[API] /user/info 获取用户信息失败", zap.Error(err))
			c.JSON(500, gin.H{"code": 500, "message": "获取用户信息失败", "data": nil})
			return
		}
		logger.Info("[API] /user/info 成功", zap.Any("user", user))
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
		c.JSON(200, resp)
	}
}
