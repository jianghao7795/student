#!/bin/bash

# RBAC中间件测试脚本

echo "=== RBAC中间件测试脚本 ==="
echo

# 设置基础URL
BASE_URL="http://localhost:8000"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local token=$3
    local expected_status=$4
    local description=$5
    
    echo -e "${YELLOW}测试: $description${NC}"
    echo "请求: $method $BASE_URL$endpoint"
    
    if [ -n "$token" ]; then
        echo "Token: ${token:0:20}..."
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $token" \
            -H "X-Request-Path: $endpoint" \
            -H "X-Request-Method: $method")
    else
        echo "Token: 无"
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
    fi
    
    # 分离响应体和状态码
    body=$(echo "$response" | head -n -1)
    status=$(echo "$response" | tail -n -1)
    
    echo "响应状态码: $status"
    echo "响应内容: $body"
    
    if [ "$status" = "$expected_status" ]; then
        echo -e "${GREEN}✓ 测试通过${NC}"
    else
        echo -e "${RED}✗ 测试失败 (期望: $expected_status, 实际: $status)${NC}"
    fi
    echo
}

# 检查服务是否运行
check_service() {
    echo "检查服务是否运行..."
    if curl -s "$BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 服务正在运行${NC}"
        return 0
    else
        echo -e "${RED}✗ 服务未运行，请先启动服务${NC}"
        echo "启动命令: make run"
        return 1
    fi
}

# 获取JWT token
get_token() {
    echo "获取JWT token..."
    
    # 尝试登录获取token
    login_response=$(curl -s -X POST "$BASE_URL/v1/user/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "admin",
            "password": "admin123"
        }')
    
    # 提取token
    token=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$token" ]; then
        echo -e "${GREEN}✓ 成功获取token${NC}"
        echo "Token: ${token:0:20}..."
        return 0
    else
        echo -e "${RED}✗ 获取token失败${NC}"
        echo "响应: $login_response"
        return 1
    fi
}

# 主测试流程
main() {
    echo "开始RBAC中间件测试..."
    echo
    
    # 检查服务状态
    if ! check_service; then
        exit 1
    fi
    
    # 获取token
    if ! get_token; then
        echo "使用空token进行测试..."
        token=""
    fi
    
    echo
    echo "=== 开始权限测试 ==="
    echo
    
    # 测试1: 无token访问受保护资源
    test_endpoint "GET" "/v1/users" "" "401" "无token访问用户列表"
    
    # 测试2: 有token访问受保护资源
    if [ -n "$token" ]; then
        test_endpoint "GET" "/v1/users" "$token" "200" "有token访问用户列表"
    fi
    
    # 测试3: 访问不需要权限的路径
    test_endpoint "POST" "/v1/user/login" "" "200" "访问登录接口（无需权限）"
    
    # 测试4: 权限检查API
    if [ -n "$token" ]; then
        test_endpoint "POST" "/v1/permissions/check" "$token" "200" "检查用户权限"
    fi
    
    # 测试5: 角色管理API
    if [ -n "$token" ]; then
        test_endpoint "GET" "/v1/roles" "$token" "200" "获取角色列表"
    fi
    
    echo
    echo "=== 测试完成 ==="
    echo
    echo "测试说明:"
    echo "1. 401状态码表示需要认证"
    echo "2. 403状态码表示权限不足"
    echo "3. 200状态码表示请求成功"
    echo "4. 500状态码表示服务器错误"
    echo
    echo "如果看到401错误，说明JWT认证中间件正常工作"
    echo "如果看到403错误，说明RBAC权限中间件正常工作"
    echo "如果看到200成功，说明用户有相应权限"
}

# 运行主函数
main
