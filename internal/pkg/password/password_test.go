package password

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	// 测试密码加密
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 验证加密后的密码不是原密码
	if hashedPassword == password {
		t.Error("Hashed password should not be equal to original password")
	}

	// 验证密码格式
	if !IsHashed(hashedPassword) {
		t.Error("Hashed password should be recognized as hashed")
	}

	// 验证密码检查
	if !CheckPassword(password, hashedPassword) {
		t.Error("Password check should pass for correct password")
	}

	// 验证错误密码
	if CheckPassword("wrongpassword", hashedPassword) {
		t.Error("Password check should fail for wrong password")
	}
}

func TestHashPasswordWithCost(t *testing.T) {
	password := "testpassword123"

	// 测试不同成本的密码加密
	costs := []int{10, 12, 14}

	for _, cost := range costs {
		hashedPassword, err := HashPasswordWithCost(password, cost)
		if err != nil {
			t.Fatalf("HashPasswordWithCost failed with cost %d: %v", cost, err)
		}

		// 验证密码检查
		if !CheckPassword(password, hashedPassword) {
			t.Errorf("Password check should pass for cost %d", cost)
		}
	}
}

func TestIsHashed(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "bcrypt hash",
			password: "$2a$12$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890",
			expected: true,
		},
		{
			name:     "plain password",
			password: "plainpassword",
			expected: false,
		},
		{
			name:     "empty string",
			password: "",
			expected: false,
		},
		{
			name:     "short string",
			password: "abc",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHashed(tt.password)
			if result != tt.expected {
				t.Errorf("IsHashed(%s) = %v, want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := HashPassword(password)

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hashedPassword,
			expected: true,
		},
		{
			name:     "wrong password",
			password: "wrongpassword",
			hash:     hashedPassword,
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hashedPassword,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			if result != tt.expected {
				t.Errorf("CheckPassword(%s, %s) = %v, want %v", tt.password, tt.hash, result, tt.expected)
			}
		})
	}
}

// 性能测试
func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword failed: %v", err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "testpassword123"
	hashedPassword, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hashedPassword)
	}
}
