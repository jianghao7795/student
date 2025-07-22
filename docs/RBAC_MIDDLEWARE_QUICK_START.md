# RBAC中间件快速开始指南

## 概述

本项目已经成功集成了RBAC中间件，提供了完整的权限控制功能。本文档将指导您如何快速使用RBAC中间件。

## 已完成的配置

### 1. 中间件集成

RBAC中间件已经集成到HTTP服务器中：

```go
// internal/server/http.go
func NewHTTPServer(c *conf.Bootstrap, student *service.StudentService, user *service.UserService, rbac *service.RBACService, rbacUC *biz.RBACUsecase, jwtUtil *jwt.JWTUtil, logger log.Logger) *http.Server {
    var opts = []http.ServerOption{
        http.Middleware(
            recovery.Recovery(),
            // JWT认证中间件
            middleware.JWTAuth(&middleware.JWTConfig{
                JWTUtil:   jwtUtil,
                SkipPaths: []string{"/v1/user/login", "/v1/user/register"},
            }),
            // RBAC权限中间件
            middleware.RBACMiddleware(rbacUC, jwtUtil),
        ),
    }
    // ...
}
```

### 2. 依赖注入配置

Wire依赖注入已经正确配置，确保RBACUsecase被注入到HTTP服务器中。

## 使用方法

### 1. 启动服务

```bash
# 编译项目
go build ./cmd/student

# 启动服务
./bin/student -conf configs/config.yaml
```

### 2. 测试RBAC中间件

使用提供的测试脚本：

```bash
# 运行RBAC中间件测试
./script/test_rbac_middleware.sh
```

### 3. 手动测试

#### 获取JWT Token

```bash
curl -X POST http://localhost:8000/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

#### 使用Token访问受保护资源

```bash
# 替换YOUR_JWT_TOKEN为实际的token
curl -X GET http://localhost:8000/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Request-Path: /v1/users" \
  -H "X-Request-Method: GET"
```

## 中间件功能

### 1. 自动权限检查

RBAC中间件会自动：
- 验证JWT token
- 提取用户信息
- 检查用户权限
- 返回相应的HTTP状态码

### 2. 错误处理

| 情况 | HTTP状态码 | 说明 |
|------|-----------|------|
| 无Authorization头 | 401 | 需要认证 |
| 无效JWT token | 401 | Token验证失败 |
| 权限不足 | 403 | 没有访问权限 |
| 权限检查失败 | 500 | 服务器错误 |

### 3. 请求头要求

- **必需**: `Authorization: Bearer <token>`
- **可选**: `X-Request-Path` 或 `X-Original-URI`
- **可选**: `X-Request-Method` 或 `X-HTTP-Method`

## 权限配置

### 1. 角色定义

系统预定义了以下角色：
- `admin`: 管理员，拥有所有权限
- `user`: 普通用户，拥有基本权限
- `guest`: 访客，只有查看权限

### 2. 权限策略

权限策略在 `rbac_policy.csv` 中定义：

```csv
p, admin, /v1/users, *
p, user, /v1/users, GET
g, 1, admin
g, 2, user
```

### 3. 权限检查

使用权限检查API：

```bash
curl -X POST http://localhost:8000/v1/permissions/check \
  -H "Content-Type: application/json" \
  -d '{
    "user": "1",
    "resource": "/v1/users",
    "action": "GET"
  }'
```

## 扩展使用

### 1. 自定义权限检查

可以在业务代码中直接调用权限检查：

```go
hasPermission, err := rbacUC.CheckPermission(ctx, userID, resource, action)
if err != nil {
    return nil, err
}
if !hasPermission {
    return nil, errors.Forbidden("FORBIDDEN", "没有访问权限")
}
```

### 2. 简化版中间件

对于特定资源，可以使用简化版中间件：

```go
userReadMiddleware := middleware.SimpleRBACMiddleware(rbacUC, jwtUtil, "/v1/users", "GET")
```

## 故障排除

### 1. 编译错误

如果遇到编译错误，请检查：
- Wire依赖注入是否正确生成
- 所有依赖是否正确导入

### 2. 权限检查不生效

检查：
- JWT token是否有效
- 请求头是否包含正确的路径和方法信息
- RBAC策略配置是否正确

### 3. 服务启动失败

检查：
- 数据库连接是否正常
- 配置文件是否正确
- 端口是否被占用

## 总结

RBAC中间件已经成功集成到项目中，提供了完整的权限控制功能。通过JWT认证和RBAC权限检查，确保只有具有适当权限的用户才能访问受保护的资源。

如需更详细的使用说明，请参考 `docs/RBAC_MIDDLEWARE_USAGE.md`。
