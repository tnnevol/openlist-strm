#!/bin/bash

# 用户管理脚本运行器
# 用法: ./run_user_management.sh [action] [userid]

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 切换到项目根目录
cd "$PROJECT_ROOT"

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go命令，请先安装Go"
    exit 1
fi

# 检查脚本文件是否存在
if [ ! -f "scripts/user_management.go" ]; then
    echo "错误: 未找到 user_management.go 脚本文件"
    exit 1
fi

# 运行脚本
echo "运行用户管理脚本..."
go run scripts/user_management.go "$@" 
