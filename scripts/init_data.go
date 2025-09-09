package main

import (
	"employee-commission-system/config"
	"employee-commission-system/models"
	"log"
)

func main() {
	// 初始化数据库
	config.InitDatabase()

	// 创建示例数据
	createSampleData()

	log.Println("示例数据初始化完成")
}

func createSampleData() {
	db := config.GetDB()

	// 1. 创建用户
	users := []models.User{
		{Username: "admin", Name: "管理员", IsAdmin: 1, EmployeeScope: "0", Status: 1},
		{Username: "user1", Name: "普通用户1", IsAdmin: 0, EmployeeScope: "1,2", Status: 1},
		{Username: "user2", Name: "普通用户2", IsAdmin: 0, EmployeeScope: "3,4,5", Status: 1},
	}

	for _, user := range users {
		var existingUser models.User
		if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
			db.Create(&user)
			log.Printf("创建用户: %s", user.Username)
		}
	}

	// 2. 创建公司
	companies := []models.Company{
		{Name: "美芦", Status: 1},
		{Name: "奥汉", Status: 1},
	}

	for _, company := range companies {
		var existing models.Company
		if err := db.Where("name = ?", company.Name).First(&existing).Error; err != nil {
			db.Create(&company)
			log.Printf("创建公司: %s", company.Name)
		}
	}

	// 3. 创建部门
	departments := []models.Department{
		{CompanyID: 1, Name: "美芦SPACE营地", Status: 1},
		{CompanyID: 2, Name: "奥汉南湖国际水上中心", Status: 1},
		{CompanyID: 2, Name: "营销策划部", Status: 1},
	}

	for _, dept := range departments {
		var existing models.Department
		if err := db.Where("company_id = ? AND name = ?", dept.CompanyID, dept.Name).First(&existing).Error; err != nil {
			db.Create(&dept)
			log.Printf("创建部门: %s", dept.Name)
		}
	}

	// 4. 创建岗位
	positions := []models.Position{
		{Name: "营地管家", Status: 1},
		{Name: "桨板教练", Status: 1},
	}

	for _, pos := range positions {
		var existing models.Position
		if err := db.Where("name = ?", pos.Name).First(&existing).Error; err != nil {
			db.Create(&pos)
			log.Printf("创建岗位: %s", pos.Name)
		}
	}

	// 5. 创建项目权限
	permissions := []models.ProjectPermission{
		{Permission: "全员通用", Status: 1},
		{Permission: "营销策划部", Status: 1},
		{Permission: "客户关系部", Status: 1},
		{Permission: "项目运营部", Status: 1},
	}

	for _, perm := range permissions {
		var existing models.ProjectPermission
		if err := db.Where("permission = ?", perm.Permission).First(&existing).Error; err != nil {
			db.Create(&perm)
			log.Printf("创建项目权限: %s", perm.Permission)
		}
	}

	// 6. 创建员工
	employees := []models.Employee{
		{CompanyID: 1, DepartmentID: 1, PositionID: 1, ProjectPermIDs: "1,2", Name: "夏子杰", Status: 1},
		{CompanyID: 1, DepartmentID: 1, PositionID: 1, ProjectPermIDs: "1,3", Name: "邵燕平", Status: 1},
		{CompanyID: 2, DepartmentID: 2, PositionID: 2, ProjectPermIDs: "1,4", Name: "胡紫阳", Status: 1},
		{CompanyID: 2, DepartmentID: 2, PositionID: 2, ProjectPermIDs: "1,2,3", Name: "廖鑫", Status: 1},
		{CompanyID: 2, DepartmentID: 2, PositionID: 2, ProjectPermIDs: "1,2,3,4", Name: "汪宇博", Status: 1},
	}

	for _, emp := range employees {
		var existing models.Employee
		if err := db.Where("name = ? AND company_id = ?", emp.Name, emp.CompanyID).First(&existing).Error; err != nil {
			db.Create(&emp)
			log.Printf("创建员工: %s", emp.Name)
		}
	}

	// 7. 创建提成项目
	commissionProjects := []models.CommissionProject{
		{FieldName: "注册会员提成", ProjectPermID: 1, Status: 1},
		{FieldName: "储值会员提成", ProjectPermID: 1, Status: 1},
		{FieldName: "储值提成", ProjectPermID: 1, Status: 1},
		{FieldName: "公众号新增数提成", ProjectPermID: 1, Status: 1},
		{FieldName: "好评数提成", ProjectPermID: 1, Status: 1},
		{FieldName: "沙龙执行提成", ProjectPermID: 1, Status: 1},
		{FieldName: "水上课程及次卡类首次成交提成", ProjectPermID: 1, Status: 1},
		{FieldName: "续课提成（首次成交）", ProjectPermID: 1, Status: 1},
		{FieldName: "续卡提成", ProjectPermID: 1, Status: 1},
		{FieldName: "线上产值提成", ProjectPermID: 2, Status: 1},
		{FieldName: "线下产值提成", ProjectPermID: 2, Status: 1},
		{FieldName: "沙龙提成", ProjectPermID: 2, Status: 1},
		{FieldName: "平面作品提成", ProjectPermID: 2, Status: 1},
		{FieldName: "视频作品提成", ProjectPermID: 2, Status: 1},
		{FieldName: "线上平台收入提成", ProjectPermID: 2, Status: 1},
		{FieldName: "客关注册会员工作提成", ProjectPermID: 3, Status: 1},
		{FieldName: "客关储值会员工作提成", ProjectPermID: 3, Status: 1},
		{FieldName: "客关储值工作提成", ProjectPermID: 3, Status: 1},
		{FieldName: "会员消费额提成", ProjectPermID: 3, Status: 1},
		{FieldName: "项目产值提成", ProjectPermID: 4, Status: 1},
		{FieldName: "水上体验执行提成（救生员）", ProjectPermID: 4, Status: 1},
		{FieldName: "水上体验执行提成（安全员）", ProjectPermID: 4, Status: 1},
		{FieldName: "续课提成（执行教练）", ProjectPermID: 4, Status: 1},
		{FieldName: "桨板课程执行提成", ProjectPermID: 4, Status: 1},
		{FieldName: "拓展活动提成", ProjectPermID: 4, Status: 1},
		{FieldName: "手工活动提成", ProjectPermID: 4, Status: 1},
		{FieldName: "策划活动执行提成", ProjectPermID: 4, Status: 1},
		{FieldName: "餐食出品提成", ProjectPermID: 4, Status: 1},
	}

	for _, cp := range commissionProjects {
		var existing models.CommissionProject
		if err := db.Where("field_name = ? AND project_perm_id = ?", cp.FieldName, cp.ProjectPermID).First(&existing).Error; err != nil {
			db.Create(&cp)
			log.Printf("创建提成项目: %s", cp.FieldName)
		}
	}

	// 8. 创建示例报表数据
	reports := []models.DailyMonthlyReport{
		{Date: "2024-01-15", EmployeeID: 1, CommissionProjectID: 1, CommissionValue: 100.0, Status: 1},
		{Date: "2024-01-15", EmployeeID: 1, CommissionProjectID: 2, CommissionValue: 200.0, Status: 1},
		{Date: "2024-01-16", EmployeeID: 2, CommissionProjectID: 3, CommissionValue: 150.0, Status: 1},
		{Date: "2024-01-16", EmployeeID: 3, CommissionProjectID: 20, CommissionValue: 300.0, Status: 1},
		{Date: "2024-01-17", EmployeeID: 4, CommissionProjectID: 21, CommissionValue: 250.0, Status: 1},
		{Date: "2024-01-17", EmployeeID: 5, CommissionProjectID: 22, CommissionValue: 180.0, Status: 1},
	}

	for _, report := range reports {
		// 获取员工信息来填充相关字段
		var employee models.Employee
		if err := db.First(&employee, report.EmployeeID).Error; err == nil {
			report.CompanyID = employee.CompanyID
			report.DepartmentID = employee.DepartmentID
			report.PositionID = employee.PositionID
			report.EmployeeName = employee.Name

			var existing models.DailyMonthlyReport
			if err := db.Where("date = ? AND employee_id = ? AND commission_project_id = ?", 
				report.Date, report.EmployeeID, report.CommissionProjectID).First(&existing).Error; err != nil {
				db.Create(&report)
				log.Printf("创建报表记录: %s - %s", report.Date, report.EmployeeName)
			}
		}
	}
}