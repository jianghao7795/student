package main

import (
	"fmt"
	"log"
	"student/internal/pkg/password"
)

func main() {
	fmt.Println("=== 密码加密演示 ===")

	// 示例1：基本密码加密
	fmt.Println("\n1. 基本密码加密:")
	plainPassword := "mypassword123"
	hashedPassword, err := password.HashPassword(plainPassword)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}
	fmt.Printf("原始密码: %s\n", plainPassword)
	fmt.Printf("加密后密码: %s\n", hashedPassword)

	// 验证密码
	if password.CheckPassword(plainPassword, hashedPassword) {
		fmt.Println("✅ 密码验证成功")
	} else {
		fmt.Println("❌ 密码验证失败")
	}

	// 示例2：不同成本的密码加密
	fmt.Println("\n2. 不同成本的密码加密:")
	costs := []int{10, 12, 14}
	for _, cost := range costs {
		hashed, err := password.HashPasswordWithCost(plainPassword, cost)
		if err != nil {
			log.Printf("成本 %d 的密码加密失败: %v", cost, err)
			continue
		}
		fmt.Printf("成本 %d: %s\n", cost, hashed)
	}

	// 示例3：密码格式检测
	fmt.Println("\n3. 密码格式检测:")
	testPasswords := []string{
		plainPassword,
		hashedPassword,
		"$2a$12$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890",
		"",
		"abc",
	}

	for _, pwd := range testPasswords {
		isHashed := password.IsHashed(pwd)
		fmt.Printf("密码: %s -> 是否已加密: %t\n", pwd, isHashed)
	}

	// 示例4：错误密码验证
	fmt.Println("\n4. 错误密码验证:")
	wrongPasswords := []string{
		"wrongpassword",
		"mypassword",
		"mypassword1234",
		"",
	}

	for _, wrongPwd := range wrongPasswords {
		isValid := password.CheckPassword(wrongPwd, hashedPassword)
		fmt.Printf("错误密码 '%s' 验证结果: %t\n", wrongPwd, isValid)
	}

	// 示例5：用户注册和登录模拟
	fmt.Println("\n5. 用户注册和登录模拟:")

	// 模拟用户注册
	userPassword := "user123456"
	userHashedPassword, _ := password.HashPassword(userPassword)
	fmt.Printf("用户注册 - 原始密码: %s\n", userPassword)
	fmt.Printf("用户注册 - 存储到数据库的密码: %s\n", userHashedPassword)

	// 模拟用户登录
	loginPassword := "user123456"
	if password.CheckPassword(loginPassword, userHashedPassword) {
		fmt.Println("✅ 用户登录成功")
	} else {
		fmt.Println("❌ 用户登录失败")
	}

	// 模拟错误登录
	wrongLoginPassword := "wrongpassword"
	if password.CheckPassword(wrongLoginPassword, userHashedPassword) {
		fmt.Println("✅ 用户登录成功")
	} else {
		fmt.Println("❌ 用户登录失败")
	}

	fmt.Println("\n=== 演示完成 ===")
}
