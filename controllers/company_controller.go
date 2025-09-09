package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCompany 创建公司
func CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := config.DB.Create(&company).Error; err != nil {
		utils.InternalError(c, "创建公司失败: "+err.Error())
		return
	}

	utils.Success(c, company)
}

// GetCompanies 获取公司列表
func GetCompanies(c *gin.Context) {
	var companies []models.Company
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")

	query := config.DB.Model(&models.Company{})
	
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&companies).Error; err != nil {
		utils.InternalError(c, "获取公司列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"companies": companies,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetCompany 获取单个公司
func GetCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	if err := config.DB.First(&company, id).Error; err != nil {
		utils.NotFound(c, "公司不存在")
		return
	}

	utils.Success(c, company)
}

// UpdateCompany 更新公司
func UpdateCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	if err := config.DB.First(&company, id).Error; err != nil {
		utils.NotFound(c, "公司不存在")
		return
	}

	var updateData models.Company
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := config.DB.Model(&company).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新公司失败: "+err.Error())
		return
	}

	utils.Success(c, company)
}

// DeleteCompany 删除公司
func DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	if err := config.DB.First(&company, id).Error; err != nil {
		utils.NotFound(c, "公司不存在")
		return
	}

	// 检查是否有关联的部门或员工
	var departmentCount int64
	config.DB.Model(&models.Department{}).Where("company_id = ?", id).Count(&departmentCount)
	if departmentCount > 0 {
		utils.BadRequest(c, "该公司下还有部门，无法删除")
		return
	}

	var employeeCount int64
	config.DB.Model(&models.Employee{}).Where("company_id = ?", id).Count(&employeeCount)
	if employeeCount > 0 {
		utils.BadRequest(c, "该公司下还有员工，无法删除")
		return
	}

	if err := config.DB.Delete(&company).Error; err != nil {
		utils.InternalError(c, "删除公司失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "公司删除成功"})
}

// ToggleCompanyStatus 切换公司状态
func ToggleCompanyStatus(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	if err := config.DB.First(&company, id).Error; err != nil {
		utils.NotFound(c, "公司不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if company.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&company).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新公司状态失败: "+err.Error())
		return
	}

	company.Status = newStatus
	utils.Success(c, company)
}