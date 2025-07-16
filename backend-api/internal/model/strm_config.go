package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

type UpdateMode string

const (
	UpdateModeIncremental UpdateMode = "incremental"
	UpdateModeFull        UpdateMode = "full"
)

type StrmConfig struct {
	ID               int
	UserID           int
	Name             string
	AlistBasePath    string
	StrmOutputPath   string
	DownloadEnabled  bool
	DownloadInterval int
	UpdateMode       UpdateMode
	ServiceID        int
	IsUseBackupUrl   bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateStrmConfig 创建Strm配置
func CreateStrmConfig(db *sql.DB, config *StrmConfig) error {
	logger.Info("[DB] CreateStrmConfig", zap.String("name", config.Name), zap.Int("service_id", config.ServiceID), zap.Int("user_id", config.UserID))

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	_, err := db.Exec(`
		INSERT INTO strm_config (user_id, name, alist_base_path, strm_output_path, download_enabled, download_interval, update_mode, service_id, is_use_backup_url, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		config.UserID, config.Name, config.AlistBasePath, config.StrmOutputPath, config.DownloadEnabled, config.DownloadInterval, config.UpdateMode, config.ServiceID, config.IsUseBackupUrl, config.CreatedAt, config.UpdatedAt)

	if err != nil {
		logger.Error("[DB] CreateStrmConfig error", zap.Error(err))
		return err
	}

	return nil
}

// GetStrmConfigByID 根据ID获取Strm配置
func GetStrmConfigByID(db *sql.DB, id int) (*StrmConfig, error) {
	var config StrmConfig
	err := db.QueryRow(`
		SELECT id, user_id, name, alist_base_path, strm_output_path, download_enabled, download_interval, update_mode, service_id, is_use_backup_url, created_at, updated_at 
		FROM strm_config WHERE id = ?`, id).
		Scan(&config.ID, &config.UserID, &config.Name, &config.AlistBasePath, &config.StrmOutputPath, &config.DownloadEnabled, &config.DownloadInterval, &config.UpdateMode, &config.ServiceID, &config.IsUseBackupUrl, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetStrmConfigsByServiceID 根据服务ID获取所有Strm配置
func GetStrmConfigsByServiceID(db *sql.DB, serviceID int) ([]*StrmConfig, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, alist_base_path, strm_output_path, download_enabled, download_interval, update_mode, service_id, is_use_backup_url, created_at, updated_at 
		FROM strm_config WHERE service_id = ? ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*StrmConfig
	for rows.Next() {
		var config StrmConfig
		err := rows.Scan(&config.ID, &config.UserID, &config.Name, &config.AlistBasePath, &config.StrmOutputPath, &config.DownloadEnabled, &config.DownloadInterval, &config.UpdateMode, &config.ServiceID, &config.IsUseBackupUrl, &config.CreatedAt, &config.UpdatedAt)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// GetAllStrmConfigs 获取所有Strm配置
func GetAllStrmConfigs(db *sql.DB) ([]*StrmConfig, error) {
	rows, err := db.Query(`SELECT id, user_id, name, alist_base_path, strm_output_path, download_enabled, download_interval, update_mode, service_id, is_use_backup_url, created_at, updated_at FROM strm_config ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var configs []*StrmConfig
	for rows.Next() {
		var config StrmConfig
		err := rows.Scan(&config.ID, &config.UserID, &config.Name, &config.AlistBasePath, &config.StrmOutputPath, &config.DownloadEnabled, &config.DownloadInterval, &config.UpdateMode, &config.ServiceID, &config.IsUseBackupUrl, &config.CreatedAt, &config.UpdatedAt)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

// UpdateStrmConfig 更新Strm配置
func UpdateStrmConfig(db *sql.DB, config *StrmConfig) error {
	logger.Info("[DB] UpdateStrmConfig", zap.Int("id", config.ID))

	config.UpdatedAt = time.Now()

	_, err := db.Exec(`
		UPDATE strm_config 
		SET user_id = ?, name = ?, alist_base_path = ?, strm_output_path = ?, download_enabled = ?, download_interval = ?, update_mode = ?, is_use_backup_url = ?, updated_at = ? 
		WHERE id = ?`,
		config.UserID, config.Name, config.AlistBasePath, config.StrmOutputPath, config.DownloadEnabled, config.DownloadInterval, config.UpdateMode, config.IsUseBackupUrl, config.UpdatedAt, config.ID)

	if err != nil {
		logger.Error("[DB] UpdateStrmConfig error", zap.Error(err))
		return err
	}

	return nil
}

// DeleteStrmConfig 删除Strm配置
func DeleteStrmConfig(db *sql.DB, id int) error {
	logger.Info("[DB] DeleteStrmConfig", zap.Int("id", id))
	
	_, err := db.Exec("DELETE FROM strm_config WHERE id = ?", id)
	if err != nil {
		logger.Error("[DB] DeleteStrmConfig error", zap.Error(err))
		return err
	}
	
	return nil
}

// ToggleStrmConfigDownloadEnabled 切换下载启用状态
func ToggleStrmConfigDownloadEnabled(db *sql.DB, id int, enabled bool) error {
	logger.Info("[DB] ToggleStrmConfigDownloadEnabled", zap.Int("id", id), zap.Bool("enabled", enabled))
	
	_, err := db.Exec("UPDATE strm_config SET download_enabled = ?, updated_at = ? WHERE id = ?", enabled, time.Now(), id)
	if err != nil {
		logger.Error("[DB] ToggleStrmConfigDownloadEnabled error", zap.Error(err))
		return err
	}
	
	return nil
} 

// 新增用于接口返回的小驼峰结构体
// swagger:model
// 用于 StrmConfigPageResult 的 list 元素
//
type StrmConfigResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	AlistBasePath    string    `json:"alistBasePath"`
	StrmOutputPath   string    `json:"strmOutputPath"`
	DownloadEnabled  bool      `json:"downloadEnabled"`
	DownloadInterval int       `json:"downloadInterval"`
	UpdateMode       string    `json:"updateMode"`
	ServiceID        int       `json:"serviceId"`
	IsUseBackupUrl   bool      `json:"isUseBackupUrl"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
} 
