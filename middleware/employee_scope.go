package middleware

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// EmployeeScopeMiddleware 员工范围权限中间件
// 用于限制用户只能查看和操作其管理范围内的员工数据
func EmployeeScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Unauthorized(c, "未授权")
			c.Abort()
			return
		}

		isAdmin, exists := c.Get("is_admin")
		if !exists {
			utils.Unauthorized(c, "未授权")
			c.Abort()
			return
		}

		// 如果是管理员，跳过权限检查
		if isAdmin.(bool) {
			c.Next()
			return
		}

		// 获取当前用户信息
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			utils.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}

		// 解析用户的员工管理范围
		employeeScope := user.EmployeeScope
		var allowedEmployeeIDs []int

		if employeeScope == "0" {
			// 0表示可以管理所有员工
			var allEmployees []models.Employee
			config.DB.Select("id").Find(&allEmployees)
			for _, emp := range allEmployees {
				allowedEmployeeIDs = append(allowedEmployeeIDs, emp.ID)
			}
		} else if employeeScope != "" {
			// 解析员工ID列表
			idStrings := strings.Split(employeeScope, ",")
			for _, idStr := range idStrings {
				if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
					allowedEmployeeIDs = append(allowedEmployeeIDs, id)
				}
			}
		}

		// 将允许的员工ID列表存储到上下文中
		c.Set("allowed_employee_ids", allowedEmployeeIDs)

		c.Next()
	}
}

// CheckEmployeeAccess 检查员工访问权限
func CheckEmployeeAccess(c *gin.Context, employeeID int) bool {
	isAdmin, exists := c.Get("is_admin")
	if exists && isAdmin.(bool) {
		return true // 管理员有所有权限
	}

	allowedIDs, exists := c.Get("allowed_employee_ids")
	if !exists {
		return false
	}

	allowedEmployeeIDs := allowedIDs.([]int)
	for _, id := range allowedEmployeeIDs {
		if id == employeeID {
			return true
		}
	}

	return false
}