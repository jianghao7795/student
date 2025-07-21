# Casbin RBAC 模型配置快速参考

## 配置文件结构

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

## 参数说明

| 参数  | 含义              | 示例                     |
| ----- | ----------------- | ------------------------ |
| `sub` | 主体（用户/角色） | `"1"`, `"admin"`         |
| `obj` | 对象（资源路径）  | `"/api/v1/users"`        |
| `act` | 动作（HTTP 方法） | `"GET"`, `"POST"`, `"*"` |

## 配置段说明

### [request_definition]

- **作用**: 定义权限检查请求格式
- **格式**: `r = sub, obj, act`
- **示例**: `r = "1", "/api/v1/users", "GET"`

### [policy_definition]

- **作用**: 定义权限策略格式
- **格式**: `p = sub, obj, act`
- **示例**: `p = "admin", "/api/v1/users", "*"`

### [role_definition]

- **作用**: 定义用户角色关系
- **格式**: `g = user, role`
- **示例**: `g = "1", "admin"`

### [policy_effect]

- **作用**: 定义权限决策规则
- **当前**: 允许优先（白名单）
- **说明**: 只要有一个策略允许就允许访问

### [matchers]

- **作用**: 定义权限匹配规则
- **逻辑**: 角色匹配 && 资源匹配 && 动作匹配

## 匹配器详解

```conf
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

### 1. 角色匹配: `g(r.sub, p.sub)`

- 检查用户是否拥有指定角色
- 通过角色分配表验证

### 2. 资源匹配: `keyMatch2(r.obj, p.obj)`

- 支持通配符的路径匹配
- 支持 `*`, `{id}`, `[id]` 等模式

### 3. 动作匹配: `(r.act == p.act || p.act == "*")`

- 检查 HTTP 方法是否匹配
- 支持 `*` 表示所有方法

## 通配符说明

| 通配符 | 含义             | 示例                 |
| ------ | ---------------- | -------------------- |
| `*`    | 匹配任意字符序列 | `/api/v1/users/*`    |
| `{id}` | 匹配路径参数     | `/api/v1/users/{id}` |
| `[id]` | 匹配可选参数     | `/api/v1/users/[id]` |

## 权限检查流程

```
1. 解析请求: r = (用户, 资源, 动作)
2. 查找角色: g(用户, 角色)
3. 匹配策略: p(角色, 资源, 动作)
4. 应用效果: 允许/拒绝
```

## 常见策略示例

### 角色分配

```sql
g, 1, admin    -- 用户1有admin角色
g, 2, user     -- 用户2有user角色
g, 3, guest    -- 用户3有guest角色
```

### 权限策略

```sql
p, admin, /api/v1/users, *      -- admin可以访问所有用户操作
p, user, /api/v1/users, GET     -- user只能GET用户列表
p, guest, /api/v1/users, GET    -- guest只能GET用户列表
```

## 测试示例

| 用户 | 角色  | 请求                 | 结果    |
| ---- | ----- | -------------------- | ------- |
| 1    | admin | GET /api/v1/users    | ✅ 允许 |
| 1    | admin | POST /api/v1/users   | ✅ 允许 |
| 2    | user  | GET /api/v1/users    | ✅ 允许 |
| 2    | user  | POST /api/v1/users   | ❌ 拒绝 |
| 3    | guest | GET /api/v1/users    | ✅ 允许 |
| 3    | guest | DELETE /api/v1/users | ❌ 拒绝 |

## 常用配置模式

### 1. 严格模式（默认）

```conf
[policy_effect]
e = some(where (p.eft == allow))
```

### 2. 拒绝优先模式

```conf
[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
```

### 3. 优先级模式

```conf
[policy_effect]
e = priority(p.eft) || deny
```

## 性能优化建议

1. **索引优化**: 为角色和策略表添加索引
2. **缓存策略**: 缓存权限检查结果
3. **批量操作**: 批量加载策略数据
4. **定期清理**: 清理无效的策略数据

## 安全建议

1. **最小权限**: 遵循最小权限原则
2. **定期审计**: 定期审查权限分配
3. **版本控制**: 对策略变更进行版本控制
4. **日志记录**: 记录权限检查日志

## 故障排除

### 常见问题

1. **权限被拒绝**

   - 检查用户是否有对应角色
   - 检查角色是否有对应权限
   - 检查资源路径是否匹配

2. **性能问题**

   - 检查策略数据量
   - 优化数据库查询
   - 启用缓存机制

3. **配置错误**
   - 检查配置文件语法
   - 验证策略数据格式
   - 检查匹配器逻辑
