package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	v1 "student/api/student/v1"
)

// StudentAPIExample 学生API使用示例
type StudentAPIExample struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewStudentAPIExample 创建学生API示例实例
func NewStudentAPIExample(baseURL string) *StudentAPIExample {
	return &StudentAPIExample{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetToken 设置认证token
func (e *StudentAPIExample) SetToken(token string) {
	e.Token = token
}

// ListStudents 获取学生列表示例
func (e *StudentAPIExample) ListStudents(page, pageSize int32) error {
	url := fmt.Sprintf("%s/v1/students?page=%d&page_size=%d", e.BaseURL, page, pageSize)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.Token)
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("list students request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("list students failed with status %d: %s", resp.StatusCode, string(body))
	}
	var listResp v1.ListStudentsReply
	if err := json.Unmarshal(body, &listResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}
	fmt.Printf("学生列表 (总数: %d):\n", listResp.Total)
	for _, student := range listResp.Data {
		fmt.Printf("  ID: %d, 姓名: %s, 年龄: %d, 状态: %d, 信息: %s\n",
			student.Id, student.Name, student.Age, student.Status, student.Info)
	}
	return nil
}

// GetStudent 获取学生详情示例
func (e *StudentAPIExample) GetStudent(id int32) error {
	url := fmt.Sprintf("%s/v1/student/%d", e.BaseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.Token)
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("get student request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get student failed with status %d: %s", resp.StatusCode, string(body))
	}
	var studentResp v1.GetStudentReply
	if err := json.Unmarshal(body, &studentResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}
	fmt.Printf("学生详情:\n")
	fmt.Printf("  ID: %d\n", studentResp.Id)
	fmt.Printf("  姓名: %s\n", studentResp.Name)
	fmt.Printf("  年龄: %d\n", studentResp.Age)
	fmt.Printf("  状态: %d\n", studentResp.Status)
	fmt.Printf("  信息: %s\n", studentResp.Info)
	fmt.Printf("  创建时间: %s\n", studentResp.CreatedAt)
	fmt.Printf("  更新时间: %s\n", studentResp.UpdatedAt)
	return nil
}

// CreateStudent 创建学生示例
func (e *StudentAPIExample) CreateStudent(name string, age int32, status int32, info string) error {
	studentData := map[string]interface{}{
		"name":   name,
		"age":    age,
		"status": status,
		"info":   info,
	}
	jsonData, err := json.Marshal(studentData)
	if err != nil {
		return fmt.Errorf("marshal student data failed: %v", err)
	}
	req, err := http.NewRequest("POST", e.BaseURL+"/v1/student", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.Token)
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("create student request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create student failed with status %d: %s", resp.StatusCode, string(body))
	}
	var createResp v1.CreateStudentReply
	if err := json.Unmarshal(body, &createResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}
	fmt.Printf("创建学生成功，消息: %s\n", createResp.Message)
	return nil
}

// UpdateStudent 更新学生示例
func (e *StudentAPIExample) UpdateStudent(id int32, name string, age int32, status int32, info string) error {
	studentData := map[string]interface{}{
		"name":   name,
		"age":    age,
		"status": status,
		"info":   info,
	}
	jsonData, err := json.Marshal(studentData)
	if err != nil {
		return fmt.Errorf("marshal student data failed: %v", err)
	}
	url := fmt.Sprintf("%s/v1/student/%d", e.BaseURL, id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.Token)
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("update student request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update student failed with status %d: %s", resp.StatusCode, string(body))
	}
	var updateResp v1.UpdateStudentReply
	if err := json.Unmarshal(body, &updateResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}
	fmt.Printf("更新学生成功，消息: %s\n", updateResp.Message)
	return nil
}

// DeleteStudent 删除学生示例
func (e *StudentAPIExample) DeleteStudent(id int32) error {
	url := fmt.Sprintf("%s/v1/student/%d", e.BaseURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.Token)
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete student request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete student failed with status %d: %s", resp.StatusCode, string(body))
	}
	var deleteResp v1.DeleteStudentReply
	if err := json.Unmarshal(body, &deleteResp); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}
	fmt.Printf("删除学生成功，消息: %s\n", deleteResp.Message)
	return nil
}

func main() {
	// 用法说明
	fmt.Println("请先通过用户API获取JWT Token，并设置到环境变量 STUDENT_API_TOKEN")
	token := os.Getenv("STUDENT_API_TOKEN")
	if token == "" {
		fmt.Println("未检测到token，跳过示例。请先设置 STUDENT_API_TOKEN 环境变量。")
		return
	}
	api := NewStudentAPIExample("http://localhost:8000")
	api.SetToken(token)

	// 示例1: 获取学生列表
	fmt.Println("=== 示例1: 获取学生列表 ===")
	if err := api.ListStudents(1, 10); err != nil {
		fmt.Printf("获取学生列表失败: %v\n", err)
	}

	// 示例2: 获取学生详情
	fmt.Println("\n=== 示例2: 获取学生详情 ===")
	if err := api.GetStudent(1); err != nil {
		fmt.Printf("获取学生详情失败: %v\n", err)
	}

	// 示例3: 创建新学生
	fmt.Println("\n=== 示例3: 创建新学生 ===")
	if err := api.CreateStudent("张三", 20, 1, "计算机科学专业学生"); err != nil {
		fmt.Printf("创建学生失败: %v\n", err)
	}

	// 示例4: 更新学生信息
	fmt.Println("\n=== 示例4: 更新学生信息 ===")
	if err := api.UpdateStudent(1, "张三", 21, 1, "计算机科学专业学生，成绩优秀"); err != nil {
		fmt.Printf("更新学生失败: %v\n", err)
	}

	// 示例5: 删除学生
	fmt.Println("\n=== 示例5: 删除学生 ===")
	if err := api.DeleteStudent(1); err != nil {
		fmt.Printf("删除学生失败: %v\n", err)
	}

	fmt.Println("\n=== 学生API示例完成 ===")
}
