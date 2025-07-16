#!/bin/bash

# 测试用户管理脚本的功能

echo "=== 用户管理脚本测试 ==="
echo

# 测试1: 显示帮助信息
echo "测试1: 显示帮助信息"
./scripts/run_user_management.sh -help
echo

# 测试2: 测试无效参数
echo "测试2: 测试无效参数"
./scripts/run_user_management.sh -action invalid_action
echo

# 测试3: 测试缺少用户ID参数
echo "测试3: 测试缺少用户ID参数"
./scripts/run_user_management.sh -action getUser
echo

echo "=== 测试完成 ==="
echo "注意: 要测试完整功能，请确保："
echo "1. 数据库中有用户数据"
echo "2. 数据库连接配置正确"
echo
echo "完整测试命令："
echo "  # 获取用户信息"
echo "  ./scripts/run_user_management.sh -action getUser -userid 1"
echo
echo "  # 添加测试数据"
echo "  ./scripts/run_user_management.sh -action addTestData" 
