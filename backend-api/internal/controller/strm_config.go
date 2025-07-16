package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/middleware"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"github.com/tnnevol/openlist-strm/backend-api/internal/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 删除本地 extractUserIDFromClaims 实现，全部调用 util.ExtractUserIDFromClaims



type StrmConfigReq struct {
	Name             string      `json:"name" binding:"required"`
	AlistBasePath    string      `json:"alistBasePath" binding:"required"`
	StrmOutputPath   string      `json:"strmOutputPath" binding:"required"`
	DownloadEnabled  model.DownloadEnabled `json:"downloadEnabled" binding:"required"`
	DownloadInterval int         `json:"downloadInterval" binding:"required"`
	UpdateMode       model.UpdateMode  `json:"updateMode" binding:"required"`
	ServiceID        int         `json:"serviceId" binding:"required"`
	IsUseBackupUrl   model.IsUseBackupUrl `json:"isUseBackupUrl" binding:"required"`
}

type StrmConfigCopyReq struct {
	IDs []int `json:"ids" binding:"required"`
}

// ListStrmConfig godoc
// @Summary      分页查询Strm配置
// @Description  分页查询Strm配置，支持serviceId筛选，返回list/total/page/pageSize
// @Tags         StrmConfig
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        serviceId query int false "服务ID"
// @Param        page query int false "页码(默认1)"
// @Param        pageSize query int false "每页条数(默认10)"
// @Success      200 {object} middleware.Response[model.PageResult[model.StrmConfigResponse]]
// @Router       /strm/config/list [get]
func ListStrmConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		logger.Info("[DEBUG] /strm/config claims内容", zap.Any("claims", claims), zap.Int("userID", userID))
		serviceID, _ := strconv.Atoi(c.Query("serviceId"))
		page, pageSize := util.GetPageParams(c)
		configs, total, err := service.ListStrmConfigs(db, userID, serviceID, page, pageSize)
		if err != nil {
			middleware.InternalServerError(c, "查询失败")
			return
		}
		list := make([]model.StrmConfigResponse, len(configs))
		for i, v := range configs {
			if v != nil {
				list[i] = model.StrmConfigResponse{
					ID: v.ID,
					Name: v.Name,
					AlistBasePath: v.AlistBasePath,
					StrmOutputPath: v.StrmOutputPath,
					DownloadEnabled: model.DownloadEnabled(strconv.Itoa(util.Bool2Int(v.DownloadEnabled))),
					DownloadInterval: v.DownloadInterval,
					UpdateMode: string(v.UpdateMode),
					ServiceID: v.ServiceID,
					IsUseBackupUrl: model.IsUseBackupUrl(strconv.Itoa(util.Bool2Int(v.IsUseBackupUrl))),
					CreatedAt: v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
				}
			}
		}
		result := model.PageResult[model.StrmConfigResponse]{
			List: list,
			Total: total,
			Page: page,
			PageSize: pageSize,
		}
		middleware.Success(c, result)
	}
}

// CopyStrmConfig godoc
// @Summary 批量复制Strm配置
// @Description 通过ids批量复制Strm配置，configName后追加“-复制”
// @Tags StrmConfig
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {{accessToken}}"
// @Param body body StrmConfigCopyReq true "要复制的id列表"
// @Success 200 {object} middleware.Response[string] // data为操作结果描述
// @Router /strm/config/copy [post]
func CopyStrmConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StrmConfigCopyReq
		if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
			middleware.ValidationError(c, "参数错误")
			return
		}
		err := service.CopyStrmConfigs(db, req.IDs)
		if err != nil {
			middleware.InternalServerError(c, "复制失败")
			return
		}
		middleware.SuccessWithMessage(c, "复制成功", nil)
	}
}

// CreateStrmConfig godoc
// @Summary      新增Strm配置
// @Description  新增Strm配置
// @Tags         StrmConfig
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        body body StrmConfigReq true "配置内容"
// @Success      200 {object} middleware.Response[string]
// @Router       /strm/config/add [post]
func CreateStrmConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StrmConfigReq
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /strm/config 新增 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "参数错误")
			return
		}
		logger.Info("[API] /strm/config 新增 参数绑定成功", zap.Any("req", req))
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /strm/config 新增 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		logger.Info("[API] /strm/config 新增 claims和userID", zap.Any("claims", claims), zap.Int("userID", userID))
		downloadEnabledBool := util.ParseEnabled(req.DownloadEnabled)
		isUseBackupUrlBool := util.ParseEnabled(req.IsUseBackupUrl)
		cfg := &model.StrmConfig{
			UserID: userID,
			Name: req.Name,
			AlistBasePath: req.AlistBasePath,
			StrmOutputPath: req.StrmOutputPath,
			DownloadEnabled: downloadEnabledBool,
			DownloadInterval: req.DownloadInterval,
			UpdateMode: model.UpdateMode(req.UpdateMode),
			ServiceID: req.ServiceID,
			IsUseBackupUrl: isUseBackupUrlBool,
		}
		err := service.CreateStrmConfig(db, cfg)
		if err != nil {
			logger.Error("[API] /strm/config 新增 service.CreateStrmConfig失败", zap.Error(err))
			middleware.InternalServerError(c, "新增失败")
			return
		}
		logger.Info("[API] /strm/config 新增成功", zap.Any("cfg", cfg))
		middleware.SuccessWithMessage(c, "新增成功", nil)
	}
}

// UpdateStrmConfig godoc
// @Summary      编辑Strm配置
// @Description  编辑指定ID的Strm配置
// @Tags         StrmConfig
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        id path int true "配置ID"
// @Param        body body StrmConfigReq true "配置内容"
// @Success      200 {object} middleware.Response[string]
// @Router       /strm/config/update/{id} [put]
func UpdateStrmConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req StrmConfigReq
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ValidationError(c, "参数错误")
			return
		}
		claims, ok := c.Get("claims")
		if !ok {
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		cfg, err := service.GetStrmConfigByID(db, id)
		if err != nil || cfg == nil || cfg.UserID != userID {
			middleware.NotFound(c, "配置不存在")
			return
		}
		downloadEnabledBool := util.ParseEnabled(req.DownloadEnabled)
		isUseBackupUrlBool := util.ParseEnabled(req.IsUseBackupUrl)
		cfg.Name = req.Name
		cfg.AlistBasePath = req.AlistBasePath
		cfg.StrmOutputPath = req.StrmOutputPath
		cfg.DownloadEnabled = downloadEnabledBool
		cfg.DownloadInterval = req.DownloadInterval
		cfg.UpdateMode = model.UpdateMode(req.UpdateMode)
		cfg.ServiceID = req.ServiceID
		cfg.IsUseBackupUrl = isUseBackupUrlBool
		err = service.UpdateStrmConfig(db, cfg)
		if err != nil {
			middleware.InternalServerError(c, "编辑失败")
			return
		}
		middleware.SuccessWithMessage(c, "编辑成功", nil)
	}
}

// DeleteStrmConfig godoc
// @Summary      删除Strm配置
// @Description  删除指定ID的Strm配置
// @Tags         StrmConfig
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        id path int true "配置ID"
// @Success      200 {object} middleware.Response[string]
// @Router       /strm/config/delete/{id} [delete]
func DeleteStrmConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		claims, ok := c.Get("claims")
		if !ok {
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		cfg, err := service.GetStrmConfigByID(db, id)
		if err != nil || cfg == nil || cfg.UserID != userID {
			middleware.NotFound(c, "配置不存在")
			return
		}
		err = service.DeleteStrmConfig(db, id)
		if err != nil {
			middleware.InternalServerError(c, "删除失败")
			return
		}
		middleware.SuccessWithMessage(c, "删除成功", nil)
	}
}

// GetStrmConfig godoc
// @Summary      Strm配置详情
// @Description  获取指定ID的Strm配置详情
// @Tags         StrmConfig
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        id path int true "配置ID"
// @Success      200 {object} middleware.Response[model.StrmConfigResponse]
// @Router       /strm/config/detail/{id} [get]
func GetStrmConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		claims, ok := c.Get("claims")
		if !ok {
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		cfg, err := service.GetStrmConfigByID(db, id)
		if err != nil || cfg == nil || cfg.UserID != userID {
			middleware.NotFound(c, "配置不存在")
			return
		}
		resp := model.StrmConfigResponse{
			ID: cfg.ID,
			Name: cfg.Name,
			AlistBasePath: cfg.AlistBasePath,
			StrmOutputPath: cfg.StrmOutputPath,
			DownloadEnabled: model.DownloadEnabled(strconv.Itoa(util.Bool2Int(cfg.DownloadEnabled))),
			DownloadInterval: cfg.DownloadInterval,
			UpdateMode: string(cfg.UpdateMode),
			ServiceID: cfg.ServiceID,
			IsUseBackupUrl: model.IsUseBackupUrl(strconv.Itoa(util.Bool2Int(cfg.IsUseBackupUrl))),
			CreatedAt: cfg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: cfg.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		middleware.Success(c, resp)
	}
}

// RegisterStrmConfigRoutes 统一注册/strm/config相关接口
func RegisterStrmConfigRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/config/list", ListStrmConfig(db))
	rg.POST("/config/add", CreateStrmConfig(db))
	rg.GET("/config/detail/:id", GetStrmConfig(db))
	rg.PUT("/config/update/:id", UpdateStrmConfig(db))
	rg.DELETE("/config/delete/:id", DeleteStrmConfig(db))
	rg.POST("/config/copy", CopyStrmConfig(db))
}
