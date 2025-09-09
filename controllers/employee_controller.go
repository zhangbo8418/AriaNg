package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEmployee 创建员工
func CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证关联数据是否存在
	var company models.Company
	if err := config.DB.First(&company, employee.CompanyID).Error; err != nil {
		utils.BadRequest(c, "所属公司不存在")
		return
	}

	var department models.Department
	if err := config.DB.First(&department, employee.DepartmentID).Error; err != nil {
		utils.BadRequest(c, "所属部门不存在")
		return
	}

	var position models.Position
	if err := config.DB.First(&position, employee.PositionID).Error; err != nil {
		utils.BadRequest(c, "所属岗位不存在")
		return
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		utils.InternalError(c, "创建员工失败: "+err.Error())
		return
	}

	// 加载关联数据
	config.DB.Preload("Company").Preload("Department").Preload("Position").First(&employee, employee.ID)

	utils.Success(c, employee)
}

// GetEmployees 获取员工列表
func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")
	companyID := c.Query("company_id")
	departmentID := c.Query("department_id")
	positionID := c.Query("position_id")

	query := config.DB.Model(&models.Employee{}).
		Preload("Company").
		Preload("Department").
		Preload("Position")
	
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}

	if positionID != "" {
		query = query.Where("position_id = ?", positionID)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&employees).Error; err != nil {
		utils.InternalError(c, "获取员工列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"employees": employees,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetEmployee 获取单个员工
func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.Preload("Company").Preload("Department").Preload("Position").First(&employee, id).Error; err != nil {
		utils.NotFound(c, "员工不存在")
		return
	}

	utils.Success(c, employee)
}

// UpdateEmployee 更新员工
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.First(&employee, id).Error; err != nil {
		utils.NotFound(c, "员工不存在")
		return
	}

	var updateData models.Employee
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证关联数据是否存在
	if updateData.CompanyID != 0 && updateData.CompanyID != employee.CompanyID {
		var company models.Company
		if err := config.DB.First(&company, updateData.CompanyID).Error; err != nil {
			utils.BadRequest(c, "所属公司不存在")
			return
		}
	}

	if updateData.DepartmentID != 0 && updateData.DepartmentID != employee.DepartmentID {
		var department models.Department
		if err := config.DB.First(&department, updateData.DepartmentID).Error; err != nil {
			utils.BadRequest(c, "所属部门不存在")
			return
		}
	}

	if updateData.PositionID != 0 && updateData.PositionID != employee.PositionID {
		var position models.Position
		if err := config.DB.First(&position, updateData.PositionID).Error; err != nil {
			utils.BadRequest(c, "所属岗位不存在")
			return
		}
	}

	if err := config.DB.Model(&employee).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新员工失败: "+err.Error())
		return
	}

	// 重新加载关联数据
	config.DB.Preload("Company").Preload("Department").Preload("Position").First(&employee, employee.ID)

	utils.Success(c, employee)
}

// DeleteEmployee 删除员工
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.First(&employee, id).Error; err != nil {
		utils.NotFound(c, "员工不存在")
		return
	}

	// 检查是否有关联的报表记录
	var reportCount int64
	config.DB.Model(&models.DailyMonthlyReport{}).Where("employee_id = ?", id).Count(&reportCount)
	if reportCount > 0 {
		utils.BadRequest(c, "该员工还有相关报表记录，无法删除")
		return
	}

	if err := config.DB.Delete(&employee).Error; err != nil {
		utils.InternalError(c, "删除员工失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "员工删除成功"})
}

// ToggleEmployeeStatus 切换员工状态
func ToggleEmployeeStatus(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.First(&employee, id).Error; err != nil {
		utils.NotFound(c, "员工不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if employee.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&employee).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新员工状态失败: "+err.Error())
		return
	}

	employee.Status = newStatus
	// 重新加载关联数据
	config.DB.Preload("Company").Preload("Department").Preload("Position").First(&employee, employee.ID)
	
	utils.Success(c, employee)
}