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

## 目录结构

```
tests/
├── README.md                 # 测试文档
├── run_tests.sh             # 测试运行脚本
├── test_config.go           # 测试配置
├── http_test_utils.go       # HTTP测试工具
├── unit/                    # 单元测试
│   ├── user_controller_test.go
│   └── user_service_test.go
├── integration/             # 集成测试
│   └── user_integration_test.go
└── fixtures/                # 测试数据
```

## 快速开始

### 1. 安装依赖

```bash
go get github.com/stretchr/testify
```

### 2. 运行测试

```bash
# 运行所有测试
./tests/run_tests.sh

# 运行单元测试
./tests/run_tests.sh unit

# 运行集成测试
./tests/run_tests.sh integration

# 生成覆盖率报告
./tests/run_tests.sh coverage

# 查看帮助
./tests/run_tests.sh help
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
