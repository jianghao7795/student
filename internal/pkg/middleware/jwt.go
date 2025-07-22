package middleware

import (
	"context"
	"log"
	"net/http"
	"slices"
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
			log.Println(config.SkipPaths)
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
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "email", claims.Email)

			// 添加调试信息
			log.Printf("JWT中间件: 用户ID=%d, 用户名=%s, 邮箱=%s", claims.UserID, claims.Username, claims.Email)

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
	// 从上下文中获取HTTP请求信息
	if httpCtx, ok := ctx.(interface {
		Request() *http.Request
	}); ok {
		req := httpCtx.Request()
		path := req.URL.Path

		// 检查当前路径是否在跳过列表中
		return slices.Contains(skipPaths, path)
	}
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
