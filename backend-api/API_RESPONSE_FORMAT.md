# API 响应格式规范

## 统一响应结构

所有API响应都使用以下统一格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## 响应字段说明

- `code`: 状态码，表示请求处理结果
- `message`: 响应消息，描述处理结果
- `data`: 响应数据，可选字段

## 状态码定义

| 状态码 | 说明                    |
| ------ | ----------------------- |
| 200    | 成功                    |
| 400    | 请求参数错误            |
| 401    | 未授权                  |
| 40101  | Token过期（特殊错误码） |
| 404    | 资源不存在              |
| 429    | 请求过于频繁            |
| 500    | 服务器内部错误          |

## Token过期处理

### 特殊错误码说明

系统对Token过期进行了特殊处理，使用 `40101` 错误码来区分Token过期和其他认证失败情况：

- `401`: 未授权（未登录或Token缺失）
- `40101`: Token过期专用错误码

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

## API 示例

### 1. 用户注册（第一步：发送验证码）

**请求:**

```http
POST /user/register
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**成功响应:**

```json
{
  "code": 200,
  "message": "验证码已发送到邮箱",
  "data": {
    "debug": "验证码：123456"
  }
}
```

**错误响应（邮箱已注册）:**

```json
{
  "code": 400,
  "message": "邮箱已注册"
}
```

**错误响应（邮箱格式错误）:**

```json
{
  "code": 400,
  "message": "参数错误：邮箱格式不正确"
}
```

### 2. 重新发送验证码

**请求:**

```http
POST /user/send-code
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**成功响应:**

```json
{
  "code": 200,
  "message": "验证码已发送",
  "data": {
    "debug": "验证码：123456"
  }
}
```

**错误响应（邮箱已注册）:**

```json
{
  "code": 400,
  "message": "邮箱已注册"
}
```

### 3. 激活账户（第二步：设置密码）

**请求:**

```http
POST /user/verify-code
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123",
  "code": "123456"
}
```

**成功响应:**

```json
{
  "code": 200,
  "message": "激活成功，请登录",
  "data": null
}
```

**错误响应（验证码错误）:**

```json
{
  "code": 400,
  "message": "验证码无效或已过期"
}
```

**错误响应（密码格式错误）:**

```json
{
  "code": 400,
  "message": "参数错误：密码需8位以上，含大小写字母和数字"
}
```

### 4. 用户登录

**请求:**

```http
POST /user/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123"
}
```

**成功响应:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应（账户未激活）:**

```json
{
  "code": 401,
  "message": "账户未激活，请先验证邮箱"
}
```

**错误响应（密码错误）:**

```json
{
  "code": 401,
  "message": "邮箱或密码错误"
}
```

### 5. Token过期处理

**使用过期Token访问需要认证的接口:**

```http
GET /user/profile
Authorization: Bearer expired_token_here
```

**响应（Token过期）:**

```json
{
  "code": 40101,
  "message": "token已过期"
}
```

**响应（无Token）:**

```json
{
  "code": 401,
  "message": "未授权访问"
}
```

### 6. 生成过期Token（测试接口）

**请求:**

```http
POST /user/generate-expired-token
Authorization: Bearer valid_token_here
```

**成功响应:**

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

## 注册流程说明

### 两步注册流程：

1. **第一步：注册邮箱**
   - 用户提供邮箱地址
   - 系统检查邮箱是否已注册
   - 如果未注册，发送验证码
   - 如果已注册，返回错误提示

2. **第二步：激活账户**
   - 用户提供邮箱、密码和验证码
   - 系统验证验证码有效性
   - 验证密码强度
   - 激活账户并设置密码

### 优势：

- 避免重复注册
- 分离邮箱验证和密码设置
- 更好的用户体验
- 减少无效注册

## 错误处理

### 参数验证错误 (400)

```json
{
  "code": 400,
  "message": "参数错误：邮箱和密码不能为空"
}
```

### 邮箱已注册 (400)

```json
{
  "code": 400,
  "message": "邮箱已注册"
}
```

### 未授权错误 (401)

```json
{
  "code": 401,
  "message": "账户未激活，请先验证邮箱"
}
```

### Token过期错误 (40101)

```json
{
  "code": 40101,
  "message": "token已过期"
}
```

### 验证码错误 (400)

```json
{
  "code": 400,
  "message": "验证码无效或已过期"
}
```

### 请求过于频繁 (429)

```json
{
  "code": 429,
  "message": "账户暂时锁定，请稍后再试"
}
```

### 服务器错误 (500)

```json
{
  "code": 500,
  "message": "发送验证码失败"
}
```

## 测试支持

### 测试接口

系统提供了专门的测试接口来支持自动化测试：

- `POST /user/generate-expired-token`: 生成过期Token用于测试
- 所有接口都支持测试模式，返回详细的调试信息

### 测试数据

测试系统使用真实的用户数据进行测试，支持：

- 正常用户测试
- 边界条件测试
- 错误处理测试
- Token过期测试

详细测试说明请参考：`tests/README_REAL_USER_TESTING.md`
