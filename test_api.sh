#!/bin/bash

# 员工提成系统API测试脚本

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_header() {
    echo -e "${YELLOW}========== $1 ==========${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 测试登录
test_login() {
    print_header "测试用户登录"
    
    # 管理员登录
    response=$(curl -s -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d '{"username": "admin", "password": "password"}')
    
    echo "登录响应: $response"
    
    TOKEN=$(echo $response | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$TOKEN" ]; then
        print_success "管理员登录成功，Token: ${TOKEN:0:20}..."
    else
        print_error "管理员登录失败"
        exit 1
    fi
}

# 测试用户管理
test_user_management() {
    print_header "测试用户管理"
    
    # 获取用户列表
    response=$(curl -s -X GET "$BASE_URL/admin/users" \
        -H "Authorization: Bearer $TOKEN")
    echo "用户列表: $response"
    
    # 创建新用户
    response=$(curl -s -X POST "$BASE_URL/admin/users" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "name": "测试用户", "is_admin": 0, "employee_scope": "1", "status": 1}')
    echo "创建用户: $response"
}

# 测试公司管理
test_company_management() {
    print_header "测试公司管理"
    
    # 获取公司列表
    response=$(curl -s -X GET "$BASE_URL/admin/companies" \
        -H "Authorization: Bearer $TOKEN")
    echo "公司列表: $response"
    
    # 创建新公司
    response=$(curl -s -X POST "$BASE_URL/admin/companies" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"name": "测试公司", "status": 1}')
    echo "创建公司: $response"
}

# 测试员工管理
test_employee_management() {
    print_header "测试员工管理"
    
    # 获取员工列表
    response=$(curl -s -X GET "$BASE_URL/admin/employees" \
        -H "Authorization: Bearer $TOKEN")
    echo "员工列表: $response"
}

# 测试报表管理
test_report_management() {
    print_header "测试报表管理"
    
    # 获取报表列表
    response=$(curl -s -X GET "$BASE_URL/admin/reports" \
        -H "Authorization: Bearer $TOKEN")
    echo "报表列表: $response"
    
    # 创建报表记录
    response=$(curl -s -X POST "$BASE_URL/reports" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"date": "2024-01-20", "employee_id": 1, "commission_project_id": 1, "commission_value": 500.0}')
    echo "创建报表: $response"
    
    # 获取月度汇总
    response=$(curl -s -X GET "$BASE_URL/admin/reports/monthly-summary?month=2024-01" \
        -H "Authorization: Bearer $TOKEN")
    echo "月度汇总: $response"
}

# 测试普通用户权限
test_user_permissions() {
    print_header "测试普通用户权限"
    
    # 普通用户登录
    response=$(curl -s -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d '{"username": "user1", "password": "password"}')
    
    USER_TOKEN=$(echo $response | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$USER_TOKEN" ]; then
        print_success "普通用户登录成功"
        
        # 尝试访问管理员接口（应该失败）
        response=$(curl -s -X GET "$BASE_URL/admin/users" \
            -H "Authorization: Bearer $USER_TOKEN")
        echo "普通用户访问管理员接口: $response"
        
        # 访问自己的报表
        response=$(curl -s -X GET "$BASE_URL/my-reports" \
            -H "Authorization: Bearer $USER_TOKEN")
        echo "普通用户查看自己的报表: $response"
    else
        print_error "普通用户登录失败"
    fi
}

# 主函数
main() {
    echo "开始测试员工提成系统API..."
    echo "请确保服务器正在运行在 $BASE_URL"
    echo ""
    
    test_login
    test_user_management
    test_company_management
    test_employee_management
    test_report_management
    test_user_permissions
    
    echo ""
    print_success "API测试完成"
}

# 运行测试
main