package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	v1 "student/api/student/v1"
)

func main() {
	// 连接 gRPC 服务器
	conn, err := grpc.Dial("localhost:9600", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	// 创建客户端
	client := v1.NewStudentClient(conn)

	// 测试场景1: 没有 Token 的请求（应该被拒绝）
	fmt.Println("=== 测试场景1: 没有 Token 的请求 ===")
	ctx1 := context.Background()
	_, err1 := client.GetStudent(ctx1, &v1.GetStudentRequest{Id: 1})
	if err1 != nil {
		fmt.Printf("预期错误: %v\n", err1)
	} else {
		fmt.Println("意外成功: 没有 Token 的请求被允许")
	}

	// 测试场景2: 无效 Token 的请求（应该被拒绝）
	fmt.Println("\n=== 测试场景2: 无效 Token 的请求 ===")
	ctx2 := context.Background()
	ctx2 = metadata.AppendToOutgoingContext(ctx2, "authorization", "Bearer invalid-token")
	_, err2 := client.GetStudent(ctx2, &v1.GetStudentRequest{Id: 1})
	if err2 != nil {
		fmt.Printf("预期错误: %v\n", err2)
	} else {
		fmt.Println("意外成功: 无效 Token 的请求被允许")
	}

	// 测试场景3: 有效 Token 的请求（需要有效的 JWT Token）
	fmt.Println("\n=== 测试场景3: 有效 Token 的请求 ===")
	// 注意：这里需要替换为实际有效的 JWT Token
	validToken := "your-valid-jwt-token-here"
	ctx3 := context.Background()
	ctx3 = metadata.AppendToOutgoingContext(ctx3, "authorization", "Bearer "+validToken)
	_, err3 := client.GetStudent(ctx3, &v1.GetStudentRequest{Id: 1})
	if err3 != nil {
		fmt.Printf("错误: %v\n", err3)
	} else {
		fmt.Println("成功: 有效 Token 的请求被允许")
	}

	fmt.Println("\n=== 测试完成 ===")
}
