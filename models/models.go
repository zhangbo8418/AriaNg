package models

import (
	"time"
)

// User 用户表
type User struct {
	ID             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       string `json:"username" gorm:"unique;not null"`
	Name           string `json:"name" gorm:"not null"`
	IsAdmin        int    `json:"is_admin" gorm:"default:0"` // 1为管理员，0为普通用户
	EmployeeScope  string `json:"employee_scope"`            // 管理员工范围，员工IDs，0为全部员工
	Status         int    `json:"status" gorm:"default:1"`   // 1启用，0停用
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Company 公司表
type Company struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	Status    int       `json:"status" gorm:"default:1"` // 1启用，0停用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Department 部门表
type Department struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID int       `json:"company_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Status    int       `json:"status" gorm:"default:1"` // 1启用，0停用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 关联
	Company Company `json:"company" gorm:"foreignKey:CompanyID"`
}

// Position 岗位表
type Position struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	Status    int       `json:"status" gorm:"default:1"` // 1启用，0停用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Employee 员工表
type Employee struct {
	ID              int    `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID       int    `json:"company_id" gorm:"not null"`
	DepartmentID    int    `json:"department_id" gorm:"not null"`
	PositionID      int    `json:"position_id" gorm:"not null"`
	ProjectPermIDs  string `json:"project_perm_ids"` // 项目权限IDs，逗号分隔
	Name            string `json:"name" gorm:"not null"`
	Status          int    `json:"status" gorm:"default:1"` // 1启用，0停用
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	
	// 关联
	Company    Company    `json:"company" gorm:"foreignKey:CompanyID"`
	Department Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Position   Position   `json:"position" gorm:"foreignKey:PositionID"`
}

// ProjectPermission 项目权限表
type ProjectPermission struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Permission string    `json:"permission" gorm:"not null"`
	Status     int       `json:"status" gorm:"default:1"` // 1启用，0停用
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CommissionProject 提成项目表
type CommissionProject struct {
	ID               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	FieldName        string `json:"field_name" gorm:"not null"`
	ProjectPermID    int    `json:"project_perm_id" gorm:"not null"`
	Status           int    `json:"status" gorm:"default:1"` // 1启用，0停用
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	
	// 关联
	ProjectPermission ProjectPermission `json:"project_permission" gorm:"foreignKey:ProjectPermID"`
}

// DailyMonthlyReport 日/月报表
type DailyMonthlyReport struct {
	ID                  int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Date                string    `json:"date" gorm:"not null"`        // 日期（年月日）
	EmployeeID          int       `json:"employee_id" gorm:"not null"`
	CompanyID           int       `json:"company_id" gorm:"not null"`
	DepartmentID        int       `json:"department_id" gorm:"not null"`
	PositionID          int       `json:"position_id" gorm:"not null"`
	EmployeeName        string    `json:"employee_name" gorm:"not null"`
	CommissionProjectID int       `json:"commission_project_id" gorm:"not null"`
	CommissionValue     float64   `json:"commission_value" gorm:"not null"` // 提成项目值
	Status              int       `json:"status" gorm:"default:1"`          // 1启用，0停用
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	
	// 关联
	Employee          Employee          `json:"employee" gorm:"foreignKey:EmployeeID"`
	Company           Company           `json:"company" gorm:"foreignKey:CompanyID"`
	Department        Department        `json:"department" gorm:"foreignKey:DepartmentID"`
	Position          Position          `json:"position" gorm:"foreignKey:PositionID"`
	CommissionProject CommissionProject `json:"commission_project" gorm:"foreignKey:CommissionProjectID"`
}

// ReportSummary 报表汇总结构（用于月报统计）
type ReportSummary struct {
	EmployeeID          int     `json:"employee_id"`
	EmployeeName        string  `json:"employee_name"`
	CompanyID           int     `json:"company_id"`
	CompanyName         string  `json:"company_name"`
	DepartmentID        int     `json:"department_id"`
	DepartmentName      string  `json:"department_name"`
	PositionID          int     `json:"position_id"`
	PositionName        string  `json:"position_name"`
	CommissionProjectID int     `json:"commission_project_id"`
	ProjectName         string  `json:"project_name"`
	TotalValue          float64 `json:"total_value"`
}