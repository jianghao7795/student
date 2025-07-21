package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClientExample HTTP客户端使用示例
type HTTPClientExample struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewHTTPClientExample 创建HTTP客户端示例实例
func NewHTTPClientExample(baseURL string) *HTTPClientExample {
	return &HTTPClientExample{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRequest GET请求示例
func (e *HTTPClientExample) GetRequest(path string) error {
	url := e.BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create GET request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "StudentSystem-Client/1.0")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	fmt.Printf("GET %s - Status: %d\n", url, resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
	return nil
}

// PostRequest POST请求示例
func (e *HTTPClientExample) PostRequest(path string, data interface{}) error {
	url := e.BaseURL + path

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create POST request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "StudentSystem-Client/1.0")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	fmt.Printf("POST %s - Status: %d\n", url, resp.StatusCode)
	fmt.Printf("Request Data: %s\n", string(jsonData))
	fmt.Printf("Response: %s\n", string(body))
	return nil
}

// PutRequest PUT请求示例
func (e *HTTPClientExample) PutRequest(path string, data interface{}) error {
	url := e.BaseURL + path

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create PUT request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "StudentSystem-Client/1.0")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("PUT request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	fmt.Printf("PUT %s - Status: %d\n", url, resp.StatusCode)
	fmt.Printf("Request Data: %s\n", string(jsonData))
	fmt.Printf("Response: %s\n", string(body))
	return nil
}

// DeleteRequest DELETE请求示例
func (e *HTTPClientExample) DeleteRequest(path string) error {
	url := e.BaseURL + path
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("create DELETE request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "StudentSystem-Client/1.0")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	fmt.Printf("DELETE %s - Status: %d\n", url, resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
	return nil
}

// RequestWithAuth 带认证的请求示例
func (e *HTTPClientExample) RequestWithAuth(method, path, token string, data interface{}) error {
	url := e.BaseURL + path

	var req *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshal data failed: %v", err)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "StudentSystem-Client/1.0")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s request failed: %v", method, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}

	fmt.Printf("%s %s - Status: %d\n", method, url, resp.StatusCode)
	if data != nil {
		jsonData, _ := json.Marshal(data)
		fmt.Printf("Request Data: %s\n", string(jsonData))
	}
	fmt.Printf("Response: %s\n", string(body))
	return nil
}

func main() {
	// 创建HTTP客户端示例实例
	client := NewHTTPClientExample("http://localhost:8000")

	// 示例1: GET请求
	fmt.Println("=== 示例1: GET请求 ===")
	if err := client.GetRequest("/v1/users"); err != nil {
		fmt.Printf("GET请求失败: %v\n", err)
	}

	// 示例2: POST请求
	fmt.Println("\n=== 示例2: POST请求 ===")
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	if err := client.PostRequest("/v1/user/login", loginData); err != nil {
		fmt.Printf("POST请求失败: %v\n", err)
	}

	// 示例3: PUT请求
	fmt.Println("\n=== 示例3: PUT请求 ===")
	updateData := map[string]interface{}{
		"name":   "张三",
		"age":    21,
		"status": 1,
		"info":   "更新后的学生信息",
	}
	if err := client.PutRequest("/v1/student/1", updateData); err != nil {
		fmt.Printf("PUT请求失败: %v\n", err)
	}

	// 示例4: DELETE请求
	fmt.Println("\n=== 示例4: DELETE请求 ===")
	if err := client.DeleteRequest("/v1/student/1"); err != nil {
		fmt.Printf("DELETE请求失败: %v\n", err)
	}

	// 示例5: 带认证的请求
	fmt.Println("\n=== 示例5: 带认证的请求 ===")
	token := "your-jwt-token-here"
	if err := client.RequestWithAuth("GET", "/v1/user/me", token, nil); err != nil {
		fmt.Printf("带认证的请求失败: %v\n", err)
	}

	fmt.Println("\n=== HTTP客户端示例完成 ===")
}
