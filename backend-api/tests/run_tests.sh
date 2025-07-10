#!/bin/bash

# 测试运行脚本
# 使用方法: ./tests/run_tests.sh [unit|integration|all]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Go是否安装
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go未安装，请先安装Go"
        exit 1
    fi
    print_message "Go版本: $(go version)"
}

# 安装测试依赖
install_deps() {
    print_message "安装测试依赖..."
    go mod tidy
    go get github.com/stretchr/testify
    print_success "依赖安装完成"
}

# 运行单元测试
run_unit_tests() {
    print_message "运行单元测试..."
    cd tests/unit
    go test -v -cover ./...
    cd ../..
    print_success "单元测试完成"
}

# 运行集成测试
run_integration_tests() {
    print_message "运行集成测试..."
    cd tests/integration
    go test -v -cover ./...
    cd ../..
    print_success "集成测试完成"
}

# 运行所有测试
run_all_tests() {
    print_message "运行所有测试..."
    
    # 运行单元测试
    run_unit_tests
    
    # 运行集成测试
    run_integration_tests
    
    print_success "所有测试完成"
}

# 生成测试覆盖率报告
generate_coverage() {
    print_message "生成测试覆盖率报告..."
    
    # 创建覆盖率目录
    mkdir -p coverage
    
    # 运行测试并生成覆盖率报告
    go test -v -coverprofile=coverage/coverage.out ./tests/...
    
    # 生成HTML报告
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    
    print_success "覆盖率报告已生成: coverage/coverage.html"
}

# 清理测试文件
cleanup() {
    print_message "清理测试文件..."
    rm -rf coverage/
    print_success "清理完成"
}

# 主函数
main() {
    check_go
    install_deps
    
    case "${1:-all}" in
        "unit")
            run_unit_tests
            ;;
        "integration")
            run_integration_tests
            ;;
        "all")
            run_all_tests
            ;;
        "coverage")
            generate_coverage
            ;;
        "clean")
            cleanup
            ;;
        "help"|"-h"|"--help")
            echo "使用方法: $0 [unit|integration|all|coverage|clean|help]"
            echo "  unit        - 运行单元测试"
            echo "  integration - 运行集成测试"
            echo "  all         - 运行所有测试（默认）"
            echo "  coverage    - 生成测试覆盖率报告"
            echo "  clean       - 清理测试文件"
            echo "  help        - 显示帮助信息"
            ;;
        *)
            print_error "未知参数: $1"
            echo "使用 '$0 help' 查看帮助信息"
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@" 
