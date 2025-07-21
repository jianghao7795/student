package main

import (
	"context"
	"fmt"
	"log"

	v1 "student/api/rbac/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接gRPC服务器
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := v1.NewRBACServiceClient(conn)
	ctx := context.Background()

	// 测试获取角色列表
	fmt.Println("=== 测试获取角色列表 ===")
	rolesResp, err := client.ListRoles(ctx, &v1.ListRolesRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("Failed to list roles: %v", err)
	} else {
		fmt.Printf("角色总数: %d\n", rolesResp.Total)
		for _, role := range rolesResp.Roles {
			fmt.Printf("角色: %s - %s (状态: %d)\n", role.Name, role.Description, role.Status)
		}
	}

	// 测试获取权限列表
	fmt.Println("\n=== 测试获取权限列表 ===")
	permissionsResp, err := client.ListPermissions(ctx, &v1.ListPermissionsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("Failed to list permissions: %v", err)
	} else {
		fmt.Printf("权限总数: %d\n", permissionsResp.Total)
		for _, permission := range permissionsResp.Permissions {
			fmt.Printf("权限: %s - %s %s (状态: %d)\n",
				permission.Name, permission.Action, permission.Resource, permission.Status)
		}
	}

	// 测试权限检查
	fmt.Println("\n=== 测试权限检查 ===")
	checkResp, err := client.CheckPermission(ctx, &v1.CheckPermissionRequest{
		User:     "1", // admin用户ID
		Resource: "/api/v1/users",
		Action:   "GET",
	})
	if err != nil {
		log.Printf("Failed to check permission: %v", err)
	} else {
		fmt.Printf("用户1是否有权限访问 /api/v1/users (GET): %t\n", checkResp.HasPermission)
	}

	// 测试获取用户角色
	fmt.Println("\n=== 测试获取用户角色 ===")
	userRolesResp, err := client.GetUserRoles(ctx, &v1.GetUserRolesRequest{
		UserId: 1,
	})
	if err != nil {
		log.Printf("Failed to get user roles: %v", err)
	} else {
		fmt.Printf("用户1的角色数量: %d\n", len(userRolesResp.UserRoles))
		for _, userRole := range userRolesResp.UserRoles {
			if userRole.Role != nil {
				fmt.Printf("角色: %s - %s\n", userRole.Role.Name, userRole.Role.Description)
			}
		}
	}

	fmt.Println("\n=== RBAC测试完成 ===")
}
