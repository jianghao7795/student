package server

import (
	"fmt"
	stdhttp "net/http"
	"net/http/httputil"
	"net/url"

	"student/internal/conf"
	"student/internal/pkg/nacos"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// GatewayProxy 网关代理
type GatewayProxy struct {
	discovery *nacos.Discovery
	config    *conf.Bootstrap
	log       *log.Helper
}

// NewGatewayProxy 创建网关代理
func NewGatewayProxy(discovery *nacos.Discovery, config *conf.Bootstrap, logger log.Logger) *GatewayProxy {
	return &GatewayProxy{
		discovery: discovery,
		config:    config,
		log:       log.NewHelper(logger),
	}
}

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

	proxy := NewGatewayProxy(discovery, c, logger)

	// 用户服务路由
	srv.HandlePrefix("/v1/user/", stdhttp.HandlerFunc(proxy.handleUserService))

	// 学生服务路由
	srv.HandlePrefix("/v1/student/", stdhttp.HandlerFunc(proxy.handleStudentService))

	// RBAC服务路由
	srv.HandlePrefix("/v1/rbac/", stdhttp.HandlerFunc(proxy.handleRBACService))

	// 健康检查端点
	srv.HandleFunc("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		w.Write([]byte("OK"))
	})

	return srv
}

// handleUserService 处理用户服务请求
func (p *GatewayProxy) handleUserService(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	instances, err := p.discovery.GetServiceInstances("user-service")
	if err != nil || len(instances) == 0 {
		p.log.Errorf("Failed to get user-service instances: %v", err)
		stdhttp.Error(w, "Service unavailable", stdhttp.StatusServiceUnavailable)
		return
	}

	// 简单的负载均衡 - 选择第一个实例
	instance := instances[0]
	target := fmt.Sprintf("http://%s:%d", instance.IP, instance.Port)

	p.proxyRequest(w, r, target)
}

// handleStudentService 处理学生服务请求
func (p *GatewayProxy) handleStudentService(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	instances, err := p.discovery.GetServiceInstances("student-service")
	if err != nil || len(instances) == 0 {
		p.log.Errorf("Failed to get student-service instances: %v", err)
		stdhttp.Error(w, "Service unavailable", stdhttp.StatusServiceUnavailable)
		return
	}

	// 简单的负载均衡 - 选择第一个实例
	instance := instances[0]
	target := fmt.Sprintf("http://%s:%d", instance.IP, instance.Port)

	p.proxyRequest(w, r, target)
}

// handleRBACService 处理RBAC服务请求
func (p *GatewayProxy) handleRBACService(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	instances, err := p.discovery.GetServiceInstances("rbac-service")
	if err != nil || len(instances) == 0 {
		p.log.Errorf("Failed to get rbac-service instances: %v", err)
		stdhttp.Error(w, "Service unavailable", stdhttp.StatusServiceUnavailable)
		return
	}

	// 简单的负载均衡 - 选择第一个实例
	instance := instances[0]
	target := fmt.Sprintf("http://%s:%d", instance.IP, instance.Port)

	p.proxyRequest(w, r, target)
}

// proxyRequest 代理请求到目标服务
func (p *GatewayProxy) proxyRequest(w stdhttp.ResponseWriter, r *stdhttp.Request, target string) {
	targetURL, err := url.Parse(target)
	if err != nil {
		p.log.Errorf("Failed to parse target URL: %v", err)
		stdhttp.Error(w, "Internal server error", stdhttp.StatusInternalServerError)
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 修改请求
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	p.log.Infof("Proxying request: %s %s to %s", r.Method, r.URL.Path, target)

	// 代理请求
	proxy.ServeHTTP(w, r)
}
