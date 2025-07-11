# 单元测试框架

本项目使用Go语言内置的`testing`包和`testify`库构建了完整的单元测试框架。

## 测试框架特性

- **分层测试**: 单元测试和集成测试分离
- **测试套件**: 使用`testify/suite`组织测试用例
- **断言库**: 使用`testify/assert`进行断言
- **HTTP测试**: 使用`httptest`测试API接口
- **内存数据库**: 使用SQLite内存数据库进行测试
- **测试工具**: 提供丰富的测试工具函数
- **覆盖率报告**: 支持生成测试覆盖率报告
- **真实用户测试**: 基于YAML配置的真实用户测试系统

## 目录结构

```
tests/
├── README.md                           # 测试文档
├── README_REAL_USER_TESTING.md         # 真实用户测试系统文档
├── README_TOKEN_EXPIRED.md             # Token过期测试文档
├── run_tests.sh                       # 测试运行脚本
├── test_config.go                     # 测试配置
├── http_test_utils.go                 # HTTP测试工具
├── unit/                              # 单元测试
│   ├── user_controller_test.go
│   └── user_service_test.go
├── integration/                       # 集成测试
│   ├── user_integration_test.go
│   ├── real_user_integration_test.go  # 真实用户集成测试
│   └── token_expired_integration_test.go
├── config/                            # 配置管理
│   └── test_config.go                 # YAML配置解析器
└── fixtures/                          # 测试数据
    └── test_users.yml                 # 用户测试配置文件
```

## 快速开始

### 1. 安装依赖

```bash
go get github.com/stretchr/testify
go get gopkg.in/yaml.v3
```

### 2. 运行测试

```bash
# 运行所有测试
./run_tests.sh

# 运行单元测试
./run_tests.sh unit

# 运行集成测试
./run_tests.sh integration

# 运行真实用户测试
./run_tests.sh real-user

# 运行token过期测试
./run_tests.sh token-expired

# 生成覆盖率报告
./run_tests.sh coverage

# 查看帮助
./run_tests.sh help
```

### 3. 直接使用go test

```bash
# 运行所有测试
go test -v ./tests/...

# 运行特定测试
go test -v ./tests/unit/...

# 生成覆盖率报告
go test -v -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html
```

## 测试体系

### 1. 真实用户测试系统

项目引入了基于 YAML 配置的真实用户测试系统，提供更贴近实际使用场景的测试：

- **配置文件驱动**: 使用 `tests/fixtures/test_users.yml` 管理测试用户数据
- **多用户类型**: 支持正常用户、边界测试用户、错误测试用户
- **场景化测试**: 预定义注册、登录、Token管理等测试场景
- **自动化管理**: 自动清理测试数据，确保测试独立性

详细说明请参考：[README_REAL_USER_TESTING.md](README_REAL_USER_TESTING.md)

### 2. 传统测试框架

- **分层测试**: 单元测试和集成测试分离
- **测试套件**: 使用 `testify/suite` 组织测试用例
- **断言库**: 使用 `testify/assert` 进行断言
- **HTTP测试**: 使用 `httptest` 测试API接口
- **内存数据库**: 使用SQLite内存数据库进行测试

## 测试工具

### TestConfig

测试配置结构，包含测试数据库连接：

```go
type TestConfig struct {
    DB *sql.DB
}
```

### 测试环境设置

```go
// 设置测试环境
config := tests.SetupTestEnvironment(t)

// 清理测试环境
tests.CleanupTestEnvironment(t)
```

### 数据库操作

```go
// 插入测试用户
userID := tests.InsertTestUser(t, db, username, email, passwordHash, isActive)

// 获取测试用户
user := tests.GetTestUser(t, db, userID)

// 清理测试数据
tests.ClearTestData(t, db)
```

### HTTP测试

```go
// 创建测试服务器
router := tests.CreateTestServer(t, db)

// 发送测试请求
req := tests.TestRequest{
    Method:      "POST",
    URL:         "/user/send-code",
    Body:        map[string]interface{}{"email": "test@example.com"},
    ExpectedCode: 200,
}
w := tests.MakeTestRequest(t, router, req)

// 断言响应
tests.AssertResponse(t, w, 200, expectedBody)
```

## 编写测试

### 单元测试示例

```go
package unit

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/tnnevol/openlist-strm/backend-api/tests"
)

type UserServiceTestSuite struct {
    suite.Suite
    config *tests.TestConfig
}

func (suite *UserServiceTestSuite) SetupSuite() {
    suite.config = tests.SetupTestEnvironment(suite.T())
}

func (suite *UserServiceTestSuite) SetupTest() {
    tests.ClearTestData(suite.T(), suite.config.DB)
}

func (suite *UserServiceTestSuite) TestCheckEmailExists() {
    userData := tests.GetDefaultTestUser()

    suite.Run("邮箱不存在", func() {
        exists, err := service.CheckEmailExists(suite.config.DB, userData.Email)
        assert.NoError(suite.T(), err)
        assert.False(suite.T(), exists)
    })
}

func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

### 集成测试示例

```go
package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/tnnevol/openlist-strm/backend-api/tests"
)

type UserIntegrationTestSuite struct {
    suite.Suite
    config *tests.TestConfig
}

func (suite *UserIntegrationTestSuite) TestUserRegistrationFlow() {
    router := tests.CreateTestServer(suite.T(), suite.config.DB)

    suite.Run("完整注册流程", func() {
        // 发送验证码
        req := tests.TestRequest{
            Method:      "POST",
            URL:         "/user/send-code",
            Body:        map[string]interface{}{"email": "test@example.com"},
            ExpectedCode: 200,
        }
        w := tests.MakeTestRequest(suite.T(), router, req)
        assert.Equal(suite.T(), 200, w.Code)
    })
}
```

### 真实用户测试示例

```go
package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/tnnevol/openlist-strm/backend-api/tests"
    "github.com/tnnevol/openlist-strm/backend-api/tests/config"
)

type RealUserIntegrationTestSuite struct {
    suite.Suite
    config *tests.TestConfig
    userConfig *config.TestConfig
}

func (suite *RealUserIntegrationTestSuite) SetupSuite() {
    suite.config = tests.SetupTestEnvironment(suite.T())

    // 加载真实用户配置
    userConfig, err := config.LoadTestConfig("")
    assert.NoError(suite.T(), err)
    suite.userConfig = userConfig
}

func (suite *RealUserIntegrationTestSuite) TestRealUserRegistration() {
    router := tests.CreateTestServer(suite.T(), suite.config.DB)

    // 获取正常用户组中的第一个用户
    normalUsers := suite.userConfig.GetNormalUsers()
    testUser := normalUsers[0]

    suite.Run("真实用户注册流程", func() {
        // 发送验证码
        req := tests.TestRequest{
            Method:      "POST",
            URL:         "/user/send-code",
            Body:        map[string]interface{}{"email": testUser.Email},
            ExpectedCode: 200,
        }
        w := tests.MakeTestRequest(suite.T(), router, req)
        assert.Equal(suite.T(), 200, w.Code)
    })
}
```

## Token过期测试

### 概述

本系统对token过期进行了特殊处理，当token过期时不会直接拒绝请求，而是返回特定的过期状态码，让前端能够区分不同的认证失败原因。

### 错误码定义

- `401`: 未授权（未登录或token缺失）
- `40101`: Token过期专用错误码

### 处理逻辑

在`AuthMiddleware`中，当检测到token过期时：

```go
// 检查token是否已过期
if exp, ok := claims["exp"].(float64); ok {
    if time.Now().Unix() > int64(exp) {
        println("[AuthMiddleware] token已过期，返回过期状态")
        TokenExpired(c, "token已过期")
        c.Abort()
        return
    }
}
```

### 响应格式

当token过期时，系统会返回以下格式的响应：

```json
{
  "code": 40101,
  "message": "token已过期"
}
```

### 前端处理建议

前端可以根据不同的错误码进行不同的处理：

```javascript
// 示例前端处理逻辑
if (response.code === 401) {
  // 未登录，跳转到登录页
  redirectToLogin();
} else if (response.code === 40101) {
  // token过期，静默刷新token或跳转登录页
  handleTokenExpired();
} else {
  // 其他错误
  handleOtherError(response);
}
```

### 测试覆盖范围

token过期测试包含以下测试场景：

1. **token过期处理流程**: 测试完整的token过期处理机制
2. **无token访问**: 测试未提供token时的处理
3. **无效token**: 测试无效token的处理
4. **token黑名单**: 测试过期token在黑名单中的处理
5. **生成过期token接口**: 测试生成过期token的API接口
6. **无认证访问**: 测试无认证访问生成过期token接口

### 测试用例说明

#### TestTokenExpiredHandling

- 创建测试用户并获取有效token
- 使用有效token访问需要认证的接口
- 生成过期token（1秒后过期）
- 等待token过期
- 使用过期token访问接口，验证返回40101错误码

#### TestNoTokenHandling

- 测试无token访问需要认证的接口
- 验证返回401错误码（不是40101）

#### TestInvalidTokenHandling

- 测试无效token访问需要认证的接口
- 验证返回401错误码（不是40101）

#### TestTokenBlacklistWithExpiredToken

- 创建测试用户并获取有效token
- 先登出，将token加入黑名单
- 使用已加入黑名单的token访问接口
- 验证返回401错误码（不是40101）

#### TestGenerateExpiredTokenEndpoint

- 测试生成过期token接口的响应格式
- 验证返回的token、expires_in、message字段

#### TestGenerateExpiredTokenWithoutAuth

- 测试无认证访问生成过期token接口
- 验证返回401错误码

### 运行token过期测试

```bash
# 运行所有token过期测试
./run_tests.sh token-expired

# 直接使用go test运行
go test -v -run TestTokenExpiredIntegrationTestSuite ./integration/...
```

### 测试接口

系统提供了测试接口来生成过期token：

```
POST /user/generate-expired-token
Authorization: Bearer <valid_token>
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "expired_token_string",
    "expires_in": 1,
    "message": "过期token已生成，1秒后过期"
  }
}
```

### 测试数据

测试使用以下默认用户数据：

```go
TestUserData{
    Username:     "testuser",
    Email:        "test@example.com",
    Password:     "TestPass123",
    IsActive:     true,
    PasswordHash: "$2a$10$test.hash.for.testing",
}
```

### 优势

1. **精确的错误分类**: 区分未登录和token过期两种情况
2. **更好的用户体验**: 前端可以根据不同情况提供不同的处理方式
3. **便于调试**: 明确的错误码便于问题定位
4. **符合RESTful规范**: 使用标准的HTTP状态码

### 注意事项

1. 过期token仍然会被加入黑名单（如果用户主动登出）
2. 系统会同时检查token的`exp`字段和`token_invalid_before`字段
3. 建议前端在收到40101错误码时，引导用户重新登录

## 测试最佳实践

### 1. 测试命名

- 测试函数名以`Test`开头
- 测试用例使用描述性的名称
- 使用`suite.Run`组织相关测试

### 2. 测试数据

- 使用`SetupTest`在每个测试前清理数据
- 使用`fixtures`目录存放测试数据
- 避免测试间的数据依赖

### 3. 断言

- 使用`assert`包进行断言
- 提供清晰的错误信息
- 测试边界条件和异常情况

### 4. 测试覆盖

- 测试正常流程
- 测试异常流程
- 测试边界条件
- 测试错误处理

### 5. 性能考虑

- 使用内存数据库提高测试速度
- 避免不必要的网络请求
- 合理使用Mock对象

## 测试覆盖率

项目目标测试覆盖率：**80%以上**

查看覆盖率报告：

```bash
./tests/run_tests.sh coverage
```

生成的HTML报告包含：

- 总体覆盖率
- 文件级别覆盖率
- 行级别覆盖率详情

## 持续集成

测试框架支持CI/CD集成：

```yaml
# GitHub Actions示例
- name: Run Tests
  run: |
    cd backend-api
    ./tests/run_tests.sh all

- name: Generate Coverage Report
  run: |
    cd backend-api
    ./tests/run_tests.sh coverage
```

## 故障排除

### 常见问题

1. **导入错误**: 确保`go.mod`文件正确配置
2. **数据库连接失败**: 检查SQLite驱动是否正确安装
3. **测试超时**: 增加测试超时时间或优化测试逻辑

### 调试技巧

1. 使用`-v`参数查看详细输出
2. 使用`-run`参数运行特定测试
3. 在测试中添加日志输出
4. 使用断点调试

## 扩展测试框架

### 添加新的测试工具

在`tests/`目录下添加新的工具函数：

```go
// tests/new_utils.go
package tests

func NewTestUtility(t *testing.T) {
    // 新的测试工具函数
}
```

### 添加新的测试类型

1. 创建新的测试目录
2. 编写测试文件
3. 更新`run_tests.sh`脚本
4. 更新文档

## 贡献指南

1. 为新功能编写测试
2. 确保测试覆盖率不降低
3. 遵循测试命名规范
4. 更新相关文档

## 相关文档

- [测试体系总览](README_TESTING_SYSTEM.md) - 完整测试体系的架构和特性总览
- [真实用户测试系统](README_REAL_USER_TESTING.md) - 基于YAML配置的真实用户测试系统
- [Token过期测试](README_TOKEN_EXPIRED.md) - Token过期处理测试详细说明
- [API响应格式](../API_RESPONSE_FORMAT.md) - API响应格式规范
