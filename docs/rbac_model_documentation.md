# Casbin RBAC 模型配置文件文档

## 文件概述

`configs/rbac_model.conf` 是 Casbin 权限管理引擎的核心配置文件，定义了基于角色的访问控制（RBAC）模型的规则和行为。该文件决定了系统如何进行权限检查和决策。

## 配置文件内容

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

## 详细说明

### 1. [request_definition] - 请求定义

**语法**: `r = sub, obj, act`

**作用**: 定义权限检查请求的格式和参数

**参数说明**:

- **r**: 请求标识符（request）
- **sub**: 主体（subject），通常是用户 ID 或用户名
- **obj**: 对象（object），通常是资源路径或 API 端点
- **act**: 动作（action），通常是 HTTP 方法或操作类型

**示例**:

```
r = "1", "/api/v1/users", "GET"
```

表示：用户 1 请求 GET 访问用户列表 API

### 2. [policy_definition] - 策略定义

**语法**: `p = sub, obj, act`

**作用**: 定义权限策略的格式和参数

**参数说明**:

- **p**: 策略标识符（policy）
- **sub**: 主体，通常是角色名称
- **obj**: 对象，资源路径
- **act**: 动作，HTTP 方法或操作类型

**示例**:

```
p = "admin", "/api/v1/users", "*"
p = "user", "/api/v1/users", "GET"
```

表示：

- admin 角色可以访问用户相关的所有操作
- user 角色只能 GET 访问用户列表

### 3. [role_definition] - 角色定义

**语法**: `g = _, _`

**作用**: 定义用户与角色的关联关系

**参数说明**:

- **g**: 分组标识符（grouping）
- 第一个参数：用户标识符
- 第二个参数：角色名称

**示例**:

```
g = "1", "admin"
g = "2", "user"
```

表示：

- 用户 1 被分配了 admin 角色
- 用户 2 被分配了 user 角色

### 4. [policy_effect] - 策略效果

**语法**: `e = some(where (p.eft == allow))`

**作用**: 定义权限决策的规则

**说明**:

- **some(where (p.eft == allow))**: 只要有一个策略允许，就允许访问
- 这是"允许优先"的策略，即白名单机制
- 如果没有匹配的允许策略，则拒绝访问

**其他常见效果**:

- `e = some(where (p.eft == allow)) && !some(where (p.eft == deny))` - 允许优先，但明确拒绝优先
- `e = priority(p.eft) || deny` - 优先级策略

### 5. [matchers] - 匹配器

**语法**: `m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")`

**作用**: 定义权限匹配的具体规则

**匹配逻辑**:

1. **角色匹配**: `g(r.sub, p.sub)`

   - 检查请求用户是否拥有策略中指定的角色
   - 通过角色分配表验证用户角色

2. **资源匹配**: `keyMatch2(r.obj, p.obj)`

   - 使用 keyMatch2 函数进行路径匹配
   - 支持通配符和模式匹配

3. **动作匹配**: `(r.act == p.act || p.act == "*")`
   - 检查请求动作是否匹配策略动作
   - 支持通配符"\*"表示所有动作

## 匹配器函数说明

### keyMatch2 函数

**功能**: 支持通配符的路径匹配

**语法**: `keyMatch2(key1, key2)`

**支持的通配符**:

- `*`: 匹配任意字符序列
- `{id}`: 匹配路径参数
- `[id]`: 匹配可选路径参数

**示例**:

```
keyMatch2("/api/v1/users/123", "/api/v1/users/*")     → true
keyMatch2("/api/v1/users/123", "/api/v1/users/{id}")  → true
keyMatch2("/api/v1/users", "/api/v1/users")           → true
keyMatch2("/api/v1/users", "/api/v1/students")        → false
```

## 权限检查流程

### 1. 请求解析

```
请求: 用户1 GET /api/v1/users
解析: r = ("1", "/api/v1/users", "GET")
```

### 2. 角色验证

```
检查: g("1", "admin") → true (用户1有admin角色)
```

### 3. 策略匹配

```
策略: p = ("admin", "/api/v1/users", "*")
匹配:
  - 角色匹配: g("1", "admin") → true
  - 资源匹配: keyMatch2("/api/v1/users", "/api/v1/users") → true
  - 动作匹配: ("GET" == "*" || "*" == "*") → true
```

### 4. 权限决策

```
结果: 所有条件满足 → 允许访问
```

## 实际应用示例

### 数据库策略示例

```sql
-- 角色分配
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
('g', '1', 'admin'),
('g', '2', 'user'),
('g', '3', 'guest');

-- 权限策略
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'admin', '/api/v1/users', '*'),
('p', 'admin', '/api/v1/students', '*'),
('p', 'user', '/api/v1/users', 'GET'),
('p', 'user', '/api/v1/students', 'GET'),
('p', 'guest', '/api/v1/users', 'GET');
```

### 权限检查示例

| 用户 | 角色  | 请求                 | 结果    | 说明                   |
| ---- | ----- | -------------------- | ------- | ---------------------- |
| 1    | admin | GET /api/v1/users    | ✅ 允许 | admin 有所有权限       |
| 1    | admin | POST /api/v1/users   | ✅ 允许 | admin 有所有权限       |
| 2    | user  | GET /api/v1/users    | ✅ 允许 | user 有 GET 权限       |
| 2    | user  | POST /api/v1/users   | ❌ 拒绝 | user 没有 POST 权限    |
| 3    | guest | GET /api/v1/users    | ✅ 允许 | guest 有 GET 权限      |
| 3    | guest | DELETE /api/v1/users | ❌ 拒绝 | guest 没有 DELETE 权限 |

## 配置优化建议

### 1. 性能优化

- 使用索引优化角色查询
- 缓存权限策略
- 定期清理无效策略

### 2. 安全性增强

- 添加审计日志
- 实现策略版本控制
- 支持策略继承

### 3. 灵活性提升

- 支持动态策略加载
- 实现策略优先级
- 添加条件策略

## 常见问题

### Q1: 如何添加新的匹配函数？

A: 可以在 matchers 中使用自定义函数，如：

```conf
[matchers]
m = g(r.sub, p.sub) && customMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

### Q2: 如何实现角色继承？

A: 可以通过添加角色间的关联关系实现：

```conf
g = "admin", "super_admin"  # admin继承super_admin的权限
```

### Q3: 如何支持更复杂的权限逻辑？

A: 可以在 matchers 中添加更多条件：

```conf
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*") && timeCheck()
```

## 总结

`rbac_model.conf` 文件是 Casbin 权限系统的核心配置，通过定义请求格式、策略规则、角色关系和匹配逻辑，实现了灵活而强大的 RBAC 权限控制。合理配置该文件可以满足大多数应用的权限管理需求。
