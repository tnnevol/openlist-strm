package tests

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"github.com/tnnevol/openlist-strm/backend-api/tests/config"
	"golang.org/x/crypto/bcrypt"
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

	// 创建测试表 - 与实际项目结构保持一致
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50),
		email VARCHAR(100) NOT NULL,
		password_hash VARCHAR(255),
		is_active BOOLEAN DEFAULT FALSE,
		code VARCHAR(10),
		code_expire_at DATETIME,
		failed_login_count INTEGER DEFAULT 0,
		locked_until DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		token_invalid_before DATETIME
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

	// 设置全局数据库连接（供中间件使用）
	service.SetGlobalDB(db)

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
	INSERT INTO user (username, email, password_hash, is_active, created_at, token_invalid_before)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, ?)
	`
	
	// 设置token_invalid_before为当前时间之前，这样新生成的token不会被拒绝
	tokenInvalidBefore := time.Now().Add(-1 * time.Hour) // 1小时前
	
	result, err := db.Exec(query, username, email, passwordHash, isActive, tokenInvalidBefore)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert id: %v", err)
	}

	return userID
}

// InsertTestUserFromConfig 从配置插入测试用户
func InsertTestUserFromConfig(t *testing.T, db *sql.DB, testUser config.TestUser) int64 {
	// 对密码进行哈希处理
	passwordHash := hashPassword(testUser.Password)
	
	return InsertTestUser(t, db, testUser.Username, testUser.Email, passwordHash, testUser.IsActive)
}

// SetupTestUsersFromConfig 从配置设置测试用户
func SetupTestUsersFromConfig(t *testing.T, db *sql.DB, userConfig *config.TestConfigFile) []int64 {
	var userIDs []int64
	
	// 插入所有活跃用户
	activeUsers := userConfig.GetActiveUsers()
	for _, user := range activeUsers {
		userID := InsertTestUserFromConfig(t, db, user)
		userIDs = append(userIDs, userID)
	}
	
	return userIDs
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

// GetTestUserByUsername 根据用户名获取测试用户
func GetTestUserByUsername(t *testing.T, db *sql.DB, username string) map[string]interface{} {
	query := `SELECT id, username, email, is_active, created_at FROM user WHERE username = ?`
	
	var user map[string]interface{}
	user = make(map[string]interface{})
	
	var id int64
	var userUsername, email sql.NullString
	var isActive bool
	var createdAt string
	
	err := db.QueryRow(query, username).Scan(&id, &userUsername, &email, &isActive, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		t.Fatalf("Failed to get test user by username: %v", err)
	}

	user["id"] = id
	user["username"] = userUsername.String
	user["email"] = email.String
	user["isActive"] = isActive
	user["createdAt"] = createdAt

	return user
}

// GetTestUserByEmail 根据邮箱获取测试用户
func GetTestUserByEmail(t *testing.T, db *sql.DB, email string) map[string]interface{} {
	query := `SELECT id, username, email, is_active, created_at FROM user WHERE email = ?`
	
	var user map[string]interface{}
	user = make(map[string]interface{})
	
	var id int64
	var username, userEmail sql.NullString
	var isActive bool
	var createdAt string
	
	err := db.QueryRow(query, email).Scan(&id, &username, &userEmail, &isActive, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		t.Fatalf("Failed to get test user by email: %v", err)
	}

	user["id"] = id
	user["username"] = username.String
	user["email"] = userEmail.String
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

// hashPassword 对密码进行哈希处理
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// 如果bcrypt失败，回退到SHA256（不推荐，但保持兼容性）
		shaHash := sha256.Sum256([]byte(password))
		return hex.EncodeToString(shaHash[:])
	}
	return string(hash)
}

// HashPasswordForTest 供外部调用的密码哈希
func HashPasswordForTest(password string) string {
	return hashPassword(password)
}

// ValidateTestUser 验证测试用户数据
func ValidateTestUser(t *testing.T, user config.TestUser) error {
	if user.Email == "" {
		return fmt.Errorf("用户 %s 缺少邮箱", user.Name)
	}
	if user.Username == "" {
		return fmt.Errorf("用户 %s 缺少用户名", user.Name)
	}
	if user.Password == "" {
		return fmt.Errorf("用户 %s 缺少密码", user.Name)
	}
	return nil
}

// GetDefaultTestUserMap 获取默认测试用户（map格式，保持向后兼容）
func GetDefaultTestUserMap() map[string]interface{} {
	return map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "TestPass123!",
		"isActive": true,
	}
}

// SetupTestUserMap 设置测试用户（map格式，保持向后兼容）
func SetupTestUserMap(t *testing.T, db *sql.DB, userData map[string]interface{}) int64 {
	username := userData["username"].(string)
	email := userData["email"].(string)
	password := userData["password"].(string)
	isActive := userData["isActive"].(bool)
	
	return InsertTestUser(t, db, username, email, hashPassword(password), isActive)
} 
