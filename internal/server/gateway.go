package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"student/internal/conf"
	"student/internal/pkg/nacos"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewGatewayHTTPServer 创建网关HTTP服务器
func NewGatewayHTTPServer(c *conf.Bootstrap, discovery *nacos.Discovery, logger log.Logger) *kratoshttp.Server {
	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Server.Http.Timeout.AsDuration()))
	}
	srv := kratoshttp.NewServer(opts...)

	// 创建反向代理
	proxy := &GatewayProxy{
		discovery: discovery,
		log:       log.NewHelper(logger),
	}

	// 注册路由
	srv.HandleFunc("/v1/user/", proxy.handleUserService)
	srv.HandleFunc("/v1/student/", proxy.handleStudentService)
	srv.HandleFunc("/v1/rbac/", proxy.handleRBACService)
	srv.HandleFunc("/v1/roles/", proxy.handleRBACService)
	srv.HandleFunc("/v1/permissions/", proxy.handleRBACService)

	// 健康检查
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return srv
}

// GatewayProxy 网关代理
type GatewayProxy struct {
	discovery *nacos.Discovery
	log       *log.Helper
}

// handleUserService 处理用户服务请求
func (p *GatewayProxy) handleUserService(w http.ResponseWriter, r *http.Request) {
	instances, err := p.discovery.GetServiceInstances("user-service")
	if err != nil {
		p.log.Error("Failed to get user service instances", "error", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	if len(instances) == 0 {
		p.log.Error("No user service instances available")
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	// 简单的负载均衡：选择第一个可用实例
	instance := instances[0]
	targetURL := fmt.Sprintf("http://%s", instance.GetServiceURL())

	p.log.Info("Routing to user service", "target", targetURL, "path", r.URL.Path)

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   instance.GetServiceURL(),
	})

	// 修改请求路径，去掉服务前缀
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/v1/user")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	proxy.ServeHTTP(w, r)
}

// handleStudentService 处理学生服务请求
func (p *GatewayProxy) handleStudentService(w http.ResponseWriter, r *http.Request) {
	instances, err := p.discovery.GetServiceInstances("student-service")
	if err != nil {
		p.log.Error("Failed to get student service instances", "error", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	if len(instances) == 0 {
		p.log.Error("No student service instances available")
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	// 简单的负载均衡：选择第一个可用实例
	instance := instances[0]
	targetURL := fmt.Sprintf("http://%s", instance.GetServiceURL())

	p.log.Info("Routing to student service", "target", targetURL, "path", r.URL.Path)

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   instance.GetServiceURL(),
	})

	// 修改请求路径，去掉服务前缀
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/v1/student")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	proxy.ServeHTTP(w, r)
}

// handleRBACService 处理RBAC服务请求
func (p *GatewayProxy) handleRBACService(w http.ResponseWriter, r *http.Request) {
	instances, err := p.discovery.GetServiceInstances("rbac-service")
	if err != nil {
		p.log.Error("Failed to get rbac service instances", "error", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	if len(instances) == 0 {
		p.log.Error("No rbac service instances available")
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	// 简单的负载均衡：选择第一个可用实例
	instance := instances[0]
	targetURL := fmt.Sprintf("http://%s", instance.GetServiceURL())

	p.log.Info("Routing to rbac service", "target", targetURL, "path", r.URL.Path)

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   instance.GetServiceURL(),
	})

	// 修改请求路径，去掉服务前缀
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/v1/rbac")
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/v1/roles")
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/v1/permissions")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	proxy.ServeHTTP(w, r)
}
