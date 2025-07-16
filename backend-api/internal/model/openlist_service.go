package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

type OpenListService struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Account   string    `json:"account"`
	Token     string    `json:"token"`
	ServiceUrl string   `json:"serviceUrl"`
	BackupUrl  string   `json:"backupUrl"`
	Enabled    bool     `json:"enabled"`
	UserID     int      `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// CreateOpenListService 创建OpenList服务
func CreateOpenListService(db *sql.DB, service *OpenListService) error {
	logger.Info("[DB] CreateOpenListService", zap.String("name", service.Name), zap.Int("user_id", service.UserID))
	
	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now
	
	_, err := db.Exec(`
		INSERT INTO openlist_service (name, account, token, service_url, backup_url, enabled, user_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		service.Name, service.Account, service.Token, service.ServiceUrl, service.BackupUrl, service.Enabled, service.UserID, service.CreatedAt, service.UpdatedAt)
	
	if err != nil {
		logger.Error("[DB] CreateOpenListService error", zap.Error(err))
		return err
	}
	
	return nil
}

// GetOpenListServiceByID 根据ID获取OpenList服务
func GetOpenListServiceByID(db *sql.DB, id int) (*OpenListService, error) {
	var service OpenListService
	err := db.QueryRow(`
		SELECT id, name, account, token, service_url, backup_url, enabled, user_id, created_at, updated_at 
		FROM openlist_service WHERE id = ?`, id).
		Scan(&service.ID, &service.Name, &service.Account, &service.Token, &service.ServiceUrl, &service.BackupUrl, &service.Enabled, &service.UserID, &service.CreatedAt, &service.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetOpenListServicesByUserID 根据用户ID获取所有OpenList服务
func GetOpenListServicesByUserID(db *sql.DB, userID int) ([]*OpenListService, error) {
	rows, err := db.Query(`
		SELECT id, name, account, token, service_url, backup_url, enabled, user_id, created_at, updated_at 
		FROM openlist_service WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var services []*OpenListService
	for rows.Next() {
		var service OpenListService
		err := rows.Scan(&service.ID, &service.Name, &service.Account, &service.Token, &service.ServiceUrl, &service.BackupUrl, &service.Enabled, &service.UserID, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	
	return services, nil
}

// UpdateOpenListService 更新OpenList服务
func UpdateOpenListService(db *sql.DB, service *OpenListService) error {
	logger.Info("[DB] UpdateOpenListService", zap.Int("id", service.ID))
	
	service.UpdatedAt = time.Now()
	
	_, err := db.Exec(`
		UPDATE openlist_service 
		SET name = ?, account = ?, token = ?, service_url = ?, backup_url = ?, enabled = ?, updated_at = ? 
		WHERE id = ?`,
		service.Name, service.Account, service.Token, service.ServiceUrl, service.BackupUrl, service.Enabled, service.UpdatedAt, service.ID)
	
	if err != nil {
		logger.Error("[DB] UpdateOpenListService error", zap.Error(err))
		return err
	}
	
	return nil
}

// DeleteOpenListService 删除OpenList服务
func DeleteOpenListService(db *sql.DB, id int) error {
	logger.Info("[DB] DeleteOpenListService", zap.Int("id", id))
	
	_, err := db.Exec("DELETE FROM openlist_service WHERE id = ?", id)
	if err != nil {
		logger.Error("[DB] DeleteOpenListService error", zap.Error(err))
		return err
	}
	
	return nil
}

// ToggleOpenListServiceEnabled 切换服务启用状态
func ToggleOpenListServiceEnabled(db *sql.DB, id int, enabled bool) error {
	logger.Info("[DB] ToggleOpenListServiceEnabled", zap.Int("id", id), zap.Bool("enabled", enabled))
	
	_, err := db.Exec("UPDATE openlist_service SET enabled = ?, updated_at = ? WHERE id = ?", enabled, time.Now(), id)
	if err != nil {
		logger.Error("[DB] ToggleOpenListServiceEnabled error", zap.Error(err))
		return err
	}
	
	return nil
} 
