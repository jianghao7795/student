# Student Management System

基于 Kratos 框架开发的学生管理系统，包含用户管理、学生管理、RBAC 权限控制等功能。

## 项目结构

```
student/
├── api/                    # API 定义文件 (protobuf)
│   ├── student/v1/        # 学生服务 API
│   ├── user/v1/           # 用户服务 API
│   ├── rbac/v1/           # RBAC 权限服务 API
│   └── errors/v1/         # 错误处理服务 API
├── cmd/                   # 应用程序入口
├── configs/               # 配置文件
├── docs/                  # 项目文档
│   ├── README.md         # 文档目录说明
│   ├── RBAC_README.md    # RBAC 权限系统文档
│   ├── API_VERSIONING_CHANGES.md  # API 版本化更改
│   └── GET_ME_API_USAGE.md       # GetMe API 使用指南
├── script/                # 脚本文件
│   ├── README.md         # 脚本说明
│   ├── test_getme.sh     # GetMe API 测试脚本
│   └── run_tests.sh      # 测试运行脚本
├── internal/              # 内部代码
│   ├── biz/              # 业务逻辑层
│   ├── data/             # 数据访问层
│   ├── service/          # 服务层
│   ├── server/           # 服务器配置
│   └── pkg/              # 公共包
├── migrate/              # 数据库迁移文件
├── tests/                # 测试文件
├── third_party/          # 第三方依赖
├── deloy.sh              # 部署脚本
├── Dockerfile            # Docker 配置
├── Makefile              # 构建脚本
└── README.md             # 项目说明
```

## 快速开始

### 环境要求

- Go 1.23+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖

```bash
# 安装 Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# 安装 Wire
go get github.com/google/wire/cmd/wire

# 下载项目依赖
make init
```

### 配置数据库

1. 创建数据库
2. 执行迁移文件：
   ```bash
   mysql -u username -p database_name < migrate/migrate.sql
   mysql -u username -p database_name < migrate/rbac_migrate.sql
   ```

### 运行项目

```bash
# 生成代码
make api

# 生成依赖注入
make wire

# 构建项目
make build

# 运行项目
./bin/student -conf ./configs
```

## API 接口

### 用户管理

- `POST /v1/user/login` - 用户登录
- `GET /v1/user/me` - 获取当前用户信息
- `GET /v1/users` - 获取用户列表
- `POST /v1/user` - 创建用户
- `GET /v1/user/{id}` - 获取用户详情
- `PUT /v1/user/{id}` - 更新用户
- `DELETE /v1/user/{id}` - 删除用户

### 学生管理

- `GET /v1/students` - 获取学生列表
- `POST /v1/student` - 创建学生
- `GET /v1/student/{id}` - 获取学生详情
- `PUT /v1/student/{id}` - 更新学生
- `DELETE /v1/student/{id}` - 删除学生

### RBAC 权限管理

- `GET /v1/roles` - 获取角色列表
- `POST /v1/roles` - 创建角色
- `GET /v1/permissions` - 获取权限列表
- `POST /v1/permissions/check` - 权限检查

## 文档

详细文档请查看 `docs/` 目录：

- [RBAC 权限系统文档](docs/RBAC_README.md)
- [API 版本化更改说明](docs/API_VERSIONING_CHANGES.md)
- [GetMe API 使用指南](docs/GET_ME_API_USAGE.md)

## 测试

```bash
# 运行测试脚本
./script/run_tests.sh

# 测试 GetMe API
./script/test_getme.sh
```

## 部署

```bash
# 使用 Docker 部署
docker build -t student-system .
docker run --rm -p 8000:8000 -p 9000:9000 -v ./configs:/data/conf student-system

# 或使用部署脚本
./deloy.sh
```

## 开发

### 生成代码

```bash
# 生成 API 文件 (pb.go, http, grpc, validate, swagger)
make api

# 生成所有文件
make all

# 生成 Wire 依赖注入
make wire
```

### 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方代码规范
- 添加必要的注释和文档

## 许可证

MIT License
