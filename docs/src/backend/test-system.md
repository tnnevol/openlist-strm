# 后端测试体系总览

> 本文档整合自 backend-api/tests/README.md、README_TESTING_SYSTEM.md、README_REAL_USER_TESTING.md、README_TOKEN_EXPIRED.md，系统梳理了 OpenList Stream 后端的测试体系、用法与最佳实践。

---

## 1. 测试体系架构

```
测试体系
├── 传统测试框架
│   ├── 单元测试 (tests/unit/)
│   ├── 集成测试 (tests/integration/)
│   └── 测试工具 (tests/http_test_utils.go)
├── 真实用户测试系统
│   ├── YAML配置驱动 (tests/fixtures/test_users.yml)
│   ├── 配置解析器 (tests/config/test_config.go)
│   └── 真实用户集成测试 (tests/integration/real_user_integration_test.go)
└── 测试管理
    ├── 统一运行脚本 (tests/run_tests.sh)
    ├── 覆盖率报告
    └── CI/CD集成
```

## 2. 测试类型与特性

### 2.1 单元测试

- 测试单个函数和模块的逻辑
- 使用 Go 内置 testing 包 + testify/suite/assert
- 支持内存数据库、HTTP接口、丰富工具函数

### 2.2 集成测试

- 测试模块间交互和 API 接口
- 结合 httptest、真实路由、数据库

### 2.3 真实用户测试

- 基于 YAML 配置，驱动多用户多场景
- 支持正常/边界/错误用户，注册、登录、Token、验证码等全流程
- 自动清理测试数据，确保独立性

### 2.4 Token 过期测试

- 专用错误码 40101 区分 token 过期
- 完整覆盖 token 过期、无 token、无效 token、黑名单等场景

---

## 3. 目录结构

```
tests/
├── README.md                           # 测试文档
├── README_TESTING_SYSTEM.md            # 测试体系总览
├── README_REAL_USER_TESTING.md         # 真实用户测试系统文档
├── README_TOKEN_EXPIRED.md             # Token过期测试文档
├── run_tests.sh                        # 测试运行脚本
├── test_config.go                      # 测试配置
├── http_test_utils.go                  # HTTP测试工具
├── unit/                               # 单元测试
│   ├── user_controller_test.go
│   └── user_service_test.go
├── integration/                        # 集成测试
│   ├── user_integration_test.go
│   ├── real_user_integration_test.go
│   └── token_expired_integration_test.go
├── config/                             # 配置管理
│   └── test_config.go                  # YAML配置解析器
└── fixtures/                           # 测试数据
    └── test_users.yml                  # 用户测试配置文件
```

---

## 4. 运行方式

### 4.1 统一脚本

```bash
# 运行所有测试
./tests/run_tests.sh

# 运行单元测试
./tests/run_tests.sh unit

# 运行集成测试
./tests/run_tests.sh integration

# 运行真实用户测试
./tests/run_tests.sh real-user

# 运行token过期测试
./tests/run_tests.sh token-expired

# 生成覆盖率报告
./tests/run_tests.sh coverage

# 清理测试数据
./tests/run_tests.sh cleanup

# 查看帮助
./tests/run_tests.sh help
```

### 4.2 直接 go test

```bash
# 运行所有测试
go test -v ./tests/...
# 运行特定测试
go test -v ./tests/unit/...
go test -v ./tests/integration/...
# 生成覆盖率报告
go test -v -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html
```

---

## 5. 真实用户测试系统

- 配置驱动，支持多用户多场景
- 详见 [test_users.yml](https://github.com/tnnevol/openlist-strm/blob/main/backend-api/tests/fixtures/test_users.yml)
- 支持注册、登录、Token、验证码、边界与异常场景
- 自动清理、批量测试、详细日志

---

## 6. Token过期测试

- 专用错误码40101，区分token过期与未登录
- 详见 [token_expired_integration_test.go](https://github.com/tnnevol/openlist-strm/blob/main/backend-api/tests/integration/token_expired_integration_test.go)
- 统一通过 run_tests.sh 管理
- 响应示例：

```json
{
  "code": 40101,
  "message": "token已过期"
}
```

---

## 7. 测试工具与最佳实践

- 丰富的测试工具函数（见 http_test_utils.go、test_config.go）
- 使用内存数据库，自动清理，避免数据污染
- 测试用例命名规范，覆盖正常/异常/边界
- 断言清晰，日志详细
- 推荐80%以上覆盖率

---

## 8. 持续集成与覆盖率

- 支持 GitHub Actions/CI
- 自动生成覆盖率报告
- 目标：核心业务90%+，API接口95%+

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

---

## 9. 故障排除与扩展

- 常见问题、调试技巧、扩展方法详见原文档
- 推荐定期运行、关注覆盖率、及时修复失败测试

---

## 10. 相关文档

- [API响应格式](../API_RESPONSE_FORMAT.md)
- [数据库实现](./database-implementation.md)
- [数据库结构](./database-schema.md)
- [Token黑名单](./token-blacklist.md)
- [日志轮转](./log-rotation.md)
