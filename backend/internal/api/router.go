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
		// courses := apiV1.Group("/courses")
		// ...

		// 成绩管理
		// grades := apiV1.Group("/grades")
		// ...

		// 考勤管理
		// attendance := apiV1.Group("/attendance")
		// ...

		// 奖惩管理
		// rewards := apiV1.Group("/rewards")
		// ...

		// 通知管理
		// notifications := apiV1.Group("/notifications")
		// ...
	}

	return r
}
