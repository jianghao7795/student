package main

import (
	"fmt"
	"student/internal/biz"
	"time"
)

func main() {
	// 创建测试数据
	now := time.Now()
	student := &biz.Student{
		ID:        1,
		Name:      "测试学生",
		Info:      "测试信息",
		Status:    1,
		Age:       20,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// 测试格式化方法
	fmt.Println("格式化前:")
	fmt.Printf("CreatedAt: %v\n", student.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", student.UpdatedAt)
	fmt.Printf("CreatedAtStr: %s\n", student.CreatedAtStr)
	fmt.Printf("UpdatedAtStr: %s\n", student.UpdatedAtStr)

	// 调用格式化方法
	student.FormatTimeFields()

	fmt.Println("\n格式化后:")
	fmt.Printf("CreatedAtStr: %s\n", student.CreatedAtStr)
	fmt.Printf("UpdatedAtStr: %s\n", student.UpdatedAtStr)

	// 测试批量格式化
	students := []*biz.Student{
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
	}

	biz.FormatTimeFieldsBatch(students)

	fmt.Println("\n批量格式化后:")
	for i, stu := range students {
		fmt.Printf("学生%d - CreatedAtStr: %s, UpdatedAtStr: %s\n", i+1, stu.CreatedAtStr, stu.UpdatedAtStr)
	}

	fmt.Println("\n优化完成！时间字段不再需要每次都进行Format转换。")
}
