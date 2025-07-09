package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/tnnevol/openlist-strm/backend-api/docs"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/controller"
	"github.com/tnnevol/openlist-strm/backend-api/internal/middleware"
)

func RegisterRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	userGroup := r.Group("/user")
	userGroup.POST("/send-code", controller.SendCode(db))
	userGroup.POST("/register", controller.Register(db))
	userGroup.POST("/login", controller.Login(db))
	userGroup.POST("/forgot-password/send-code", controller.ForgotPasswordSendCode(db))
	userGroup.POST("/forgot-password/reset", controller.ForgotPasswordReset(db))
	userGroup.GET("/info", controller.UserInfo(db))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
} 
