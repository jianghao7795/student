package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Age      int32  `json:"age"`
	Avatar   string `json:"avatar"`
}

// RegisterReply 注册响应结构
type RegisterReply struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	UserInfo UserInfo `json:"user_info,omitempty"`
}

// UserInfo 用户信息结构
type UserInfo struct {
	Id        int32  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    int32  `json:"status"`
	Age       int32  `json:"age"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func main() {
	// 注册请求数据
	registerData := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Password: "password123",
		Age:      25,
		Avatar:   "https://example.com/avatar.jpg",
	}

	// 将数据转换为JSON
	jsonData, err := json.Marshal(registerData)
	if err != nil {
		fmt.Printf("JSON编码错误: %v\n", err)
		return
	}

	// 发送POST请求
	resp, err := http.Post("http://localhost:8000/v1/user/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应错误: %v\n", err)
		return
	}

	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))

	// 解析响应
	var registerReply RegisterReply
	if err := json.Unmarshal(body, &registerReply); err != nil {
		fmt.Printf("JSON解析错误: %v\n", err)
		return
	}

	if registerReply.Success {
		fmt.Printf("注册成功！用户ID: %d, 用户名: %s\n", registerReply.UserInfo.Id, registerReply.UserInfo.Username)
	} else {
		fmt.Printf("注册失败: %s\n", registerReply.Message)
	}
}
