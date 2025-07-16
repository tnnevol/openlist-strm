package app

import (
	"github.com/gin-gonic/gin"
	_ "github.com/tnnevol/openlist-strm/backend-api/docs"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/controller"
	"github.com/tnnevol/openlist-strm/backend-api/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	// 全局注册AuthMiddleware
	r.Use(middleware.AuthMiddleware(db))

	userGroup := r.Group("/user")
	controller.RegisterUserRoutes(userGroup, db)

	openlistGroup := r.Group("/openlist")
	controller.RegisterOpenListServiceRoutes(openlistGroup, db)

	strmGroup := r.Group("/strm")
	controller.RegisterStrmConfigRoutes(strmGroup, db)

	dictGroup := r.Group("/dict")
	controller.RegisterDictRoutes(dictGroup, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
} 
