package config

import (
	"fmt"
	"log"
	"os"

	"student-management-system/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("警告: .env 文件未找到，将使用默认配置")
	}
	// 从环境变量读取配置
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "your_password")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "student_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	})

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移 - 按依赖关系排序
	err = DB.AutoMigrate(
		// 1. 基础表（无外键依赖）
		&models.Role{},
		&models.User{},
		&models.Teacher{},
		&models.Class{},
		&models.Course{},
		&models.Notification{},

		// 2. 依赖基础表的表
		&models.Student{}, // 依赖 Class
		&models.Parent{},  // 依赖 Student

		// 3. 依赖多个表的关联表
		&models.Enrollment{},       // 依赖 Student, Course
		&models.Grade{},            // 依赖 Enrollment
		&models.Attendance{},       // 依赖 Student
		&models.RewardPunishment{}, // 依赖 Student
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库表迁移成功")

	// 初始化默认数据
	initDefaultData()
}

// initDefaultData 初始化默认数据
func initDefaultData() {
	// 创建默认角色
	roles := []models.Role{
		{RoleName: "admin"},
		{RoleName: "teacher"},
		{RoleName: "student"},
		{RoleName: "parent"},
	}

	for _, role := range roles {
		var count int64
		DB.Model(&models.Role{}).Where("role_name = ?", role.RoleName).Count(&count)
		if count == 0 {
			DB.Create(&role)
		}
	}

	log.Println("默认角色初始化完成")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
