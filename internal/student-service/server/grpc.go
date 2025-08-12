package server

import (
	v1 "student/api/student/v1"
	"student/internal/biz"
	"student/internal/conf"
	"student/internal/pkg/jwt"
	"student/internal/pkg/middleware"
	"student/internal/student-service/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, student *service.StudentService, rbacUC *biz.RBACUsecase, jwtUtil *jwt.JWTUtil, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			// JWT认证中间件
			middleware.JWTAuth(&middleware.JWTConfig{
				JWTUtil: jwtUtil,
				SkipPaths: []string{
					"/student.v1.Student/HealthCheck",
				},
			}),
			// RBAC权限中间件
			middleware.RBACMiddleware(&middleware.RBACConfig{
				RBACUC:  rbacUC,
				JWTUtil: jwtUtil,
				SkipPaths: []string{
					"/student.v1.Student/HealthCheck",
				},
			}),
		),
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
	return srv
}
