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

// RequireAuthorizationHeader 是一个中间件，要求请求必须带 Authorization 头
func RequireAuthorizationHeader() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing Authorization header"))
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware(jwtUtil *jwt.JWTUtil, skipPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否需要跳过认证
			if shouldSkipPath(r.URL.Path, skipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 检查 Authorization 头
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendErrorResponse(w, http.StatusUnauthorized, "缺少 Authorization 头")
				return
			}

			// 从 Authorization 头中提取 token
			token, err := extractTokenFromHeader(authHeader)
			if err != nil {
				sendErrorResponse(w, http.StatusUnauthorized, "Authorization 头格式错误")
				return
			}

			// 验证 JWT token
			claims, err := jwtUtil.ValidateToken(token)
			if err != nil {
				sendErrorResponse(w, http.StatusUnauthorized, "Token 验证失败")
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

// 从 Authorization 头中提取 token
func extractTokenFromHeader(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New(400, "INVALID_TOKEN_FORMAT", "Authorization 头格式错误")
	}
	return parts[1], nil
}

// 检查是否需要跳过认证
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// 发送错误响应
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}
