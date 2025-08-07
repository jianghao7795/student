#!/bin/bash

# 微服务部署脚本
set -e

echo "🚀 开始部署学生管理系统微服务架构..."

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 构建项目
echo "📦 构建微服务..."
make build-microservices

# 启动服务
echo "🔧 启动微服务架构..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 30

# 检查服务状态
echo "🔍 检查服务状态..."
docker-compose ps

# 检查Nacos是否可用
echo "🔍 检查Nacos服务..."
if curl -s http://localhost:8848/nacos/v1/console/health/readiness | grep -q "UP"; then
    echo "✅ Nacos服务启动成功"
else
    echo "❌ Nacos服务启动失败"
    exit 1
fi

# 检查各个微服务是否可用
echo "🔍 检查微服务状态..."

# 检查网关服务
if curl -s http://localhost:8600/health | grep -q "OK"; then
    echo "✅ 网关服务启动成功"
else
    echo "❌ 网关服务启动失败"
fi

# 检查用户服务
if curl -s http://localhost:8601/health | grep -q "OK"; then
    echo "✅ 用户服务启动成功"
else
    echo "❌ 用户服务启动失败"
fi

# 检查学生服务
if curl -s http://localhost:8602/health | grep -q "OK"; then
    echo "✅ 学生服务启动成功"
else
    echo "❌ 学生服务启动失败"
fi

# 检查RBAC服务
if curl -s http://localhost:8603/health | grep -q "OK"; then
    echo "✅ RBAC服务启动成功"
else
    echo "❌ RBAC服务启动失败"
fi

echo ""
echo "🎉 微服务架构部署完成！"
echo ""
echo "📋 服务访问地址："
echo "   Nacos控制台: http://localhost:8848/nacos (用户名/密码: nacos/nacos)"
echo "   API网关: http://localhost:8600"
echo "   用户服务: http://localhost:8601"
echo "   学生服务: http://localhost:8602"
echo "   RBAC服务: http://localhost:8603"
echo ""
echo "🔧 管理命令："
echo "   查看服务状态: docker-compose ps"
echo "   查看服务日志: docker-compose logs -f [service-name]"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart [service-name]"
echo "" 