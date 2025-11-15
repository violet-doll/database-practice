package v1

import (
	"net/http"

	"student-management-system/config"
	"student-management-system/internal/models"
	"student-management-system/internal/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string                 `json:"token"`
	User     models.User            `json:"user"`
	UserInfo map[string]interface{} `json:"user_info,omitempty"` // 根据角色返回的详细信息（学生/教师/家长信息）
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查询用户
	var user models.User
	db := config.GetDB()
	if err := db.Where("username = ?", req.Username).Preload("Role").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "用户已被禁用",
		})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Token生成失败",
		})
		return
	}

	// 获取权限列表
	var permissionList []string
	var role models.Role
	if err := db.Preload("Permissions").First(&role, user.RoleID).Error; err == nil {
		for _, perm := range role.Permissions {
			permissionList = append(permissionList, perm.Permission)
		}
	}

	// 根据角色获取详细信息
	response := LoginResponse{
		Token: token,
		User:  user,
	}

	// 根据用户类型和角色返回对应的详细信息
	if user.UserID > 0 {
		switch user.UserType {
		case "student":
			var student models.Student
			if err := db.Preload("Class").Preload("Class.Teacher").First(&student, user.UserID).Error; err == nil {
				response.UserInfo = map[string]interface{}{
					"student": student,
				}
			}
		case "teacher":
			var teacher models.Teacher
			if err := db.First(&teacher, user.UserID).Error; err == nil {
				response.UserInfo = map[string]interface{}{
					"teacher": teacher,
				}
			}
		case "parent":
			var parent models.Parent
			if err := db.Preload("Student").Preload("Student.Class").First(&parent, user.UserID).Error; err == nil {
				response.UserInfo = map[string]interface{}{
					"parent":  parent,
					"student": parent.Student, // 包含关联的学生信息
				}
			}
		}
	}

	// 将权限列表添加到响应中
	responseData := map[string]interface{}{
		"token":       response.Token,
		"user":        response.User,
		"user_info":   response.UserInfo,
		"permissions": permissionList,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    responseData,
	})
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	db := config.GetDB()
	if err := db.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 获取权限列表（从中间件设置的context中获取，如果没有则从数据库查询）
	permissions, exists := c.Get("permissions")
	var permissionList []string
	if exists {
		permissionMap := permissions.(map[string]bool)
		for perm := range permissionMap {
			permissionList = append(permissionList, perm)
		}
	} else {
		// 如果context中没有，从数据库查询
		var role models.Role
		if err := db.Preload("Permissions").First(&role, user.RoleID).Error; err == nil {
			for _, perm := range role.Permissions {
				permissionList = append(permissionList, perm.Permission)
			}
		}
	}

	// 构建响应数据
	response := map[string]interface{}{
		"user":        user,
		"permissions": permissionList,
	}

	// 根据用户类型返回对应的详细信息
	if user.UserID > 0 {
		switch user.UserType {
		case "student":
			var student models.Student
			if err := db.Preload("Class").Preload("Class.Teacher").Preload("Parents").First(&student, user.UserID).Error; err == nil {
				response["student"] = student
			}
		case "teacher":
			var teacher models.Teacher
			if err := db.First(&teacher, user.UserID).Error; err == nil {
				response["teacher"] = teacher
			}
		case "parent":
			var parent models.Parent
			if err := db.Preload("Student").Preload("Student.Class").Preload("Student.Class.Teacher").First(&parent, user.UserID).Error; err == nil {
				response["parent"] = parent
				response["student"] = parent.Student // 包含关联的学生信息
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// JWT是无状态的，登出通常在前端处理（删除token）
	// 如果需要服务端登出，可以使用Redis黑名单机制
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
	})
}
