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
		}

		// 学生管理（需要认证）
		students := apiV1.Group("/students")
		students.Use(middleware.AuthMiddleware())
		{
			students.GET("", v1.GetStudents)
			students.GET("/:id", v1.GetStudent)
			students.POST("", v1.CreateStudent)
			students.PUT("/:id", v1.UpdateStudent)
			students.DELETE("/:id", v1.DeleteStudent)
		}

		// TODO: 添加其他路由
		// 班级管理
		// classes := apiV1.Group("/classes")
		// classes.Use(middleware.AuthMiddleware())
		// {
		//     classes.GET("", v1.GetClasses)
		//     ...
		// }

		// 课程管理
		courses := apiV1.Group("/courses")
		courses.Use(middleware.AuthMiddleware())
		{
			courses.GET("", v1.GetCourses)
			courses.GET("/:id", v1.GetCourse)
			courses.POST("", v1.CreateCourse)
			courses.PUT("/:id", v1.UpdateCourse)
			courses.DELETE("/:id", v1.DeleteCourse)
		}

		// 成绩管理
		grades := apiV1.Group("/grades")
		grades.Use(middleware.AuthMiddleware())
		{
			grades.GET("", v1.GetGrades) // 新增：全部成绩分页列表（可筛选）
			grades.POST("", v1.CreateGrade)
			grades.GET("/student/:id", v1.GetGradesByStudent)
			grades.GET("/course/:id", v1.GetGradesByCourse)
		}

		// 考勤管理
		attendance := apiV1.Group("/attendance")
		attendance.Use(middleware.AuthMiddleware())
		{
			attendance.GET("", v1.GetAttendance)                      // 列表（分页+筛选）
			attendance.POST("", v1.CreateAttendance)                   // 新增记录（签到/缺席/请假/迟到）
			attendance.GET("/student/:id", v1.GetAttendanceByStudent)  // 按学生查询
			attendance.GET("/stats", v1.GetAttendanceStats)            // 统计
			attendance.DELETE(":id", v1.DeleteAttendance)              // 删除记录
		}

		// 奖惩管理
		rewards := apiV1.Group("/rewards")
		rewards.Use(middleware.AuthMiddleware())
		{
			rewards.GET("", v1.GetRewards)                    // 列表（分页+筛选）
			rewards.POST("", v1.CreateReward)                  // 新增奖惩记录
			rewards.GET("/student/:id", v1.GetRewardsByStudent) // 按学生查询
			rewards.DELETE(":id", v1.DeleteReward)             // 删除记录
		}


		// 家长联系方式管理
		parents := apiV1.Group("/parents")
		parents.Use(middleware.AuthMiddleware())
		{
			parents.GET("", v1.GetParents)
			parents.POST("", v1.CreateParent)
			parents.PUT(":id", v1.UpdateParent)
			parents.DELETE(":id", v1.DeleteParent)
		}

		// 通知管理
		notifications := apiV1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware())
		{
			notifications.GET("", v1.GetNotifications)
			notifications.POST("", v1.CreateNotification)
		}
	}

	return r
}
