package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDepartment 创建部门
func CreateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证公司是否存在
	var company models.Company
	if err := config.DB.First(&company, department.CompanyID).Error; err != nil {
		utils.BadRequest(c, "所属公司不存在")
		return
	}

	if err := config.DB.Create(&department).Error; err != nil {
		utils.InternalError(c, "创建部门失败: "+err.Error())
		return
	}

	// 加载关联数据
	config.DB.Preload("Company").First(&department, department.ID)

	utils.Success(c, department)
}

// GetDepartments 获取部门列表
func GetDepartments(c *gin.Context) {
	var departments []models.Department
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")
	companyID := c.Query("company_id")

	query := config.DB.Model(&models.Department{}).Preload("Company")
	
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&departments).Error; err != nil {
		utils.InternalError(c, "获取部门列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"departments": departments,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

// GetDepartment 获取单个部门
func GetDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.Preload("Company").First(&department, id).Error; err != nil {
		utils.NotFound(c, "部门不存在")
		return
	}

	utils.Success(c, department)
}

// UpdateDepartment 更新部门
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.First(&department, id).Error; err != nil {
		utils.NotFound(c, "部门不存在")
		return
	}

	var updateData models.Department
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 如果更新了公司ID，验证公司是否存在
	if updateData.CompanyID != 0 && updateData.CompanyID != department.CompanyID {
		var company models.Company
		if err := config.DB.First(&company, updateData.CompanyID).Error; err != nil {
			utils.BadRequest(c, "所属公司不存在")
			return
		}
	}

	if err := config.DB.Model(&department).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新部门失败: "+err.Error())
		return
	}

	// 重新加载关联数据
	config.DB.Preload("Company").First(&department, department.ID)

	utils.Success(c, department)
}

// DeleteDepartment 删除部门
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.First(&department, id).Error; err != nil {
		utils.NotFound(c, "部门不存在")
		return
	}

	// 检查是否有关联的员工
	var employeeCount int64
	config.DB.Model(&models.Employee{}).Where("department_id = ?", id).Count(&employeeCount)
	if employeeCount > 0 {
		utils.BadRequest(c, "该部门下还有员工，无法删除")
		return
	}

	if err := config.DB.Delete(&department).Error; err != nil {
		utils.InternalError(c, "删除部门失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "部门删除成功"})
}

// ToggleDepartmentStatus 切换部门状态
func ToggleDepartmentStatus(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.First(&department, id).Error; err != nil {
		utils.NotFound(c, "部门不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if department.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&department).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新部门状态失败: "+err.Error())
		return
	}

	department.Status = newStatus
	// 重新加载关联数据
	config.DB.Preload("Company").First(&department, department.ID)
	
	utils.Success(c, department)
}