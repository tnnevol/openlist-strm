package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tnnevol/openlist-strm/backend-api/tests"
)

// TokenExpiredIntegrationTestSuite token过期集成测试套件
type TokenExpiredIntegrationTestSuite struct {
	suite.Suite
	config *tests.TestConfig
}

// SetupSuite 设置测试套件
func (suite *TokenExpiredIntegrationTestSuite) SetupSuite() {
	suite.config = tests.SetupTestEnvironment(suite.T())
}

// TearDownSuite 清理测试套件
func (suite *TokenExpiredIntegrationTestSuite) TearDownSuite() {
	tests.CleanupTestEnvironment(suite.T())
}

// SetupTest 每个测试前的设置
func (suite *TokenExpiredIntegrationTestSuite) SetupTest() {
	tests.ClearTestData(suite.T(), suite.config.DB)
}

// TestTokenExpiredHandling 测试token过期处理
func (suite *TokenExpiredIntegrationTestSuite) TestTokenExpiredHandling() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	suite.Run("token过期处理流程", func() {
		// 步骤1: 创建测试用户并获取有效token
		userID := tests.InsertTestUser(suite.T(), suite.config.DB, userData.Username, userData.Email, userData.PasswordHash, userData.IsActive)
		assert.Greater(suite.T(), userID, int64(0))

		// 创建有效token
		validToken := suite.createValidToken(userData.Username, userData.Email, userID)

		// 步骤2: 使用有效token访问需要认证的接口
		validReq := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			Headers:     map[string]string{"Authorization": "Bearer " + validToken},
			ExpectedCode: 200,
		}

		w := tests.MakeTestRequest(suite.T(), router, validReq)
		assert.Equal(suite.T(), 200, w.Code)

		// 步骤3: 生成过期token
		expiredTokenReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/generate-expired-token",
			Headers:     map[string]string{"Authorization": "Bearer " + validToken},
			ExpectedCode: 200,
		}

		w = tests.MakeTestRequest(suite.T(), router, expiredTokenReq)
		assert.Equal(suite.T(), 200, w.Code)

		// 解析响应获取过期token
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(suite.T(), ok)
		expiredToken, ok := data["token"].(string)
		assert.True(suite.T(), ok)
		assert.NotEmpty(suite.T(), expiredToken)

		// 步骤4: 等待token过期（1秒）
		time.Sleep(2 * time.Second)

		// 步骤5: 使用过期token访问接口
		expiredReq := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			Headers:     map[string]string{"Authorization": "Bearer " + expiredToken},
			ExpectedCode: 401,
		}

		w = tests.MakeTestRequest(suite.T(), router, expiredReq)
		assert.Equal(suite.T(), 401, w.Code)

		// 验证返回的是40101错误码
		var expiredResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &expiredResponse)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(40101), expiredResponse["code"])
		assert.Equal(suite.T(), "token已过期", expiredResponse["message"])
	})
}

// TestNoTokenHandling 测试无token访问
func (suite *TokenExpiredIntegrationTestSuite) TestNoTokenHandling() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	suite.Run("无token访问处理", func() {
		// 测试无token访问需要认证的接口
		req := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			ExpectedCode: 401,
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), 401, w.Code)

		// 验证返回的是401错误码（不是40101）
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(401), response["code"])
		assert.Equal(suite.T(), "未登录或token缺失", response["message"])
	})
}

// TestInvalidTokenHandling 测试无效token
func (suite *TokenExpiredIntegrationTestSuite) TestInvalidTokenHandling() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	suite.Run("无效token处理", func() {
		// 测试无效token访问需要认证的接口
		req := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			Headers:     map[string]string{"Authorization": "Bearer invalid-token"},
			ExpectedCode: 401,
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), 401, w.Code)

		// 验证返回的是401错误码（不是40101）
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(401), response["code"])
		assert.Equal(suite.T(), "token无效或已过期", response["message"])
	})
}

// TestTokenBlacklistWithExpiredToken 测试过期token在黑名单中的处理
func (suite *TokenExpiredIntegrationTestSuite) TestTokenBlacklistWithExpiredToken() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	suite.Run("过期token黑名单处理", func() {
		// 创建测试用户
		userID := tests.InsertTestUser(suite.T(), suite.config.DB, userData.Username, userData.Email, userData.PasswordHash, userData.IsActive)
		assert.Greater(suite.T(), userID, int64(0))

		// 创建有效token
		validToken := suite.createValidToken(userData.Username, userData.Email, userID)

		// 先登出，将token加入黑名单
		logoutReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/logout",
			Headers:     map[string]string{"Authorization": "Bearer " + validToken},
			ExpectedCode: 200,
		}

		w := tests.MakeTestRequest(suite.T(), router, logoutReq)
		assert.Equal(suite.T(), 200, w.Code)

		// 使用已加入黑名单的token访问接口
		blacklistedReq := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			Headers:     map[string]string{"Authorization": "Bearer " + validToken},
			ExpectedCode: 401,
		}

		w = tests.MakeTestRequest(suite.T(), router, blacklistedReq)
		assert.Equal(suite.T(), 401, w.Code)

		// 验证返回的是401错误码（不是40101）
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(401), response["code"])
		assert.Equal(suite.T(), "token已失效，请重新登录", response["message"])
	})
}

// TestGenerateExpiredTokenEndpoint 测试生成过期token接口
func (suite *TokenExpiredIntegrationTestSuite) TestGenerateExpiredTokenEndpoint() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	suite.Run("生成过期token接口", func() {
		// 创建测试用户
		userID := tests.InsertTestUser(suite.T(), suite.config.DB, userData.Username, userData.Email, userData.PasswordHash, userData.IsActive)
		assert.Greater(suite.T(), userID, int64(0))

		// 创建有效token
		validToken := suite.createValidToken(userData.Username, userData.Email, userID)

		// 测试生成过期token接口
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/generate-expired-token",
			Headers:     map[string]string{"Authorization": "Bearer " + validToken},
			ExpectedCode: 200,
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), 200, w.Code)

		// 验证响应格式
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(suite.T(), ok)
		
		// 验证返回的字段
		assert.Contains(suite.T(), data, "token")
		assert.Contains(suite.T(), data, "expires_in")
		assert.Contains(suite.T(), data, "message")
		
		token, ok := data["token"].(string)
		assert.True(suite.T(), ok)
		assert.NotEmpty(suite.T(), token)
		
		expiresIn, ok := data["expires_in"].(float64)
		assert.True(suite.T(), ok)
		assert.Equal(suite.T(), float64(1), expiresIn)
		
		message, ok := data["message"].(string)
		assert.True(suite.T(), ok)
		assert.Equal(suite.T(), "过期token已生成，1秒后过期", message)
	})
}

// TestGenerateExpiredTokenWithoutAuth 测试无认证访问生成过期token接口
func (suite *TokenExpiredIntegrationTestSuite) TestGenerateExpiredTokenWithoutAuth() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	suite.Run("无认证访问生成过期token接口", func() {
		// 测试无token访问生成过期token接口
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/generate-expired-token",
			ExpectedCode: 401,
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), 401, w.Code)

		// 验证返回的是401错误码
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(401), response["code"])
		assert.Equal(suite.T(), "未登录或token缺失", response["message"])
	})
}

// createValidToken 创建有效的JWT token
func (suite *TokenExpiredIntegrationTestSuite) createValidToken(username, email string, userID int64) string {
	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(31 * 24 * time.Hour).Unix(),
		"iat":      now,
	})

	// 使用与测试环境相同的JWT密钥
	jwtKey := []byte("test-secret-key")
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		suite.T().Fatalf("Failed to create token: %v", err)
	}

	return tokenString
}

// TestTokenExpiredIntegrationTestSuite 运行测试套件
func TestTokenExpiredIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TokenExpiredIntegrationTestSuite))
} 
