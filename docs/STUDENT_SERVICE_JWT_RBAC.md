# Student服务 JWT和RBAC功能说明

## 概述

Student微服务已经成功集成了JWT（JSON Web Token）认证和RBAC（基于角色的访问控制）功能，确保API的安全性和权限管理。

## 功能特性

### 1. JWT认证
- **自动验证**: 所有API请求都会自动验证JWT token
- **跳过路径**: 健康检查等公开API可以跳过JWT验证
- **Token格式**: 支持Bearer token格式
- **配置化**: JWT密钥和过期时间可通过配置文件设置

### 2. RBAC权限控制
- **基于角色**: 用户通过角色获得权限
- **资源控制**: 对不同的API资源进行权限控制
- **动态检查**: 实时检查用户是否有访问特定资源的权限
- **Casbin集成**: 使用Casbin作为权限引擎

### 3. 中间件集成
- **JWT中间件**: 自动提取和验证JWT token
- **RBAC中间件**: 检查用户权限
- **错误处理**: 统一的错误响应格式

## 配置说明

### JWT配置
```yaml
jwt:
  secret_key: "your-secret-key-here-make-it-long-and-secure"
  expire: 86400s  # 24小时
```

### RBAC配置
```yaml
rbac:
  model_path: "rbac_model.conf"
  policy_path: "rbac_policy.csv"
  enabled: true
```

## API端点

### 公开端点（无需认证）
- `GET /v1/students/health` - 健康检查
- `GET /health` - 基础健康检查

### 受保护端点（需要JWT认证和RBAC权限）
- `GET /v1/student/{id}` - 获取学生信息
- `POST /v1/student` - 创建学生
- `PUT /v1/student/{id}` - 更新学生信息
- `DELETE /v1/student/{id}` - 删除学生
- `GET /v1/students` - 获取学生列表

## 权限配置示例

### RBAC模型配置 (rbac_model.conf)
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

### 权限策略配置 (rbac_policy.csv)
```csv
p, admin, /v1/student/*, GET
p, admin, /v1/student/*, POST
p, admin, /v1/student/*, PUT
p, admin, /v1/student/*, DELETE
p, teacher, /v1/student/*, GET
p, teacher, /v1/students, GET
p, student, /v1/student/*, GET
g, user1, admin
g, user2, teacher
g, user3, student
```

## 使用示例

### 1. 获取JWT Token
```bash
# 通过用户服务登录获取token
curl -X POST "http://localhost:8601/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

### 2. 使用Token访问API
```bash
# 使用获取到的token访问学生API
curl -X GET "http://localhost:8602/v1/student/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. 健康检查
```bash
# 无需token的健康检查
curl -X GET "http://localhost:8602/v1/students/health"
```

## 错误响应

### 401 Unauthorized
```json
{
  "code": 401,
  "reason": "UNAUTHORIZED",
  "message": "未提供有效的认证token"
}
```

### 403 Forbidden
```json
{
  "code": 403,
  "reason": "FORBIDDEN",
  "message": "没有访问权限"
}
```

## 部署说明

### 1. 启动依赖服务
```bash
# 启动数据库和Redis
docker-compose up -d mysql redis

# 启动RBAC服务
./bin/rbac-service -conf configs/rbac-service.yaml

# 启动用户服务
./bin/user-service -conf configs/user-service.yaml
```

### 2. 启动Student服务
```bash
./bin/student-service -conf configs/student-service.yaml
```

### 3. 测试功能
```bash
# 运行测试脚本
chmod +x test-student-service.sh
./test-student-service.sh
```

## 安全建议

1. **JWT密钥**: 使用强密钥，定期更换
2. **Token过期**: 设置合理的过期时间
3. **HTTPS**: 生产环境使用HTTPS
4. **权限最小化**: 遵循最小权限原则
5. **日志监控**: 记录认证和授权日志

## 故障排除

### 常见问题

1. **401错误**: 检查JWT token是否有效
2. **403错误**: 检查用户是否有相应权限
3. **服务启动失败**: 检查数据库和Redis连接
4. **权限不生效**: 检查RBAC配置文件

### 调试方法

1. 查看服务日志
2. 检查JWT token内容
3. 验证RBAC策略配置
4. 测试健康检查端点

## 总结

Student服务现在具备了完整的JWT认证和RBAC权限控制功能，可以安全地处理学生相关的API请求。通过合理的权限配置，可以确保不同角色的用户只能访问其被授权的资源。
