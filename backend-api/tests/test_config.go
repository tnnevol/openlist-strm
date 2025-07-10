package tests

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
)

// TestConfig 测试配置
type TestConfig struct {
	DB *sql.DB
}

var testConfig *TestConfig

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T) *sql.DB {
	// 使用内存数据库进行测试
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 创建测试表
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255),
		is_active BOOLEAN DEFAULT FALSE,
		code VARCHAR(10),
		code_expire_at DATETIME,
		failed_login_count INTEGER DEFAULT 0,
		locked_until DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	return db
}

// SetupTestEnvironment 设置测试环境
func SetupTestEnvironment(t *testing.T) *TestConfig {
	// 初始化日志（测试模式）
	logger.Init()

	// 设置测试数据库
	db := SetupTestDB(t)

	// 设置环境变量
	os.Setenv("JWT_SECRET", "test-secret-key")

	testConfig = &TestConfig{
		DB: db,
	}

	return testConfig
}

// CleanupTestEnvironment 清理测试环境
func CleanupTestEnvironment(t *testing.T) {
	if testConfig != nil && testConfig.DB != nil {
		testConfig.DB.Close()
	}
}

// InsertTestUser 插入测试用户数据
func InsertTestUser(t *testing.T, db *sql.DB, username, email, passwordHash string, isActive bool) int64 {
	query := `
	INSERT INTO user (username, email, password_hash, is_active, created_at, updated_at)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	
	result, err := db.Exec(query, username, email, passwordHash, isActive)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert id: %v", err)
	}

	return userID
}

// GetTestUser 获取测试用户
func GetTestUser(t *testing.T, db *sql.DB, userID int64) map[string]interface{} {
	query := `SELECT id, username, email, is_active, created_at FROM user WHERE id = ?`
	
	var user map[string]interface{}
	user = make(map[string]interface{})
	
	var id int64
	var username, email sql.NullString
	var isActive bool
	var createdAt string
	
	err := db.QueryRow(query, userID).Scan(&id, &username, &email, &isActive, &createdAt)
	if err != nil {
		t.Fatalf("Failed to get test user: %v", err)
	}

	user["id"] = id
	user["username"] = username.String
	user["email"] = email.String
	user["isActive"] = isActive
	user["createdAt"] = createdAt

	return user
}

// ClearTestData 清理测试数据
func ClearTestData(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM user")
	if err != nil {
		t.Fatalf("Failed to clear test data: %v", err)
	}
} 
