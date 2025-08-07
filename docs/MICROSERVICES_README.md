# 学生管理系统微服务架构

## 概述

本项目已成功改造为基于微服务架构的学生管理系统，使用 Nacos 作为服务注册与配置中心。

## 架构设计

### 服务拆分

1. **user-service** (用户服务)

   - 端口: HTTP 8601, gRPC 9601
   - 功能: 用户管理、认证、注册登录

2. **student-service** (学生服务)

   - 端口: HTTP 8602, gRPC 9602
   - 功能: 学生信息管理

3. **rbac-service** (权限服务)

   - 端口: HTTP 8603, gRPC 9603
   - 功能: 角色权限管理

4. **gateway-service** (API 网关)
   - 端口: HTTP 8600, gRPC 9600
   - 功能: 统一入口、路由转发、负载均衡

### 技术栈

- **框架**: Kratos v2
- **服务注册**: Nacos
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **容器化**: Docker + Docker Compose
- **API**: gRPC + HTTP

## 快速开始

### 环境要求

- Go 1.23+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 6.0+

### 本地开发

1. **安装依赖**

```bash
make init
```

2. **生成代码**

```bash
make api
make wire
```

3. **构建微服务**

```bash
make build-microservices
```

4. **启动 Nacos**

```bash
# 使用Docker启动Nacos
docker run -d --name nacos-standalone \
  -p 8848:8848 \
  -p 9848:9848 \
  -e MODE=standalone \
  nacos/nacos-server:v2.2.3
```

5. **启动微服务**

```bash
# 启动网关服务
./bin/gateway-service -conf ./configs/gateway-service.yaml

# 启动用户服务
./bin/user-service -conf ./configs/user-service.yaml

# 启动学生服务
./bin/student-service -conf ./configs/student-service.yaml

# 启动RBAC服务
./bin/rbac-service -conf ./configs/rbac-service.yaml
```

### Docker 部署

使用一键部署脚本：

```bash
./deploy-microservices.sh
```

或手动部署：

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]
```

## 服务访问

### 服务地址

- **Nacos 控制台**: http://localhost:8848/nacos (用户名/密码: nacos/nacos)
- **API 网关**: http://localhost:8600
- **用户服务**: http://localhost:8601
- **学生服务**: http://localhost:8602
- **RBAC 服务**: http://localhost:8603

### API 接口

所有 API 通过网关统一访问：

```
# 用户管理
POST   /v1/user/login          # 用户登录
POST   /v1/user/register       # 用户注册
GET    /v1/user/me             # 获取当前用户信息
GET    /v1/users               # 获取用户列表
POST   /v1/user                # 创建用户
GET    /v1/user/{id}           # 获取用户详情
PUT    /v1/user/{id}           # 更新用户
DELETE /v1/user/{id}           # 删除用户

# 学生管理
GET    /v1/students            # 获取学生列表
POST   /v1/student             # 创建学生
GET    /v1/student/{id}        # 获取学生详情
PUT    /v1/student/{id}        # 更新学生
DELETE /v1/student/{id}        # 删除学生

# 权限管理
GET    /v1/roles               # 获取角色列表
POST   /v1/roles               # 创建角色
GET    /v1/permissions         # 获取权限列表
POST   /v1/permissions/check   # 权限检查
```

## 配置说明

### Nacos 配置

每个服务都有独立的配置文件：

- `configs/user-service.yaml` - 用户服务配置
- `configs/student-service.yaml` - 学生服务配置
- `configs/rbac-service.yaml` - RBAC 服务配置
- `configs/gateway-service.yaml` - 网关服务配置

### 服务注册配置

```yaml
nacos:
  discovery:
    ip: "192.168.56.162" # Nacos服务器IP
    port: 8848 # Nacos服务器端口
    namespace_id: "public" # 命名空间
    group: "DEFAULT_GROUP" # 分组
    cluster_name: "DEFAULT" # 集群名称
    weight: 10 # 权重
    metadata:
      version: "1.0.0" # 版本
      zone: "zone1" # 区域
```

## 监控与管理

### 健康检查

每个服务都提供健康检查接口：

```bash
curl http://localhost:8600/health  # 网关服务
curl http://localhost:8601/health  # 用户服务
curl http://localhost:8602/health  # 学生服务
curl http://localhost:8603/health  # RBAC服务
```

### 服务发现

通过 Nacos 控制台可以查看所有注册的服务实例：

1. 访问 http://localhost:8848/nacos
2. 使用 nacos/nacos 登录
3. 在"服务管理"中查看服务列表

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f user-service
docker-compose logs -f student-service
docker-compose logs -f rbac-service
docker-compose logs -f gateway-service
```

## 扩展性

### 水平扩展

可以通过增加服务实例来实现水平扩展：

```bash
# 扩展用户服务实例
docker-compose up -d --scale user-service=3

# 扩展学生服务实例
docker-compose up -d --scale student-service=2
```

### 负载均衡

网关服务会自动进行负载均衡，将请求分发到不同的服务实例。

## 故障排除

### 常见问题

1. **服务无法注册到 Nacos**

   - 检查 Nacos 服务是否正常运行
   - 检查网络连接
   - 检查配置文件中的 Nacos 地址

2. **服务间无法通信**

   - 检查服务是否正常注册
   - 检查网络配置
   - 查看服务日志

3. **数据库连接失败**
   - 检查 MySQL 服务状态
   - 检查数据库连接配置
   - 检查网络连接

### 调试命令

```bash
# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]

# 进入容器调试
docker-compose exec [service-name] sh

# 重启服务
docker-compose restart [service-name]

# 停止所有服务
docker-compose down
```

## 开发指南

### 添加新服务

1. 创建服务目录结构
2. 添加配置文件
3. 实现服务逻辑
4. 注册到 Nacos
5. 更新网关路由

### 服务间调用

使用 Nacos 服务发现进行服务间调用：

```go
// 获取服务实例
instances, err := discovery.GetServiceInstances("service-name")
if err != nil {
    return err
}

// 选择实例进行调用
instance := instances[0]
targetURL := fmt.Sprintf("http://%s", instance.GetServiceURL())
```

## 总结

通过微服务架构改造，系统具备了以下优势：

1. **高可用性**: 服务独立部署，故障隔离
2. **可扩展性**: 支持水平扩展和垂直扩展
3. **可维护性**: 服务职责单一，便于维护
4. **技术多样性**: 不同服务可以使用不同的技术栈
5. **团队协作**: 不同团队可以独立开发不同服务

微服务架构为系统的长期发展奠定了坚实的基础。
