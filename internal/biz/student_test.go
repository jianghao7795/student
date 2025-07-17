package biz

import (
	"testing"
	"time"
)

func TestStudent_FormatTimeFields(t *testing.T) {
	// 创建测试时间
	now := time.Date(2025, 7, 17, 18, 18, 43, 0, time.Local)

	tests := []struct {
		name     string
		student  *Student
		expected struct {
			createdAtStr string
			updatedAtStr string
		}
	}{
		{
			name: "正常时间格式化",
			student: &Student{
				ID:        1,
				Name:      "测试学生",
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
			student: &Student{
				ID:   2,
				Name: "空时间学生",
			},
			expected: struct {
				createdAtStr string
				updatedAtStr string
			}{
				createdAtStr: "",
				updatedAtStr: "",
			},
		},
		{
			name: "部分时间字段为空",
			student: &Student{
				ID:        3,
				Name:      "部分时间学生",
				CreatedAt: &now,
				// UpdatedAt 为空
			},
			expected: struct {
				createdAtStr string
				updatedAtStr string
			}{
				createdAtStr: "2025-07-17 18:18:43",
				updatedAtStr: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用格式化方法
			tt.student.FormatTimeFields()

			// 验证结果
			if tt.student.CreatedAtStr != tt.expected.createdAtStr {
				t.Errorf("CreatedAtStr = %v, want %v", tt.student.CreatedAtStr, tt.expected.createdAtStr)
			}
			if tt.student.UpdatedAtStr != tt.expected.updatedAtStr {
				t.Errorf("UpdatedAtStr = %v, want %v", tt.student.UpdatedAtStr, tt.expected.updatedAtStr)
			}
		})
	}
}

func TestFormatTimeFieldsBatch(t *testing.T) {
	// 创建测试时间
	now := time.Date(2025, 7, 17, 18, 18, 43, 0, time.Local)

	students := []*Student{
		{
			ID:        1,
			Name:      "学生1",
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:        2,
			Name:      "学生2",
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:   3,
			Name: "学生3",
			// 时间字段为空
		},
	}

	// 调用批量格式化方法
	FormatTimeFieldsBatch(students)

	// 验证结果
	expectedStr := "2025-07-17 18:18:43"

	for i, student := range students {
		if i < 2 {
			// 前两个学生应该有格式化后的时间
			if student.CreatedAtStr != expectedStr {
				t.Errorf("学生%d CreatedAtStr = %v, want %v", i+1, student.CreatedAtStr, expectedStr)
			}
			if student.UpdatedAtStr != expectedStr {
				t.Errorf("学生%d UpdatedAtStr = %v, want %v", i+1, student.UpdatedAtStr, expectedStr)
			}
		} else {
			// 第三个学生的时间字段应该为空
			if student.CreatedAtStr != "" {
				t.Errorf("学生%d CreatedAtStr = %v, want empty string", i+1, student.CreatedAtStr)
			}
			if student.UpdatedAtStr != "" {
				t.Errorf("学生%d UpdatedAtStr = %v, want empty string", i+1, student.UpdatedAtStr)
			}
		}
	}
}

func TestTimeFormatConstant(t *testing.T) {
	// 测试时间格式常量
	if TimeFormat != "2006-01-02 15:04:05" {
		t.Errorf("TimeFormat = %v, want %v", TimeFormat, "2006-01-02 15:04:05")
	}

	// 测试常量是否能正确格式化时间
	now := time.Date(2025, 7, 17, 18, 18, 43, 0, time.Local)
	formatted := now.Format(TimeFormat)
	expected := "2025-07-17 18:18:43"

	if formatted != expected {
		t.Errorf("格式化结果 = %v, want %v", formatted, expected)
	}
}

// 性能测试：比较优化前后的性能差异
func BenchmarkFormatTimeFields(b *testing.B) {
	now := time.Now()
	student := &Student{
		ID:        1,
		Name:      "性能测试学生",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		student.FormatTimeFields()
	}
}

func BenchmarkFormatTimeFieldsBatch(b *testing.B) {
	now := time.Now()
	students := make([]*Student, 100)
	for i := 0; i < 100; i++ {
		students[i] = &Student{
			ID:        uint(i + 1),
			Name:      "性能测试学生",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatTimeFieldsBatch(students)
	}
}
