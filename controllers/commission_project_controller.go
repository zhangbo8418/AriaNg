package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCommissionProject 创建提成项目
func CreateCommissionProject(c *gin.Context) {
	var project models.CommissionProject
	if err := c.ShouldBindJSON(&project); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证项目权限是否存在
	var permission models.ProjectPermission
	if err := config.DB.First(&permission, project.ProjectPermID).Error; err != nil {
		utils.BadRequest(c, "项目权限不存在")
		return
	}

	if err := config.DB.Create(&project).Error; err != nil {
		utils.InternalError(c, "创建提成项目失败: "+err.Error())
		return
	}

	// 加载关联数据
	config.DB.Preload("ProjectPermission").First(&project, project.ID)

	utils.Success(c, project)
}

// GetCommissionProjects 获取提成项目列表
func GetCommissionProjects(c *gin.Context) {
	var projects []models.CommissionProject
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")
	projectPermID := c.Query("project_perm_id")

	query := config.DB.Model(&models.CommissionProject{}).Preload("ProjectPermission")
	
	if search != "" {
		query = query.Where("field_name LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if projectPermID != "" {
		query = query.Where("project_perm_id = ?", projectPermID)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
		utils.InternalError(c, "获取提成项目列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"projects":  projects,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetCommissionProject 获取单个提成项目
func GetCommissionProject(c *gin.Context) {
	id := c.Param("id")
	var project models.CommissionProject

	if err := config.DB.Preload("ProjectPermission").First(&project, id).Error; err != nil {
		utils.NotFound(c, "提成项目不存在")
		return
	}

	utils.Success(c, project)
}

// UpdateCommissionProject 更新提成项目
func UpdateCommissionProject(c *gin.Context) {
	id := c.Param("id")
	var project models.CommissionProject

	if err := config.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "提成项目不存在")
		return
	}

	var updateData models.CommissionProject
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 如果更新了项目权限ID，验证是否存在
	if updateData.ProjectPermID != 0 && updateData.ProjectPermID != project.ProjectPermID {
		var permission models.ProjectPermission
		if err := config.DB.First(&permission, updateData.ProjectPermID).Error; err != nil {
			utils.BadRequest(c, "项目权限不存在")
			return
		}
	}

	if err := config.DB.Model(&project).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新提成项目失败: "+err.Error())
		return
	}

	// 重新加载关联数据
	config.DB.Preload("ProjectPermission").First(&project, project.ID)

	utils.Success(c, project)
}

// DeleteCommissionProject 删除提成项目
func DeleteCommissionProject(c *gin.Context) {
	id := c.Param("id")
	var project models.CommissionProject

	if err := config.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "提成项目不存在")
		return
	}

	// 检查是否有关联的报表记录
	var reportCount int64
	config.DB.Model(&models.DailyMonthlyReport{}).Where("commission_project_id = ?", id).Count(&reportCount)
	if reportCount > 0 {
		utils.BadRequest(c, "该提成项目还有相关报表记录，无法删除")
		return
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		utils.InternalError(c, "删除提成项目失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "提成项目删除成功"})
}

// ToggleCommissionProjectStatus 切换提成项目状态
func ToggleCommissionProjectStatus(c *gin.Context) {
	id := c.Param("id")
	var project models.CommissionProject

	if err := config.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "提成项目不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if project.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&project).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新提成项目状态失败: "+err.Error())
		return
	}

	project.Status = newStatus
	// 重新加载关联数据
	config.DB.Preload("ProjectPermission").First(&project, project.ID)
	
	utils.Success(c, project)
}