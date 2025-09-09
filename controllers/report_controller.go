package controllers

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"employee-commission-system/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateReport 创建报表记录
func CreateReport(c *gin.Context) {
	var report models.DailyMonthlyReport
	if err := c.ShouldBindJSON(&report); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证关联数据是否存在
	var employee models.Employee
	if err := config.DB.First(&employee, report.EmployeeID).Error; err != nil {
		utils.BadRequest(c, "员工不存在")
		return
	}

	var commissionProject models.CommissionProject
	if err := config.DB.First(&commissionProject, report.CommissionProjectID).Error; err != nil {
		utils.BadRequest(c, "提成项目不存在")
		return
	}

	// 自动填充员工相关信息
	report.CompanyID = employee.CompanyID
	report.DepartmentID = employee.DepartmentID
	report.PositionID = employee.PositionID
	report.EmployeeName = employee.Name

	if err := config.DB.Create(&report).Error; err != nil {
		utils.InternalError(c, "创建报表记录失败: "+err.Error())
		return
	}

	// 加载关联数据
	config.DB.Preload("Employee").
		Preload("Company").
		Preload("Department").
		Preload("Position").
		Preload("CommissionProject").
		First(&report, report.ID)

	utils.Success(c, report)
}

// GetReports 获取报表记录列表
func GetReports(c *gin.Context) {
	var reports []models.DailyMonthlyReport
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 搜索参数
	search := c.Query("search")
	status := c.Query("status")
	employeeID := c.Query("employee_id")
	companyID := c.Query("company_id")
	departmentID := c.Query("department_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := config.DB.Model(&models.DailyMonthlyReport{}).
		Preload("Employee").
		Preload("Company").
		Preload("Department").
		Preload("Position").
		Preload("CommissionProject")
	
	if search != "" {
		query = query.Where("employee_name LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}

	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Order("date DESC").Find(&reports).Error; err != nil {
		utils.InternalError(c, "获取报表记录列表失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"reports":   reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetReport 获取单个报表记录
func GetReport(c *gin.Context) {
	id := c.Param("id")
	var report models.DailyMonthlyReport

	if err := config.DB.Preload("Employee").
		Preload("Company").
		Preload("Department").
		Preload("Position").
		Preload("CommissionProject").
		First(&report, id).Error; err != nil {
		utils.NotFound(c, "报表记录不存在")
		return
	}

	utils.Success(c, report)
}

// UpdateReport 更新报表记录
func UpdateReport(c *gin.Context) {
	id := c.Param("id")
	var report models.DailyMonthlyReport

	if err := config.DB.First(&report, id).Error; err != nil {
		utils.NotFound(c, "报表记录不存在")
		return
	}

	var updateData models.DailyMonthlyReport
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 如果更新了员工ID，需要同步更新相关信息
	if updateData.EmployeeID != 0 && updateData.EmployeeID != report.EmployeeID {
		var employee models.Employee
		if err := config.DB.First(&employee, updateData.EmployeeID).Error; err != nil {
			utils.BadRequest(c, "员工不存在")
			return
		}
		updateData.CompanyID = employee.CompanyID
		updateData.DepartmentID = employee.DepartmentID
		updateData.PositionID = employee.PositionID
		updateData.EmployeeName = employee.Name
	}

	// 验证提成项目是否存在
	if updateData.CommissionProjectID != 0 && updateData.CommissionProjectID != report.CommissionProjectID {
		var commissionProject models.CommissionProject
		if err := config.DB.First(&commissionProject, updateData.CommissionProjectID).Error; err != nil {
			utils.BadRequest(c, "提成项目不存在")
			return
		}
	}

	if err := config.DB.Model(&report).Updates(updateData).Error; err != nil {
		utils.InternalError(c, "更新报表记录失败: "+err.Error())
		return
	}

	// 重新加载关联数据
	config.DB.Preload("Employee").
		Preload("Company").
		Preload("Department").
		Preload("Position").
		Preload("CommissionProject").
		First(&report, report.ID)

	utils.Success(c, report)
}

// DeleteReport 删除报表记录
func DeleteReport(c *gin.Context) {
	id := c.Param("id")
	var report models.DailyMonthlyReport

	if err := config.DB.First(&report, id).Error; err != nil {
		utils.NotFound(c, "报表记录不存在")
		return
	}

	if err := config.DB.Delete(&report).Error; err != nil {
		utils.InternalError(c, "删除报表记录失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "报表记录删除成功"})
}

// ToggleReportStatus 切换报表记录状态
func ToggleReportStatus(c *gin.Context) {
	id := c.Param("id")
	var report models.DailyMonthlyReport

	if err := config.DB.First(&report, id).Error; err != nil {
		utils.NotFound(c, "报表记录不存在")
		return
	}

	// 切换状态
	newStatus := 1
	if report.Status == 1 {
		newStatus = 0
	}

	if err := config.DB.Model(&report).Update("status", newStatus).Error; err != nil {
		utils.InternalError(c, "更新报表记录状态失败: "+err.Error())
		return
	}

	report.Status = newStatus
	// 重新加载关联数据
	config.DB.Preload("Employee").
		Preload("Company").
		Preload("Department").
		Preload("Position").
		Preload("CommissionProject").
		First(&report, report.ID)
	
	utils.Success(c, report)
}

// GetMonthlySummary 获取月度汇总报表
func GetMonthlySummary(c *gin.Context) {
	month := c.Query("month") // 格式：2024-01
	if month == "" {
		utils.BadRequest(c, "请提供月份参数")
		return
	}

	// 构建日期范围
	startDate := month + "-01"
	// 计算月末日期
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		utils.BadRequest(c, "日期格式错误")
		return
	}
	endDate := t.AddDate(0, 1, -1).Format("2006-01-02")

	var summaries []models.ReportSummary
	
	// 查询月度汇总数据
	query := `
		SELECT 
			r.employee_id,
			r.employee_name,
			r.company_id,
			c.name as company_name,
			r.department_id,
			d.name as department_name,
			r.position_id,
			p.name as position_name,
			r.commission_project_id,
			cp.field_name as project_name,
			SUM(r.commission_value) as total_value
		FROM daily_monthly_reports r
		LEFT JOIN companies c ON r.company_id = c.id
		LEFT JOIN departments d ON r.department_id = d.id
		LEFT JOIN positions p ON r.position_id = p.id
		LEFT JOIN commission_projects cp ON r.commission_project_id = cp.id
		WHERE r.date >= ? AND r.date <= ? AND r.status = 1
		GROUP BY r.employee_id, r.commission_project_id
		ORDER BY r.employee_name, cp.field_name
	`

	if err := config.DB.Raw(query, startDate, endDate).Scan(&summaries).Error; err != nil {
		utils.InternalError(c, "获取月度汇总失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"summaries": summaries,
		"month":     month,
		"start_date": startDate,
		"end_date":   endDate,
	})
}

// BatchCreateReports 批量创建报表记录
func BatchCreateReports(c *gin.Context) {
	var reports []models.DailyMonthlyReport
	if err := c.ShouldBindJSON(&reports); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(reports) == 0 {
		utils.BadRequest(c, "报表记录不能为空")
		return
	}

	// 验证和填充数据
	for i := range reports {
		// 验证员工是否存在
		var employee models.Employee
		if err := config.DB.First(&employee, reports[i].EmployeeID).Error; err != nil {
			utils.BadRequest(c, "员工ID "+strconv.Itoa(reports[i].EmployeeID)+" 不存在")
			return
		}

		// 验证提成项目是否存在
		var commissionProject models.CommissionProject
		if err := config.DB.First(&commissionProject, reports[i].CommissionProjectID).Error; err != nil {
			utils.BadRequest(c, "提成项目ID "+strconv.Itoa(reports[i].CommissionProjectID)+" 不存在")
			return
		}

		// 自动填充员工相关信息
		reports[i].CompanyID = employee.CompanyID
		reports[i].DepartmentID = employee.DepartmentID
		reports[i].PositionID = employee.PositionID
		reports[i].EmployeeName = employee.Name
	}

	// 批量创建
	if err := config.DB.Create(&reports).Error; err != nil {
		utils.InternalError(c, "批量创建报表记录失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"message": "批量创建成功",
		"count":   len(reports),
	})
}