# Backend API (Root Level)

基于 Gin 框架的 Go 语言后端服务

## 开发

```bash
# 启动开发服务器
go run main.go
```

---

如遇已激活用户 password_hash 为空导致无法登录，可用如下 SQL 修复：

UPDATE user SET is_active=0 WHERE is_active=1 AND (password_hash IS NULL OR password_hash='');

---

## 单元测试与集成测试规范

### 1. 目录结构与管理

- **所有测试代码、脚本、文档必须统一放在 `tests` 目录下**，不得在 backend-api 根目录或其它目录新建测试脚本。
- 单元测试放在 `tests/unit/`，集成测试放在 `tests/integration/`。
- 所有测试相关的 shell 脚本、说明文档也应放在 `tests/` 目录下。

### 2. 测试运行方式

- 推荐在 backend-api 目录下直接运行：

```bash
# 运行所有测试
./tests/run_tests.sh

# 运行单元测试
./tests/run_tests.sh unit

# 运行集成测试
./tests/run_tests.sh integration

# 运行 token 过期相关测试
./tests/run_tests.sh token-expired

# 生成覆盖率报告
./tests/run_tests.sh coverage
```

- 也可直接用 Go 命令：

```bash
# 运行所有测试
go test -v ./tests/...
```

### 3. 测试书写要求

- **测试文件命名**：单元测试以 `_test.go` 结尾，建议以被测模块命名，如 `user_service_test.go`。
- **测试函数命名**：以 `TestXxx` 开头，推荐分组使用 testify/suite。
- **测试内容**：应覆盖正常流程、异常流程、边界条件。
- **断言**：统一使用 testify/assert。
- **Mock/依赖隔离**：如需 mock，优先使用内存数据库、mock 对象等。
- **测试数据**：如需初始化数据，使用 `tests/fixtures/` 或测试工具函数。
- **覆盖率**：建议所有核心业务逻辑覆盖率 > 80%。
- **文档**：如有特殊测试说明，请写在 `tests/README.md`。

### 4. 新增/维护测试用例流程

- 新增功能时同步补充/完善对应的单元测试和集成测试。
- 测试用例需能重复运行、互不影响。
- 所有测试通过后方可提交代码。
- 如需新增测试脚本/工具/文档，统一放在 `tests/` 目录，并在 `tests/README.md` 说明。

### 5. 其它说明

- **严禁**在 backend-api 根目录新建测试脚本、测试文档。
- CI/CD 可直接调用 `./tests/run_tests.sh all` 或 `go test ./tests/...`。
- 详细测试用法、token 过期测试等见 `tests/README.md`。

---

如有疑问请联系项目维护者或查阅 `tests/README.md`。
