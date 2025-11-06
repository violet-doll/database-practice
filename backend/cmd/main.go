package main

import (
	"fmt"
	"log"

	"student-management-system/config"
	"student-management-system/internal/api"
	"student-management-system/internal/utils"
)

func main() {
	// 设置JWT密钥
	utils.SetJWTSecret("your_jwt_secret_key_change_in_production")

	// 初始化数据库
	config.InitDB()

	// 设置路由
	router := api.SetupRouter()

	// 启动服务器
	port := ":8080"
	fmt.Printf("服务器启动成功，监听端口 %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
