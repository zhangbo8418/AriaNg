package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProjectPermission 创建项目权限
func CreateProjectPermission(c *gin.Context) {
	var permission models.ProjectPermission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := config.DB.Create(&permission).Error; err != nil {
		utils.InternalError(c, "创建项目权限失败: "+err.Error())
		return
	}

	utils.Success(c, permission)
}

// GetProjectPermissions 获取项目权限列表
func GetProjectPermissions(c *gin.Context) {
	var permissions []models.ProjectPermission
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")

	query := config.DB.Model(&models.ProjectPermission{})
	
	if search != "" {
		query = query.Where("permission LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		utils.InternalError(c, "获取项目权限列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"permissions": permissions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

// GetProjectPermission 获取单个项目权限
func GetProjectPermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.ProjectPermission

	if err := config.DB.First(&permission, id).Error; err != nil {
		utils.NotFound(c, "项目权限不存在")
		return
	}

	utils.Success(c, permission)
}

// UpdateProjectPermission 更新项目权限
func UpdateProjectPermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.ProjectPermission

	if err := config.DB.First(&permission, id).Error; err != nil {
		utils.NotFound(c, "项目权限不存在")
		return
	}

	var updateData models.ProjectPermission
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := config.DB.Model(&permission).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新项目权限失败: "+err.Error())
		return
	}

	utils.Success(c, permission)
}

// DeleteProjectPermission 删除项目权限
func DeleteProjectPermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.ProjectPermission

	if err := config.DB.First(&permission, id).Error; err != nil {
		utils.NotFound(c, "项目权限不存在")
		return
	}

	// 检查是否有关联的提成项目
	var commissionCount int64
	config.DB.Model(&models.CommissionProject{}).Where("project_perm_id = ?", id).Count(&commissionCount)
	if commissionCount > 0 {
		utils.BadRequest(c, "该项目权限下还有提成项目，无法删除")
		return
	}

	if err := config.DB.Delete(&permission).Error; err != nil {
		utils.InternalError(c, "删除项目权限失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "项目权限删除成功"})
}

// ToggleProjectPermissionStatus 切换项目权限状态
func ToggleProjectPermissionStatus(c *gin.Context) {
	id := c.Param("id")
	var permission models.ProjectPermission

	if err := config.DB.First(&permission, id).Error; err != nil {
		utils.NotFound(c, "项目权限不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if permission.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&permission).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新项目权限状态失败: "+err.Error())
		return
	}

	permission.Status = newStatus
	utils.Success(c, permission)
}