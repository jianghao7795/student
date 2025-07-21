# RBAC 策略文件文档总览

## 概述

本文档集合涵盖了 `rbac_policy.csv` 文件的完整使用指南，包括详细文档、快速参考、配置模板和最佳实践。

## 文档列表

### 📖 详细文档

- **[rbac_policy_documentation.md](./rbac_policy_documentation.md)** - 完整的配置文档
  - 文件格式和语法说明
  - 配置详解和示例
  - 最佳实践和安全建议
  - 故障排除和调试方法

### ⚡ 快速参考

- **[rbac_policy_quick_reference.md](./rbac_policy_quick_reference.md)** - 快速参考指南
  - 常用配置模式
  - 验证和测试命令
  - 故障排除步骤

### 📋 配置模板

- **[rbac_policy_template.csv](./rbac_policy_template.csv)** - 配置模板文件
  - 完整的配置示例
  - 不同角色的权限配置
  - 高级配置选项

## 文件结构

```
docs/
├── rbac_policy_documentation.md      # 详细文档
├── rbac_policy_quick_reference.md    # 快速参考
├── rbac_policy_template.csv          # 配置模板
└── RBAC_POLICY_README.md             # 本文档
```

## 快速开始

### 1. 了解文件格式

`rbac_policy.csv` 使用 CSV 格式，包含两种规则类型：

```csv
# 权限策略
p, 角色, 资源路径, 操作

# 角色分配
g, 用户ID, 角色名
```

### 2. 当前配置

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

### 3. 权限矩阵

| 用户 | 角色  | 用户 API | 学生 API |
| ---- | ----- | -------- | -------- |
| 1    | admin | 所有操作 | 所有操作 |
| 2    | user  | 仅 GET   | 仅 GET   |
| 3    | guest | 仅 GET   | 无权限   |

## 常用操作

### 验证配置

```bash
# 语法检查
cat rbac_policy.csv | grep -v "^#" | grep -v "^$" | while read line; do
    echo "$line" | awk -F',' '{print NF}' | xargs -I {} test {} -ge 3 || echo "Invalid: $line"
done

# 权限测试
curl -X POST http://localhost:8000/api/v1/permissions/check \
  -H "Content-Type: application/json" \
  -d '{"user": "1", "resource": "/api/v1/users", "action": "GET"}'
```

### 重新加载策略

```bash
# 重新加载策略
curl -X POST http://localhost:8000/api/v1/rbac/policies/reload

# 检查用户角色
curl -X GET http://localhost:8000/api/v1/users/1/roles
```

## 配置最佳实践

### 1. 权限设计原则

- **最小权限原则**: 只授予必要权限
- **角色分离原则**: 不同角色承担不同职责
- **权限继承原则**: 保持权限层次清晰

### 2. 命名规范

```csv
# 推荐的角色命名
admin          # 系统管理员
manager        # 部门经理
user           # 普通用户
guest          # 访客

# 推荐的资源路径
/api/v1/users          # 用户管理
/api/v1/students       # 学生管理
/api/v1/roles          # 角色管理
```

### 3. 配置结构

```csv
# 按功能分组
# 用户管理权限
p, admin, /api/v1/users, *
p, user, /api/v1/users, GET
p, guest, /api/v1/users, GET

# 学生管理权限
p, admin, /api/v1/students, *
p, user, /api/v1/students, GET

# 角色分配
g, 1, admin
g, 2, user
g, 3, guest
```

## 常见配置模式

### 基础权限模式

```csv
# 管理员权限
p, admin, /api/v1/users, *
p, admin, /api/v1/students, *
p, admin, /api/v1/roles, *
p, admin, /api/v1/permissions, *

# 普通用户权限
p, user, /api/v1/users, GET
p, user, /api/v1/students, GET

# 访客权限
p, guest, /api/v1/users, GET
```

### 细粒度权限模式

```csv
# 用户管理细粒度权限
p, admin, /api/v1/users, GET
p, admin, /api/v1/users, POST
p, admin, /api/v1/users, PUT
p, admin, /api/v1/users, DELETE

# 学生管理细粒度权限
p, manager, /api/v1/students, GET
p, manager, /api/v1/students, POST
p, manager, /api/v1/students, PUT
```

## 故障排除

### 常见问题

1. **权限检查失败**

   - 检查策略是否正确加载
   - 验证用户角色分配
   - 查看系统日志

2. **策略不生效**

   - 重新加载策略
   - 检查文件格式
   - 验证配置语法

3. **系统启动失败**
   - 检查文件是否存在
   - 验证文件权限
   - 查看错误日志

### 调试方法

```bash
# 启用调试模式
# configs/config.yaml
rbac:
  enabled: true
  debug: true

# 查看权限检查日志
tail -f logs/app.log | grep "casbin"

# 检查当前策略
curl -X GET http://localhost:8000/api/v1/rbac/policies
```

## 安全建议

### 1. 文件安全

- 设置适当的文件权限：`chmod 600 rbac_policy.csv`
- 定期备份策略文件
- 使用版本控制管理策略变更

### 2. 权限审计

- 定期审查权限分配
- 记录权限变更日志
- 实施权限最小化原则

### 3. 测试验证

- 编写权限测试用例
- 定期进行权限渗透测试
- 验证权限配置的有效性

## 扩展功能

### 1. 条件权限

```csv
# 支持时间条件的权限（需要自定义匹配器）
p, user, /api/v1/students, GET, 9:00-18:00
p, guest, /api/v1/users, GET, 8:00-20:00
```

### 2. 域权限

```csv
# 支持多租户的域权限
p, admin, domain1, /api/v1/users, *
p, user, domain2, /api/v1/users, GET
```

### 3. 优先级权限

```csv
# 支持权限优先级
p, admin, /api/v1/users, *, 1
p, user, /api/v1/users, GET, 2
```

## 相关文档

- [RBAC 系统使用说明](./RBAC_README.md)
- [RBAC 模型配置文档](./rbac_model_documentation.md)
- [RBAC 实现总结](./RBAC_IMPLEMENTATION_SUMMARY.md)

## 总结

`rbac_policy.csv` 文件是 RBAC 权限系统的核心配置文件，通过合理配置可以实现灵活的访问控制。建议：

1. **仔细阅读详细文档**了解完整功能
2. **使用快速参考**进行日常操作
3. **参考配置模板**创建新的配置
4. **遵循最佳实践**确保安全性
5. **定期审查和更新**权限配置

如有问题，请参考相关文档或联系系统管理员。
