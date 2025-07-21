# RBAC 策略文件快速参考

## 文件格式

```csv
规则类型, 参数1, 参数2, 参数3
```

## 规则类型

### 权限策略 (Policy)

```csv
p, 角色, 资源路径, 操作
```

### 角色分配 (Grouping)

```csv
g, 用户ID, 角色名
```

## 当前配置

```csv
# 权限策略
p, admin, /api/v1/users, *
p, admin, /api/v1/students, *
p, user, /api/v1/users, GET
p, user, /api/v1/students, GET
p, guest, /api/v1/users, GET

# 角色分配
g, 1, admin
g, 2, user
g, 3, guest
```

## 常用操作

| 操作     | 说明     |
| -------- | -------- |
| `GET`    | 获取资源 |
| `POST`   | 创建资源 |
| `PUT`    | 更新资源 |
| `DELETE` | 删除资源 |
| `*`      | 所有操作 |

## 权限矩阵

| 用户 | 角色  | 用户 API | 学生 API |
| ---- | ----- | -------- | -------- |
| 1    | admin | 所有操作 | 所有操作 |
| 2    | user  | 仅 GET   | 仅 GET   |
| 3    | guest | 仅 GET   | 无权限   |

## 常见配置模式

### 基础权限

```csv
# 管理员
p, admin, /api/v1/users, *
p, admin, /api/v1/students, *

# 普通用户
p, user, /api/v1/users, GET
p, user, /api/v1/students, GET

# 访客
p, guest, /api/v1/users, GET
```

### 细粒度权限

```csv
# 用户管理
p, admin, /api/v1/users, GET
p, admin, /api/v1/users, POST
p, admin, /api/v1/users, PUT
p, admin, /api/v1/users, DELETE

# 学生管理
p, manager, /api/v1/students, GET
p, manager, /api/v1/students, POST
p, manager, /api/v1/students, PUT
```

### 路径参数权限

```csv
p, user, /api/v1/users/{id}, GET
p, user, /api/v1/users/{id}, PUT
p, admin, /api/v1/users/{id}, DELETE
```

## 验证命令

### 语法检查

```bash
cat rbac_policy.csv | grep -v "^#" | grep -v "^$" | while read line; do
    echo "$line" | awk -F',' '{print NF}' | xargs -I {} test {} -ge 3 || echo "Invalid: $line"
done
```

### 权限测试

```bash
curl -X POST http://localhost:8000/api/v1/permissions/check \
  -H "Content-Type: application/json" \
  -d '{"user": "1", "resource": "/api/v1/users", "action": "GET"}'
```

## 故障排除

### 权限不生效

```bash
# 重新加载策略
curl -X POST http://localhost:8000/api/v1/rbac/policies/reload

# 检查用户角色
curl -X GET http://localhost:8000/api/v1/users/1/roles
```

### 策略检查

```bash
# 查看当前策略
curl -X GET http://localhost:8000/api/v1/rbac/policies

# 检查日志
tail -f logs/app.log | grep "casbin"
```

## 安全建议

1. **最小权限**: 只授予必要权限
2. **定期审查**: 定期检查权限分配
3. **版本控制**: 使用 Git 管理策略变更
4. **备份**: 定期备份策略文件

## 相关文件

- `configs/rbac_model.conf`: 模型配置
- `migrate/rbac_migrate.sql`: 数据库迁移
- `docs/rbac_policy_documentation.md`: 详细文档
