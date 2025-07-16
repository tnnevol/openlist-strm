package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskMode string

const (
	TaskModeCreate TaskMode = "create"
	TaskModeCheck  TaskMode = "check"
)

type StrmTask struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int       `json:"userId"`
	Name          string    `json:"name" gorm:"type:varchar(128);index"`
	ScheduledTime time.Time `json:"scheduledTime"`
	TaskMode      TaskMode  `json:"taskMode" gorm:"type:varchar(32)"`
	Enabled       bool      `json:"enabled"`
	ServiceID     int       `json:"serviceId"`
	ConfigID      int       `json:"configId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func CreateStrmTask(db *gorm.DB, task *StrmTask) error {
	logger.Info("[DB] CreateStrmTask", zap.String("name", task.Name), zap.Int("service_id", task.ServiceID), zap.Int("config_id", task.ConfigID), zap.Int("user_id", task.UserID))
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	if err := db.Create(task).Error; err != nil {
		logger.Error("[DB] CreateStrmTask error", zap.Error(err))
		return err
	}
	return nil
}

func GetStrmTaskByID(db *gorm.DB, id int) (*StrmTask, error) {
	var task StrmTask
	if err := db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func GetStrmTasksByServiceID(db *gorm.DB, serviceID int) ([]*StrmTask, error) {
	var tasks []*StrmTask
	if err := db.Where("service_id = ?", serviceID).Order("scheduled_time ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetStrmTasksByConfigID(db *gorm.DB, configID int) ([]*StrmTask, error) {
	var tasks []*StrmTask
	if err := db.Where("config_id = ?", configID).Order("scheduled_time ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetEnabledStrmTasks(db *gorm.DB) ([]*StrmTask, error) {
	var tasks []*StrmTask
	if err := db.Where("enabled = ?", true).Order("scheduled_time ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func UpdateStrmTask(db *gorm.DB, task *StrmTask) error {
	logger.Info("[DB] UpdateStrmTask", zap.Int("id", task.ID))
	task.UpdatedAt = time.Now()
	if err := db.Model(&StrmTask{}).Where("id = ?", task.ID).Updates(task).Error; err != nil {
		logger.Error("[DB] UpdateStrmTask error", zap.Error(err))
		return err
	}
	return nil
}

func DeleteStrmTask(db *gorm.DB, id int) error {
	logger.Info("[DB] DeleteStrmTask", zap.Int("id", id))
	if err := db.Delete(&StrmTask{}, id).Error; err != nil {
		logger.Error("[DB] DeleteStrmTask error", zap.Error(err))
		return err
	}
	return nil
}

func ToggleStrmTaskEnabled(db *gorm.DB, id int, enabled bool) error {
	logger.Info("[DB] ToggleStrmTaskEnabled", zap.Int("id", id), zap.Bool("enabled", enabled))
	if err := db.Model(&StrmTask{}).Where("id = ?", id).Updates(map[string]interface{}{"enabled": enabled, "updated_at": time.Now()}).Error; err != nil {
		logger.Error("[DB] ToggleStrmTaskEnabled error", zap.Error(err))
		return err
	}
	return nil
}

func UpdateStrmTaskScheduledTime(db *gorm.DB, id int, scheduledTime time.Time) error {
	logger.Info("[DB] UpdateStrmTaskScheduledTime", zap.Int("id", id), zap.Time("scheduled_time", scheduledTime))
	if err := db.Model(&StrmTask{}).Where("id = ?", id).Updates(map[string]interface{}{"scheduled_time": scheduledTime, "updated_at": time.Now()}).Error; err != nil {
		logger.Error("[DB] UpdateStrmTaskScheduledTime error", zap.Error(err))
		return err
	}
	return nil
} 
