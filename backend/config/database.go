package config

import (
	"fmt"
	"log"
	"os"

	"student-management-system/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	})

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移 - 按依赖关系排序
	err = DB.AutoMigrate(
		// 1. 基础表（无外键依赖）
		&models.Permission{},     // 权限表
		&models.Role{},           // 角色表
		&models.RolePermission{}, // 角色-权限关联表
		&models.User{},
		&models.Teacher{},
		&models.Class{},
		&models.Course{},
		&models.CoursePrerequisite{}, // 课程先修关系表
		&models.Notification{},

		// 2. 依赖基础表的表
		&models.Student{}, // 依赖 Class
		&models.Parent{},  // 依赖 Student

		// 3. 依赖多个表的关联表
		&models.Enrollment{},       // 依赖 Student, Course
		&models.Grade{},            // 依赖 Enrollment
		&models.GradeAuditLog{},    // 依赖 Grade - 成绩审计日志（配合触发器使用）
		&models.Attendance{},       // 依赖 Student
		&models.RewardPunishment{}, // 依赖 Student
		&models.Schedule{},         // 依赖 Course, Class, Teacher
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
	// 创建默认角色（幂等）
	roles := []models.Role{
		{RoleName: "admin"},
		{RoleName: "teacher"},
		{RoleName: "student"},
		{RoleName: "parent"},
	}

	// 使用 ON CONFLICT DO NOTHING（MySQL 下会使用 INSERT IGNORE）避免唯一键冲突报错
	DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles)

	log.Println("默认角色初始化完成")

	// 初始化权限并为各角色分配默认权限
	initPermissionsAndAssignRoles()
}

// 初始化权限并为各角色分配默认权限
func initPermissionsAndAssignRoles() {
	// 定义所有权限（与 admin_permissions.go 中的定义保持一致）
	allPermissions := []models.Permission{
		// 管理员权限
		{Name: "查看用户", Permission: "admin:user:read", Group: "admin"},
		{Name: "创建用户", Permission: "admin:user:create", Group: "admin"},
		{Name: "修改用户", Permission: "admin:user:update", Group: "admin"},
		{Name: "删除用户", Permission: "admin:user:delete", Group: "admin"},
		{Name: "查看角色", Permission: "admin:role:read", Group: "admin"},
		{Name: "创建角色", Permission: "admin:role:create", Group: "admin"},
		{Name: "修改角色", Permission: "admin:role:update", Group: "admin"},
		{Name: "删除角色", Permission: "admin:role:delete", Group: "admin"},
		{Name: "查看统计", Permission: "admin:stats:read", Group: "admin"},
	}

	// 初始化权限（如果不存在则创建）
	for _, perm := range allPermissions {
		var existing models.Permission
		if err := DB.Where("permission = ?", perm.Permission).First(&existing).Error; err != nil {
			// 权限不存在，创建它
			if err := DB.Create(&perm).Error; err != nil {
				log.Printf("创建权限失败: %v, error: %v", perm.Permission, err)
			}
		}
	}

	// 获取 admin 角色
	var adminRole models.Role
	if err := DB.Where("role_name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Printf("获取 admin 角色失败: %v", err)
		return
	}

	// 为 admin 角色分配所有权限
	count := DB.Model(&adminRole).Association("Permissions").Count()
	if count == 0 {
		var permissions []models.Permission
		if err := DB.Find(&permissions).Error; err != nil {
			log.Printf("获取权限列表失败: %v", err)
			return
		}
		if err := DB.Model(&adminRole).Association("Permissions").Replace(permissions); err != nil {
			log.Printf("为 admin 角色分配权限失败: %v", err)
		} else {
			log.Printf("已为 admin 角色分配 %d 个权限", len(permissions))
		}
	} else {
		log.Printf("admin 角色已有 %d 个权限，跳过初始化", count)
	}
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
