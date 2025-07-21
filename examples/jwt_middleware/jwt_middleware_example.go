package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTMiddleware JWT中间件示例
type JWTMiddleware struct {
	SecretKey []byte
}

// NewJWTMiddleware 创建JWT中间件实例
func NewJWTMiddleware(secretKey string) *JWTMiddleware {
	return &JWTMiddleware{
		SecretKey: []byte(secretKey),
	}
}

// GenerateToken 生成JWT token
func (m *JWTMiddleware) GenerateToken(userID uint, username string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "student-system",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.SecretKey)
}

// ValidateToken 验证JWT token
func (m *JWTMiddleware) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token failed: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractTokenFromHeader 从请求头中提取token
func (m *JWTMiddleware) ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	// 检查Bearer前缀
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

// AuthMiddleware 认证中间件
func (m *JWTMiddleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 提取token
		tokenString, err := m.ExtractTokenFromHeader(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// 验证token
		claims, err := m.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// 将用户信息添加到请求上下文
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "claims", claims)

		// 调用下一个处理器
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func (m *JWTMiddleware) OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 尝试提取token
		tokenString, err := m.ExtractTokenFromHeader(r)
		if err != nil {
			// token不存在或格式错误，继续执行但不添加用户信息
			next.ServeHTTP(w, r)
			return
		}

		// 尝试验证token
		claims, err := m.ValidateToken(tokenString)
		if err != nil {
			// token无效，继续执行但不添加用户信息
			next.ServeHTTP(w, r)
			return
		}

		// token有效，将用户信息添加到请求上下文
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "claims", claims)

		// 调用下一个处理器
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// LoggingMiddleware 日志中间件示例
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 记录请求信息
		fmt.Printf("[%s] %s %s - Started\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		// 记录响应时间
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s - Completed in %v\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, duration)
	}
}

// CORSMiddleware CORS中间件示例
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	}
}

// RateLimitMiddleware 简单的速率限制中间件示例
type RateLimitMiddleware struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimitMiddleware 创建速率限制中间件
func NewRateLimitMiddleware(limit int, window time.Duration) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit 速率限制中间件
func (m *RateLimitMiddleware) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取客户端IP（简化版本）
		clientIP := r.RemoteAddr
		now := time.Now()

		// 清理过期的请求记录
		if requests, exists := m.requests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < m.window {
					validRequests = append(validRequests, reqTime)
				}
			}
			m.requests[clientIP] = validRequests
		}

		// 检查请求数量
		if len(m.requests[clientIP]) >= m.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// 记录当前请求
		m.requests[clientIP] = append(m.requests[clientIP], now)

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	}
}

// 示例处理器
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	username := r.Context().Value("username")

	fmt.Fprintf(w, "Protected endpoint - User ID: %v, Username: %v\n", userID, username)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID != nil {
		fmt.Fprintf(w, "Public endpoint - Authenticated user: %v\n", userID)
	} else {
		fmt.Fprintf(w, "Public endpoint - Anonymous user\n")
	}
}

func main() {
	// 创建JWT中间件
	jwtMiddleware := NewJWTMiddleware("your-secret-key")

	// 创建速率限制中间件
	rateLimitMiddleware := NewRateLimitMiddleware(10, time.Minute) // 每分钟最多10个请求

	// 示例1: 生成token
	fmt.Println("=== 示例1: 生成JWT Token ===")
	token, err := jwtMiddleware.GenerateToken(1, "admin")
	if err != nil {
		fmt.Printf("生成token失败: %v\n", err)
	} else {
		fmt.Printf("生成的token: %s\n", token)
	}

	// 示例2: 验证token
	fmt.Println("\n=== 示例2: 验证JWT Token ===")
	if token != "" {
		claims, err := jwtMiddleware.ValidateToken(token)
		if err != nil {
			fmt.Printf("验证token失败: %v\n", err)
		} else {
			fmt.Printf("Token验证成功 - 用户ID: %d, 用户名: %s\n", claims.UserID, claims.Username)
		}
	}

	// 示例3: 设置HTTP服务器（仅作示例，不实际启动）
	fmt.Println("\n=== 示例3: HTTP中间件链示例 ===")
	fmt.Println("中间件链: CORS -> Logging -> RateLimit -> Auth -> Handler")

	// 创建中间件链
	handler := http.HandlerFunc(protectedHandler)
	handler = jwtMiddleware.AuthMiddleware(handler)
	handler = rateLimitMiddleware.RateLimit(handler)
	handler = LoggingMiddleware(handler)
	handler = CORSMiddleware(handler)

	fmt.Println("中间件链创建完成")

	// 示例4: 可选认证中间件
	fmt.Println("\n=== 示例4: 可选认证中间件示例 ===")
	publicHandler := http.HandlerFunc(publicHandler)
	publicHandler = jwtMiddleware.OptionalAuthMiddleware(publicHandler)
	publicHandler = LoggingMiddleware(publicHandler)
	publicHandler = CORSMiddleware(publicHandler)

	fmt.Println("可选认证中间件链创建完成")

	fmt.Println("\n=== JWT中间件示例完成 ===")
}
