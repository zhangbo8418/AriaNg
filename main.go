package main

import (
	"employee-commission-system/config"
	"employee-commission-system/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDatabase()

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}