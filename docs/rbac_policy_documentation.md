# RBAC 策略配置文件文档

## 文件概述

`rbac_policy.csv` 是 Casbin 权限管理引擎的策略配置文件，定义了基于角色的访问控制（RBAC）系统的具体权限规则。该文件是 RBAC 权限系统的核心配置，决定了"谁可以做什么"的权限策略。

## 文件位置

```
student/
└── rbac_policy.csv
```

## 文件格式

### 基本语法

文件采用 CSV 格式，每行表示一条策略规则，字段之间用逗号分隔：

```csv
规则类型, 参数1, 参数2, 参数3, ...
```

### 规则类型

#### 1. 权限策略 (Policy Rules) - 以 `p` 开头

**格式**: `p, 角色, 资源路径, 操作`

**作用**: 定义某个角色对特定资源的操作权限

**示例**:

```csv
p, admin, /api/v1/users, *
p, user, /api/v1/users, GET
p, guest, /api/v1/students, GET
```

#### 2. 角色分配 (Grouping Rules) - 以 `g` 开头

**格式**: `g, 用户ID, 角色名`

**作用**: 定义用户与角色的关联关系

**示例**:

```csv
g, 1, admin
g, 2, user
g, 3, guest
```

## 当前配置详解

### 权限策略配置

```csv
p, admin, /api/v1/users, *
p, admin, /api/v1/students, *
p, user, /api/v1/users, GET
p, user, /api/v1/students, GET
p, guest, /api/v1/users, GET
```

#### 详细说明

| 规则                             | 角色  | 资源               | 操作  | 说明                                 |
| -------------------------------- | ----- | ------------------ | ----- | ------------------------------------ |
| `p, admin, /api/v1/users, *`     | admin | `/api/v1/users`    | `*`   | admin 角色可以访问用户相关的所有操作 |
| `p, admin, /api/v1/students, *`  | admin | `/api/v1/students` | `*`   | admin 角色可以访问学生相关的所有操作 |
| `p, user, /api/v1/users, GET`    | user  | `/api/v1/users`    | `GET` | user 角色只能 GET 访问用户列表       |
| `p, user, /api/v1/students, GET` | user  | `/api/v1/students` | `GET` | user 角色只能 GET 访问学生列表       |
| `p, guest, /api/v1/users, GET`   | guest | `/api/v1/users`    | `GET` | guest 角色只能 GET 访问用户列表      |

### 角色分配配置

```csv
g, 1, admin
g, 2, user
g, 3, guest
```

#### 详细说明

| 规则          | 用户 ID | 角色  | 说明                                   |
| ------------- | ------- | ----- | -------------------------------------- |
| `g, 1, admin` | 1       | admin | 用户 ID 为 1 的用户被分配了 admin 角色 |
| `g, 2, user`  | 2       | user  | 用户 ID 为 2 的用户被分配了 user 角色  |
| `g, 3, guest` | 3       | guest | 用户 ID 为 3 的用户被分配了 guest 角色 |

## 权限矩阵

基于当前配置的权限矩阵：

| 用户 ID | 角色  | 用户 API | 学生 API | 说明         |
| ------- | ----- | -------- | -------- | ------------ |
| 1       | admin | 所有操作 | 所有操作 | 完全权限     |
| 2       | user  | 仅 GET   | 仅 GET   | 只读权限     |
| 3       | guest | 仅 GET   | 无权限   | 有限查看权限 |

## 操作类型说明

### HTTP 方法映射

| 操作     | HTTP 方法 | 说明               |
| -------- | --------- | ------------------ |
| `GET`    | GET       | 获取资源列表或详情 |
| `POST`   | POST      | 创建新资源         |
| `PUT`    | PUT       | 更新资源           |
| `DELETE` | DELETE    | 删除资源           |
| `*`      | 所有方法  | 支持所有 HTTP 方法 |

### 通配符说明

| 通配符 | 含义                 | 示例                                    |
| ------ | -------------------- | --------------------------------------- |
| `*`    | 匹配所有操作         | `/api/v1/users, *` 表示所有用户相关操作 |
| `/*`   | 匹配路径下的所有资源 | `/api/v1/users/*` 表示所有用户资源      |

## 配置最佳实践

### 1. 权限设计原则

#### 最小权限原则

- 只授予用户完成工作所需的最小权限
- 避免过度授权

#### 角色分离原则

- 不同角色承担不同职责
- 避免角色权限重叠

#### 权限继承原则

- 高级角色继承低级角色权限
- 保持权限层次清晰

### 2. 命名规范

#### 角色命名

```csv
# 推荐
admin          # 系统管理员
manager        # 部门经理
user           # 普通用户
guest          # 访客

# 不推荐
role1          # 无意义的名称
admin_role     # 冗余后缀
```

#### 资源路径命名

```csv
# 推荐
/api/v1/users          # 用户管理
/api/v1/students       # 学生管理
/api/v1/roles          # 角色管理

# 不推荐
/api/users             # 缺少版本号
/api/v1/user           # 单数形式不一致
```

### 3. 配置结构

#### 按功能分组

```csv
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

### 1. 基础权限模式

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

### 2. 细粒度权限模式

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

### 3. 路径参数权限模式

```csv
# 支持路径参数的权限
p, user, /api/v1/users/{id}, GET
p, user, /api/v1/users/{id}, PUT
p, admin, /api/v1/users/{id}, DELETE
```

## 配置验证

### 1. 语法检查

```bash
# 检查 CSV 格式
cat rbac_policy.csv | grep -v "^#" | grep -v "^$" | while read line; do
    echo "$line" | awk -F',' '{print NF}' | xargs -I {} test {} -ge 3 || echo "Invalid line: $line"
done
```

### 2. 权限测试

```bash
# 测试权限检查
curl -X POST http://localhost:8000/api/v1/permissions/check \
  -H "Content-Type: application/json" \
  -d '{
    "user": "1",
    "resource": "/api/v1/users",
    "action": "GET"
  }'
```

## 动态管理

### 1. 通过 API 管理

```bash
# 添加权限策略
curl -X POST http://localhost:8000/api/v1/rbac/policies \
  -H "Content-Type: application/json" \
  -d '{
    "role": "manager",
    "resource": "/api/v1/students",
    "action": "POST"
  }'

# 分配用户角色
curl -X POST http://localhost:8000/api/v1/users/4/roles \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 2
  }'
```

### 2. 数据库同步

策略会自动同步到数据库的 `casbin_rule` 表中：

```sql
-- 查看当前策略
SELECT * FROM casbin_rule;

-- 手动添加策略
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'manager', '/api/v1/students', 'POST'),
('g', '4', 'manager');
```

## 故障排除

### 1. 常见问题

#### 权限检查失败

```bash
# 检查策略是否正确加载
curl -X GET http://localhost:8000/api/v1/rbac/policies

# 检查用户角色
curl -X GET http://localhost:8000/api/v1/users/1/roles
```

#### 策略不生效

```bash
# 重新加载策略
curl -X POST http://localhost:8000/api/v1/rbac/policies/reload

# 检查日志
tail -f logs/app.log | grep "casbin"
```

### 2. 调试方法

#### 启用调试模式

```yaml
# configs/config.yaml
rbac:
  enabled: true
  debug: true
```

#### 查看权限检查日志

```bash
# 查看权限检查详情
curl -X POST http://localhost:8000/api/v1/permissions/check \
  -H "Content-Type: application/json" \
  -H "X-Debug: true" \
  -d '{
    "user": "1",
    "resource": "/api/v1/users",
    "action": "GET"
  }'
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

## 扩展配置

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

## 总结

`rbac_policy.csv` 文件是 RBAC 权限系统的核心配置文件，通过定义权限策略和角色分配，实现了灵活的访问控制。合理配置该文件可以确保系统的安全性和可用性。

### 关键要点

1. **文件格式**: CSV 格式，支持权限策略和角色分配
2. **权限粒度**: 支持资源级别的细粒度权限控制
3. **动态管理**: 支持通过 API 和数据库进行动态管理
4. **安全原则**: 遵循最小权限原则和角色分离原则
5. **维护性**: 定期审查和更新权限配置

### 相关文件

- `configs/rbac_model.conf`: Casbin 模型配置文件
- `migrate/rbac_migrate.sql`: RBAC 数据库迁移文件
- `docs/RBAC_README.md`: RBAC 系统使用说明
