# 数据库表结构文档

本文档描述了OpenList Strm项目的数据库表结构。

## 表结构概览

项目包含以下5个主要表：

1. **user** - 用户表
2. **openlist_service** - OpenList服务表
3. **strm_config** - Strm配置表
4. **strm_task** - Strm任务表
5. **log_record** - 日志记录表

## 表关系图

```
User "1" --> "*" OpenListService : 拥有
OpenListService "1" --> "*" StrmConfig : 包含配置
OpenListService "1" --> "*" StrmTask : 执行任务
StrmConfig "*" --> "1" StrmTask : 任务关联
StrmTask "1" --> "*" LogRecord : 生成日志
```

## 详细表结构

### 1. user 表

用户信息表，存储系统用户的基本信息。

| 字段名               | 类型     | 约束                      | 说明           |
| -------------------- | -------- | ------------------------- | -------------- |
| id                   | INTEGER  | PRIMARY KEY AUTOINCREMENT | 用户ID         |
| username             | TEXT     | UNIQUE                    | 用户名         |
| email                | TEXT     | UNIQUE                    | 邮箱地址       |
| password_hash        | TEXT     | NOT NULL                  | 密码哈希       |
| is_active            | INTEGER  | DEFAULT 0                 | 是否激活       |
| code                 | TEXT     |                           | 验证码         |
| code_expire_at       | DATETIME |                           | 验证码过期时间 |
| failed_login_count   | INTEGER  | DEFAULT 0                 | 登录失败次数   |
| locked_until         | DATETIME |                           | 锁定截止时间   |
| created_at           | DATETIME |                           | 创建时间       |
| token_invalid_before | DATETIME |                           | Token失效时间  |

### 2. openlist_service 表

OpenList服务表，存储用户配置的OpenList服务信息。

| 字段名       | 类型     | 约束                      | 说明     |
| ------------ | -------- | ------------------------- | -------- |
| id           | INTEGER  | PRIMARY KEY AUTOINCREMENT | 服务ID   |
| service_name | TEXT     | NOT NULL                  | 服务名称 |
| account      | TEXT     | NOT NULL                  | 账户名   |
| token        | TEXT     | NOT NULL                  | 访问令牌 |
| service_url  | TEXT     | NOT NULL                  | 服务URL  |
| backup_url   | TEXT     |                           | 备用URL  |
| enabled      | INTEGER  | DEFAULT 1                 | 是否启用 |
| user_id      | INTEGER  | NOT NULL, FK              | 用户ID   |
| created_at   | DATETIME | NOT NULL                  | 创建时间 |
| updated_at   | DATETIME | NOT NULL                  | 更新时间 |

**外键约束：**

- `user_id` 引用 `user(id)` ON DELETE CASCADE

### 3. strm_config 表

Strm配置表，存储每个服务的Strm相关配置。

| 字段名            | 类型     | 约束                      | 说明           |
| ----------------- | -------- | ------------------------- | -------------- |
| id                | INTEGER  | PRIMARY KEY AUTOINCREMENT | 配置ID         |
| config_name       | TEXT     | NOT NULL                  | 配置名称       |
| alist_base_path   | TEXT     | NOT NULL                  | Alist基础路径  |
| strm_output_path  | TEXT     | NOT NULL                  | Strm输出路径   |
| download_enabled  | INTEGER  | DEFAULT 0                 | 是否启用下载   |
| download_interval | INTEGER  | DEFAULT 3600              | 下载间隔（秒） |
| update_mode       | TEXT     | DEFAULT 'incremental'     | 更新模式       |
| service_id        | INTEGER  | NOT NULL, FK              | 服务ID         |
| created_at        | DATETIME | NOT NULL                  | 创建时间       |
| updated_at        | DATETIME | NOT NULL                  | 更新时间       |

**外键约束：**

- `service_id` 引用 `openlist_service(id)` ON DELETE CASCADE

**枚举值：**

- `update_mode`: 'incremental' | 'full'

### 4. strm_task 表

Strm任务表，存储定时任务的配置信息。

| 字段名         | 类型     | 约束                      | 说明           |
| -------------- | -------- | ------------------------- | -------------- |
| id             | INTEGER  | PRIMARY KEY AUTOINCREMENT | 任务ID         |
| task_name      | TEXT     | NOT NULL                  | 任务名称       |
| scheduled_time | DATETIME | NOT NULL                  | 计划执行时间   |
| task_mode      | TEXT     | NOT NULL                  | 任务模式       |
| enabled        | INTEGER  | DEFAULT 1                 | 是否启用       |
| service_id     | INTEGER  | NOT NULL, FK              | 服务ID（冗余） |
| config_id      | INTEGER  | NOT NULL, FK              | 配置ID         |
| created_at     | DATETIME | NOT NULL                  | 创建时间       |
| updated_at     | DATETIME | NOT NULL                  | 更新时间       |

**外键约束：**

- `service_id` 引用 `openlist_service(id)` ON DELETE CASCADE
- `config_id` 引用 `strm_config(id)` ON DELETE CASCADE

**枚举值：**

- `task_mode`: 'create' | 'check'

### 5. log_record 表

日志记录表，存储任务执行的日志信息。

| 字段名      | 类型     | 约束                      | 说明         |
| ----------- | -------- | ------------------------- | ------------ |
| id          | INTEGER  | PRIMARY KEY AUTOINCREMENT | 日志ID       |
| log_name    | TEXT     | NOT NULL                  | 日志名称     |
| log_path    | TEXT     | NOT NULL                  | 日志文件路径 |
| created_at  | DATETIME | NOT NULL                  | 创建时间     |
| task_status | TEXT     | NOT NULL                  | 任务状态     |
| task_id     | INTEGER  | NOT NULL, FK              | 任务ID       |

**外键约束：**

- `task_id` 引用 `strm_task(id)` ON DELETE CASCADE

**枚举值：**

- `log_name`: 'create' | 'check'
- `task_status`: 'running' | 'error' | 'completed'

## 索引建议

为了提高查询性能，建议在以下字段上创建索引：

1. `user.email` - 用户邮箱查询
2. `user.username` - 用户名查询
3. `openlist_service.user_id` - 用户服务查询
4. `strm_config.service_id` - 服务配置查询
5. `strm_task.service_id` - 服务任务查询
6. `strm_task.config_id` - 配置任务查询
7. `strm_task.scheduled_time` - 任务调度查询
8. `log_record.task_id` - 任务日志查询
9. `log_record.created_at` - 日志时间查询

## 数据完整性

### 级联删除

- 删除用户时，自动删除该用户的所有服务
- 删除服务时，自动删除该服务的所有配置和任务
- 删除配置时，自动删除该配置的所有任务
- 删除任务时，自动删除该任务的所有日志记录

### 约束检查

- 枚举字段使用CHECK约束确保数据有效性
- 外键约束确保数据引用完整性
- 唯一约束防止重复数据

## 使用示例

### 创建服务

```go
service := &model.OpenListService{
    ServiceName: "我的服务",
    Account:     "myaccount",
    Token:       "mytoken",
    ServiceUrl:  "http://localhost:5244",
    Enabled:     true,
    UserID:      1,
}
err := model.CreateOpenListService(db, service)
```

### 查询用户的所有服务

```go
services, err := model.GetOpenListServicesByUserID(db, userID)
```

### 创建任务

```go
task := &model.StrmTask{
    TaskName:      "每日同步",
    ScheduledTime: time.Now().Add(24 * time.Hour),
    TaskMode:      model.TaskModeCreate,
    Enabled:       true,
    ServiceID:     1,
    ConfigID:      1,
}
err := model.CreateStrmTask(db, task)
```

## 测试

运行数据库测试：

```bash
cd backend-api
go run cmd/test_db/main.go
```

这将创建所有表并插入测试数据，验证数据库功能是否正常。
