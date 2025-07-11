package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tnnevol/openlist-strm/backend-api/tests"
	"github.com/tnnevol/openlist-strm/backend-api/tests/config"
)

// RealUserIntegrationTestSuite 真实用户集成测试套件
type RealUserIntegrationTestSuite struct {
	suite.Suite
	testConfig *tests.TestConfig
	userConfig *config.TestConfigFile
}

// SetupSuite 设置测试套件
func (suite *RealUserIntegrationTestSuite) SetupSuite() {
	suite.testConfig = tests.SetupTestEnvironment(suite.T())
	
	// 加载用户测试配置
	userConfig, err := config.LoadTestConfig("")
	if err != nil {
		suite.T().Fatalf("加载用户测试配置失败: %v", err)
	}
	suite.userConfig = userConfig
	
	// 验证配置
	err = suite.userConfig.ValidateConfig()
	if err != nil {
		suite.T().Fatalf("验证用户测试配置失败: %v", err)
	}
}

// TearDownSuite 清理测试套件
func (suite *RealUserIntegrationTestSuite) TearDownSuite() {
	tests.CleanupTestEnvironment(suite.T())
}

// SetupTest 每个测试前的设置
func (suite *RealUserIntegrationTestSuite) SetupTest() {
	tests.ClearTestData(suite.T(), suite.testConfig.DB)
}

// TestNormalUserRegistration 测试正常用户注册流程
func (suite *RealUserIntegrationTestSuite) TestNormalUserRegistration() {
	router := tests.CreateTestServer(suite.T(), suite.testConfig.DB)
	
	// 获取正常用户组中的第一个用户
	normalUsers := suite.userConfig.GetActiveUsers()
	if len(normalUsers) == 0 {
		suite.T().Skip("没有可用的正常用户进行测试")
	}
	
	testUser := normalUsers[0]
	
	suite.Run(fmt.Sprintf("正常用户注册: %s", testUser.Name), func() {
		// 步骤1: 发送验证码
		sendCodeReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/send-code",
			Body:        map[string]interface{}{"email": testUser.Email},
			ExpectedCode: 200, // HTTP状态码始终为200
		}

		w := tests.MakeTestRequest(suite.T(), router, sendCodeReq)
		assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

		// 验证响应体中的业务状态码
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(200), response["code"]) // 业务状态码
		assert.Equal(suite.T(), "验证码已发送", response["message"])

		// 步骤2: 测试注册接口（由于验证码是随机生成的，这里只测试接口响应格式）
		// 注意：实际注册需要正确的验证码，这里只测试参数验证
		registerReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/register",
			Body: map[string]interface{}{
				"email":           testUser.Email,
				"username":        testUser.Username,
				"password":        testUser.Password,
				"confirmPassword": testUser.Password,
				"code":            "000000", // 使用错误的验证码测试
			},
			ExpectedCode: 200, // HTTP状态码始终为200
		}

		w = tests.MakeTestRequest(suite.T(), router, registerReq)
		assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

		// 验证注册失败（验证码错误）
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(400), response["code"]) // 业务状态码
		assert.Contains(suite.T(), response["message"], "验证码")
	})
}

// TestNormalUserLogin 测试正常用户登录流程
func (suite *RealUserIntegrationTestSuite) TestNormalUserLogin() {
	router := tests.CreateTestServer(suite.T(), suite.testConfig.DB)
	
	// 获取正常用户组中的第一个用户
	normalUsers := suite.userConfig.GetActiveUsers()
	if len(normalUsers) == 0 {
		suite.T().Skip("没有可用的正常用户进行测试")
	}
	
	testUser := normalUsers[0]
	
	suite.Run(fmt.Sprintf("正常用户登录: %s", testUser.Name), func() {
		// 先创建测试用户
		userID := tests.InsertTestUserFromConfig(suite.T(), suite.testConfig.DB, testUser)
		assert.Greater(suite.T(), userID, int64(0))

		// 登录
		loginReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/login",
			Body: map[string]interface{}{
				"username": testUser.Username,
				"password": testUser.Password,
			},
			ExpectedCode: 200, // HTTP状态码始终为200
		}

		w := tests.MakeTestRequest(suite.T(), router, loginReq)
		assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

		// 验证登录成功
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(200), response["code"]) // 业务状态码
		
		// 检查accessToken在data字段中
		if response["data"] != nil {
			data := response["data"].(map[string]interface{})
			assert.Contains(suite.T(), data, "accessToken")
		} else {
			// 如果没有data字段，检查根级别是否有accessToken
			assert.Contains(suite.T(), response, "accessToken")
		}
	})
}

// TestEdgeCaseUsers 测试边界情况用户
func (suite *RealUserIntegrationTestSuite) TestEdgeCaseUsers() {
	router := tests.CreateTestServer(suite.T(), suite.testConfig.DB)
	
	edgeCaseUsers := suite.userConfig.UserGroups.EdgeCaseUsers
	if len(edgeCaseUsers) == 0 {
		suite.T().Skip("没有可用的边界测试用户")
	}
	
	for _, testUser := range edgeCaseUsers {
		suite.Run(fmt.Sprintf("边界测试用户: %s", testUser.Name), func() {
			// 测试发送验证码
			sendCodeReq := tests.TestRequest{
				Method:      "POST",
				URL:         "/user/send-code",
				Body:        map[string]interface{}{"email": testUser.Email},
				ExpectedCode: 200, // HTTP状态码始终为200
			}

			w := tests.MakeTestRequest(suite.T(), router, sendCodeReq)
			assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

			// 验证响应体中的业务状态码
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), float64(200), response["code"]) // 业务状态码
		})
	}
}

// TestErrorCaseUsers 测试错误情况用户
func (suite *RealUserIntegrationTestSuite) TestErrorCaseUsers() {
	router := tests.CreateTestServer(suite.T(), suite.testConfig.DB)
	
	errorTestUsers := suite.userConfig.UserGroups.ErrorTestUsers
	if len(errorTestUsers) == 0 {
		suite.T().Skip("没有可用的错误测试用户")
	}
	
	for _, testUser := range errorTestUsers {
		suite.Run(fmt.Sprintf("错误测试用户: %s", testUser.Name), func() {
			// 测试发送验证码
			// 注意：接口只验证邮箱格式，空邮箱会返回400业务码，无效邮箱格式会返回400业务码
			expectedBusinessCode := float64(200)
			if testUser.Email == "" {
				expectedBusinessCode = 400 // 空邮箱
			} else if testUser.Email == "invalid-email" {
				expectedBusinessCode = 400 // 无效邮箱格式
			}
			
			sendCodeReq := tests.TestRequest{
				Method:      "POST",
				URL:         "/user/send-code",
				Body:        map[string]interface{}{"email": testUser.Email},
				ExpectedCode: 200, // HTTP状态码始终为200
			}

			w := tests.MakeTestRequest(suite.T(), router, sendCodeReq)
			assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

			// 验证响应体中的业务状态码
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), expectedBusinessCode, response["code"]) // 业务状态码
		})
	}
}

// TestUserScenarios 测试预定义的测试场景
func (suite *RealUserIntegrationTestSuite) TestUserScenarios() {
	router := tests.CreateTestServer(suite.T(), suite.testConfig.DB)
	
	// 测试注册场景
	for _, scenario := range suite.userConfig.TestScenarios.Registration {
		suite.Run(fmt.Sprintf("注册场景: %s", scenario.Name), func() {
			suite.runTestScenario(router, scenario)
		})
	}
	
	// 测试登录场景
	for _, scenario := range suite.userConfig.TestScenarios.Login {
		suite.Run(fmt.Sprintf("登录场景: %s", scenario.Name), func() {
			suite.runTestScenario(router, scenario)
		})
	}
}

// runTestScenario 运行测试场景
func (suite *RealUserIntegrationTestSuite) runTestScenario(router *gin.Engine, scenario config.TestScenario) {
	for _, step := range scenario.Steps {
		suite.Run(step.Step, func() {
			// 如果是登录接口，自动插入测试用户
			if step.Endpoint == "/user/login" && step.Method == "POST" {
				username, _ := step.Body["username"].(string)
				password, _ := step.Body["password"].(string)
				if username != "" && password != "" {
					// 检查用户是否已存在
					db := suite.testConfig.DB
					user := tests.GetTestUserByUsername(suite.T(), db, username)
					if user == nil {
						// 自动插入测试用户，email用username拼接@autotest.com
						tests.InsertTestUser(suite.T(), db, username, username+"@autotest.com", tests.HashPasswordForTest(password), true)
					}
				}
			}

			// 构建请求
			req := tests.TestRequest{
				Method:      step.Method,
				URL:         step.Endpoint,
				Body:        step.Body,
				Headers:     step.Headers,
				ExpectedCode: 200, // HTTP状态码始终为200
			}

			w := tests.MakeTestRequest(suite.T(), router, req)
			assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

			// 验证响应体中的业务状态码
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			// 验证期望的业务状态码
			if step.Expected.Code != 0 {
				assert.Equal(suite.T(), float64(step.Expected.Code), response["code"])
			}

			// 验证期望的消息
			if step.Expected.Message != "" {
				assert.Equal(suite.T(), step.Expected.Message, response["message"])
			}

			// 验证期望包含的内容
			if step.Expected.Contains != "" {
				// 先检查data字段中是否包含
				if response["data"] != nil {
					data := response["data"].(map[string]interface{})
					dataStr := fmt.Sprintf("%v", data)
					if strings.Contains(dataStr, step.Expected.Contains) {
						return // 在data字段中找到，验证通过
					}
				}
				
				// 如果data字段中没有，检查根级别
				responseStr := fmt.Sprintf("%v", response)
				assert.Contains(suite.T(), responseStr, step.Expected.Contains)
			}
		})
	}
}

// TestUserTokenExpired 测试用户Token过期场景
func (suite *RealUserIntegrationTestSuite) TestUserTokenExpired() {
	router := tests.CreateTestServer(suite.T(), suite.testConfig.DB)
	
	// 获取正常用户
	normalUsers := suite.userConfig.GetActiveUsers()
	if len(normalUsers) == 0 {
		suite.T().Skip("没有可用的正常用户进行测试")
	}
	
	testUser := normalUsers[0]
	
	suite.Run(fmt.Sprintf("用户Token过期测试: %s", testUser.Name), func() {
		// 先创建用户并登录获取token
		userID := tests.InsertTestUserFromConfig(suite.T(), suite.testConfig.DB, testUser)
		assert.Greater(suite.T(), userID, int64(0))

		loginReq := tests.TestRequest{
			Method:      "POST",
			URL:         "/user/login",
			Body: map[string]interface{}{
				"username": testUser.Username,
				"password": testUser.Password,
			},
			ExpectedCode: 200, // HTTP状态码始终为200
		}

		w := tests.MakeTestRequest(suite.T(), router, loginReq)
		assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

		var loginResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
		assert.NoError(suite.T(), err)
		
		// 检查是否有token
		if loginResponse["data"] != nil {
			data := loginResponse["data"].(map[string]interface{})
			if data["accessToken"] != nil {
				token := data["accessToken"].(string)
				assert.NotEmpty(suite.T(), token)

				// 生成过期token
				expiredTokenReq := tests.TestRequest{
					Method:      "POST",
					URL:         "/user/generate-expired-token",
					Headers:     map[string]string{"Authorization": "Bearer " + token},
					ExpectedCode: 200, // HTTP状态码始终为200
				}

				w = tests.MakeTestRequest(suite.T(), router, expiredTokenReq)
				assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

				var expiredResponse map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &expiredResponse)
				assert.NoError(suite.T(), err)
				
				if expiredResponse["data"] != nil {
					expiredData := expiredResponse["data"].(map[string]interface{})
					if expiredData["token"] != nil {
						expiredToken := expiredData["token"].(string)
						assert.NotEmpty(suite.T(), expiredToken)

						// 等待token过期
						time.Sleep(2 * time.Second)

						// 使用过期token访问接口
						infoReq := tests.TestRequest{
							Method:      "GET",
							URL:         "/user/info",
							Headers:     map[string]string{"Authorization": "Bearer " + expiredToken},
							ExpectedCode: 200, // HTTP状态码始终为200
						}

						w = tests.MakeTestRequest(suite.T(), router, infoReq)
						assert.Equal(suite.T(), 200, w.Code) // HTTP状态码检查

						var infoResponse map[string]interface{}
						err = json.Unmarshal(w.Body.Bytes(), &infoResponse)
						assert.NoError(suite.T(), err)
						// 检查业务码是否为40101（token过期）
						assert.Equal(suite.T(), float64(40101), infoResponse["code"])
						assert.Equal(suite.T(), "token已过期", infoResponse["message"])
					}
				}
			}
		}
	})
}

// TestRealUserIntegrationTestSuite 运行测试套件
func TestRealUserIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RealUserIntegrationTestSuite))
} 
