package middleware

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"student/internal/biz"
	"student/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// RBACConfig RBAC中间件配置
type RBACConfig struct {
	RBACUC  *biz.RBACUsecase
	JWTUtil *jwt.JWTUtil
	// 不需要进行RBAC权限检查的路径
	SkipPaths []string
}

// RBACMiddleware RBAC权限中间件
func RBACMiddleware(config *RBACConfig) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 获取HTTP请求信息
			if tr, ok := transport.FromServerContext(ctx); ok {
				// 获取请求路径和方法
				path := tr.RequestHeader().Get("X-Request-Path")
				if path == "" {
					path = tr.RequestHeader().Get("X-Original-URI")
				}
				method := tr.RequestHeader().Get("X-Request-Method")
				if method == "" {
					method = tr.RequestHeader().Get("X-HTTP-Method")
				}

				// 如果无法获取路径和方法，跳过权限检查
				if path == "" || method == "" {
					return handler(ctx, req)
				}

				// 检查是否需要跳过RBAC权限检查
				if shouldSkipRBAC(path, config.SkipPaths) {
					return handler(ctx, req)
				}

				// 从请求头获取JWT token
				token := tr.RequestHeader().Get("Authorization")
				if token == "" {
					return nil, errors.Unauthorized("UNAUTHORIZED", "未提供认证token")
				}

				// 移除Bearer前缀
				token = strings.TrimPrefix(token, "Bearer ")

				// 验证JWT token
				claims, err := config.JWTUtil.ValidateToken(token)
				if err != nil {
					return nil, errors.Unauthorized("UNAUTHORIZED", "无效的token")
				}

				// 检查权限
				userID := strconv.Itoa(int(claims.UserID))
				hasPermission, err := config.RBACUC.CheckPermission(ctx, userID, path, method)
				if err != nil {
					return nil, errors.InternalServer("INTERNAL_ERROR", "权限检查失败")
				}

				if !hasPermission {
					return nil, errors.Forbidden("FORBIDDEN", "没有访问权限")
				}
			}

			return handler(ctx, req)
		}
	}
}

// 检查是否需要跳过RBAC权限检查
func shouldSkipRBAC(path string, skipPaths []string) bool {
	return slices.Contains(skipPaths, path)
}

// SimpleRBACMiddleware 简化版RBAC中间件，用于特定路径的权限检查
func SimpleRBACMiddleware(rbacUC *biz.RBACUsecase, jwtUtil *jwt.JWTUtil, resource, action string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 获取HTTP请求信息
			if tr, ok := transport.FromServerContext(ctx); ok {
				// 从请求头获取JWT token
				token := tr.RequestHeader().Get("Authorization")
				if token == "" {
					return nil, errors.Unauthorized("UNAUTHORIZED", "未提供认证token")
				}

				// 移除Bearer前缀
				token = strings.TrimPrefix(token, "Bearer ")

				// 验证JWT token
				claims, err := jwtUtil.ValidateToken(token)
				if err != nil {
					return nil, errors.Unauthorized("UNAUTHORIZED", "无效的token")
				}

				// 检查权限
				userID := strconv.Itoa(int(claims.UserID))
				hasPermission, err := rbacUC.CheckPermission(ctx, userID, resource, action)
				if err != nil {
					return nil, errors.InternalServer("INTERNAL_ERROR", "权限检查失败")
				}

				if !hasPermission {
					return nil, errors.Forbidden("FORBIDDEN", "没有访问权限")
				}
			}

			return handler(ctx, req)
		}
	}
}
