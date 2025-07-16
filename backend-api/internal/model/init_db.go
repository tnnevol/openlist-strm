package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitTestData 初始化测试数据
func InitTestData(db *gorm.DB) error {
	logger.Info("[DB] 开始初始化测试数据")

	// 创建测试用户
	testUser := &User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$test_hash_for_testing",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	err := CreateUser(db, testUser)
	if err != nil {
		logger.Error("[DB] 创建测试用户失败", zap.Error(err))
		return err
	}

	// 获取用户ID（假设是第一个用户）
	var user User
	err = db.Where("email = ?", testUser.Email).First(&user).Error
	if err != nil {
		logger.Error("[DB] 获取用户ID失败", zap.Error(err))
		return err
	}
	userID := user.ID

	// 创建测试OpenList服务
	testService := &OpenListService{
		Name:      "测试服务",
		Account:   "test_account",
		Token:     "test_token_123",
		ServiceUrl:  "http://localhost:5244",
		BackupUrl:   "http://backup.localhost:5244",
		Enabled:     true,
		UserID:      userID,
	}

	err = CreateOpenListService(db, testService)
	if err != nil {
		logger.Error("[DB] 创建测试服务失败", zap.Error(err))
		return err
	}

	// 获取服务ID
	var service OpenListService
	err = db.Where("name = ?", testService.Name).First(&service).Error
	if err != nil {
		logger.Error("[DB] 获取服务ID失败", zap.Error(err))
		return err
	}
	serviceID := service.ID

	// 创建测试Strm配置
	testConfig := &StrmConfig{
		Name:       "测试配置",
		AlistBasePath:    "/test/path",
		StrmOutputPath:   "/output/strm",
		DownloadEnabled:  true,
		DownloadInterval: 3600,
		UpdateMode:       UpdateModeIncremental,
		ServiceID:        serviceID,
		UserID:           userID,
		IsUseBackupUrl:   true,
	}

	err = CreateStrmConfig(db, testConfig)
	if err != nil {
		logger.Error("[DB] 创建测试配置失败", zap.Error(err))
		return err
	}

	// 获取配置ID
	var config StrmConfig
	err = db.Where("name = ?", testConfig.Name).First(&config).Error
	if err != nil {
		logger.Error("[DB] 获取配置ID失败", zap.Error(err))
		return err
	}
	configID := config.ID

	// 创建测试Strm任务
	testTask := &StrmTask{
		Name:      "测试任务",
		ScheduledTime: time.Now().Add(time.Hour),
		TaskMode:      TaskModeCreate,
		Enabled:       true,
		ServiceID:     serviceID,
		ConfigID:      configID,
		UserID:        userID,
	}

	err = CreateStrmTask(db, testTask)
	if err != nil {
		logger.Error("[DB] 创建测试任务失败", zap.Error(err))
		return err
	}

	// 获取任务ID
	var task StrmTask
	err = db.Where("name = ?", testTask.Name).First(&task).Error
	if err != nil {
		logger.Error("[DB] 获取任务ID失败", zap.Error(err))
		return err
	}
	taskID := task.ID

	// 创建测试日志记录
	testLog := &LogRecord{
		Name:    LogNameCreate,
		LogPath:    "/logs/test.log",
		TaskStatus: TaskStatusCompleted,
		TaskID:     taskID,
		UserID:     userID,
	}

	err = CreateLogRecord(db, testLog)
	if err != nil {
		logger.Error("[DB] 创建测试日志失败", zap.Error(err))
		return err
	}

	logger.Info("[DB] 测试数据初始化完成", 
		zap.Int("user_id", userID),
		zap.Int("service_id", serviceID),
		zap.Int("config_id", configID),
		zap.Int("task_id", taskID))

	return nil
}

// CheckTablesExist 检查表是否存在（GORM 迁移后可省略，建议用 AutoMigrate）
func CheckTablesExist(db *gorm.DB) error {
	tables := []interface{}{&User{}, &OpenListService{}, &StrmConfig{}, &StrmTask{}, &LogRecord{}}
	for _, table := range tables {
		err := db.AutoMigrate(table)
		if err != nil {
			logger.Error("[DB] 表迁移失败", zap.Error(err))
			return err
		}
		logger.Info("[DB] 表已存在或迁移成功", zap.String("table", db.Migrator().CurrentDatabase()))
	}
	return nil
} 
