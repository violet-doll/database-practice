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
    "gorm.io/gorm/clause"
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
		&models.Permission{},      // 权限表
		&models.Role{},            // 角色表
		&models.RolePermission{},  // 角色-权限关联表
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

// initPermissionsAndAssignRoles 初始化权限并为各角色分配默认权限
func initPermissionsAndAssignRoles() {
	// 定义所有权限（与 admin_permissions.go 中的定义保持一致）
	allPermissions := []models.Permission{
		// 学生管理权限
		{Name: "查看学生", Permission: "student:read", Group: "student"},
		{Name: "创建学生", Permission: "student:create", Group: "student"},
		{Name: "修改学生", Permission: "student:update", Group: "student"},
		{Name: "删除学生", Permission: "student:delete", Group: "student"},

		// 班级管理权限
		{Name: "查看班级", Permission: "class:read", Group: "class"},
		{Name: "创建班级", Permission: "class:create", Group: "class"},
		{Name: "修改班级", Permission: "class:update", Group: "class"},
		{Name: "删除班级", Permission: "class:delete", Group: "class"},

		// 课程管理权限
		{Name: "查看课程", Permission: "course:read", Group: "course"},
		{Name: "创建课程", Permission: "course:create", Group: "course"},
		{Name: "修改课程", Permission: "course:update", Group: "course"},
		{Name: "删除课程", Permission: "course:delete", Group: "course"},

		// 排课管理权限
		{Name: "查看排课", Permission: "schedule:read", Group: "schedule"},
		{Name: "创建排课", Permission: "schedule:create", Group: "schedule"},
		{Name: "修改排课", Permission: "schedule:update", Group: "schedule"},
		{Name: "删除排课", Permission: "schedule:delete", Group: "schedule"},

		// 选课管理权限
		{Name: "查看选课", Permission: "enrollment:read", Group: "enrollment"},
		{Name: "创建选课", Permission: "enrollment:create", Group: "enrollment"},
		{Name: "删除选课", Permission: "enrollment:delete", Group: "enrollment"},

		// 成绩管理权限
		{Name: "查看成绩", Permission: "grade:read", Group: "grade"},
		{Name: "创建成绩", Permission: "grade:create", Group: "grade"},
		{Name: "修改成绩", Permission: "grade:update", Group: "grade"},

		// 考勤管理权限
		{Name: "查看考勤", Permission: "attendance:read", Group: "attendance"},
		{Name: "创建考勤", Permission: "attendance:create", Group: "attendance"},
		{Name: "删除考勤", Permission: "attendance:delete", Group: "attendance"},

		// 奖惩管理权限
		{Name: "查看奖惩", Permission: "reward:read", Group: "reward"},
		{Name: "创建奖惩", Permission: "reward:create", Group: "reward"},
		{Name: "删除奖惩", Permission: "reward:delete", Group: "reward"},

		// 家长管理权限
		{Name: "查看家长", Permission: "parent:read", Group: "parent"},
		{Name: "创建家长", Permission: "parent:create", Group: "parent"},
		{Name: "修改家长", Permission: "parent:update", Group: "parent"},
		{Name: "删除家长", Permission: "parent:delete", Group: "parent"},

		// 通知管理权限
		{Name: "查看通知", Permission: "notification:read", Group: "notification"},
		{Name: "创建通知", Permission: "notification:create", Group: "notification"},

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

	// 为 teacher 角色分配权限
	var teacherRole models.Role
	if err := DB.Where("role_name = ?", "teacher").First(&teacherRole).Error; err == nil {
		count := DB.Model(&teacherRole).Association("Permissions").Count()
		if count == 0 {
			// 教师权限：查看和创建课程、成绩、考勤、奖惩等
			teacherPerms := []string{
				"course:read", "course:create", "course:update",
				"schedule:read", "schedule:create", "schedule:update",
				"enrollment:read", "enrollment:create", "enrollment:delete",
				"grade:read", "grade:create", "grade:update",
				"attendance:read", "attendance:create", "attendance:delete",
				"reward:read", "reward:create", "reward:delete",
				"student:read", "class:read",
				"notification:read", "notification:create",
			}
			assignPermissionsToRole(&teacherRole, teacherPerms, "teacher")
		}
	}

	// 为学生角色分配权限
	var studentRole models.Role
	if err := DB.Where("role_name = ?", "student").First(&studentRole).Error; err == nil {
		count := DB.Model(&studentRole).Association("Permissions").Count()
		if count == 0 {
			// 学生权限：查看自己的信息、课程、成绩、考勤、奖惩
			studentPerms := []string{
				"student:read",      // 只能查看自己的信息（需要在业务逻辑中限制）
				"course:read",
				"enrollment:read",
				"grade:read",        // 只能查看自己的成绩（需要在业务逻辑中限制）
				"attendance:read",  // 只能查看自己的考勤（需要在业务逻辑中限制）
				"reward:read",      // 只能查看自己的奖惩（需要在业务逻辑中限制）
				"notification:read",
			}
			assignPermissionsToRole(&studentRole, studentPerms, "student")
		}
	}

	// 为家长角色分配权限
	var parentRole models.Role
	if err := DB.Where("role_name = ?", "parent").First(&parentRole).Error; err == nil {
		count := DB.Model(&parentRole).Association("Permissions").Count()
		if count == 0 {
			// 家长权限：查看关联学生的成绩、考勤、奖惩
			parentPerms := []string{
				"student:read",     // 只能查看关联学生的信息（需要在业务逻辑中限制）
				"grade:read",        // 只能查看关联学生的成绩（需要在业务逻辑中限制）
				"attendance:read",   // 只能查看关联学生的考勤（需要在业务逻辑中限制）
				"reward:read",       // 只能查看关联学生的奖惩（需要在业务逻辑中限制）
				"notification:read",
			}
			assignPermissionsToRole(&parentRole, parentPerms, "parent")
		}
	}
}

// assignPermissionsToRole 为角色分配权限
func assignPermissionsToRole(role *models.Role, permissionStrings []string, roleName string) {
	var permissions []models.Permission
	for _, permStr := range permissionStrings {
		var perm models.Permission
		if err := DB.Where("permission = ?", permStr).First(&perm).Error; err == nil {
			permissions = append(permissions, perm)
		}
	}
	if len(permissions) > 0 {
		if err := DB.Model(role).Association("Permissions").Replace(permissions); err != nil {
			log.Printf("为 %s 角色分配权限失败: %v", roleName, err)
		} else {
			log.Printf("已为 %s 角色分配 %d 个权限", roleName, len(permissions))
		}
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
