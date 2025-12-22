package api

import (
	v1 "student-management-system/internal/api/v1"
	"student-management-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS中间件
	r.Use(middleware.CORSMiddleware())

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关（无需认证）
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", v1.Login)
			auth.POST("/logout", v1.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), v1.GetCurrentUser)
			auth.PUT("/password", middleware.AuthMiddleware(), v1.UpdatePassword)
		}

		// 学生管理（需要认证 和 特定权限）
		students := apiV1.Group("/students")
		students.Use(middleware.AuthMiddleware()) // AuthMiddleware 必须在前面, 它负责设置权限列表
		{
			students.GET("", middleware.PermissionMiddleware("student:read"), v1.GetStudents)
			students.GET("/:id", middleware.PermissionMiddleware("student:read"), v1.GetStudent)
			students.POST("", middleware.PermissionMiddleware("student:create"), v1.CreateStudent)
			students.PUT("/:id", middleware.PermissionMiddleware("student:update"), v1.UpdateStudent)
			students.DELETE("/:id", middleware.PermissionMiddleware("student:delete"), v1.DeleteStudent)
		}

		// 班级管理（需要认证 和 特定权限）
		classes := apiV1.Group("/classes")
		classes.Use(middleware.AuthMiddleware())
		{
			classes.GET("", middleware.PermissionMiddleware("class:read"), v1.GetClasses)
			classes.GET(":id", middleware.PermissionMiddleware("class:read"), v1.GetClass)
			classes.POST("", middleware.PermissionMiddleware("class:create"), v1.CreateClass)
			classes.PUT(":id", middleware.PermissionMiddleware("class:update"), v1.UpdateClass)
			classes.DELETE(":id", middleware.PermissionMiddleware("class:delete"), v1.DeleteClass)
		}

		// 课程管理（需要认证 和 特定权限）
		courses := apiV1.Group("/courses")
		courses.Use(middleware.AuthMiddleware())
		{
			courses.GET("", middleware.PermissionMiddleware("course:read"), v1.GetCourses)
			courses.GET("/:id", middleware.PermissionMiddleware("course:read"), v1.GetCourse)
			courses.POST("", middleware.PermissionMiddleware("course:create"), v1.CreateCourse)
			courses.PUT("/:id", middleware.PermissionMiddleware("course:update"), v1.UpdateCourse)
			courses.DELETE("/:id", middleware.PermissionMiddleware("course:delete"), v1.DeleteCourse)
		}

		// 排课管理（需要认证 和 特定权限）
		schedules := apiV1.Group("/schedules")
		schedules.Use(middleware.AuthMiddleware())
		{
			schedules.GET("", middleware.PermissionMiddleware("schedule:read"), v1.GetSchedules)
			schedules.GET("/:id", middleware.PermissionMiddleware("schedule:read"), v1.GetSchedule)
			schedules.POST("", middleware.PermissionMiddleware("schedule:create"), v1.CreateSchedule)
			schedules.PUT("/:id", middleware.PermissionMiddleware("schedule:update"), v1.UpdateSchedule)
			schedules.DELETE("/:id", middleware.PermissionMiddleware("schedule:delete"), v1.DeleteSchedule)
		}

		// 选课管理 (独立于成绩)（需要认证 和 特定权限）
		enrollments := apiV1.Group("/enrollments")
		enrollments.Use(middleware.AuthMiddleware())
		{
			enrollments.GET("", middleware.PermissionMiddleware("enrollment:read"), v1.GetEnrollments)
			enrollments.POST("", middleware.PermissionMiddleware("enrollment:create"), v1.CreateEnrollment)
			enrollments.DELETE("/:id", middleware.PermissionMiddleware("enrollment:delete"), v1.DeleteEnrollment)
		}

		// 成绩管理（需要认证 和 特定权限）
		grades := apiV1.Group("/grades")
		grades.Use(middleware.AuthMiddleware())
		{
			grades.GET("", middleware.PermissionMiddleware("grade:read"), v1.GetGrades) // 新增：全部成绩分页列表（可筛选）
			grades.GET("/student/:id", middleware.PermissionMiddleware("grade:read"), v1.GetGradesByStudent)
			grades.GET("/course/:id", middleware.PermissionMiddleware("grade:read"), v1.GetGradesByCourse)
			grades.POST("", middleware.PermissionMiddleware("grade:create"), v1.CreateGrade)
		}

		// 考勤管理（需要认证 和 特定权限）
		attendance := apiV1.Group("/attendance")
		attendance.Use(middleware.AuthMiddleware())
		{
			attendance.GET("", middleware.PermissionMiddleware("attendance:read"), v1.GetAttendance)                      // 列表（分页+筛选）
			attendance.GET("/student/:id", middleware.PermissionMiddleware("attendance:read"), v1.GetAttendanceByStudent) // 按学生查询
			attendance.GET("/stats", middleware.PermissionMiddleware("attendance:read"), v1.GetAttendanceStats)           // 统计
			attendance.POST("", middleware.PermissionMiddleware("attendance:create"), v1.CreateAttendance)                // 新增记录（签到/缺席/请假/迟到）
			attendance.DELETE(":id", middleware.PermissionMiddleware("attendance:delete"), v1.DeleteAttendance)           // 删除记录
		}

		// 奖惩管理（需要认证 和 特定权限）
		rewards := apiV1.Group("/rewards")
		rewards.Use(middleware.AuthMiddleware())
		{
			rewards.GET("", middleware.PermissionMiddleware("reward:read"), v1.GetRewards)                      // 列表（分页+筛选）
			rewards.GET("/student/:id", middleware.PermissionMiddleware("reward:read"), v1.GetRewardsByStudent) // 按学生查询
			rewards.POST("", middleware.PermissionMiddleware("reward:create"), v1.CreateReward)                 // 新增奖惩记录
			rewards.DELETE(":id", middleware.PermissionMiddleware("reward:delete"), v1.DeleteReward)            // 删除记录
		}

		// 家长联系方式管理（需要认证 和 特定权限）
		parents := apiV1.Group("/parents")
		parents.Use(middleware.AuthMiddleware())
		{
			parents.GET("", middleware.PermissionMiddleware("parent:read"), v1.GetParents)
			parents.POST("", middleware.PermissionMiddleware("parent:create"), v1.CreateParent)
			parents.PUT(":id", middleware.PermissionMiddleware("parent:update"), v1.UpdateParent)
			parents.DELETE(":id", middleware.PermissionMiddleware("parent:delete"), v1.DeleteParent)
		}

		// 通知管理（需要认证 和 特定权限）
		notifications := apiV1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware())
		{
			notifications.GET("", middleware.PermissionMiddleware("notification:read"), v1.GetNotifications)
			notifications.POST("", middleware.PermissionMiddleware("notification:create"), v1.CreateNotification)
		}

		// 数据统计（需要认证 和 特定权限，通常开放给管理员或教师）
		stats := apiV1.Group("/stats")
		stats.Use(middleware.AuthMiddleware())
		{
			stats.GET("/dashboard", middleware.PermissionMiddleware("admin:stats:read"), v1.GetDashboardStats)
		}

		// 管理员模块（需要认证 和 特定权限）
		admin := apiV1.Group("/admin")
		admin.Use(middleware.AuthMiddleware()) // AuthMiddleware 必须在前面
		{
			// 用户管理
			admin.GET("/users", middleware.PermissionMiddleware("admin:user:read"), v1.AdminListUsers)
			admin.POST("/users", middleware.PermissionMiddleware("admin:user:create"), v1.AdminCreateUser)
			admin.PUT("/users/:id", middleware.PermissionMiddleware("admin:user:update"), v1.AdminUpdateUser)
			admin.DELETE("/users/:id", middleware.PermissionMiddleware("admin:user:delete"), v1.AdminDeleteUser)

			// 角色管理
			admin.GET("/roles", middleware.PermissionMiddleware("admin:role:read"), v1.AdminListRoles)
			admin.POST("/roles", middleware.PermissionMiddleware("admin:role:create"), v1.AdminCreateRole)
			admin.PUT("/roles/:id", middleware.PermissionMiddleware("admin:role:update"), v1.AdminUpdateRole)
			admin.DELETE("/roles/:id", middleware.PermissionMiddleware("admin:role:delete"), v1.AdminDeleteRole)

			// 权限管理 (新 API)
			admin.GET("/permissions", middleware.PermissionMiddleware("admin:role:read"), v1.AdminListPermissions)
			admin.GET("/roles/:id/permissions", middleware.PermissionMiddleware("admin:role:read"), v1.AdminGetRolePermissions)
			admin.POST("/roles/:id/permissions", middleware.PermissionMiddleware("admin:role:update"), v1.AdminUpdateRolePermissions)

			// 统计概览
			admin.GET("/stats/overview", middleware.PermissionMiddleware("admin:stats:read"), v1.AdminOverviewStats)
		}

		// 数据库管理（需要认证和管理员权限）
		database := apiV1.Group("/database")
		database.Use(middleware.AuthMiddleware())
		{
			database.GET("/tables", middleware.PermissionMiddleware("admin:user:read"), v1.GetTableList)
			database.GET("/tables/:table", middleware.PermissionMiddleware("admin:user:read"), v1.GetTableData)
			database.GET("/tables/:table/schema", middleware.PermissionMiddleware("admin:user:read"), v1.GetTableSchema)
			database.POST("/tables/:table", middleware.PermissionMiddleware("admin:user:create"), v1.CreateTableData)
			database.PUT("/tables/:table/:id", middleware.PermissionMiddleware("admin:user:update"), v1.UpdateTableData)
			database.DELETE("/tables/:table/:id", middleware.PermissionMiddleware("admin:user:delete"), v1.DeleteTableData)
			database.GET("/tables/:table/export", middleware.PermissionMiddleware("admin:user:read"), v1.ExportTableData)
			database.POST("/execute", middleware.PermissionMiddleware("admin:user:create"), v1.ExecuteSQL)
		}
	}

	return r
}
