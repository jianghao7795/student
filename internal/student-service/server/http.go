package server

import (
	stdhttp "net/http"
	v1 "student/api/student/v1"
	"student/internal/biz"
	"student/internal/conf"
	"student/internal/pkg/jwt"
	"student/internal/pkg/middleware"
	"student/internal/student-service/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, student *service.StudentService, rbacUC *biz.RBACUsecase, jwtUtil *jwt.JWTUtil, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			// JWT认证中间件
			middleware.JWTAuth(&middleware.JWTConfig{
				JWTUtil: jwtUtil,
				SkipPaths: []string{
					"/health",
					"/v1/students/health",
				},
			}),
			// RBAC权限中间件
			middleware.RBACMiddleware(&middleware.RBACConfig{
				RBACUC:  rbacUC,
				JWTUtil: jwtUtil,
				SkipPaths: []string{
					"/health",
					"/v1/students/health",
				},
			}),
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

	// 注册学生服务
	v1.RegisterStudentHTTPServer(srv, student)

	// 添加健康检查端点
	srv.HandleFunc("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		w.Write([]byte("OK"))
	})

	return srv
}
