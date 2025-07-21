#!/bin/bash

# GetMe API 测试脚本
# 使用方法: ./test_getme.sh

BASE_URL="http://localhost:8000"

echo "=== GetMe API 测试 ==="
echo ""

# 1. 测试登录获取 token
echo "1. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/user/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }')

echo "登录响应: $LOGIN_RESPONSE"
echo ""

# 提取 token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取 token"
    exit 1
fi

echo "✅ 登录成功，获取到 token: ${TOKEN:0:20}..."
echo ""

# 2. 测试 GetMe 接口
echo "2. 测试 GetMe 接口..."
ME_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/user/me" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "GetMe 响应: $ME_RESPONSE"
echo ""

# 检查响应
if echo "$ME_RESPONSE" | grep -q '"success":true'; then
    echo "✅ GetMe 接口测试成功"
else
    echo "❌ GetMe 接口测试失败"
    exit 1
fi

echo ""
echo "=== 测试完成 ===" 