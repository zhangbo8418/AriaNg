# 员工提成数据提交系统

一个基于Go语言开发的员工提成数据管理系统，支持用户管理、公司组织架构管理、员工管理、提成项目管理和日/月报表管理。

## 功能特性

### 用户管理
- 用户增删改查
- 用户启用/停用
- 管理员和普通用户权限区分
- 员工管理范围控制

### 组织架构管理
- 公司管理（增删改查，启用停用）
- 部门管理（增删改查，启用停用）
- 岗位管理（增删改查，启用停用）
- 员工管理（增删改查，启用停用）

### 权限和项目管理
- 项目权限管理（增删改查，启用停用）
- 提成项目管理（增删改查，启用停用）

### 报表管理
- 日/月报表管理（增删改查，启用停用）
- 批量数据提交
- 月度汇总统计
- 权限范围内的数据查看

## 技术栈

- **后端框架**: Gin (Go)
- **数据库**: SQLite
- **ORM**: GORM
- **认证**: JWT
- **API**: RESTful API

## 项目结构

```
employee-commission-system/
├── config/                 # 配置文件
│   └── database.go         # 数据库配置
├── controllers/            # 控制器
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── company_controller.go
│   ├── department_controller.go
│   ├── position_controller.go
│   ├── employee_controller.go
│   ├── project_permission_controller.go
│   ├── commission_project_controller.go
│   └── report_controller.go
├── middleware/             # 中间件
│   ├── auth.go
│   └── employee_scope.go
├── models/                 # 数据模型
│   └── models.go
├── routes/                 # 路由
│   └── routes.go
├── scripts/                # 脚本
│   └── init_data.go
├── utils/                  # 工具函数
│   ├── auth.go
│   └── response.go
├── main.go                 # 入口文件
├── go.mod                  # Go模块文件
└── README.md              # 说明文档
```

## 数据库表结构

### 用户表 (users)
- ID, 用户名, 名称, 后台管理权限, 管理员工范围, 状态

### 公司表 (companies)
- ID, 公司名, 状态

### 部门表 (departments)
- ID, 公司ID, 部门名, 状态

### 岗位表 (positions)
- ID, 名称, 状态

### 员工表 (employees)
- ID, 公司ID, 部门ID, 岗位ID, 项目权限IDs, 姓名, 状态

### 项目权限表 (project_permissions)
- ID, 权限名称, 状态

### 提成项目表 (commission_projects)
- ID, 字段名, 项目权限ID, 状态

### 日/月报表 (daily_monthly_reports)
- ID, 日期, 员工ID, 公司ID, 部门ID, 岗位ID, 姓名, 提成项目ID, 提成项目值, 状态

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 初始化示例数据

```bash
go run scripts/init_data.go
```

### 3. 启动服务器

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 接口

### 认证接口

#### 用户登录
- **POST** `/api/v1/login`
- **请求体**:
```json
{
  "username": "admin",
  "password": "password"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "jwt_token_here",
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

### 管理员接口

所有管理员接口都需要在请求头中包含 `Authorization: Bearer <token>`

#### 用户管理
- **GET** `/api/v1/admin/users` - 获取用户列表
- **POST** `/api/v1/admin/users` - 创建用户
- **GET** `/api/v1/admin/users/:id` - 获取单个用户
- **PUT** `/api/v1/admin/users/:id` - 更新用户
- **DELETE** `/api/v1/admin/users/:id` - 删除用户
- **PUT** `/api/v1/admin/users/:id/toggle-status` - 切换用户状态

#### 公司管理
- **GET** `/api/v1/admin/companies` - 获取公司列表
- **POST** `/api/v1/admin/companies` - 创建公司
- **GET** `/api/v1/admin/companies/:id` - 获取单个公司
- **PUT** `/api/v1/admin/companies/:id` - 更新公司
- **DELETE** `/api/v1/admin/companies/:id` - 删除公司
- **PUT** `/api/v1/admin/companies/:id/toggle-status` - 切换公司状态

#### 部门管理
- **GET** `/api/v1/admin/departments` - 获取部门列表
- **POST** `/api/v1/admin/departments` - 创建部门
- **GET** `/api/v1/admin/departments/:id` - 获取单个部门
- **PUT** `/api/v1/admin/departments/:id` - 更新部门
- **DELETE** `/api/v1/admin/departments/:id` - 删除部门
- **PUT** `/api/v1/admin/departments/:id/toggle-status` - 切换部门状态

#### 岗位管理
- **GET** `/api/v1/admin/positions` - 获取岗位列表
- **POST** `/api/v1/admin/positions` - 创建岗位
- **GET** `/api/v1/admin/positions/:id` - 获取单个岗位
- **PUT** `/api/v1/admin/positions/:id` - 更新岗位
- **DELETE** `/api/v1/admin/positions/:id` - 删除岗位
- **PUT** `/api/v1/admin/positions/:id/toggle-status` - 切换岗位状态

#### 员工管理
- **GET** `/api/v1/admin/employees` - 获取员工列表
- **POST** `/api/v1/admin/employees` - 创建员工
- **GET** `/api/v1/admin/employees/:id` - 获取单个员工
- **PUT** `/api/v1/admin/employees/:id` - 更新员工
- **DELETE** `/api/v1/admin/employees/:id` - 删除员工
- **PUT** `/api/v1/admin/employees/:id/toggle-status` - 切换员工状态

#### 项目权限管理
- **GET** `/api/v1/admin/project-permissions` - 获取项目权限列表
- **POST** `/api/v1/admin/project-permissions` - 创建项目权限
- **GET** `/api/v1/admin/project-permissions/:id` - 获取单个项目权限
- **PUT** `/api/v1/admin/project-permissions/:id` - 更新项目权限
- **DELETE** `/api/v1/admin/project-permissions/:id` - 删除项目权限
- **PUT** `/api/v1/admin/project-permissions/:id/toggle-status` - 切换项目权限状态

#### 提成项目管理
- **GET** `/api/v1/admin/commission-projects` - 获取提成项目列表
- **POST** `/api/v1/admin/commission-projects` - 创建提成项目
- **GET** `/api/v1/admin/commission-projects/:id` - 获取单个提成项目
- **PUT** `/api/v1/admin/commission-projects/:id` - 更新提成项目
- **DELETE** `/api/v1/admin/commission-projects/:id` - 删除提成项目
- **PUT** `/api/v1/admin/commission-projects/:id/toggle-status` - 切换提成项目状态

#### 报表管理
- **GET** `/api/v1/admin/reports` - 获取报表列表
- **GET** `/api/v1/admin/reports/:id` - 获取单个报表
- **PUT** `/api/v1/admin/reports/:id` - 更新报表
- **DELETE** `/api/v1/admin/reports/:id` - 删除报表
- **PUT** `/api/v1/admin/reports/:id/toggle-status` - 切换报表状态
- **GET** `/api/v1/admin/reports/monthly-summary?month=2024-01` - 获取月度汇总

### 普通用户接口

普通用户接口受员工管理范围限制，只能操作权限范围内的数据。

#### 报表提交
- **POST** `/api/v1/reports` - 创建报表记录
- **POST** `/api/v1/reports/batch` - 批量创建报表记录

#### 数据查看
- **GET** `/api/v1/my-reports` - 获取我的报表列表
- **GET** `/api/v1/my-reports/:id` - 获取单个报表
- **GET** `/api/v1/my-reports/monthly-summary?month=2024-01` - 获取月度汇总

#### 基础数据查看（只读）
- **GET** `/api/v1/companies` - 获取公司列表
- **GET** `/api/v1/departments` - 获取部门列表
- **GET** `/api/v1/positions` - 获取岗位列表
- **GET** `/api/v1/employees` - 获取员工列表
- **GET** `/api/v1/project-permissions` - 获取项目权限列表
- **GET** `/api/v1/commission-projects` - 获取提成项目列表

## 默认账号

系统初始化后会创建以下默认账号：

- **管理员账号**: 
  - 用户名: `admin`
  - 密码: `password`
  - 权限: 管理员，可管理所有数据

- **普通用户1**: 
  - 用户名: `user1`
  - 密码: `password`
  - 权限: 可管理员工ID为1,2的数据

- **普通用户2**: 
  - 用户名: `user2`
  - 密码: `password`
  - 权限: 可管理员工ID为3,4,5的数据

## 查询参数

大部分列表接口支持以下查询参数：

- `page`: 页码，默认为1
- `page_size`: 每页大小，默认为10
- `search`: 搜索关键词
- `status`: 状态筛选（1为启用，0为停用）
- 其他特定筛选参数（如company_id, department_id等）

## 权限说明

1. **管理员权限**: 
   - 可以管理所有用户、公司、部门、岗位、员工、项目权限、提成项目
   - 可以查看和管理所有报表数据

2. **普通用户权限**: 
   - 只能查看基础数据（公司、部门、岗位等）
   - 只能在自己的员工管理范围内提交和查看报表数据
   - 员工管理范围由用户的`employee_scope`字段控制

3. **员工管理范围**: 
   - `0`: 表示可以管理所有员工
   - `1,2,3`: 表示只能管理员工ID为1、2、3的数据

## 注意事项

1. 所有API接口都返回统一格式的JSON响应
2. 删除操作会检查关联数据，有关联时不允许删除
3. 状态字段：1表示启用，0表示停用
4. 日期格式：YYYY-MM-DD
5. JWT令牌有效期为24小时
6. 系统使用SQLite数据库，数据文件为`commission.db`