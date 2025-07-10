package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tnnevol/openlist-strm/backend-api/internal/service"
	"github.com/tnnevol/openlist-strm/backend-api/tests"
)

// UserServiceTestSuite 用户服务测试套件
type UserServiceTestSuite struct {
	suite.Suite
	config *tests.TestConfig
}

// SetupSuite 设置测试套件
func (suite *UserServiceTestSuite) SetupSuite() {
	suite.config = tests.SetupTestEnvironment(suite.T())
}

// TearDownSuite 清理测试套件
func (suite *UserServiceTestSuite) TearDownSuite() {
	tests.CleanupTestEnvironment(suite.T())
}

// SetupTest 每个测试前的设置
func (suite *UserServiceTestSuite) SetupTest() {
	tests.ClearTestData(suite.T(), suite.config.DB)
}

// TestCheckEmailExists 测试邮箱存在性检查
func (suite *UserServiceTestSuite) TestCheckEmailExists() {
	userData := tests.GetDefaultTestUser()

	// 测试用例1: 邮箱不存在
	suite.Run("邮箱不存在", func() {
		exists, err := service.CheckEmailExists(suite.config.DB, userData.Email)
		assert.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
	})

	// 测试用例2: 邮箱存在
	suite.Run("邮箱存在", func() {
		// 先插入用户
		userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
		assert.Greater(suite.T(), userID, int64(0))

		// 检查邮箱是否存在
		exists, err := service.CheckEmailExists(suite.config.DB, userData.Email)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
	})
}

// TestCheckUsernameExists 测试用户名存在性检查
func (suite *UserServiceTestSuite) TestCheckUsernameExists() {
	userData := tests.GetDefaultTestUser()

	// 测试用例1: 用户名不存在
	suite.Run("用户名不存在", func() {
		exists, err := service.CheckUsernameExists(suite.config.DB, userData.Username)
		assert.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
	})

	// 测试用例2: 用户名存在
	suite.Run("用户名存在", func() {
		// 先插入用户
		userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
		assert.Greater(suite.T(), userID, int64(0))

		// 检查用户名是否存在
		exists, err := service.CheckUsernameExists(suite.config.DB, userData.Username)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
	})
}

// TestSendCode 测试发送验证码
func (suite *UserServiceTestSuite) TestSendCode() {
	userData := tests.GetDefaultTestUser()

	// 测试用例1: 发送验证码给新邮箱
	suite.Run("发送验证码给新邮箱", func() {
		code, err := service.SendCode(suite.config.DB, userData.Email)
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), code)
		assert.Len(suite.T(), code, 6) // 验证码应该是6位
	})

	// 测试用例2: 发送验证码给已存在但未激活的用户
	suite.Run("发送验证码给已存在但未激活的用户", func() {
		// 先发送一次验证码
		code1, err := service.SendCode(suite.config.DB, userData.Email)
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), code1)

		// 再次发送验证码
		code2, err := service.SendCode(suite.config.DB, userData.Email)
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), code2)
		// 两次验证码可能不同（因为会更新）
	})
}

// TestLoginUser 测试用户登录
func (suite *UserServiceTestSuite) TestLoginUser() {
	userData := tests.GetDefaultTestUser()

	// 测试用例1: 用户名不存在
	suite.Run("用户名不存在", func() {
		token, err := service.LoginUser(suite.config.DB, "nonexistent", userData.Password)
		assert.Error(suite.T(), err)
		assert.Empty(suite.T(), token)
	})

	// 测试用例2: 用户未激活
	suite.Run("用户未激活", func() {
		// 创建未激活用户
		userData.IsActive = false
		userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
		assert.Greater(suite.T(), userID, int64(0))

		token, err := service.LoginUser(suite.config.DB, userData.Username, userData.Password)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "not_activated", token)
	})
}

// TestGetUserBaseInfo 测试获取用户基础信息
func (suite *UserServiceTestSuite) TestGetUserBaseInfo() {
	userData := tests.GetDefaultTestUser()

	// 测试用例1: 用户名不存在
	suite.Run("用户名不存在", func() {
		userInfo, err := service.GetUserBaseInfo(suite.config.DB, "nonexistent")
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), userInfo)
	})

	// 测试用例2: 用户名存在
	suite.Run("用户名存在", func() {
		// 先创建用户
		userID := tests.SetupTestUser(suite.T(), suite.config.DB, userData)
		assert.Greater(suite.T(), userID, int64(0))

		// 获取用户信息
		userInfo, err := service.GetUserBaseInfo(suite.config.DB, userData.Username)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), userInfo)
		assert.Equal(suite.T(), userData.Username, userInfo["username"])
		assert.Equal(suite.T(), userData.Email, userInfo["email"])
	})
}

// TestTokenBlacklist 测试token黑名单功能
func (suite *UserServiceTestSuite) TestTokenBlacklist() {
	// 测试用例1: 添加token到黑名单
	suite.Run("添加token到黑名单", func() {
		blacklist := service.GetTokenBlacklist()
		assert.NotNil(suite.T(), blacklist)

		// 检查初始大小
		initialSize := blacklist.GetBlacklistSize()
		assert.Equal(suite.T(), 0, initialSize)

		// 添加token到黑名单
		testToken := "test-token-123"
		blacklist.AddToBlacklist(testToken, time.Now().Add(1*time.Hour))

		// 检查是否在黑名单中
		isBlacklisted := blacklist.IsBlacklisted(testToken)
		assert.True(suite.T(), isBlacklisted)

		// 检查黑名单大小
		newSize := blacklist.GetBlacklistSize()
		assert.Equal(suite.T(), 1, newSize)
	})

	// 测试用例2: 检查不存在的token
	suite.Run("检查不存在的token", func() {
		blacklist := service.GetTokenBlacklist()
		isBlacklisted := blacklist.IsBlacklisted("non-existent-token")
		assert.False(suite.T(), isBlacklisted)
	})
}

// TestUserServiceTestSuite 运行测试套件
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
} 
