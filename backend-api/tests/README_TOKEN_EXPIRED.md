# Token过期测试说明

## 文件重组

根据项目规范，所有测试相关的文件已统一管理在`tests`目录下：

### 移动的文件

1. `backend-api/test_token_expired_unit.sh` → `tests/test_token_expired_unit.sh`
2. `backend-api/TOKEN_EXPIRED_HANDLING.md` → 内容已整合到 `tests/README.md`

### 统一管理

所有token过期测试现在都通过 `tests/run_tests.sh` 进行管理：

```bash
# 运行token过期测试
./tests/run_tests.sh token-expired

# 运行所有测试
./tests/run_tests.sh all

# 查看帮助
./tests/run_tests.sh help
```

## 测试文件结构

```
tests/
├── README.md                           # 主要测试文档（包含token过期说明）
├── run_tests.sh                        # 统一测试运行脚本
├── test_token_expired_unit.sh          # token过期专用脚本
├── integration/
│   └── token_expired_integration_test.go  # token过期测试套件
└── ...
```

## 运行方式

### 1. 统一脚本（推荐）

```bash
cd tests
./run_tests.sh token-expired
```

### 2. 专用脚本

```bash
cd tests
./test_token_expired_unit.sh
```

### 3. 直接运行

```bash
cd tests/integration
go test -v -run TestTokenExpiredIntegrationTestSuite ./...
```

## 注意事项

- 所有测试脚本现在都在`tests`目录下
- 通过`run_tests.sh`统一管理所有测试
- 详细的token过期处理说明已整合到`tests/README.md`
