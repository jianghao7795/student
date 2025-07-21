#!/bin/bash
echo "=== 运行测试套件 ==="
echo "1. 单元测试..."
go test ./tests/unit/ -v
echo "2. 集成测试..."
go test ./tests/integration/ -v
echo "3. 性能测试..."
go test -bench=. ./tests/benchmark/ -v
echo "4. 端到端测试..."
go test ./tests/e2e/ -v
echo "=== 测试完成 ==="
