package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateMode string

const (
	UpdateModeIncremental UpdateMode = "incremental"
	UpdateModeFull        UpdateMode = "full"
)

type StrmConfig struct {
	ID               int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int        `json:"userId"`
	Name             string     `json:"name" gorm:"type:varchar(128);index"`
	AlistBasePath    string     `json:"alistBasePath" gorm:"type:varchar(255)"`
	StrmOutputPath   string     `json:"strmOutputPath" gorm:"type:varchar(255)"`
	DownloadEnabled  bool       `json:"downloadEnabled"`
	DownloadInterval int        `json:"downloadInterval"`
	UpdateMode       UpdateMode `json:"updateMode" gorm:"type:varchar(32)"`
	ServiceID        int        `json:"serviceId"`
	IsUseBackupUrl   bool       `json:"isUseBackupUrl"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

func CreateStrmConfig(db *gorm.DB, config *StrmConfig) error {
	logger.Info("[DB] CreateStrmConfig", zap.String("name", config.Name), zap.Int("service_id", config.ServiceID), zap.Int("user_id", config.UserID))
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now
	if err := db.Create(config).Error; err != nil {
		logger.Error("[DB] CreateStrmConfig error", zap.Error(err))
		return err
	}
	return nil
}

func GetStrmConfigByID(db *gorm.DB, id int) (*StrmConfig, error) {
	var config StrmConfig
	if err := db.First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// 分页获取指定 serviceID 的 StrmConfig，返回数据和总数
func GetStrmConfigsByServiceID(db *gorm.DB, serviceID, userID, page, pageSize int) ([]*StrmConfig, int64, error) {
	var configs []*StrmConfig
	var total int64
	db = db.Model(&StrmConfig{}).Where("service_id = ? AND user_id = ?", serviceID, userID)
	db.Count(&total)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}

// 分页获取所有 StrmConfig，返回数据和总数
func GetAllStrmConfigs(db *gorm.DB, userID, page, pageSize int) ([]*StrmConfig, int64, error) {
	var configs []*StrmConfig
	var total int64
	db = db.Model(&StrmConfig{}).Where("user_id = ?", userID)
	db.Count(&total)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}

func UpdateStrmConfig(db *gorm.DB, config *StrmConfig) error {
	logger.Info("[DB] UpdateStrmConfig", zap.Int("id", config.ID))
	config.UpdatedAt = time.Now()
	if err := db.Model(&StrmConfig{}).Where("id = ?", config.ID).Updates(config).Error; err != nil {
		logger.Error("[DB] UpdateStrmConfig error", zap.Error(err))
		return err
	}
	return nil
}

func DeleteStrmConfig(db *gorm.DB, id int) error {
	logger.Info("[DB] DeleteStrmConfig", zap.Int("id", id))
	if err := db.Delete(&StrmConfig{}, id).Error; err != nil {
		logger.Error("[DB] DeleteStrmConfig error", zap.Error(err))
		return err
	}
	return nil
}

func ToggleStrmConfigDownloadEnabled(db *gorm.DB, id int, enabled bool) error {
	logger.Info("[DB] ToggleStrmConfigDownloadEnabled", zap.Int("id", id), zap.Bool("enabled", enabled))
	if err := db.Model(&StrmConfig{}).Where("id = ?", id).Updates(map[string]interface{}{"download_enabled": enabled, "updated_at": time.Now()}).Error; err != nil {
		logger.Error("[DB] ToggleStrmConfigDownloadEnabled error", zap.Error(err))
		return err
	}
	return nil
} 
