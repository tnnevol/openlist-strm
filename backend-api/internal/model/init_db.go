package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

// InitTestData 初始化测试数据
func InitTestData(db *sql.DB) error {
	logger.Info("[DB] 开始初始化测试数据")
	
	// 创建测试用户
	testUser := &User{
		Username:     sql.NullString{String: "testuser", Valid: true},
		Email:        sql.NullString{String: "test@example.com", Valid: true},
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
	var userID int
	err = db.QueryRow("SELECT id FROM user WHERE email = ?", testUser.Email.String).Scan(&userID)
	if err != nil {
		logger.Error("[DB] 获取用户ID失败", zap.Error(err))
		return err
	}
	
	// 创建测试OpenList服务
	testService := &OpenListService{
		Name: "测试服务",
		Account:     "test_account",
		Token:       "test_token_123",
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
	var serviceID int
	err = db.QueryRow("SELECT id FROM openlist_service WHERE name = ?", testService.Name).Scan(&serviceID)
	if err != nil {
		logger.Error("[DB] 获取服务ID失败", zap.Error(err))
		return err
	}
	
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
	var configID int
	err = db.QueryRow("SELECT id FROM strm_config WHERE name = ?", testConfig.Name).Scan(&configID)
	if err != nil {
		logger.Error("[DB] 获取配置ID失败", zap.Error(err))
		return err
	}
	
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
	var taskID int
	err = db.QueryRow("SELECT id FROM strm_task WHERE name = ?", testTask.Name).Scan(&taskID)
	if err != nil {
		logger.Error("[DB] 获取任务ID失败", zap.Error(err))
		return err
	}
	
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

// CheckTablesExist 检查表是否存在
func CheckTablesExist(db *sql.DB) error {
	tables := []string{"user", "openlist_service", "strm_config", "strm_task", "log_record"}
	
	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		if err != nil {
			logger.Error("[DB] 检查表失败", zap.String("table", table), zap.Error(err))
			return err
		}
		
		if count == 0 {
			logger.Error("[DB] 表不存在", zap.String("table", table))
			return err
		}
		
		logger.Info("[DB] 表存在", zap.String("table", table))
	}
	
	return nil
} 
