#!/bin/bash

# 测试student服务的JWT和RBAC功能

echo "=== 测试Student服务的JWT和RBAC功能 ==="

# 服务地址
STUDENT_SERVICE_URL="http://localhost:8602"
GATEWAY_URL="http://localhost:8600"

echo "1. 测试健康检查端点（无需JWT验证）"
echo "----------------------------------------"
curl -X GET "${STUDENT_SERVICE_URL}/v1/students/health" -H "Content-Type: application/json"
echo -e "\n"

echo "2. 测试健康检查端点（通过网关，无需JWT验证）"
echo "----------------------------------------"
curl -X GET "${GATEWAY_URL}/v1/students/health" -H "Content-Type: application/json"
echo -e "\n"

echo "3. 测试需要JWT验证的端点（无token，应该返回401）"
echo "----------------------------------------"
curl -X GET "${STUDENT_SERVICE_URL}/v1/student/1" -H "Content-Type: application/json"
echo -e "\n"

echo "4. 测试需要JWT验证的端点（通过网关，无token，应该返回401）"
echo "----------------------------------------"
curl -X GET "${GATEWAY_URL}/v1/student/1" -H "Content-Type: application/json"
echo -e "\n"

echo "5. 测试需要JWT验证的端点（无效token，应该返回401）"
echo "----------------------------------------"
curl -X GET "${STUDENT_SERVICE_URL}/v1/student/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid-token"
echo -e "\n"

echo "6. 测试需要JWT验证的端点（通过网关，无效token，应该返回401）"
echo "----------------------------------------"
curl -X GET "${GATEWAY_URL}/v1/student/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid-token"
echo -e "\n"

echo "=== 测试完成 ==="
echo ""
echo "注意："
echo "- 健康检查端点应该返回200状态码"
echo "- 需要认证的端点应该返回401状态码（未授权）"
echo "- 如果服务未启动，请先启动相关服务"
echo "- 要测试完整的JWT功能，需要先通过用户服务获取有效的token"
