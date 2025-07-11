#!/bin/bash

# Token过期测试单元运行脚本
echo "=== 运行Token过期测试单元 ==="

# 检查是否在正确的目录
if [ ! -f "../go.mod" ]; then
    echo "❌ 请在tests目录下运行此脚本"
    exit 1
fi

# 检查run_tests.sh是否存在
if [ ! -f "run_tests.sh" ]; then
    echo "❌ run_tests.sh不存在"
    exit 1
fi

# 通过run_tests.sh运行token过期测试
echo "通过run_tests.sh运行token过期测试..."
./run_tests.sh token-expired

# 检查测试结果
if [ $? -eq 0 ]; then
    echo "✅ Token过期测试通过"
else
    echo "❌ Token过期测试失败"
    exit 1
fi

echo ""
echo "=== 测试完成 ==="
echo "如需运行所有测试，请使用: ./run_tests.sh all"
echo "如需运行其他测试，请使用: ./run_tests.sh help" 
