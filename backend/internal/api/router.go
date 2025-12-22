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
