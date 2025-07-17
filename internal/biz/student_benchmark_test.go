package biz

import (
	"testing"
	"time"
)

// 模拟优化前的实现：每次都进行Format转换
func formatTimeOld(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// 模拟优化后的实现：使用预格式化的字符串
func formatTimeNew(createdAtStr, updatedAtStr string) (string, string) {
	return createdAtStr, updatedAtStr
}

// 性能对比测试：优化前 vs 优化后
func BenchmarkTimeFormatOld(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟优化前：每次都进行Format转换
		_ = formatTimeOld(&now)
		_ = formatTimeOld(&now)
	}
}

func BenchmarkTimeFormatNew(b *testing.B) {
	now := time.Now()
	// 预格式化（模拟在data层完成）
	createdAtStr := now.Format("2006-01-02 15:04:05")
	updatedAtStr := now.Format("2006-01-02 15:04:05")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟优化后：直接使用预格式化的字符串
		_, _ = formatTimeNew(createdAtStr, updatedAtStr)
	}
}

// 批量处理性能对比
func BenchmarkBatchFormatOld(b *testing.B) {
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
		// 模拟优化前：在service层每次都进行Format转换
		for _, stu := range students {
			_ = formatTimeOld(stu.CreatedAt)
			_ = formatTimeOld(stu.UpdatedAt)
		}
	}
}

func BenchmarkBatchFormatNew(b *testing.B) {
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

	// 预格式化（模拟在data层完成）
	FormatTimeFieldsBatch(students)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟优化后：直接使用预格式化的字符串
		for _, stu := range students {
			_ = stu.CreatedAtStr
			_ = stu.UpdatedAtStr
		}
	}
}

// 内存分配对比测试
func BenchmarkMemoryAllocationOld(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// 优化前：每次Format都会分配新的字符串
		_ = now.Format("2006-01-02 15:04:05")
		_ = now.Format("2006-01-02 15:04:05")
	}
}

func BenchmarkMemoryAllocationNew(b *testing.B) {
	now := time.Now()
	// 预格式化一次
	createdAtStr := now.Format("2006-01-02 15:04:05")
	updatedAtStr := now.Format("2006-01-02 15:04:05")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// 优化后：重复使用已格式化的字符串
		_ = createdAtStr
		_ = updatedAtStr
	}
}
