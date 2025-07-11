#!/bin/bash

# 测试运行脚本
# 支持运行不同类型的测试

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Go版本
check_go_version() {
    log_info "Go版本: $(go version)"
}

# 安装测试依赖
install_dependencies() {
    log_info "安装测试依赖..."
    
    # 安装yaml解析库
    go get gopkg.in/yaml.v3
    
    # 安装其他测试依赖
    go mod tidy
    
    log_success "依赖安装完成"
}

# 运行单元测试
run_unit_tests() {
    log_info "运行单元测试..."
    
    # 运行所有单元测试
    go test -v ./tests/unit/... -coverprofile=coverage/unit.out
    
    log_success "单元测试完成"
}

# 运行集成测试
run_integration_tests() {
    log_info "运行集成测试..."
    
    # 运行所有集成测试
    go test -v ./tests/integration/... -coverprofile=coverage/integration.out
    
    log_success "集成测试完成"
}

# 运行token过期测试
run_token_expired_tests() {
    log_info "运行token过期测试..."
    
    # 运行token过期相关测试
    go test -v ./tests/integration/ -run TestTokenExpiredIntegrationTestSuite
    
    log_success "token过期测试完成"
}

# 运行真实用户测试
run_real_user_tests() {
    log_info "运行真实用户测试..."
    
    # 检查配置文件是否存在
    if [ ! -f "tests/fixtures/test_users.yml" ]; then
        log_error "测试配置文件 tests/fixtures/test_users.yml 不存在"
        exit 1
    fi
    
    # 运行真实用户测试
    go test -v ./tests/integration/ -run TestRealUserIntegrationTestSuite
    
    log_success "真实用户测试完成"
}

# 运行所有测试
run_all_tests() {
    log_info "运行所有测试..."
    
    # 运行单元测试
    run_unit_tests
    
    # 运行集成测试
    run_integration_tests
    
    # 运行真实用户测试
    run_real_user_tests
    
    log_success "所有测试完成"
}

# 生成测试覆盖率报告
generate_coverage_report() {
    log_info "生成测试覆盖率报告..."
    
    # 创建覆盖率目录
    mkdir -p coverage
    
    # 合并覆盖率文件
    if [ -f "coverage/unit.out" ] && [ -f "coverage/integration.out" ]; then
        go tool cover -html=coverage/unit.out -o coverage/coverage.html
        log_success "覆盖率报告已生成: coverage/coverage.html"
    else
        log_warning "覆盖率文件不存在，跳过覆盖率报告生成"
    fi
}

# 清理测试数据
cleanup_test_data() {
    log_info "清理测试数据..."
    
    # 删除测试生成的数据库文件
    rm -f tests/integration/db/test.db
    rm -f tests/unit/db/test.db
    
    # 删除测试日志
    rm -rf tests/integration/logs/*
    rm -rf tests/unit/logs/*
    
    log_success "测试数据清理完成"
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  unit              运行单元测试"
    echo "  integration       运行集成测试"
    echo "  token-expired     运行token过期测试"
    echo "  real-user         运行真实用户测试"
    echo "  all               运行所有测试"
    echo "  coverage          生成测试覆盖率报告"
    echo "  cleanup           清理测试数据"
    echo "  help              显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 unit                    # 运行单元测试"
    echo "  $0 integration             # 运行集成测试"
    echo "  $0 real-user               # 运行真实用户测试"
    echo "  $0 all                     # 运行所有测试"
    echo "  $0 coverage                # 生成覆盖率报告"
}

# 主函数
main() {
    # 检查是否在正确的目录
    if [ ! -f "go.mod" ]; then
        log_error "请在项目根目录运行此脚本"
        exit 1
    fi
    
    # 检查Go版本
    check_go_version
    
    # 安装依赖
    install_dependencies
    
    # 根据参数执行相应操作
    case "${1:-help}" in
        "unit")
            run_unit_tests
            ;;
        "integration")
            run_integration_tests
            ;;
        "token-expired")
            run_token_expired_tests
            ;;
        "real-user")
            run_real_user_tests
            ;;
        "all")
            run_all_tests
            ;;
        "coverage")
            generate_coverage_report
            ;;
        "cleanup")
            cleanup_test_data
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# 执行主函数
main "$@" 
