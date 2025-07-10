#!/bin/bash

echo "测试日志轮转功能..."

# 检查日志目录
echo "1. 检查日志目录结构..."
ls -la logs/

echo -e "\n2. 测试日志写入..."
echo "当前时间: $(date)"

# 发送一些测试请求来生成日志
echo "3. 发送测试请求生成日志..."

# 测试注册接口
curl -X POST http://localhost:8890/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "TestPass123",
    "confirmPassword": "TestPass123",
    "code": "123456"
  }' -s > /dev/null

# 测试登录接口
curl -X POST http://localhost:8890/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "TestPass123"
  }' -s > /dev/null

echo "4. 检查日志文件..."
echo "main.log 内容（最后10行）:"
tail -10 logs/main.log

echo -e "\n5. 检查日志文件大小和修改时间..."
ls -lh logs/

echo -e "\n6. 手动触发日志轮转..."
# 这里可以添加手动轮转的逻辑

echo -e "\n测试完成！"
echo "日志轮转说明："
echo "- 当天的日志会写入 logs/main.log"
echo "- 历史日志会按日期归档为 logs/main.log.YYYY-MM-DD"
echo "- 日志文件超过100MB会自动分割"
echo "- 保留最近30天的日志文件"
echo "- 旧日志文件会被压缩" 
