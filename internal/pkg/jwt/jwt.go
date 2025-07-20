package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT配置
type Config struct {
	SecretKey string        `json:"secret_key"`
	Expire    time.Duration `json:"expire"`
}

// 自定义Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWT工具
type JWTUtil struct {
	config *Config
}

// 创建JWT工具实例
func NewJWTUtil(config *Config) *JWTUtil {
	return &JWTUtil{
		config: config,
	}
}

// 生成JWT Token
func (j *JWTUtil) GenerateToken(userID uint, username, email string) (string, error) {
	// 创建Claims
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "student-system",
			Subject:   username,
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 验证JWT Token
func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token是否有效
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 获取Claims
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid claims")
}

// 刷新Token
func (j *JWTUtil) RefreshToken(tokenString string) (string, error) {
	// 验证当前Token
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的Token
	return j.GenerateToken(claims.UserID, claims.Username, claims.Email)
}

// 从Token中获取用户ID
func (j *JWTUtil) GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// 从Token中获取用户名
func (j *JWTUtil) GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}
