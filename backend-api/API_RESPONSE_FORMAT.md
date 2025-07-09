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

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

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