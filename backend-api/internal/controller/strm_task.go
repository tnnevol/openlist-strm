package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterStrmTaskRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/list", ListStrmTasks(db))
	rg.POST("/add", CreateStrmTask(db))
	rg.PUT("/update/:id", UpdateStrmTask(db))
	rg.DELETE("/delete/:id", DeleteStrmTask(db))
	rg.GET("/detail/:id", GetStrmTask(db))
	rg.POST("/copy", CopyStrmTask(db))
}

func ListStrmTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}

func CreateStrmTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}

func UpdateStrmTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}

func DeleteStrmTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}

func GetStrmTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}

func CopyStrmTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}

func ExecuteStrmTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现
	}
}
