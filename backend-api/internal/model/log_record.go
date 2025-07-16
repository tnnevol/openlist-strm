package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
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
	ID         int
	UserID     int
	Name       LogName
	LogPath    string
	CreatedAt  time.Time
	TaskStatus TaskStatus
	TaskID     int
}

// CreateLogRecord 创建日志记录
func CreateLogRecord(db *sql.DB, record *LogRecord) error {
	logger.Info("[DB] CreateLogRecord", zap.String("name", string(record.Name)), zap.Int("task_id", record.TaskID), zap.String("task_status", string(record.TaskStatus)), zap.Int("user_id", record.UserID))
	
	record.CreatedAt = time.Now()
	
	_, err := db.Exec(`
		INSERT INTO log_record (user_id, name, log_path, created_at, task_status, task_id) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		record.UserID, record.Name, record.LogPath, record.CreatedAt, record.TaskStatus, record.TaskID)
	
	if err != nil {
		logger.Error("[DB] CreateLogRecord error", zap.Error(err))
		return err
	}
	
	return nil
}

// GetLogRecordByID 根据ID获取日志记录
func GetLogRecordByID(db *sql.DB, id int) (*LogRecord, error) {
	var record LogRecord
	err := db.QueryRow(`
		SELECT id, user_id, name, log_path, created_at, task_status, task_id 
		FROM log_record WHERE id = ?`, id).
		Scan(&record.ID, &record.UserID, &record.Name, &record.LogPath, &record.CreatedAt, &record.TaskStatus, &record.TaskID)
	
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetLogRecordsByTaskID 根据任务ID获取所有日志记录
func GetLogRecordsByTaskID(db *sql.DB, taskID int) ([]*LogRecord, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, log_path, created_at, task_status, task_id 
		FROM log_record WHERE task_id = ? ORDER BY created_at DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*LogRecord
	for rows.Next() {
		var record LogRecord
		err := rows.Scan(&record.ID, &record.UserID, &record.Name, &record.LogPath, &record.CreatedAt, &record.TaskStatus, &record.TaskID)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	
	return records, nil
}

// GetLogRecordsByServiceID 根据服务ID获取所有日志记录（通过任务关联）
func GetLogRecordsByServiceID(db *sql.DB, serviceID int) ([]*LogRecord, error) {
	rows, err := db.Query(`
		SELECT lr.id, lr.user_id, lr.name, lr.log_path, lr.created_at, lr.task_status, lr.task_id 
		FROM log_record lr 
		JOIN strm_task st ON lr.task_id = st.id 
		WHERE st.service_id = ? 
		ORDER BY lr.created_at DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*LogRecord
	for rows.Next() {
		var record LogRecord
		err := rows.Scan(&record.ID, &record.UserID, &record.Name, &record.LogPath, &record.CreatedAt, &record.TaskStatus, &record.TaskID)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	
	return records, nil
}

// GetLogRecordsByStatus 根据状态获取日志记录
func GetLogRecordsByStatus(db *sql.DB, status TaskStatus) ([]*LogRecord, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, log_path, created_at, task_status, task_id 
		FROM log_record WHERE task_status = ? ORDER BY created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*LogRecord
	for rows.Next() {
		var record LogRecord
		err := rows.Scan(&record.ID, &record.UserID, &record.Name, &record.LogPath, &record.CreatedAt, &record.TaskStatus, &record.TaskID)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	
	return records, nil
}

// UpdateLogRecordStatus 更新日志记录状态
func UpdateLogRecordStatus(db *sql.DB, id int, status TaskStatus) error {
	logger.Info("[DB] UpdateLogRecordStatus", zap.Int("id", id), zap.String("status", string(status)))
	
	_, err := db.Exec("UPDATE log_record SET task_status = ? WHERE id = ?", status, id)
	if err != nil {
		logger.Error("[DB] UpdateLogRecordStatus error", zap.Error(err))
		return err
	}
	
	return nil
}

// DeleteLogRecord 删除日志记录
func DeleteLogRecord(db *sql.DB, id int) error {
	logger.Info("[DB] DeleteLogRecord", zap.Int("id", id))
	
	_, err := db.Exec("DELETE FROM log_record WHERE id = ?", id)
	if err != nil {
		logger.Error("[DB] DeleteLogRecord error", zap.Error(err))
		return err
	}
	
	return nil
}

// DeleteLogRecordsByTaskID 根据任务ID删除所有日志记录
func DeleteLogRecordsByTaskID(db *sql.DB, taskID int) error {
	logger.Info("[DB] DeleteLogRecordsByTaskID", zap.Int("task_id", taskID))
	
	_, err := db.Exec("DELETE FROM log_record WHERE task_id = ?", taskID)
	if err != nil {
		logger.Error("[DB] DeleteLogRecordsByTaskID error", zap.Error(err))
		return err
	}
	
	return nil
}

// GetRecentLogRecords 获取最近的日志记录
func GetRecentLogRecords(db *sql.DB, limit int) ([]*LogRecord, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, log_path, created_at, task_status, task_id 
		FROM log_record ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*LogRecord
	for rows.Next() {
		var record LogRecord
		err := rows.Scan(&record.ID, &record.UserID, &record.Name, &record.LogPath, &record.CreatedAt, &record.TaskStatus, &record.TaskID)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	
	return records, nil
} 
