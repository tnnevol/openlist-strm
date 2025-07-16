package controller

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/middleware"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"github.com/tnnevol/openlist-strm/backend-api/internal/util"
	"go.uber.org/zap"
)

type OpenListServiceReq struct {
	Name        string `json:"name" binding:"required"`
	Account     string `json:"account" binding:"required"`
	Token       string `json:"token" binding:"required"`
	ServiceUrl  string `json:"serviceUrl" binding:"required"`
	BackupUrl   string `json:"backupUrl"`
	Enabled     model.Enabled `json:"enabled"` // 支持字符串或数字
}

// convertToResponse 将 OpenListService 转换为 OpenListServiceResponse
func convertToResponse(service *model.OpenListService) *model.OpenListServiceResponse {
	return &model.OpenListServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Account:     service.Account,
		Token:       service.Token,
		ServiceUrl:  service.ServiceUrl,
		BackupUrl:   service.BackupUrl,
		Enabled:     model.Enabled(strconv.Itoa(util.Bool2Int(service.Enabled))),
		UpdatedAt:   service.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// convertToResponseList 将 OpenListService 列表转换为 OpenListServiceResponse 列表
func convertToResponseList(services []*model.OpenListService) []*model.OpenListServiceResponse {
	var responses []*model.OpenListServiceResponse
	for _, service := range services {
		responses = append(responses, convertToResponse(service))
	}
	return responses
}

// ListOpenListService godoc
// @Summary      分页查询OpenList服务
// @Description  查询当前用户所有OpenList服务，支持分页，返回list/total/page/pageSize
// @Tags         OpenListService
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        page query int false "页码(默认1)"
// @Param        pageSize query int false "每页条数(默认10)"
// @Success      200 {object} middleware.Response[model.PageResult[model.OpenListServiceResponse]]
// @Router       /openlist/list [get]
func ListOpenListService(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /openlist/service [GET] called - 请求入口")
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /openlist/service [GET] 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		if userID == 0 {
			logger.Error("[API] /openlist/service [GET] 用户信息无效")
			middleware.Unauthorized(c, "用户信息无效")
			return
		}
		services, err := service.GetOpenListServicesByUserID(db, userID)
		if err != nil {
			logger.Error("[API] /openlist/service [GET] 查询失败", zap.Error(err))
			middleware.InternalServerError(c, "查询失败")
			return
		}
		if services == nil {
			services = []*model.OpenListService{}
		}
		// 转换为响应格式
		responses := convertToResponseList(services)
		// 分页参数
		page, pageSize := util.GetPageParams(c)
		pagedPtr, total := util.Paginate(responses, page, pageSize)
		paged := make([]model.OpenListServiceResponse, len(pagedPtr))
		for i, v := range pagedPtr {
			if v != nil {
				paged[i] = *v
			}
		}
		logger.Info("[API] /openlist/service [GET] 查询成功", zap.Int("user_id", userID), zap.Int("count", len(services)), zap.Int("page", page), zap.Int("pageSize", pageSize))
		result := model.PageResult[model.OpenListServiceResponse]{
			List: paged,
			Total: total,
			Page: page,
			PageSize: pageSize,
		}
		middleware.Success(c, result)
	}
}

// CreateOpenListService godoc
// @Summary      新增OpenList服务
// @Description  新增OpenList服务，字段：name、account、token、serviceUrl、backupUrl、enabled
// @Tags         OpenListService
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        body body OpenListServiceReq true "服务信息"
// @Success      200 {object} middleware.Response[string]
// @Router       /openlist/add [post]
func CreateOpenListService(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /openlist/service [POST] called - 请求入口")
		var req OpenListServiceReq
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /openlist/service [POST] 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "参数错误")
			return
		}
		logger.Info("[API] /openlist/service [POST] 参数绑定成功", zap.Any("req", req))
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /openlist/service [POST] 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		if userID == 0 {
			logger.Error("[API] /openlist/service [POST] 用户信息无效")
			middleware.Unauthorized(c, "用户信息无效")
			return
		}
		enabledBool := util.ParseEnabled(req.Enabled)
		serviceObj := &model.OpenListService{
			Name: req.Name,
			Account: req.Account,
			Token: req.Token,
			ServiceUrl: req.ServiceUrl,
			BackupUrl: req.BackupUrl,
			Enabled: enabledBool,
			UserID: userID,
		}
		err := service.CreateOpenListService(db, serviceObj)
		if err != nil {
			logger.Error("[API] /openlist/service [POST] 创建失败", zap.Error(err))
			middleware.InternalServerError(c, "创建服务失败")
			return
		}
		logger.Info("[API] /openlist/service [POST] 创建成功", zap.Int("user_id", userID), zap.String("name", req.Name))
		middleware.SuccessWithMessage(c, "创建成功", nil)
	}
}

// GetOpenListService godoc
// @Summary      OpenList服务详情
// @Description  查询单个OpenList服务
// @Tags         OpenListService
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        id path int true "服务ID"
// @Success      200 {object} middleware.Response[model.OpenListServiceResponse]
// @Router       /openlist/detail/{id} [get]
func GetOpenListService(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /openlist/service/:id [GET] called - 请求入口")
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error("[API] /openlist/service/:id [GET] 参数错误", zap.String("id", idStr))
			middleware.ValidationError(c, "参数错误")
			return
		}
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /openlist/service/:id [GET] 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		if userID == 0 {
			logger.Error("[API] /openlist/service/:id [GET] 用户信息无效")
			middleware.Unauthorized(c, "用户信息无效")
			return
		}
		serviceObj, err := service.GetOpenListServiceByID(db, id)
		if err != nil || serviceObj == nil || serviceObj.UserID != userID {
			logger.Error("[API] /openlist/service/:id [GET] 服务不存在", zap.Int("id", id))
			middleware.NotFound(c, "服务不存在")
			return
		}
		logger.Info("[API] /openlist/service/:id [GET] 查询成功", zap.Int("id", id), zap.Int("user_id", userID))
		middleware.Success(c, convertToResponse(serviceObj))
	}
}

// UpdateOpenListService godoc
// @Summary      编辑OpenList服务
// @Description  编辑指定ID的OpenList服务
// @Tags         OpenListService
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        id path int true "服务ID"
// @Param        body body OpenListServiceReq true "服务信息"
// @Success      200 {object} middleware.Response[string]
// @Router       /openlist/update/{id} [put]
func UpdateOpenListService(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /openlist/service/:id [PUT] called - 请求入口")
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error("[API] /openlist/service/:id [PUT] 参数错误", zap.String("id", idStr))
			middleware.ValidationError(c, "参数错误")
			return
		}
		var req OpenListServiceReq
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("[API] /openlist/service/:id [PUT] 参数绑定失败", zap.Error(err))
			middleware.ValidationError(c, "参数错误")
			return
		}
		logger.Info("[API] /openlist/service/:id [PUT] 参数绑定成功", zap.Any("req", req))
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /openlist/service/:id [PUT] 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		if userID == 0 {
			logger.Error("[API] /openlist/service/:id [PUT] 用户信息无效")
			middleware.Unauthorized(c, "用户信息无效")
			return
		}
		serviceObj, err := service.GetOpenListServiceByID(db, id)
		if err != nil || serviceObj == nil || serviceObj.UserID != userID {
			logger.Error("[API] /openlist/service/:id [PUT] 服务不存在", zap.Int("id", id))
			middleware.NotFound(c, "服务不存在")
			return
		}
		enabledBool := util.ParseEnabled(req.Enabled)
		serviceObj.Enabled = enabledBool
		err = service.UpdateOpenListService(db, serviceObj)
		if err != nil {
			logger.Error("[API] /openlist/service/:id [PUT] 更新失败", zap.Error(err))
			middleware.InternalServerError(c, "更新失败")
			return
		}
		logger.Info("[API] /openlist/service/:id [PUT] 更新成功", zap.Int("id", id), zap.Int("user_id", userID))
		middleware.SuccessWithMessage(c, "更新成功", nil)
	}
}

// DeleteOpenListService godoc
// @Summary      删除OpenList服务
// @Description  删除指定ID的OpenList服务
// @Tags         OpenListService
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {accessToken}"
// @Param        id path int true "服务ID"
// @Success      200 {object} middleware.Response[string]
// @Router       /openlist/delete/{id} [delete]
func DeleteOpenListService(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("[API] /openlist/service/:id [DELETE] called - 请求入口")
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error("[API] /openlist/service/:id [DELETE] 参数错误", zap.String("id", idStr))
			middleware.ValidationError(c, "参数错误")
			return
		}
		claims, ok := c.Get("claims")
		if !ok {
			logger.Error("[API] /openlist/service/:id [DELETE] 未获取到claims")
			middleware.Unauthorized(c, "未登录或token缺失")
			return
		}
		userID := util.ExtractUserIDFromClaims(claims)
		if userID == 0 {
			logger.Error("[API] /openlist/service/:id [DELETE] 用户信息无效")
			middleware.Unauthorized(c, "用户信息无效")
			return
		}
		serviceObj, err := service.GetOpenListServiceByID(db, id)
		if err != nil || serviceObj == nil || serviceObj.UserID != userID {
			logger.Error("[API] /openlist/service/:id [DELETE] 服务不存在", zap.Int("id", id))
			middleware.NotFound(c, "服务不存在")
			return
		}
		err = service.DeleteOpenListService(db, id)
		if err != nil {
			logger.Error("[API] /openlist/service/:id [DELETE] 删除失败", zap.Error(err))
			middleware.InternalServerError(c, "删除失败")
			return
		}
		logger.Info("[API] /openlist/service/:id [DELETE] 删除成功", zap.Int("id", id), zap.Int("user_id", userID))
		middleware.SuccessWithMessage(c, "删除成功", nil)
	}
}

// RegisterOpenListServiceRoutes 统一注册/openlist/service相关接口
func RegisterOpenListServiceRoutes(rg *gin.RouterGroup, db *sql.DB) {
	rg.GET("/list", ListOpenListService(db))
	rg.POST("/add", CreateOpenListService(db))
	rg.GET("/detail/:id", GetOpenListService(db))
	rg.PUT("/update/:id", UpdateOpenListService(db))
	rg.DELETE("/delete/:id", DeleteOpenListService(db))
} 
