package server

import (
	v1 "student/api/student/v1"
	userV1 "student/api/user/v1"
	"student/internal/biz"
	"student/internal/conf"
	"student/internal/pkg/jwt"
	"student/internal/service"
	"student/internal/pkg/middleware"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, student *service.StudentService, user *service.UserService, rbacUC *biz.RBACUsecase, jwtUtil *jwt.JWTUtil, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{}
	
	// 如果启用了 RBAC，添加 RBAC 中间件
	if c.Rbac != nil && c.Rbac.Enabled {
		// 创建 RBAC 中间件配置
		rbacConfig := &middleware.RBACConfig{
			RBACUC: rbacUC,
			JWTUtil: jwtUtil,
			SkipPaths: []string{
				// 可以在这里添加不需要权限检查的 gRPC 方法路径
				// 例如："/student.v1.Student/GetStudent",
			},
		}
		
		// 添加 RBAC 中间件到 gRPC 中间件链
		opts = append(opts, grpc.Middleware(
			recovery.Recovery(),
			middleware.RBACMiddleware(rbacConfig),
		))
	} else {
		// 如果没有启用 RBAC，只使用 recovery 中间件
		opts = append(opts, grpc.Middleware(
			recovery.Recovery(),
		))
	}
	
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterStudentServer(srv, student)
	userV1.RegisterUserServer(srv, user)
	return srv
}
