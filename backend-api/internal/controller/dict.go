package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/middleware"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"go.uber.org/zap"
)

type DictReq struct {
	Type        string `json:"type" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
	ParentID    int    `json:"parentId"`
}

type DictPageResult struct {
	List     []*model.Dict `json:"list"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

func RegisterDictRoutes(rg *gin.RouterGroup, db *sql.DB) {
	rg.GET("/list", ListDicts(db))
	rg.POST("/add", CreateDict(db))
	rg.GET("/detail/:id", GetDict(db))
	rg.PUT("/update/:id", UpdateDict(db))
	rg.DELETE("/delete/:id", DeleteDict(db))
}

// CreateDict godoc
// @Summary      新增字典项
// @Description  新增一条字典数据
// @Tags         字典
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {{accessToken}}"
// @Param        dict  body  DictReq  true  "字典信息"
// @Success      200   {object} model.Response
// @Router       /dict/add [post]
func CreateDict(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DictReq
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /dict 新增 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "参数错误")
			return
		}
		logger.Info("[API] /dict 新增 参数绑定成功", zap.Any("req", req))
		d := &model.Dict{
			Type: req.Type,
			Key: req.Key,
			Value: req.Value,
			Description: req.Description,
			ParentID: req.ParentID,
		}
		err := service.CreateDict(db, d)
		if err != nil {
			logger.Error("[API] /dict 新增 service.CreateDict失败", zap.Error(err))
			middleware.InternalServerError(c, "新增失败")
			return
		}
		logger.Info("[API] /dict 新增成功", zap.Any("dict", d))
		middleware.SuccessWithMessage(c, "新增成功", nil)
	}
}

// UpdateDict godoc
// @Summary      编辑字典项
// @Description  编辑指定ID的字典项
// @Tags         字典
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {{accessToken}}"
// @Param        id    path  int     true  "字典ID"
// @Param        dict  body  DictReq true  "字典信息"
// @Success      200   {object} model.Response
// @Router       /dict/update/{id} [put]
func UpdateDict(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req DictReq
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /dict 编辑 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "参数错误")
			return
		}
		d, err := service.GetDictByID(db, id)
		if err != nil || d == nil {
			middleware.NotFound(c, "字典项不存在")
			return
		}
		d.Type = req.Type
		d.Key = req.Key
		d.Value = req.Value
		d.Description = req.Description
		d.ParentID = req.ParentID
		d.UpdatedAt = time.Now()
		err = service.UpdateDict(db, d)
		if err != nil {
			logger.Error("[API] /dict 编辑 service.UpdateDict失败", zap.Error(err))
			middleware.InternalServerError(c, "编辑失败")
			return
		}
		logger.Info("[API] /dict 编辑成功", zap.Any("dict", d))
		middleware.SuccessWithMessage(c, "编辑成功", nil)
	}
}

// DeleteDict godoc
// @Summary      删除字典项
// @Description  删除指定ID的字典项
// @Tags         字典
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {{accessToken}}"
// @Param        id   path  int  true  "字典ID"
// @Success      200  {object} model.Response
// @Router       /dict/delete/{id} [delete]
func DeleteDict(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		d, err := service.GetDictByID(db, id)
		if err != nil || d == nil {
			middleware.NotFound(c, "字典项不存在")
			return
		}
		err = service.DeleteDict(db, id)
		if err != nil {
			logger.Error("[API] /dict 删除 service.DeleteDict失败", zap.Error(err))
			middleware.InternalServerError(c, "删除失败")
			return
		}
		logger.Info("[API] /dict 删除成功", zap.Int("id", id))
		middleware.SuccessWithMessage(c, "删除成功", nil)
	}
}

// GetDict godoc
// @Summary      字典详情
// @Description  获取指定ID的字典项详情
// @Tags         字典
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {{accessToken}}"
// @Param        id   path  int  true  "字典ID"
// @Success      200  {object} model.Dict
// @Router       /dict/detail/{id} [get]
func GetDict(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		d, err := service.GetDictByID(db, id)
		if err != nil || d == nil {
			middleware.NotFound(c, "字典项不存在")
			return
		}
		middleware.Success(c, d)
	}
}

// ListDicts godoc
// @Summary      字典列表
// @Description  查询全部字典项，可按类型和父级筛选，无分页
// @Tags         字典
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {{accessToken}}"
// @Param        type     query    string  false  "字典类型"
// @Param        parentId query    int     false  "父级ID"
// @Success      200      {object} []model.Dict
// @Router       /dict/list [get]
func ListDicts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dictType := c.Query("type")
		parentIdStr := c.Query("parentId")
		var parentId int
		if parentIdStr != "" {
			parentId, _ = strconv.Atoi(parentIdStr)
		}
		all, err := service.ListDicts(db, dictType)
		if err != nil {
			logger.Error("[API] /dict 查询 service.ListDicts失败", zap.Error(err))
			middleware.InternalServerError(c, "查询失败")
			return
		}
		var filtered []*model.Dict
		if parentIdStr != "" {
			for _, d := range all {
				if d.ParentID == parentId {
					filtered = append(filtered, d)
				}
			}
		} else {
			filtered = all
		}
		logger.Info("[API] /dict 查询成功", zap.Int("count", len(filtered)))
		middleware.Success(c, filtered)
	}
} 
