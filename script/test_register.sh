#!/bin/bash
echo "测试用户注册接口..."
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser","email":"test@example.com","phone":"13800138000","password":"password123","age":25,"avatar":"https://example.com/avatar.jpg"}' http://localhost:8000/v1/user/register
