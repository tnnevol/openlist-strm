package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LogName string
type TaskStatus string

const (
	LogNameCreate LogName = "create"
	LogNameCheck  LogName = "check"
)

const (
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusError     TaskStatus = "error"
	TaskStatusCompleted TaskStatus = "completed"
)

type LogRecord struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int        `json:"userId"`
	Name       LogName    `json:"name" gorm:"type:varchar(32)"`
	LogPath    string     `json:"logPath" gorm:"type:varchar(255)"`
	CreatedAt  time.Time  `json:"createdAt"`
	TaskStatus TaskStatus `json:"taskStatus" gorm:"type:varchar(32)"`
	TaskID     int        `json:"taskId"`
}

func CreateLogRecord(db *gorm.DB, record *LogRecord) error {
	logger.Info("[DB] CreateLogRecord", zap.String("name", string(record.Name)), zap.Int("task_id", record.TaskID), zap.String("task_status", string(record.TaskStatus)), zap.Int("user_id", record.UserID))
	record.CreatedAt = time.Now()
	if err := db.Create(record).Error; err != nil {
		logger.Error("[DB] CreateLogRecord error", zap.Error(err))
		return err
	}
	return nil
}

func GetLogRecordByID(db *gorm.DB, id int) (*LogRecord, error) {
	var record LogRecord
	if err := db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func GetLogRecordsByTaskID(db *gorm.DB, taskID int) ([]*LogRecord, error) {
	var records []*LogRecord
	if err := db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func GetLogRecordsByServiceID(db *gorm.DB, serviceID int) ([]*LogRecord, error) {
	var records []*LogRecord
	if err := db.Joins("JOIN strm_task st ON log_record.task_id = st.id").Where("st.service_id = ?", serviceID).Order("log_record.created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func GetLogRecordsByStatus(db *gorm.DB, status TaskStatus) ([]*LogRecord, error) {
	var records []*LogRecord
	if err := db.Where("task_status = ?", status).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func UpdateLogRecordStatus(db *gorm.DB, id int, status TaskStatus) error {
	logger.Info("[DB] UpdateLogRecordStatus", zap.Int("id", id), zap.String("status", string(status)))
	if err := db.Model(&LogRecord{}).Where("id = ?", id).Update("task_status", status).Error; err != nil {
		logger.Error("[DB] UpdateLogRecordStatus error", zap.Error(err))
		return err
	}
	return nil
}

func DeleteLogRecord(db *gorm.DB, id int) error {
	logger.Info("[DB] DeleteLogRecord", zap.Int("id", id))
	if err := db.Delete(&LogRecord{}, id).Error; err != nil {
		logger.Error("[DB] DeleteLogRecord error", zap.Error(err))
		return err
	}
	return nil
}

func DeleteLogRecordsByTaskID(db *gorm.DB, taskID int) error {
	logger.Info("[DB] DeleteLogRecordsByTaskID", zap.Int("task_id", taskID))
	if err := db.Where("task_id = ?", taskID).Delete(&LogRecord{}).Error; err != nil {
		logger.Error("[DB] DeleteLogRecordsByTaskID error", zap.Error(err))
		return err
	}
	return nil
}

func GetRecentLogRecords(db *gorm.DB, limit int) ([]*LogRecord, error) {
	var records []*LogRecord
	if err := db.Order("created_at DESC").Limit(limit).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
} 
