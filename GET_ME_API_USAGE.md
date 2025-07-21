# GetMe API 使用指南

## 接口概述

`GET /v1/user/me` 接口用于获取当前登录用户的详细信息，该接口需要在请求头中包含有效的 JWT token。

## 接口详情

### 请求信息

- **方法**: `GET`
- **路径**: `/v1/user/me`
- **认证**: 需要 JWT token（Authorization 头）

### 请求头

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

### 响应格式

#### 成功响应 (200 OK)

```json
{
  "success": true,
  "message": "获取用户信息成功",
  "user_info": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "phone": "13800138000",
    "status": 1,
    "age": 25,
    "avatar": "https://example.com/avatar.jpg",
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 12:00:00"
  }
}
```

#### 失败响应 (401 Unauthorized)

```json
{
  "success": false,
  "message": "未找到用户信息"
}
```

## 使用示例

### 1. 使用 curl 测试

```bash
# 首先登录获取 token
curl -X POST http://localhost:8000/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'

# 使用返回的 token 调用 GetMe 接口
curl -X GET http://localhost:8000/v1/user/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### 2. 使用 JavaScript/Fetch API

```javascript
// 登录获取 token
const loginResponse = await fetch("/v1/user/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    username: "admin",
    password: "password123",
  }),
});

const loginData = await loginResponse.json();
const token = loginData.token;

// 使用 token 获取当前用户信息
const meResponse = await fetch("/v1/user/me", {
  method: "GET",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  },
});

const userData = await meResponse.json();
console.log("当前用户信息:", userData.user_info);
```

### 3. 使用 Python requests

```python
import requests

# 登录获取 token
login_data = {
    "username": "admin",
    "password": "password123"
}

login_response = requests.post(
    "http://localhost:8000/v1/user/login",
    json=login_data
)

token = login_response.json()["token"]

# 使用 token 获取当前用户信息
headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}

me_response = requests.get(
    "http://localhost:8000/v1/user/me",
    headers=headers
)

user_info = me_response.json()
print("当前用户信息:", user_info["user_info"])
```

## 错误处理

### 常见错误码

| 状态码 | 错误信息               | 说明                                |
| ------ | ---------------------- | ----------------------------------- |
| 401    | 未提供有效的认证 token | 缺少 Authorization 头               |
| 401    | token 验证失败         | JWT token 无效或已过期              |
| 401    | 未找到用户信息         | 用户不存在或 token 中的用户 ID 无效 |

### 错误响应示例

```json
{
  "success": false,
  "message": "token验证失败"
}
```

## 安全注意事项

1. **Token 安全**: JWT token 包含敏感信息，请妥善保管，不要泄露给第三方
2. **HTTPS**: 生产环境请使用 HTTPS 传输，确保 token 安全
3. **Token 过期**: JWT token 有过期时间，过期后需要重新登录获取新 token
4. **权限控制**: 该接口返回当前登录用户的完整信息，确保只有用户本人可以访问

## 业务逻辑

### 实现流程

1. **Token 验证**: 从 Authorization 头中提取 JWT token
2. **用户信息获取**: 从 token 中解析用户 ID，查询数据库获取用户详细信息
3. **角色信息获取**: 通过 RBAC 系统获取用户的角色信息
4. **响应返回**: 返回用户完整信息（不包含密码）

### 数据来源

- **用户基本信息**: 从 `users` 表获取
- **角色信息**: 从 RBAC 系统获取
- **时间信息**: 自动格式化为可读字符串

## 相关接口

- `POST /v1/user/login` - 用户登录，获取 JWT token
- `GET /v1/user/{id}` - 获取指定用户信息（需要权限）
- `PUT /v1/user/{id}` - 更新用户信息（需要权限）

## 更新日志

- **v1.0.0**: 初始版本，支持获取当前用户基本信息
- **v1.1.0**: 添加角色信息支持
- **v1.2.0**: 优化错误处理和响应格式
