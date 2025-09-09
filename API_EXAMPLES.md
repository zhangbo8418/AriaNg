# 员工提成系统 API 使用示例

## 快速开始

### 1. 启动系统

```bash
# 初始化示例数据（只需执行一次）
go run scripts/init_data.go

# 启动服务器
go run main.go
```

服务器启动后会监听在 `http://localhost:8080`

### 2. 默认账号

- **管理员**: `admin` / `password`
- **普通用户1**: `user1` / `password` (可管理员工ID: 1,2)
- **普通用户2**: `user2` / `password` (可管理员工ID: 3,4,5)

## API 使用示例

### 用户登录

```bash
curl -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "name": "管理员",
      "is_admin": 1,
      "employee_scope": "0",
      "status": 1
    }
  }
}
```

### 管理员接口示例

**获取公司列表**:
```bash
curl -X GET "http://localhost:8080/api/v1/admin/companies" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**创建员工**:
```bash
curl -X POST "http://localhost:8080/api/v1/admin/employees" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": 1,
    "department_id": 1,
    "position_id": 1,
    "project_perm_ids": "1,2",
    "name": "新员工",
    "status": 1
  }'
```

**获取月度汇总**:
```bash
curl -X GET "http://localhost:8080/api/v1/admin/reports/monthly-summary?month=2024-01" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 普通用户接口示例

**提交报表数据**:
```bash
curl -X POST "http://localhost:8080/api/v1/reports" \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-20",
    "employee_id": 1,
    "commission_project_id": 1,
    "commission_value": 500.0
  }'
```

**批量提交报表数据**:
```bash
curl -X POST "http://localhost:8080/api/v1/reports/batch" \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "date": "2024-01-20",
      "employee_id": 1,
      "commission_project_id": 1,
      "commission_value": 100.0
    },
    {
      "date": "2024-01-20",
      "employee_id": 1,
      "commission_project_id": 2,
      "commission_value": 200.0
    }
  ]'
```

**查看自己管理的报表**:
```bash
curl -X GET "http://localhost:8080/api/v1/my-reports" \
  -H "Authorization: Bearer USER_TOKEN"
```

**查看月度汇总**:
```bash
curl -X GET "http://localhost:8080/api/v1/my-reports/monthly-summary?month=2024-01" \
  -H "Authorization: Bearer USER_TOKEN"
```

## 权限说明

### 管理员权限
- 可以管理所有用户、公司、部门、岗位、员工
- 可以管理项目权限和提成项目
- 可以查看和管理所有报表数据

### 普通用户权限
- 只能查看基础数据（公司、部门、岗位等）
- 只能在自己的员工管理范围内提交和查看报表数据
- 员工管理范围由用户的 `employee_scope` 字段控制
  - `0`: 可以管理所有员工
  - `1,2,3`: 只能管理员工ID为1、2、3的数据

## 查询参数

大部分列表接口支持以下查询参数：

- `page`: 页码，默认为1
- `page_size`: 每页大小，默认为10
- `search`: 搜索关键词
- `status`: 状态筛选（1为启用，0为停用）
- `company_id`: 公司ID筛选
- `department_id`: 部门ID筛选
- `employee_id`: 员工ID筛选
- `start_date`: 开始日期筛选
- `end_date`: 结束日期筛选

**示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/admin/employees?company_id=1&page=1&page_size=5" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 错误处理

所有API都返回统一格式的响应：

**成功响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "错误信息"
}
```

常见错误码：
- `400`: 请求参数错误
- `401`: 未授权（token无效或过期）
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

## 测试脚本

系统提供了自动化测试脚本：

```bash
# 运行API测试
./test_api.sh
```

该脚本会测试主要的API功能，包括：
- 用户登录
- 管理员权限接口
- 普通用户权限接口
- 数据创建和查询
- 权限控制验证