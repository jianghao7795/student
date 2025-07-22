# gRPC RBAC 权限过滤使用指南

## 概述

本文档介绍如何在 gRPC 服务器中配置和使用 RBAC（基于角色的访问控制）权限过滤功能。

## 功能特性

- 支持基于 JWT Token 的用户认证
- 支持基于 RBAC 模型的权限检查
- 可配置跳过权限检查的路径
- 与 HTTP 服务器共享相同的 RBAC 配置

## 配置说明

### 1. 配置文件设置

在 `configs/config.yaml` 中确保 RBAC 配置已启用：

```yaml
rbac:
  model_path: "rbac_model.conf"
  policy_path: "rbac_policy.csv"
  enabled: true
```

### 2. gRPC 服务器配置

gRPC 服务器会自动根据配置决定是否启用 RBAC 中间件：

- 当 `rbac.enabled = true` 时，gRPC 服务器会添加 RBAC 权限检查中间件
- 当 `rbac.enabled = false` 时，gRPC 服务器只使用基础的 recovery 中间件

## 使用方法

### 1. 客户端调用

gRPC 客户端在调用需要权限验证的方法时，需要在请求头中包含 JWT Token：

```go
// 创建 gRPC 连接
conn, err := grpc.Dial("localhost:9600", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// 创建客户端
client := v1.NewStudentClient(conn)

// 创建带认证的上下文
ctx := context.Background()
ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer your-jwt-token-here")

// 调用 gRPC 方法
response, err := client.GetStudent(ctx, &v1.GetStudentRequest{Id: 1})
```

### 2. 跳过权限检查

如果需要跳过某些 gRPC 方法的权限检查，可以在 `internal/server/grpc.go` 中的 `SkipPaths` 配置中添加：

```go
rbacConfig := &middleware.RBACConfig{
    RBACUC: rbacUC,
    JWTUtil: jwtUtil,
    SkipPaths: []string{
        "/student.v1.Student/GetStudent",  // 跳过获取学生信息的权限检查
        "/user.v1.User/Login",             // 跳过登录接口的权限检查
    },
}
```

## 权限检查流程

1. **请求拦截**: gRPC 中间件拦截所有请求
2. **路径检查**: 检查请求路径是否在跳过列表中
3. **Token 验证**: 从请求头获取并验证 JWT Token
4. **权限验证**: 使用 RBAC 模型检查用户是否有访问权限
5. **请求处理**: 如果权限验证通过，继续处理请求；否则返回错误

## 错误处理

RBAC 中间件可能返回以下错误：

- `UNAUTHORIZED`: 未提供认证 token 或 token 无效
- `FORBIDDEN`: 用户没有访问权限
- `INTERNAL_ERROR`: 权限检查过程中发生内部错误

## 示例

### 完整的 gRPC 客户端示例

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    v1 "student/api/student/v1"
)

func main() {
    // 连接 gRPC 服务器
    conn, err := grpc.Dial("localhost:9600", grpc.WithInsecure())
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer conn.Close()
    
    // 创建客户端
    client := v1.NewStudentClient(conn)
    
    // 创建带认证的上下文
    ctx := context.Background()
    ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer your-jwt-token-here")
    
    // 调用需要权限验证的方法
    response, err := client.CreateStudent(ctx, &v1.CreateStudentRequest{
        Name:  "张三",
        Age:   20,
        Grade: "大一",
    })
    
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    log.Printf("Student created: %v", response)
}
```

## 注意事项

1. **Token 格式**: JWT Token 必须以 "Bearer " 开头
2. **权限配置**: 确保 RBAC 模型和策略文件配置正确
3. **性能考虑**: RBAC 检查会增加一定的请求延迟
4. **错误处理**: 客户端需要正确处理权限相关的错误响应

## 故障排除

### 常见问题

1. **权限被拒绝**: 检查用户角色和权限配置
2. **Token 无效**: 确认 JWT Token 格式和有效期
3. **配置错误**: 验证 RBAC 模型和策略文件路径

### 调试方法

1. 检查服务器日志中的权限检查信息
2. 验证 RBAC 配置是否正确加载
3. 确认客户端传递的 Token 格式正确
