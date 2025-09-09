package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePosition 创建岗位
func CreatePosition(c *gin.Context) {
	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := config.DB.Create(&position).Error; err != nil {
		utils.InternalError(c, "创建岗位失败: "+err.Error())
		return
	}

	utils.Success(c, position)
}

// GetPositions 获取岗位列表
func GetPositions(c *gin.Context) {
	var positions []models.Position
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")

	query := config.DB.Model(&models.Position{})
	
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&positions).Error; err != nil {
		utils.InternalError(c, "获取岗位列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"positions": positions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPosition 获取单个岗位
func GetPosition(c *gin.Context) {
	id := c.Param("id")
	var position models.Position

	if err := config.DB.First(&position, id).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	utils.Success(c, position)
}

// UpdatePosition 更新岗位
func UpdatePosition(c *gin.Context) {
	id := c.Param("id")
	var position models.Position

	if err := config.DB.First(&position, id).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	var updateData models.Position
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := config.DB.Model(&position).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新岗位失败: "+err.Error())
		return
	}

	utils.Success(c, position)
}

// DeletePosition 删除岗位
func DeletePosition(c *gin.Context) {
	id := c.Param("id")
	var position models.Position

	if err := config.DB.First(&position, id).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	// 检查是否有关联的员工
	var employeeCount int64
	config.DB.Model(&models.Employee{}).Where("position_id = ?", id).Count(&employeeCount)
	if employeeCount > 0 {
		utils.BadRequest(c, "该岗位下还有员工，无法删除")
		return
	}

	if err := config.DB.Delete(&position).Error; err != nil {
		utils.InternalError(c, "删除岗位失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "岗位删除成功"})
}

// TogglePositionStatus 切换岗位状态
func TogglePositionStatus(c *gin.Context) {
	id := c.Param("id")
	var position models.Position

	if err := config.DB.First(&position, id).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if position.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&position).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新岗位状态失败: "+err.Error())
		return
	}

	position.Status = newStatus
	utils.Success(c, position)
}