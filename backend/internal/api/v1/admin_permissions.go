package v1

import (
	"net/http"
	"strconv"

	"student-management-system/config"
	"student-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// 定义系统中所有可用的权限
var AllPermissions = []models.Permission{
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

// AdminListPermissions 获取所有可用的权限列表
func AdminListPermissions(c *gin.Context) {
	db := config.GetDB()

	// 从数据库获取所有权限（如果数据库中没有，则初始化）
	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 如果数据库中没有权限，则初始化
	if len(permissions) == 0 {
		for _, perm := range AllPermissions {
			var existing models.Permission
			if err := db.Where("permission = ?", perm.Permission).First(&existing).Error; err != nil {
				// 权限不存在，创建它
				db.Create(&perm)
				permissions = append(permissions, perm)
			} else {
				permissions = append(permissions, existing)
			}
		}
		// 重新查询
		db.Find(&permissions)
	}

	// 按分组组织权限
	grouped := make(map[string][]models.Permission)
	for _, perm := range permissions {
		grouped[perm.Group] = append(grouped[perm.Group], perm)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":    permissions,
			"grouped": grouped,
		},
	})
}

// AdminGetRolePermissions 获取特定角色的权限列表
func AdminGetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	db := config.GetDB()
	var role models.Role
	if err := db.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	// 返回权限标识列表
	permissions := make([]string, 0)
	for _, perm := range role.Permissions {
		permissions = append(permissions, perm.Permission)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"role_id":     role.ID,
			"role_name":   role.RoleName,
			"permissions": permissions,
		},
	})
}

// UpdateRolePermissionsRequest 更新角色权限的请求体
type UpdateRolePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"` // 权限标识数组
}

// AdminUpdateRolePermissions 更新特定角色的权限
func AdminUpdateRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var req UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	db := config.GetDB()

	// 检查角色是否存在
	var role models.Role
	if err := db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	// 验证权限标识是否存在
	var permissions []models.Permission
	if err := db.Where("permission IN ?", req.Permissions).Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询权限失败", "error": err.Error()})
		return
	}

	if len(permissions) != len(req.Permissions) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "部分权限标识不存在"})
		return
	}

	// 更新角色的权限（使用 Association 替换）
	if err := db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新权限失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data": gin.H{
			"role_id":     role.ID,
			"role_name":   role.RoleName,
			"permissions": req.Permissions,
		},
	})
}
