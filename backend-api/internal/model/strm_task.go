package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

type TaskMode string

const (
	TaskModeCreate TaskMode = "create"
	TaskModeCheck  TaskMode = "check"
)

type StrmTask struct {
	ID            int
	UserID        int
	Name          string
	ScheduledTime time.Time
	TaskMode      TaskMode
	Enabled       bool
	ServiceID     int
	ConfigID      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CreateStrmTask 创建Strm任务
func CreateStrmTask(db *sql.DB, task *StrmTask) error {
	logger.Info("[DB] CreateStrmTask", zap.String("name", task.Name), zap.Int("service_id", task.ServiceID), zap.Int("config_id", task.ConfigID), zap.Int("user_id", task.UserID))
	
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	
	_, err := db.Exec(`
		INSERT INTO strm_task (user_id, name, scheduled_time, task_mode, enabled, service_id, config_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.UserID, task.Name, task.ScheduledTime, task.TaskMode, task.Enabled, task.ServiceID, task.ConfigID, task.CreatedAt, task.UpdatedAt)
	
	if err != nil {
		logger.Error("[DB] CreateStrmTask error", zap.Error(err))
		return err
	}
	
	return nil
}

// GetStrmTaskByID 根据ID获取Strm任务
func GetStrmTaskByID(db *sql.DB, id int) (*StrmTask, error) {
	var task StrmTask
	err := db.QueryRow(`
		SELECT id, user_id, name, scheduled_time, task_mode, enabled, service_id, config_id, created_at, updated_at 
		FROM strm_task WHERE id = ?`, id).
		Scan(&task.ID, &task.UserID, &task.Name, &task.ScheduledTime, &task.TaskMode, &task.Enabled, &task.ServiceID, &task.ConfigID, &task.CreatedAt, &task.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetStrmTasksByServiceID 根据服务ID获取所有Strm任务
func GetStrmTasksByServiceID(db *sql.DB, serviceID int) ([]*StrmTask, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, scheduled_time, task_mode, enabled, service_id, config_id, created_at, updated_at 
		FROM strm_task WHERE service_id = ? ORDER BY scheduled_time ASC`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []*StrmTask
	for rows.Next() {
		var task StrmTask
		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.ScheduledTime, &task.TaskMode, &task.Enabled, &task.ServiceID, &task.ConfigID, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	
	return tasks, nil
}

// GetStrmTasksByConfigID 根据配置ID获取所有Strm任务
func GetStrmTasksByConfigID(db *sql.DB, configID int) ([]*StrmTask, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, scheduled_time, task_mode, enabled, service_id, config_id, created_at, updated_at 
		FROM strm_task WHERE config_id = ? ORDER BY scheduled_time ASC`, configID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []*StrmTask
	for rows.Next() {
		var task StrmTask
		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.ScheduledTime, &task.TaskMode, &task.Enabled, &task.ServiceID, &task.ConfigID, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	
	return tasks, nil
}

// GetEnabledStrmTasks 获取所有启用的任务
func GetEnabledStrmTasks(db *sql.DB) ([]*StrmTask, error) {
	rows, err := db.Query(`
		SELECT id, user_id, name, scheduled_time, task_mode, enabled, service_id, config_id, created_at, updated_at 
		FROM strm_task WHERE enabled = 1 ORDER BY scheduled_time ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []*StrmTask
	for rows.Next() {
		var task StrmTask
		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.ScheduledTime, &task.TaskMode, &task.Enabled, &task.ServiceID, &task.ConfigID, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	
	return tasks, nil
}

// UpdateStrmTask 更新Strm任务
func UpdateStrmTask(db *sql.DB, task *StrmTask) error {
	logger.Info("[DB] UpdateStrmTask", zap.Int("id", task.ID))
	
	task.UpdatedAt = time.Now()
	
	_, err := db.Exec(`
		UPDATE strm_task 
		SET user_id = ?, name = ?, scheduled_time = ?, task_mode = ?, enabled = ?, service_id = ?, config_id = ?, updated_at = ? 
		WHERE id = ?`,
		task.UserID, task.Name, task.ScheduledTime, task.TaskMode, task.Enabled, task.ServiceID, task.ConfigID, task.UpdatedAt, task.ID)
	
	if err != nil {
		logger.Error("[DB] UpdateStrmTask error", zap.Error(err))
		return err
	}
	
	return nil
}

// DeleteStrmTask 删除Strm任务
func DeleteStrmTask(db *sql.DB, id int) error {
	logger.Info("[DB] DeleteStrmTask", zap.Int("id", id))
	
	_, err := db.Exec("DELETE FROM strm_task WHERE id = ?", id)
	if err != nil {
		logger.Error("[DB] DeleteStrmTask error", zap.Error(err))
		return err
	}
	
	return nil
}

// ToggleStrmTaskEnabled 切换任务启用状态
func ToggleStrmTaskEnabled(db *sql.DB, id int, enabled bool) error {
	logger.Info("[DB] ToggleStrmTaskEnabled", zap.Int("id", id), zap.Bool("enabled", enabled))
	
	_, err := db.Exec("UPDATE strm_task SET enabled = ?, updated_at = ? WHERE id = ?", enabled, time.Now(), id)
	if err != nil {
		logger.Error("[DB] ToggleStrmTaskEnabled error", zap.Error(err))
		return err
	}
	
	return nil
}

// UpdateStrmTaskScheduledTime 更新任务调度时间
func UpdateStrmTaskScheduledTime(db *sql.DB, id int, scheduledTime time.Time) error {
	logger.Info("[DB] UpdateStrmTaskScheduledTime", zap.Int("id", id), zap.Time("scheduled_time", scheduledTime))
	
	_, err := db.Exec("UPDATE strm_task SET scheduled_time = ?, updated_at = ? WHERE id = ?", scheduledTime, time.Now(), id)
	if err != nil {
		logger.Error("[DB] UpdateStrmTaskScheduledTime error", zap.Error(err))
		return err
	}
	
	return nil
} 
