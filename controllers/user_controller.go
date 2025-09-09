package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, "创建用户失败: "+err.Error())
		return
	}

	utils.Success(c, user)
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	var users []models.User
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")

	query := config.DB.Model(&models.User{})
	
	if search != "" {
		query = query.Where("username LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		utils.InternalError(c, "获取用户列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 如果要更新用户名，检查是否重复
	if updateData.Username != "" && updateData.Username != user.Username {
		var existingUser models.User
		if err := config.DB.Where("username = ? AND id != ?", updateData.Username, id).First(&existingUser).Error; err == nil {
			utils.BadRequest(c, "用户名已存在")
			return
		}
	}

	if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新用户失败: "+err.Error())
		return
	}

	utils.Success(c, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.InternalError(c, "删除用户失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "用户删除成功"})
}

// ToggleUserStatus 切换用户状态
func ToggleUserStatus(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if user.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&user).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新用户状态失败: "+err.Error())
		return
	}

	user.Status = newStatus
	utils.Success(c, user)
}