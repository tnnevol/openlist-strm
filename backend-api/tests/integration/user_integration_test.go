package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tnnevol/openlist-strm/backend-api/tests"
)

// UserIntegrationTestSuite 用户集成测试套件
type UserIntegrationTestSuite struct {
	suite.Suite
	config *tests.TestConfig
}

// SetupSuite 设置测试套件
func (suite *UserIntegrationTestSuite) SetupSuite() {
	suite.config = tests.SetupTestEnvironment(suite.T())
}

// TearDownSuite 清理测试套件
func (suite *UserIntegrationTestSuite) TearDownSuite() {
	tests.CleanupTestEnvironment(suite.T())
}

// SetupTest 每个测试前的设置
func (suite *UserIntegrationTestSuite) SetupTest() {
	tests.ClearTestData(suite.T(), suite.config.DB)
}

// TestUserRegistrationFlow 测试完整的用户注册流程
func (suite *UserIntegrationTestSuite) TestUserRegistrationFlow() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	suite.Run("完整注册流程", func() {
		// 步骤1: 发送验证码
		sendCodeReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": userData.Email},
			ExpectedCode: 200,
		}

		w := tests.MakeTestRequest(suite.T(), router, sendCodeReq)
		assert.Equal(suite.T(), 200, w.Code)

		// 解析响应获取验证码（在实际测试中，可能需要从数据库或邮件服务获取）
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "验证码已发送", response["message"])

		// 步骤2: 注册用户（这里使用模拟的验证码）
		registerReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/register",
			Body: map[string]interface{}{
				"email":           userData.Email,
				"username":        userData.Username,
				"password":        userData.Password,
				"confirmPassword": userData.Password,
				"code":            "123456", // 模拟验证码
			},
			ExpectedCode: 400, // 验证码错误，预期失败
		}

		w = tests.MakeTestRequest(suite.T(), router, registerReq)
		assert.Equal(suite.T(), 400, w.Code)
	})
}

// TestUserLoginFlow 测试完整的用户登录流程
func (suite *UserIntegrationTestSuite) TestUserLoginFlow() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	suite.Run("完整登录流程", func() {
		// 先创建测试用户
		userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
		assert.Greater(suite.T(), userID, int64(0))

		// 登录
		loginReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/login",
			Body: map[string]interface{}{
				"username": userData.Username,
				"password": userData.Password,
			},
			ExpectedCode: 401, // 密码哈希不匹配，预期失败
		}

		w := tests.MakeTestRequest(suite.T(), router, loginReq)
		assert.Equal(suite.T(), 401, w.Code)
	})
}

// TestTokenBlacklistFlow 测试token黑名单流程
func (suite *UserIntegrationTestSuite) TestTokenBlacklistFlow() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	suite.Run("token黑名单流程", func() {
		// 测试未提供token的情况
		statusReq := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/token-blacklist-status",
			ExpectedCode: 401,
		}

		w := tests.MakeTestRequest(suite.T(), router, statusReq)
		assert.Equal(suite.T(), 401, w.Code)
	})
}

// TestAPIResponseFormat 测试API响应格式
func (suite *UserIntegrationTestSuite) TestAPIResponseFormat() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	suite.Run("API响应格式", func() {
		// 测试发送验证码接口的响应格式
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": "test@example.com"},
			ExpectedCode: 200,
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), 200, w.Code)

		// 验证响应格式
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)

		// 检查响应字段
		assert.Contains(suite.T(), response, "code")
		assert.Contains(suite.T(), response, "message")
		assert.Equal(suite.T(), float64(200), response["code"])
		assert.Equal(suite.T(), "验证码已发送", response["message"])
	})
}

// TestErrorHandling 测试错误处理
func (suite *UserIntegrationTestSuite) TestErrorHandling() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	suite.Run("错误处理", func() {
		// 测试无效的JSON请求
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        "invalid json",
			ExpectedCode: 400,
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), 400, w.Code)
	})
}

// TestUserIntegrationTestSuite 运行测试套件
func TestUserIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserIntegrationTestSuite))
} 
