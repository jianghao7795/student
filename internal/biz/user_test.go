package biz

import (
	"testing"
	"time"
)

func TestUser_FormatTimeFields(t *testing.T) {
	// 创建测试时间
	now := time.Date(2025, 7, 17, 18, 18, 43, 0, time.Local)

	tests := []struct {
		name     string
		user     *User
		expected struct {
			createdAtStr string
			updatedAtStr string
		}
	}{
		{
			name: "正常时间格式化",
			user: &User{
				ID:        1,
				Username:  "testuser",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			expected: struct {
				createdAtStr string
				updatedAtStr string
			}{
				createdAtStr: "2025-07-17 18:18:43",
				updatedAtStr: "2025-07-17 18:18:43",
			},
		},
		{
			name: "空时间字段",
			user: &User{
				ID:       2,
				Username: "emptyuser",
			},
			expected: struct {
				createdAtStr string
				updatedAtStr string
			}{
				createdAtStr: "",
				updatedAtStr: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用格式化方法
			tt.user.FormatTimeFields()

			// 验证结果
			if tt.user.CreatedAtStr != tt.expected.createdAtStr {
				t.Errorf("CreatedAtStr = %v, want %v", tt.user.CreatedAtStr, tt.expected.createdAtStr)
			}
			if tt.user.UpdatedAtStr != tt.expected.updatedAtStr {
				t.Errorf("UpdatedAtStr = %v, want %v", tt.user.UpdatedAtStr, tt.expected.updatedAtStr)
			}
		})
	}
}

func TestFormatUserTimeFieldsBatch(t *testing.T) {
	// 创建测试时间
	now := time.Date(2025, 7, 17, 18, 18, 43, 0, time.Local)

	users := []*User{
		{
			ID:        1,
			Username:  "user1",
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:        2,
			Username:  "user2",
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:       3,
			Username: "user3",
			// 时间字段为空
		},
	}

	// 调用批量格式化方法
	FormatUserTimeFieldsBatch(users)

	// 验证结果
	expectedStr := "2025-07-17 18:18:43"

	for i, user := range users {
		if i < 2 {
			// 前两个用户应该有格式化后的时间
			if user.CreatedAtStr != expectedStr {
				t.Errorf("用户%d CreatedAtStr = %v, want %v", i+1, user.CreatedAtStr, expectedStr)
			}
			if user.UpdatedAtStr != expectedStr {
				t.Errorf("用户%d UpdatedAtStr = %v, want %v", i+1, user.UpdatedAtStr, expectedStr)
			}
		} else {
			// 第三个用户的时间字段应该为空
			if user.CreatedAtStr != "" {
				t.Errorf("用户%d CreatedAtStr = %v, want empty string", i+1, user.CreatedAtStr)
			}
			if user.UpdatedAtStr != "" {
				t.Errorf("用户%d UpdatedAtStr = %v, want empty string", i+1, user.UpdatedAtStr)
			}
		}
	}
}

// 性能测试
func BenchmarkUserFormatTimeFields(b *testing.B) {
	now := time.Now()
	user := &User{
		ID:        1,
		Username:  "性能测试用户",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.FormatTimeFields()
	}
}

func BenchmarkUserFormatTimeFieldsBatch(b *testing.B) {
	now := time.Now()
	users := make([]*User, 100)
	for i := 0; i < 100; i++ {
		users[i] = &User{
			ID:        uint(i + 1),
			Username:  "性能测试用户",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatUserTimeFieldsBatch(users)
	}
}
