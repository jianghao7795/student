package middleware

import (
	"context"
	"net/http"
	"strings"

	"student/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// SimpleJWTConfig JWT中间件配置
type SimpleJWTConfig struct {
	JWTUtil *jwt.JWTUtil
}

// SimpleJWTAuth 简单的JWT中间件
func SimpleJWTAuth(config *SimpleJWTConfig) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 从HTTP请求中获取token
			token, err := extractTokenFromHTTPContext(ctx)
			if err != nil {
				return nil, errors.Unauthorized("UNAUTHORIZED", "未提供有效的认证token")
			}

			// 验证token
			claims, err := config.JWTUtil.ValidateToken(token)
			if err != nil {
				return nil, errors.Unauthorized("UNAUTHORIZED", "token验证失败")
			}

			// 将用户信息存储到上下文中
			ctx = context.WithValue(ctx, userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, usernameKey, claims.Username)
			ctx = context.WithValue(ctx, emailKey, claims.Email)

			return handler(ctx, req)
		}
	}
}

// extractTokenFromHTTPContext 从HTTP上下文中提取token
func extractTokenFromHTTPContext(ctx context.Context) (string, error) {
	// 对于Kratos，我们需要通过HTTP请求头获取
	if httpCtx, ok := ctx.(interface {
		Request() *http.Request
	}); ok {
		req := httpCtx.Request()
		return extractTokenFromHTTPRequest(req)
	}
	return "", errors.New(400, "INVALID_REQUEST", "无法获取请求信息")
}

// extractTokenFromHTTPRequest 从HTTP请求中提取token
func extractTokenFromHTTPRequest(req *http.Request) (string, error) {
	// 从Authorization头中获取token
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New(400, "MISSING_TOKEN", "缺少Authorization头")
	}

	// 检查Authorization头的格式
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New(400, "INVALID_TOKEN_FORMAT", "Authorization头格式错误")
	}

	return parts[1], nil
}
