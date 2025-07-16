package tests

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
)

func TestDatabaseTables(t *testing.T) {
	// 连接测试数据库
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 执行数据库迁移
	err = model.AutoMigrateAll(db)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 检查表是否存在
	err = model.CheckTablesExist(db)
	if err != nil {
		t.Fatalf("表检查失败: %v", err)
	}

	// 初始化测试数据
	err = model.InitTestData(db)
	if err != nil {
		t.Fatalf("测试数据初始化失败: %v", err)
	}

	// 验证数据查询
	services, err := model.GetOpenListServicesByUserID(db, 1)
	if err != nil {
		t.Fatalf("查询服务失败: %v", err)
	}
	if len(services) == 0 {
		t.Error("应该查询到至少一个服务")
	}

	configs, err := model.GetStrmConfigsByServiceID(db, 1)
	if err != nil {
		t.Fatalf("查询配置失败: %v", err)
	}
	if len(configs) == 0 {
		t.Error("应该查询到至少一个配置")
	}

	tasks, err := model.GetStrmTasksByServiceID(db, 1)
	if err != nil {
		t.Fatalf("查询任务失败: %v", err)
	}
	if len(tasks) == 0 {
		t.Error("应该查询到至少一个任务")
	}

	logs, err := model.GetLogRecordsByTaskID(db, 1)
	if err != nil {
		t.Fatalf("查询日志失败: %v", err)
	}
	if len(logs) == 0 {
		t.Error("应该查询到至少一条日志记录")
	}

	t.Log("✅ 所有数据库表测试通过")
}

func TestOpenListServiceCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	err = model.AutoMigrateAll(db)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建测试用户
	testUser := &model.User{
		Username:     sql.NullString{String: "testuser", Valid: true},
		Email:        sql.NullString{String: "test@example.com", Valid: true},
		PasswordHash: "$2a$10$test_hash",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	err = model.CreateUser(db, testUser)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}

	// 获取用户ID
	var userID int
	err = db.QueryRow("SELECT id FROM user WHERE email = ?", testUser.Email.String).Scan(&userID)
	if err != nil {
		t.Fatalf("获取用户ID失败: %v", err)
	}

	// 测试创建服务
	service := &model.OpenListService{
		Name: "测试服务",
		Account:     "test_account",
		Token:       "test_token",
		ServiceUrl:  "http://localhost:5244",
		Enabled:     true,
		UserID:      userID,
	}

	err = model.CreateOpenListService(db, service)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}

	// 测试查询服务
	services, err := model.GetOpenListServicesByUserID(db, userID)
	if err != nil {
		t.Fatalf("查询服务失败: %v", err)
	}
	if len(services) == 0 {
		t.Error("应该查询到至少一个服务")
	}

	// 测试更新服务
	service = services[0]
	service.Name = "更新后的服务名"
	err = model.UpdateOpenListService(db, service)
	if err != nil {
		t.Fatalf("更新服务失败: %v", err)
	}

	// 测试切换服务状态
	err = model.ToggleOpenListServiceEnabled(db, service.ID, false)
	if err != nil {
		t.Fatalf("切换服务状态失败: %v", err)
	}

	t.Log("✅ OpenListService CRUD测试通过")
}

func TestStrmConfigCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	err = model.AutoMigrateAll(db)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建测试数据
	err = model.InitTestData(db)
	if err != nil {
		t.Fatalf("初始化测试数据失败: %v", err)
	}

	// 测试创建配置
	config := &model.StrmConfig{
		Name:       "新配置",
		AlistBasePath:    "/new/path",
		StrmOutputPath:   "/new/output",
		DownloadEnabled:  true,
		DownloadInterval: 7200,
		UpdateMode:       model.UpdateModeFull,
		ServiceID:        1,
	}

	err = model.CreateStrmConfig(db, config)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	// 测试查询配置
	configs, err := model.GetStrmConfigsByServiceID(db, 1)
	if err != nil {
		t.Fatalf("查询配置失败: %v", err)
	}
	if len(configs) < 2 {
		t.Error("应该查询到至少2个配置")
	}

	// 测试更新配置
	config = configs[0]
	config.Name = "更新后的配置名"
	err = model.UpdateStrmConfig(db, config)
	if err != nil {
		t.Fatalf("更新配置失败: %v", err)
	}

	// 测试切换下载状态
	err = model.ToggleStrmConfigDownloadEnabled(db, config.ID, false)
	if err != nil {
		t.Fatalf("切换下载状态失败: %v", err)
	}

	t.Log("✅ StrmConfig CRUD测试通过")
}

func TestStrmTaskCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	err = model.AutoMigrateAll(db)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建测试数据
	err = model.InitTestData(db)
	if err != nil {
		t.Fatalf("初始化测试数据失败: %v", err)
	}

	// 测试创建任务
	task := &model.StrmTask{
		Name:      "新任务",
		ScheduledTime: time.Now().Add(2 * time.Hour),
		TaskMode:      model.TaskModeCheck,
		Enabled:       true,
		ServiceID:     1,
		ConfigID:      1,
	}

	err = model.CreateStrmTask(db, task)
	if err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}

	// 测试查询任务
	tasks, err := model.GetStrmTasksByServiceID(db, 1)
	if err != nil {
		t.Fatalf("查询任务失败: %v", err)
	}
	if len(tasks) < 2 {
		t.Error("应该查询到至少2个任务")
	}

	// 测试更新任务
	task = tasks[0]
	task.Name = "更新后的任务名"
	err = model.UpdateStrmTask(db, task)
	if err != nil {
		t.Fatalf("更新任务失败: %v", err)
	}

	// 测试切换任务状态
	err = model.ToggleStrmTaskEnabled(db, task.ID, false)
	if err != nil {
		t.Fatalf("切换任务状态失败: %v", err)
	}

	// 测试更新调度时间
	err = model.UpdateStrmTaskScheduledTime(db, task.ID, time.Now().Add(3*time.Hour))
	if err != nil {
		t.Fatalf("更新调度时间失败: %v", err)
	}

	t.Log("✅ StrmTask CRUD测试通过")
}

func TestLogRecordCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	err = model.AutoMigrateAll(db)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建测试数据
	err = model.InitTestData(db)
	if err != nil {
		t.Fatalf("初始化测试数据失败: %v", err)
	}

	// 测试创建日志记录
	logRecord := &model.LogRecord{
		Name:    model.LogNameCheck,
		LogPath:    "/logs/new_test.log",
		TaskStatus: model.TaskStatusRunning,
		TaskID:     1,
	}

	err = model.CreateLogRecord(db, logRecord)
	if err != nil {
		t.Fatalf("创建日志记录失败: %v", err)
	}

	// 测试查询日志记录
	logs, err := model.GetLogRecordsByTaskID(db, 1)
	if err != nil {
		t.Fatalf("查询日志记录失败: %v", err)
	}
	if len(logs) < 2 {
		t.Error("应该查询到至少2条日志记录")
	}

	// 测试更新日志状态
	logRecord = logs[0]
	err = model.UpdateLogRecordStatus(db, logRecord.ID, model.TaskStatusCompleted)
	if err != nil {
		t.Fatalf("更新日志状态失败: %v", err)
	}

	// 测试按状态查询
	completedLogs, err := model.GetLogRecordsByStatus(db, model.TaskStatusCompleted)
	if err != nil {
		t.Fatalf("按状态查询日志失败: %v", err)
	}
	if len(completedLogs) == 0 {
		t.Error("应该查询到至少一条已完成的日志记录")
	}

	t.Log("✅ LogRecord CRUD测试通过")
} 
