# API 版本化更改总结

## 概述

为了统一 API 接口的版本管理，所有接口路径已更新为以 `/v1` 开头，实现了一致的版本化策略。

## 更改内容

### 1. 学生服务 (Student Service)

**原路径** → **新路径**

- `GET /student/{id}` → `GET /v1/student/{id}`
- `POST /student` → `POST /v1/student`
- `PUT /student/{id}` → `PUT /v1/student/{id}`
- `DELETE /student/{id}` → `DELETE /v1/student/{id}`
- `GET /students` → `GET /v1/students`

### 2. 用户服务 (User Service)

**原路径** → **新路径**

- `GET /user/{id}` → `GET /v1/user/{id}`
- `POST /user` → `POST /v1/user`
- `PUT /user/{id}` → `PUT /v1/user/{id}`
- `DELETE /user/{id}` → `DELETE /v1/user/{id}`
- `GET /users` → `GET /v1/users`
- `POST /user/login` → `POST /v1/user/login`
- 新增: `GET /v1/user/me` - 获取当前用户信息（需要 JWT 认证）

### 3. 错误服务 (Error Service)

**原路径** → **新路径**

- `GET /errors/{error_code}` → `GET /v1/errors/{error_code}`
- `GET /errors` → `GET /v1/errors`
- `POST /errors/custom` → `POST /v1/errors/custom`

### 4. RBAC 服务 (RBAC Service)

**原路径** → **新路径**

- `GET /api/v1/roles/{id}` → `GET /v1/roles/{id}`
- `POST /api/v1/roles` → `POST /v1/roles`
- `PUT /api/v1/roles/{id}` → `PUT /v1/roles/{id}`
- `DELETE /api/v1/roles/{id}` → `DELETE /v1/roles/{id}`
- `GET /api/v1/roles` → `GET /v1/roles`
- `GET /api/v1/permissions/{id}` → `GET /v1/permissions/{id}`
- `POST /api/v1/permissions` → `POST /v1/permissions`
- `PUT /api/v1/permissions/{id}` → `PUT /v1/permissions/{id}`
- `DELETE /api/v1/permissions/{id}` → `DELETE /v1/permissions/{id}`
- `GET /api/v1/permissions` → `GET /v1/permissions`
- `GET /api/v1/users/{user_id}/roles` → `GET /v1/users/{user_id}/roles`
- `POST /api/v1/users/{user_id}/roles` → `POST /v1/users/{user_id}/roles`
- `DELETE /api/v1/users/{user_id}/roles/{role_id}` → `DELETE /v1/users/{user_id}/roles/{role_id}`
- `GET /api/v1/roles/{role_id}/permissions` → `GET /v1/roles/{role_id}/permissions`
- `POST /api/v1/roles/{role_id}/permissions` → `POST /v1/roles/{role_id}/permissions`
- `DELETE /api/v1/roles/{role_id}/permissions/{permission_id}` → `DELETE /v1/roles/{role_id}/permissions/{permission_id}`
- `POST /api/v1/permissions/check` → `POST /v1/permissions/check`

## 修改的文件

### Proto 文件

- `api/student/v1/student.proto` - 更新学生服务的 HTTP 注解
- `api/user/v1/user.proto` - 更新用户服务的 HTTP 注解，添加 GetMe 接口
- `api/errors/v1/errors.proto` - 更新错误服务的 HTTP 注解
- `api/rbac/v1/rbac.proto` - 更新 RBAC 服务的 HTTP 注解

### 生成的代码

- `api/student/v1/student_http.pb.go` - 重新生成的学生服务 HTTP 路由
- `api/user/v1/user_http.pb.go` - 重新生成的用户服务 HTTP 路由，包含 GetMe 接口
- `api/errors/v1/errors_http.pb.go` - 重新生成的错误服务 HTTP 路由
- `api/rbac/v1/rbac_http.pb.go` - 重新生成的 RBAC 服务 HTTP 路由

### 服务器配置

- `internal/server/http.go` - 添加 RBAC 服务注册到 HTTP 服务器
- `cmd/student/wire_gen.go` - 更新依赖注入配置

### 业务逻辑

- `internal/biz/user.go` - 添加 GetMe 业务逻辑方法
- `internal/service/user.go` - 添加 GetMe 服务层实现

### 配置文件

- `rbac_policy.csv` - 更新 RBAC 策略中的路径

### 文档文件

- `docs/RBAC_README.md` - 更新 API 接口文档
- `RBAC_IMPLEMENTATION_SUMMARY.md` - 更新实现总结文档
- `openapi.yaml` - 重新生成的 OpenAPI 文档
- `GET_ME_API_USAGE.md` - 新增 GetMe 接口使用指南

## 影响范围

### 客户端调用

所有使用旧路径的客户端代码需要更新为新的 `/v1` 路径。

### 权限策略

RBAC 权限策略中的路径已更新，确保权限检查正常工作。

### 文档

所有相关文档已更新，反映新的 API 路径。

## 兼容性

这是一个破坏性更改，旧版本的 API 路径将不再可用。建议：

1. 通知所有 API 使用者更新客户端代码
2. 提供迁移指南
3. 考虑在过渡期间提供重定向或代理

## 验证

- ✅ 项目编译成功
- ✅ 所有 protobuf 文件重新生成
- ✅ OpenAPI 文档更新
- ✅ RBAC 策略更新
- ✅ 文档更新
- ✅ GetMe 接口实现完成

## 下一步

1. 测试所有 API 端点功能正常
2. 更新客户端代码示例
3. 部署到测试环境验证
4. 通知相关团队进行迁移
