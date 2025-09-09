package config

import (
	"employee-commission-system/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("commission.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Department{},
		&models.Position{},
		&models.Employee{},
		&models.ProjectPermission{},
		&models.CommissionProject{},
		&models.DailyMonthlyReport{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully")
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}