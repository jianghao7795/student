# RBAC 权限控制系统使用说明

## 概述

本项目基于 Casbin 实现了完整的 RBAC（基于角色的访问控制）权限管理系统，支持角色管理、权限管理、用户角色分配和权限检查等功能。

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP/gRPC     │    │   RBAC Service  │    │   Casbin        │
│   API Layer     │───▶│   Business      │───▶│   Enforcer      │
└─────────────────┘    │   Logic         │    └─────────────────┘
                       └─────────────────┘
                                │
                       ┌─────────────────┐
                       │   Database      │
                       │   (MySQL)       │
                       └─────────────────┘
```

## 数据库表结构

### 1. 角色表 (roles)

```sql
CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '角色名称',
  `description` varchar(500) DEFAULT NULL COMMENT '角色描述',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_name` (`name`)
);
```

### 2. 权限表 (permissions)

```sql
CREATE TABLE `permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `resource` varchar(200) NOT NULL COMMENT '资源路径',
  `action` varchar(50) NOT NULL COMMENT '操作类型：GET,POST,PUT,DELETE等',
  `description` varchar(500) DEFAULT NULL COMMENT '权限描述',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_permission_resource_action` (`resource`, `action`)
);
```

### 3. 用户角色关联表 (user_roles)

```sql
CREATE TABLE `user_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
);
```

### 4. 角色权限关联表 (role_permissions)

```sql
CREATE TABLE `role_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `permission_id` int(11) NOT NULL COMMENT '权限ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
);
```

## 预定义角色和权限

### 角色

1. **admin** - 系统管理员，拥有所有权限
2. **manager** - 部门经理，拥有部门管理权限
3. **user** - 普通用户，拥有基本权限
4. **guest** - 访客，只有查看权限

### 权限

- **用户管理权限**: user:read, user:create, user:update, user:delete, user:detail
- **学生管理权限**: student:read, student:create, student:update, student:delete, student:detail
- **角色管理权限**: role:read, role:create, role:update, role:delete, role:detail
- **权限管理权限**: permission:read, permission:create, permission:update, permission:delete, permission:detail

## API 接口

### 角色管理

- `GET /v1/roles` - 获取角色列表
- `GET /v1/roles/{id}` - 获取角色详情
- `POST /v1/roles` - 创建角色
- `PUT /v1/roles/{id}` - 更新角色
- `DELETE /v1/roles/{id}` - 删除角色

### 权限管理

- `GET /v1/permissions` - 获取权限列表
- `GET /v1/permissions/{id}` - 获取权限详情
- `POST /v1/permissions` - 创建权限
- `PUT /v1/permissions/{id}` - 更新权限
- `DELETE /v1/permissions/{id}` - 删除权限

### 用户角色管理

- `GET /v1/users/{user_id}/roles` - 获取用户角色
- `POST /v1/users/{user_id}/roles` - 分配用户角色
- `DELETE /v1/users/{user_id}/roles/{role_id}` - 移除用户角色

### 角色权限管理

- `GET /v1/roles/{role_id}/permissions` - 获取角色权限
- `POST /v1/roles/{role_id}/permissions` - 分配角色权限
- `DELETE /v1/roles/{role_id}/permissions/{permission_id}` - 移除角色权限

### 权限检查

- `POST /v1/permissions/check` - 检查用户权限

## 使用示例

### 1. 用户登录获取角色信息

```go
// 用户登录后，返回的用户信息包含角色列表
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "roles": ["admin"],
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. 权限检查

```bash
# 检查用户是否有权限访问用户列表
curl -X POST http://localhost:8000/v1/permissions/check \
  -H "Content-Type: application/json" \
  -d '{
    "user": "1",
    "resource": "/v1/users",
    "action": "GET"
  }'
```

### 3. 分配用户角色

```bash
# 为用户分配角色
curl -X POST http://localhost:8000/v1/users/2/roles \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 2
  }'
```

### 4. 分配角色权限

```bash
# 为角色分配权限
curl -X POST http://localhost:8000/v1/roles/2/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "permission_id": 1
  }'
```

## 中间件使用

### RBAC 权限中间件

```go
// 在HTTP路由中使用RBAC中间件
rbacMiddleware := middleware.RBACMiddleware(rbacUC, jwtUtil)
router.Use(rbacMiddleware)

// 或者针对特定路径使用简化版中间件
simpleRBACMiddleware := middleware.SimpleRBACMiddleware(rbacUC, jwtUtil, "/v1/users", "GET")
router.Use(simpleRBACMiddleware)
```

## 配置说明

### Casbin 模型配置 (configs/rbac_model.conf)

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

## 部署步骤

1. **执行数据库迁移**

```bash
mysql -u username -p database_name < migrate/rbac_migrate.sql
```

2. **启动服务**

```bash
make run
```

3. **测试 RBAC 功能**

```bash
go run examples/rbac_test.go
```

## 注意事项

1. **权限缓存**: Casbin 会自动缓存权限策略，修改权限后需要重新加载策略
2. **角色继承**: 当前版本支持简单的角色分配，不支持角色继承
3. **权限粒度**: 权限检查支持通配符匹配，如 `/v1/users/*` 可以匹配所有用户相关操作
4. **安全性**: 所有 API 接口都需要 JWT 认证，权限检查在认证之后进行

## 扩展功能

1. **角色继承**: 可以实现角色的层级继承关系
2. **动态权限**: 支持运行时动态添加和修改权限
3. **权限组**: 可以将多个权限组合成权限组
4. **审计日志**: 记录权限检查和权限变更的审计日志
