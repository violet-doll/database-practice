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
			"list":   permissions,
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
			"role_id":    role.ID,
			"role_name":  role.RoleName,
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
			"role_id":    role.ID,
			"role_name":  role.RoleName,
			"permissions": req.Permissions,
		},
	})
}

