package middleware

import (
    "net/http"
    "strings"

    "student-management-system/config"
    "student-management-system/internal/models"
    "student-management-system/internal/utils"

    "github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请求头中auth为空",
			})
			c.Abort()
			return
		}

		// Authorization: Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请求头中auth格式有误",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的Token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)

		// 查询该角色拥有的所有权限
		var role models.Role
		if err := config.GetDB().Preload("Permissions").First(&role, claims.RoleID).Error; err == nil {
			// 构建权限映射表（map[string]bool）用于快速查找
			permissionMap := make(map[string]bool)
			for _, perm := range role.Permissions {
				permissionMap[perm.Permission] = true
			}
			c.Set("permissions", permissionMap)
		} else {
			// 如果查询失败，设置空权限映射
			c.Set("permissions", make(map[string]bool))
		}

		c.Next()
	}
}

// RoleMiddleware 角色权限中间件（保留用于向后兼容）
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限访问",
			})
			c.Abort()
			return
		}

        // 查询Role表获取角色名称，避免依赖固定ID
        var role models.Role
        if err := config.GetDB().First(&role, roleID.(uint)).Error; err != nil {
            c.JSON(http.StatusForbidden, gin.H{
                "code":    403,
                "message": "无权限访问",
            })
            c.Abort()
            return
        }

        userRole := role.RoleName
		allowed := false
		for _, role := range allowedRoles {
            if role == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionMiddleware 检查用户是否具有特定权限
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 AuthMiddleware 设置的上下文中获取权限列表
		perms, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "获取权限失败"})
			c.Abort()
			return
		}

		// 类型断言 (假设存的是 map[string]bool)
		permissionMap, ok := perms.(map[string]bool)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限格式错误"})
			c.Abort()
			return
		}

		// 检查权限是否存在
		if _, ok := permissionMap[requiredPermission]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限访问"})
			c.Abort()
			return
		}

		c.Next()
	}
}
