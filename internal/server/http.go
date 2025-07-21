package server

import (
	rbacV1 "student/api/rbac/v1"
	v1 "student/api/student/v1"
	userV1 "student/api/user/v1"
	"student/internal/conf"
	"student/internal/pkg/jwt"
	"student/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, student *service.StudentService, user *service.UserService, rbac *service.RBACService, jwtUtil *jwt.JWTUtil, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	// 注册服务
	v1.RegisterStudentHTTPServer(srv, student)
	userV1.RegisterUserHTTPServer(srv, user)
	rbacV1.RegisterRBACServiceHTTPServer(srv, rbac)

	return srv
}
