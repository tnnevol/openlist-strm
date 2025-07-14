# Token黑名单功能说明

## 功能概述

Token黑名单功能用于实现安全的用户登出机制，确保用户登出后token立即失效，无法继续使用。

## 核心特性

### 1. 立即失效

- 用户登出后，token立即加入黑名单
- 后续使用该token的请求会被拒绝
- 返回"token已失效，请重新登录"错误

### 2. 安全存储

- 使用SHA256哈希存储token，避免明文存储
- 内存存储，高性能访问
- 线程安全的读写操作

### 3. 自动清理

- 每小时自动清理过期的token
- 减少内存占用
- 避免黑名单无限增长

### 4. 监控功能

- 提供黑名单状态查询接口
- 实时监控黑名单大小
- 便于系统监控和调试

## 技术实现

### 1. 黑名单管理

```go
type TokenBlacklist struct {
    blacklist map[string]time.Time  // token哈希 -> 过期时间
    mutex     sync.RWMutex         // 读写锁
}
```

### 2. 核心方法

- `AddToBlacklist(token, expireTime)`: 添加token到黑名单
- `IsBlacklisted(token)`: 检查token是否在黑名单中
- `GetBlacklistSize()`: 获取黑名单大小
- `cleanupExpiredTokens()`: 清理过期token

### 3. 中间件集成

在 `AuthMiddleware` 中增加黑名单检查：

```go
// 检查token是否在黑名单中
blacklist := service.GetTokenBlacklist()
if blacklist.IsBlacklisted(tokenStr) {
    Unauthorized(c, "token已失效，请重新登录")
    c.Abort()
    return
}
```

## 接口说明

### 1. 登出接口

```http
POST /user/logout
Authorization: Bearer {token}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

### 2. 黑名单状态接口

```http
GET /user/token-blacklist-status
Authorization: Bearer {token}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "blacklistSize": 5,
    "timestamp": 1640995200
  }
}
```

## 使用流程

### 1. 正常登录

```bash
# 1. 用户登录获取token
curl -X POST http://localhost:8890/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "TestPass123"}'

# 2. 使用token访问受保护接口
curl -X GET http://localhost:8890/user/info \
  -H "Authorization: Bearer {token}"
```

### 2. 用户登出

```bash
# 3. 用户登出
curl -X POST http://localhost:8890/user/logout \
  -H "Authorization: Bearer {token}"

# 4. 再次使用相同token访问接口（会失败）
curl -X GET http://localhost:8890/user/info \
  -H "Authorization: Bearer {token}"
# 响应: {"code": 401, "message": "token已失效，请重新登录"}
```

## 测试验证

运行测试脚本验证功能：

```bash
./test_token_blacklist.sh
```

测试脚本会：

1. 登录获取token
2. 使用token访问接口（成功）
3. 执行登出
4. 再次使用相同token访问接口（失败）
5. 验证响应包含"token已失效"

## 性能考虑

### 1. 内存使用

- 每个token哈希占用64字节（SHA256）
- 时间戳占用8字节
- 建议监控黑名单大小，避免内存泄漏

### 2. 访问性能

- 使用map存储，O(1)查找时间
- 读写锁保证并发安全
- 适合高并发场景

### 3. 清理策略

- 每小时自动清理过期token
- 访问时检查并清理过期token
- 可配置清理间隔

## 监控建议

### 1. 黑名单大小监控

```bash
# 定期检查黑名单大小
curl -X GET http://localhost:8890/user/token-blacklist-status \
  -H "Authorization: Bearer {token}"
```

### 2. 日志监控

- 关注 `[TokenBlacklist]` 相关日志
- 监控清理过期token的频率
- 观察黑名单增长趋势

### 3. 告警设置

- 黑名单大小超过阈值时告警
- 清理频率异常时告警
- 内存使用过高时告警

## 注意事项

1. **重启影响**: 服务重启后黑名单会清空，已登出的token可能重新有效
2. **内存限制**: 大量并发登出可能导致内存使用增加
3. **时钟同步**: 确保服务器时间准确，避免token过期时间计算错误
4. **分布式部署**: 多实例部署时黑名单不共享，需要考虑分布式解决方案

## 扩展方案

### 1. 持久化存储

- 使用Redis存储黑名单
- 支持分布式部署
- 服务重启后黑名单不丢失

### 2. 数据库存储

- 使用数据库存储黑名单
- 支持复杂查询和统计
- 适合大规模部署

### 3. 分布式缓存

- 使用Memcached或Redis集群
- 高可用和高性能
- 支持自动过期
