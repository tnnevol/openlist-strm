# 用户管理脚本模块

这个脚本模块提供了用户管理和测试数据管理的功能。

## 功能特性

1. **通过用户ID获取用户信息** - 查询并显示指定用户的详细信息
2. **为openlist_service表添加测试数据** - 自动创建测试用的服务配置数据

## 文件结构

```
scripts/
├── user_management.go      # 主要的Go脚本文件
├── run_user_management.sh  # Shell运行脚本
└── README.md              # 使用说明文档
```

## 使用方法

### 方法1: 使用Shell脚本（推荐）

```bash
# 显示帮助信息
./scripts/run_user_management.sh -help

# 获取用户ID为1的用户信息
./scripts/run_user_management.sh -action getUser -userid 1

# 为用户ID为1添加openlist_service测试数据
./scripts/run_user_management.sh -action addTestData -userid 1

 ./scripts/run_user_management.sh -action createStrmConfigTestData -userid 1
```

### 方法2: 直接使用Go命令

```bash
# 切换到backend-api目录
cd backend-api

# 显示帮助信息
go run scripts/user_management.go -help

# 获取用户ID为1的用户信息
go run scripts/user_management.go -action getUser -userid 1

# 为用户ID为1添加openlist_service测试数据
go run scripts/user_management.go -action addTestData -userid 1
```

## 操作说明

### getUser 操作

通过用户ID获取用户的详细信息，包括：

- 用户名
- 邮箱地址
- 账户状态（是否激活）
- 验证码信息
- 登录失败次数
- 账户锁定状态
- 创建时间
- Token失效时间

**参数：**

- `-userid`: 要查询的用户ID（必需）

**示例：**

```bash
./scripts/run_user_management.sh -action getUser -userid 1
```

### addTestData 操作

为openlist_service表添加测试数据，包括：

- 测试服务1（启用状态）
- 测试服务2（禁用状态）
- 生产环境服务（启用状态）

**参数：**

- `-userid`: 要添加测试数据的用户ID（必需）

**注意：** 此操作需要指定有效的用户ID，脚本会验证用户是否存在。

**示例：**

```bash
./scripts/run_user_management.sh -action addTestData -userid 1
```

## 测试数据说明

addTestData操作会批量创建50条测试服务数据：

- 服务名: 测试服务1, 测试服务2, ..., 测试服务50
- 账户: test_account_1, test_account_2, ..., test_account_50
- Token: test_token_000001, test_token_000002, ..., test_token_000050
- URL: https://api.example1.com, https://api.example2.com, ..., https://api.example50.com
- 启用状态: 奇数编号启用，偶数编号禁用

## 错误处理

脚本包含完善的错误处理机制：

- 数据库连接失败时会显示错误信息
- 用户不存在时会提示错误
- 测试数据添加失败时会继续处理其他数据
- 参数错误时会显示帮助信息

## 依赖要求

- Go 1.16+
- 数据库连接配置正确
- 项目依赖已安装（go mod tidy）

## 注意事项

1. 运行脚本前确保数据库服务正在运行
2. 确保数据库连接配置正确
3. 添加测试数据前确保至少有一个用户存在
4. 脚本会自动处理数据库连接和关闭
5. 所有操作都会记录详细的日志信息
