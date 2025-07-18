package middleware

import (
	"context"
	"net/http"
	"strings"

	"student/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// JWT中间件配置
type JWTConfig struct {
	JWTUtil *jwt.JWTUtil
	// 不需要验证JWT的路径
	SkipPaths []string
}

// JWT中间件
func JWTAuth(config *JWTConfig) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 检查是否需要跳过JWT验证
			if shouldSkipJWT(ctx, config.SkipPaths) {
				return handler(ctx, req)
			}

			// 从HTTP请求中获取token
			token, err := extractTokenFromContext(ctx)
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

// 从上下文中提取token
func extractTokenFromContext(ctx context.Context) (string, error) {
	// 对于Kratos，我们需要通过HTTP请求头获取
	if httpCtx, ok := ctx.(interface {
		Request() *http.Request
	}); ok {
		req := httpCtx.Request()
		return extractTokenFromRequest(req)
	}
	return "", errors.New(400, "INVALID_REQUEST", "无法获取请求信息")
}

// 从HTTP请求中提取token
func extractTokenFromRequest(req *http.Request) (string, error) {
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

// 检查是否需要跳过JWT验证
func shouldSkipJWT(ctx context.Context, skipPaths []string) bool {
	// 这里可以根据请求路径来判断是否需要跳过JWT验证
	// 对于登录等公开接口，可以跳过JWT验证
	// 暂时返回false，后续可以根据需要扩展
	return false
}

// 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	if userID, ok := ctx.Value(userIDKey).(uint); ok {
		return userID, true
	}
	return 0, false
}

// 从上下文中获取用户名
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	if username, ok := ctx.Value(usernameKey).(string); ok {
		return username, true
	}
	return "", false
}

// 从上下文中获取邮箱
func GetEmailFromContext(ctx context.Context) (string, bool) {
	if email, ok := ctx.Value(emailKey).(string); ok {
		return email, true
	}
	return "", false
}
