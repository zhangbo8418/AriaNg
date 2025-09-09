package routes

import (
	"employee-commission-system/controllers"
	"employee-commission-system/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// API版本前缀
	api := r.Group("/api/v1")

	// 公开路由（不需要认证）
	public := api.Group("/")
	{
		public.POST("/login", controllers.Login)
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关路由
		protected.GET("/profile", controllers.GetProfile)
		protected.PUT("/profile", controllers.UpdateProfile)

		// 管理员路由
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			// 用户管理
			admin.POST("/users", controllers.CreateUser)
			admin.GET("/users", controllers.GetUsers)
			admin.GET("/users/:id", controllers.GetUser)
			admin.PUT("/users/:id", controllers.UpdateUser)
			admin.DELETE("/users/:id", controllers.DeleteUser)
			admin.PUT("/users/:id/toggle-status", controllers.ToggleUserStatus)

			// 公司管理
			admin.POST("/companies", controllers.CreateCompany)
			admin.GET("/companies", controllers.GetCompanies)
			admin.GET("/companies/:id", controllers.GetCompany)
			admin.PUT("/companies/:id", controllers.UpdateCompany)
			admin.DELETE("/companies/:id", controllers.DeleteCompany)
			admin.PUT("/companies/:id/toggle-status", controllers.ToggleCompanyStatus)

			// 部门管理
			admin.POST("/departments", controllers.CreateDepartment)
			admin.GET("/departments", controllers.GetDepartments)
			admin.GET("/departments/:id", controllers.GetDepartment)
			admin.PUT("/departments/:id", controllers.UpdateDepartment)
			admin.DELETE("/departments/:id", controllers.DeleteDepartment)
			admin.PUT("/departments/:id/toggle-status", controllers.ToggleDepartmentStatus)

			// 岗位管理
			admin.POST("/positions", controllers.CreatePosition)
			admin.GET("/positions", controllers.GetPositions)
			admin.GET("/positions/:id", controllers.GetPosition)
			admin.PUT("/positions/:id", controllers.UpdatePosition)
			admin.DELETE("/positions/:id", controllers.DeletePosition)
			admin.PUT("/positions/:id/toggle-status", controllers.TogglePositionStatus)

			// 员工管理
			admin.POST("/employees", controllers.CreateEmployee)
			admin.GET("/employees", controllers.GetEmployees)
			admin.GET("/employees/:id", controllers.GetEmployee)
			admin.PUT("/employees/:id", controllers.UpdateEmployee)
			admin.DELETE("/employees/:id", controllers.DeleteEmployee)
			admin.PUT("/employees/:id/toggle-status", controllers.ToggleEmployeeStatus)

			// 项目权限管理
			admin.POST("/project-permissions", controllers.CreateProjectPermission)
			admin.GET("/project-permissions", controllers.GetProjectPermissions)
			admin.GET("/project-permissions/:id", controllers.GetProjectPermission)
			admin.PUT("/project-permissions/:id", controllers.UpdateProjectPermission)
			admin.DELETE("/project-permissions/:id", controllers.DeleteProjectPermission)
			admin.PUT("/project-permissions/:id/toggle-status", controllers.ToggleProjectPermissionStatus)

			// 提成项目管理
			admin.POST("/commission-projects", controllers.CreateCommissionProject)
			admin.GET("/commission-projects", controllers.GetCommissionProjects)
			admin.GET("/commission-projects/:id", controllers.GetCommissionProject)
			admin.PUT("/commission-projects/:id", controllers.UpdateCommissionProject)
			admin.DELETE("/commission-projects/:id", controllers.DeleteCommissionProject)
			admin.PUT("/commission-projects/:id/toggle-status", controllers.ToggleCommissionProjectStatus)

			// 报表管理（管理员可以查看所有报表）
			admin.GET("/reports", controllers.GetReports)
			admin.GET("/reports/:id", controllers.GetReport)
			admin.PUT("/reports/:id", controllers.UpdateReport)
			admin.DELETE("/reports/:id", controllers.DeleteReport)
			admin.PUT("/reports/:id/toggle-status", controllers.ToggleReportStatus)
			admin.GET("/reports/monthly-summary", controllers.GetMonthlySummary)
		}

		// 普通用户路由（受员工范围限制）
		user := protected.Group("/")
		user.Use(middleware.EmployeeScopeMiddleware())
		{
			// 报表提交
			user.POST("/reports", controllers.CreateReport)
			user.POST("/reports/batch", controllers.BatchCreateReports)
			
			// 查看自己管理范围内的报表
			user.GET("/my-reports", controllers.GetReports)
			user.GET("/my-reports/:id", controllers.GetReport)
			user.GET("/my-reports/monthly-summary", controllers.GetMonthlySummary)

			// 查看基础数据（只读）
			user.GET("/companies", controllers.GetCompanies)
			user.GET("/departments", controllers.GetDepartments)
			user.GET("/positions", controllers.GetPositions)
			user.GET("/employees", controllers.GetEmployees)
			user.GET("/project-permissions", controllers.GetProjectPermissions)
			user.GET("/commission-projects", controllers.GetCommissionProjects)
		}
	}
}