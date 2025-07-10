package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tnnevol/openlist-strm/backend-api/internal/app"
)

// TestRequest 测试请求结构
type TestRequest struct {
	Method      string
	URL         string
	Body        interface{}
	Headers     map[string]string
	ExpectedCode int
	ExpectedBody interface{}
}

// CreateTestServer 创建测试服务器
func CreateTestServer(t *testing.T, db *sql.DB) *gin.Engine {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	
	// 创建路由
	r := app.RegisterRouter(db)
	
	return r
}

// MakeTestRequest 发送测试请求
func MakeTestRequest(t *testing.T, router *gin.Engine, req TestRequest) *httptest.ResponseRecorder {
	// 准备请求体
	var bodyBytes []byte
	var err error
	
	if req.Body != nil {
		bodyBytes, err = json.Marshal(req.Body)
		assert.NoError(t, err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, httpReq)

	return w
}

// AssertResponse 断言响应
func AssertResponse(t *testing.T, w *httptest.ResponseRecorder, expectedCode int, expectedBody interface{}) {
	// 断言状态码
	assert.Equal(t, expectedCode, w.Code)

	// 如果有期望的响应体，进行断言
	if expectedBody != nil {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 如果期望的响应体是map，进行深度比较
		if expectedMap, ok := expectedBody.(map[string]interface{}); ok {
			for key, expectedValue := range expectedMap {
				// 处理JSON中数字类型为float64的情况
				if key == "code" {
					if expectedInt, ok := expectedValue.(int); ok {
						assert.Equal(t, float64(expectedInt), response[key], "Field: %s", key)
					} else {
						assert.Equal(t, expectedValue, response[key], "Field: %s", key)
					}
				} else {
					assert.Equal(t, expectedValue, response[key], "Field: %s", key)
				}
			}
		}
	}
}

// CreateTestToken 创建测试用的JWT token
func CreateTestToken(t *testing.T, userID int, username, email string) string {
	// 这里可以创建一个简单的测试token
	// 在实际项目中，你可能需要使用真实的JWT库
	return "test-token-" + username
}

// TestUserData 测试用户数据
type TestUserData struct {
	Username     string
	Email        string
	Password     string
	IsActive     bool
	PasswordHash string
}

// GetDefaultTestUser 获取默认测试用户数据
func GetDefaultTestUser() TestUserData {
	return TestUserData{
		Username:     "testuser",
		Email:        "test@example.com",
		Password:     "TestPass123",
		IsActive:     true,
		PasswordHash: "$2a$10$test.hash.for.testing",
	}
}

// SetupTestUser 设置测试用户
func SetupTestUser(t *testing.T, db *sql.DB, userData TestUserData) int64 {
	return InsertTestUser(t, db, userData.Username, userData.Email, userData.PasswordHash, userData.IsActive)
}

// CleanupTestUser 清理测试用户
func CleanupTestUser(t *testing.T, db *sql.DB) {
	ClearTestData(t, db)
} 
