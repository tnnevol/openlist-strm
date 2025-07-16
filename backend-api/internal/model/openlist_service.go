package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OpenListService struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"type:varchar(128);index"`
	Account    string    `json:"account" gorm:"type:varchar(128)"`
	Token      string    `json:"token" gorm:"type:varchar(255)"`
	ServiceUrl string    `json:"serviceUrl" gorm:"type:varchar(255)"`
	BackupUrl  string    `json:"backupUrl" gorm:"type:varchar(255)"`
	Enabled    bool      `json:"enabled"`
	UserID     int       `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// CreateOpenListService 创建OpenList服务
func CreateOpenListService(db *gorm.DB, service *OpenListService) error {
	logger.Info("[DB] CreateOpenListService", zap.String("name", service.Name), zap.Int("user_id", service.UserID))
	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now
	if err := db.Create(service).Error; err != nil {
		logger.Error("[DB] CreateOpenListService error", zap.Error(err))
		return err
	}
	return nil
}

// GetOpenListServiceByID 根据ID获取OpenList服务
func GetOpenListServiceByID(db *gorm.DB, id int) (*OpenListService, error) {
	var service OpenListService
	if err := db.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

// GetOpenListServicesByUserID 根据用户ID分页获取OpenList服务，返回数据和总数
func GetOpenListServicesByUserID(db *gorm.DB, userID, page, pageSize int) ([]*OpenListService, int64, error) {
	var services []*OpenListService
	var total int64
	db = db.Model(&OpenListService{}).Where("user_id = ?", userID)
	db.Count(&total)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&services).Error; err != nil {
		return nil, 0, err
	}
	return services, total, nil
}

// UpdateOpenListService 更新OpenList服务
func UpdateOpenListService(db *gorm.DB, service *OpenListService) error {
	logger.Info("[DB] UpdateOpenListService", zap.Int("id", service.ID))
	service.UpdatedAt = time.Now()
	if err := db.Model(&OpenListService{}).Where("id = ?", service.ID).Updates(service).Error; err != nil {
		logger.Error("[DB] UpdateOpenListService error", zap.Error(err))
		return err
	}
	return nil
}

// DeleteOpenListService 删除OpenList服务
func DeleteOpenListService(db *gorm.DB, id int) error {
	logger.Info("[DB] DeleteOpenListService", zap.Int("id", id))
	if err := db.Delete(&OpenListService{}, id).Error; err != nil {
		logger.Error("[DB] DeleteOpenListService error", zap.Error(err))
		return err
	}
	return nil
}

// ToggleOpenListServiceEnabled 切换服务启用状态
func ToggleOpenListServiceEnabled(db *gorm.DB, id int, enabled bool) error {
	logger.Info("[DB] ToggleOpenListServiceEnabled", zap.Int("id", id), zap.Bool("enabled", enabled))
	if err := db.Model(&OpenListService{}).Where("id = ?", id).Updates(map[string]interface{}{"enabled": enabled, "updated_at": time.Now()}).Error; err != nil {
		logger.Error("[DB] ToggleOpenListServiceEnabled error", zap.Error(err))
		return err
	}
	return nil
} 
