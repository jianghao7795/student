#!/bin/bash

# gRPC RBAC 测试脚本

echo "=== gRPC RBAC 权限过滤测试 ==="

# 检查服务器是否正在运行
echo "1. 检查服务器状态..."
if pgrep -f "student" > /dev/null; then
    echo "   服务器正在运行"
else
    echo "   服务器未运行，请先启动服务器"
    echo "   启动命令: ./student -conf configs/config.yaml"
    exit 1
fi

# 等待服务器完全启动
echo "2. 等待服务器启动完成..."
sleep 3

# 测试 gRPC 连接
echo "3. 测试 gRPC 连接..."
if command -v grpcurl > /dev/null; then
    echo "   使用 grpcurl 测试连接..."
    grpcurl -plaintext localhost:9600 list
else
    echo "   grpcurl 未安装，跳过连接测试"
fi

# 运行 Go 测试程序
echo "4. 运行 gRPC RBAC 测试程序..."
if [ -f "examples/grpc/grpc_rbac/main.go" ]; then
    cd examples/grpc/grpc_rbac
    go run grpc_rbac_test.go
    cd ../..
else
    echo "   测试文件不存在: examples/grpc/grpc_rbac/main.go"
fi

echo "=== 测试完成 ==="
echo ""
echo "注意事项："
echo "1. 确保 configs/config.yaml 中 rbac.enabled = true"
echo "2. 确保 rbac_model.conf 和 rbac_policy.csv 文件存在"
echo "3. 测试前需要先获取有效的 JWT Token"
echo "4. 修改测试文件中的 Token 为实际有效的 Token"
