package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"student/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
)

// 定义自定义的上下文键类型
type contextKey string

const (
	userIDKey   contextKey = "user_id"
	usernameKey contextKey = "username"
	emailKey    contextKey = "email"
)

// JWT HTTP中间件
func JWTHTTPMiddleware(jwtUtil *jwt.JWTUtil, skipPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否需要跳过JWT验证
			if shouldSkipPath(r.URL.Path, skipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 从请求头中获取token
			token, err := extractTokenFromHeader(r)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"code":    401,
					"message": "未提供有效的认证token",
				})
				return
			}

			// 验证token
			claims, err := jwtUtil.ValidateToken(token)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"code":    401,
					"message": "token验证失败",
				})
				return
			}

			// 将用户信息存储到请求上下文中
			ctx := r.Context()
			ctx = context.WithValue(ctx, userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, usernameKey, claims.Username)
			ctx = context.WithValue(ctx, emailKey, claims.Email)

			// 更新请求上下文
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// 从请求头中提取token
func extractTokenFromHeader(r *http.Request) (string, error) {
	// 从Authorization头中获取token
	authHeader := r.Header.Get("Authorization")
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
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
