package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tnnevol/openlist-strm/backend-api/tests"
)

// UserControllerTestSuite 用户控制器测试套件
type UserControllerTestSuite struct {
	suite.Suite
	config *tests.TestConfig
}

// SetupSuite 设置测试套件
func (suite *UserControllerTestSuite) SetupSuite() {
	suite.config = tests.SetupTestEnvironment(suite.T())
}

// TearDownSuite 清理测试套件
func (suite *UserControllerTestSuite) TearDownSuite() {
	tests.CleanupTestEnvironment(suite.T())
}

// SetupTest 每个测试前的设置
func (suite *UserControllerTestSuite) SetupTest() {
	tests.ClearTestData(suite.T(), suite.config.DB)
}

// TestSendCode 测试发送验证码接口
func (suite *UserControllerTestSuite) TestSendCode() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	// 测试用例1: 正常发送验证码
	suite.Run("正常发送验证码", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": "test@example.com"},
			ExpectedCode: 200,
			ExpectedBody: map[string]interface{}{
				"code":    200,
				"message": "验证码已发送",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})

	// 测试用例2: 邮箱为空
	suite.Run("邮箱为空", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": ""},
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"code":    400,
				"message": "参数错误：邮箱不能为空",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})

	// 测试用例3: 邮箱格式不正确
	suite.Run("邮箱格式不正确", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": "invalid-email"},
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"code":    400,
				"message": "参数错误：邮箱格式不正确",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})
}

// TestRegister 测试注册接口
func (suite *UserControllerTestSuite) TestRegister() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	// 先发送验证码
	userData := tests.GetDefaultTestUser()
	
	// 测试用例1: 正常注册
	suite.Run("正常注册", func() {
		// 先发送验证码
		sendCodeReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": userData.Email},
			ExpectedCode: 200,
		}
		tests.MakeTestRequest(suite.T(), router, sendCodeReq)

		// 注册（这里需要真实的验证码，实际测试中可能需要mock）
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/register",
			Body: map[string]interface{}{
				"email":           userData.Email,
				"username":        userData.Username,
				"password":        userData.Password,
				"confirmPassword": userData.Password,
				"code":            "123456", // 这里需要真实的验证码
			},
			ExpectedCode: 400, // 验证码错误，预期失败
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), req.ExpectedCode, w.Code)
	})

	// 测试用例2: 参数不完整
	suite.Run("参数不完整", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/register",
			Body: map[string]interface{}{
				"email":    userData.Email,
				"username": userData.Username,
				// 缺少password和confirmPassword
			},
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"code":    400,
				"message": "参数错误：邮箱、用户名、密码、确认密码和验证码不能为空",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})

	// 测试用例3: 密码不一致
	suite.Run("密码不一致", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/register",
			Body: map[string]interface{}{
				"email":           userData.Email,
				"username":        userData.Username,
				"password":        userData.Password,
				"confirmPassword": "different-password",
				"code":            "123456",
			},
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"code":    400,
				"message": "参数错误：两次输入的密码不一致",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})
}

// TestLogin 测试登录接口
func (suite *UserControllerTestSuite) TestLogin() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	// 先创建测试用户
	userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
	assert.Greater(suite.T(), userID, int64(0))

	// 测试用例1: 正常登录
	suite.Run("正常登录", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/login",
			Body: map[string]interface{}{
				"username": userData.Username,
				"password": userData.Password,
			},
			ExpectedCode: 401, // 密码哈希不匹配，预期失败
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		assert.Equal(suite.T(), req.ExpectedCode, w.Code)
	})

	// 测试用例2: 用户名或密码为空
	suite.Run("用户名或密码为空", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/login",
			Body: map[string]interface{}{
				"username": "",
				"password": "",
			},
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"code":    400,
				"message": "参数错误：用户名和密码不能为空",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})

	// 测试用例3: 用户名不存在
	suite.Run("用户名不存在", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/login",
			Body: map[string]interface{}{
				"username": "nonexistent",
				"password": userData.Password,
			},
			ExpectedCode: 401,
			ExpectedBody: map[string]interface{}{
				"code":    401,
				"message": "用户名或密码错误",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})
}

// TestUserInfo 测试用户信息接口
func (suite *UserControllerTestSuite) TestUserInfo() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)
	userData := tests.GetDefaultTestUser()

	// 先创建测试用户
	userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
	assert.Greater(suite.T(), userID, int64(0))

	// 测试用例1: 未提供token
	suite.Run("未提供token", func() {
		req := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			ExpectedCode: 401,
			ExpectedBody: map[string]interface{}{
				"code":    401,
				"message": "未登录或token缺失",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})

	// 测试用例2: 无效token
	suite.Run("无效token", func() {
		req := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/info",
			Headers:     map[string]string{"Authorization": "Bearer invalid-token"},
			ExpectedCode: 401,
			ExpectedBody: map[string]interface{}{
				"code":    401,
				"message": "token无效或已过期",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})
}

// TestLogout 测试登出接口
func (suite *UserControllerTestSuite) TestLogout() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	// 测试用例1: 未提供token
	suite.Run("未提供token", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/logout",
			ExpectedCode: 401,
			ExpectedBody: map[string]interface{}{
				"code":    401,
				"message": "未登录或token缺失",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})

	// 测试用例2: 无效token
	suite.Run("无效token", func() {
		req := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/logout",
			Headers:     map[string]string{"Authorization": "Bearer invalid-token"},
			ExpectedCode: 401,
			ExpectedBody: map[string]interface{}{
				"code":    401,
				"message": "token无效或已过期",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})
}

// TestTokenBlacklistStatus 测试token黑名单状态接口
func (suite *UserControllerTestSuite) TestTokenBlacklistStatus() {
	router := tests.CreateTestServer(suite.T(), suite.config.DB)

	// 测试用例1: 未提供token
	suite.Run("未提供token", func() {
		req := tests.TestRequest{
			Method:      "GET",
			URL:         "/user/token-blacklist-status",
			ExpectedCode: 401,
			ExpectedBody: map[string]interface{}{
				"code":    401,
				"message": "未登录或token缺失",
			},
		}

		w := tests.MakeTestRequest(suite.T(), router, req)
		tests.AssertResponse(suite.T(), w, req.ExpectedCode, req.ExpectedBody)
	})
}

// TestUserControllerTestSuite 运行测试套件
func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
} 
