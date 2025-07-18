package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost 默认的加密成本
	DefaultCost = 12
)

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	return string(bytes), err
}

// HashPasswordWithCost 使用指定成本加密密码
func HashPasswordWithCost(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// CheckPassword 验证密码是否匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsHashed 检查密码是否已经加密
func IsHashed(password string) bool {
	// bcrypt 加密后的密码通常以 $2a$, $2b$, $2x$, $2y$ 开头
	return len(password) >= 4 && (password[:4] == "$2a$" || password[:4] == "$2b$" || password[:4] == "$2x$" || password[:4] == "$2y$")
}
