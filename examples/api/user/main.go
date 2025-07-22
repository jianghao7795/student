package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	v1 "student/api/user/v1"
)

// UserAPIExample 用户API使用示例
type UserAPIExample struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewUserAPIExample 创建用户API示例实例
func NewUserAPIExample(baseURL string) *UserAPIExample {
	return &UserAPIExample{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Login 用户登录示例
func (e *UserAPIExample) Login(username, password string) error {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("marshal login data failed: %v", err)
	}

	req, err := http.NewRequest("POST", e.BaseURL+"/v1/user/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp v1.LoginReply
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}

	e.Token = loginResp.Token
	fmt.Printf("登录成功，Token: %s\n", e.Token)
	return nil
}

// GetMe 获取当前用户信息示例
func (e *UserAPIExample) GetMe() error {
	if e.Token == "" {
		return fmt.Errorf("token is required, please login first")
	}

	req, err := http.NewRequest("GET", e.BaseURL+"/v1/account/me", nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+e.Token)

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("get me request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get me failed with status %d: %s", resp.StatusCode, string(body))
	}

	var meResp v1.GetMeReply
	if err := json.Unmarshal(body, &meResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}

	fmt.Printf("当前用户信息:\n")
	fmt.Printf("  成功: %t\n", meResp.Success)
	fmt.Printf("  消息: %s\n", meResp.Message)
	if meResp.UserInfo != nil {
		fmt.Printf("  用户ID: %d\n", meResp.UserInfo.Id)
		fmt.Printf("  用户名: %s\n", meResp.UserInfo.Username)
		fmt.Printf("  邮箱: %s\n", meResp.UserInfo.Email)
		fmt.Printf("  状态: %d\n", meResp.UserInfo.Status)
	}
	return nil
}

// ListUsers 获取用户列表示例
func (e *UserAPIExample) ListUsers(page, pageSize int32) error {
	if e.Token == "" {
		return fmt.Errorf("token is required, please login first")
	}

	url := fmt.Sprintf("%s/v1/users?page=%d&page_size=%d", e.BaseURL, page, pageSize)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+e.Token)

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("list users request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("list users failed with status %d: %s", resp.StatusCode, string(body))
	}

	var listResp v1.ListUsersReply
	if err := json.Unmarshal(body, &listResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}

	fmt.Printf("用户列表 (总数: %d):\n", listResp.Total)
	for _, user := range listResp.Data {
		fmt.Printf("  ID: %d, 用户名: %s, 邮箱: %s, 状态: %d\n",
			user.Id, user.Username, user.Email, user.Status)
	}
	return nil
}

// CreateUser 创建用户示例
func (e *UserAPIExample) CreateUser(username, password, email string) error {
	if e.Token == "" {
		return fmt.Errorf("token is required, please login first")
	}

	userData := map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("marshal user data failed: %v", err)
	}

	req, err := http.NewRequest("POST", e.BaseURL+"/v1/user", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.Token)

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("create user request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createResp v1.CreateUserReply
	if err := json.Unmarshal(body, &createResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}

	fmt.Printf("创建用户成功，消息: %s\n", createResp.Message)
	return nil
}

func main() {
	// 创建API示例实例
	api := NewUserAPIExample("http://localhost:8000")

	// 示例1: 用户登录
	fmt.Println("=== 示例1: 用户登录 ===")
	if err := api.Login("admin", "admin123"); err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}

	// 示例2: 获取当前用户信息
	fmt.Println("\n=== 示例2: 获取当前用户信息 ===")
	if err := api.GetMe(); err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
	}

	// 示例3: 获取用户列表
	fmt.Println("\n=== 示例3: 获取用户列表 ===")
	if err := api.ListUsers(1, 10); err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
	}

	// 示例4: 创建新用户
	fmt.Println("\n=== 示例4: 创建新用户 ===")
	if err := api.CreateUser("testuser", "testpass123", "test@example.com"); err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	fmt.Println("\n=== 用户API示例完成 ===")
}
