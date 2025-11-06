package main

import (
	"fmt"
	"log"

	"student-management-system/config"
	"student-management-system/internal/models"
	"student-management-system/internal/utils"
)

func main() {
	fmt.Println("=== 学生管理系统 - 创建管理员账号 ===")

	// 初始化数据库
	config.InitDB()
	db := config.GetDB()

	// 创建管理员账号
	username := "admin"
	password := "admin123"

	// 检查用户是否已存在
	var existingUser models.User
	result := db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		fmt.Printf("用户 '%s' 已存在，无需重复创建\n", username)
		return
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 获取管理员角色ID
	var adminRole models.Role
	if err := db.Where("role_name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Fatalf("获取管理员角色失败: %v\n提示: 请先启动主程序以初始化角色数据", err)
	}

	// 创建管理员用户
	user := models.User{
		Username: username,
		Password: hashedPassword,
		RoleID:   adminRole.ID,
		IsActive: true,
		UserType: "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("创建管理员账号失败: %v", err)
	}

	fmt.Println("\n✓ 管理员账号创建成功!")
	fmt.Printf("  用户名: %s\n", username)
	fmt.Printf("  密码: %s\n", password)
	fmt.Println("\n请妥善保管账号信息，并在首次登录后修改密码。")
}
