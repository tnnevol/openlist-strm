# 日志轮转功能说明

## 功能特性

### 1. 按天自动分割

- 当天的日志写入 `logs/main.log`
- 历史日志按日期归档为 `logs/main.log.YYYY-MM-DD`
- 项目运行时自动检测日期变化并执行轮转

### 2. 文件大小控制

- 单个日志文件最大 100MB
- 超过大小限制时自动分割
- 保留最近 30 个备份文件

### 3. 时间控制

- 保留最近 30 天的日志文件
- 超过保留期限的日志文件自动删除

### 4. 压缩功能

- 历史日志文件自动压缩
- 节省磁盘空间

## 配置说明

### 默认配置

```go
var defaultLogConfig = LogConfig{
    Filename:    "logs/main.log",  // 日志文件路径
    MaxSize:     100,              // 单个文件最大100MB
    MaxBackups:  30,               // 保留30个备份文件
    MaxAge:      30,               // 保留30天
    Compress:    true,             // 压缩旧文件
    Level:       "info",           // 日志级别
}
```

### 自定义配置

```go
import "github.com/tnnevol/openlist-strm/backend-api/internal/logger"

// 设置自定义配置
config := logger.LogConfig{
    Filename:    "logs/custom.log",
    MaxSize:     50,               // 50MB
    MaxBackups:  10,               // 保留10个备份
    MaxAge:      7,                // 保留7天
    Compress:    true,
    Level:       "debug",
}

logger.InitWithConfig(config)
```

## 日志级别

- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息
- `fatal`: 致命错误

## 使用方法

### 基本使用

```go
import "github.com/tnnevol/openlist-strm/backend-api/internal/logger"

// 初始化日志系统
logger.Init()

// 记录日志
logger.Info("这是一条信息日志")
logger.Error("这是一条错误日志")
logger.Debug("这是一条调试日志")
logger.Warn("这是一条警告日志")
```

### 带字段的日志

```go
logger.Info("用户登录",
    zap.String("username", "testuser"),
    zap.String("ip", "192.168.1.1"),
    zap.Time("loginTime", time.Now()))
```

## 文件结构

```
logs/
├── main.log                    # 当前日志文件
├── main.log.2025-01-09         # 历史日志文件
├── main.log.2025-01-08
└── main.log.2025-01-07
```

## 测试

运行测试脚本验证日志轮转功能：

```bash
./test_log_rotation.sh
```

## 注意事项

1. 日志轮转是线程安全的，使用互斥锁保护
2. 每次写入日志时都会检查日期变化
3. 轮转时会自动重新初始化 logger
4. 建议在生产环境中定期清理过期的日志文件
5. 可以通过配置文件 `config/log.yaml` 管理日志设置

## 故障排除

### 问题1: 日志文件不轮转

- 检查系统时间是否正确
- 确认日志目录有写权限
- 查看是否有其他进程占用日志文件

### 问题2: 日志文件过大

- 调整 `MaxSize` 参数
- 检查是否有大量重复日志
- 考虑调整日志级别

### 问题3: 磁盘空间不足

- 减少 `MaxBackups` 和 `MaxAge` 参数
- 启用压缩功能
- 定期手动清理旧日志文件
