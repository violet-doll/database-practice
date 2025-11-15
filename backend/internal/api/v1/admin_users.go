package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"
    "student-management-system/internal/utils"

    "github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    RoleID   uint   `json:"role_id" binding:"required"`
    IsActive *bool  `json:"is_active"`
}

type UpdateUserRequest struct {
    Password *string `json:"password"`
    RoleID   *uint   `json:"role_id"`
    IsActive *bool   `json:"is_active"`
}

// AdminListUsers 列出用户（分页、可筛选）
func AdminListUsers(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 || pageSize > 100 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize

    var total int64
    var users []models.User
    query := db.Model(&models.User{}).Preload("Role")

    if username := c.Query("username"); username != "" {
        query = query.Where("username LIKE ?", "%"+username+"%")
    }
    if roleID := c.Query("role_id"); roleID != "" {
        query = query.Where("role_id = ?", roleID)
    }
    if isActive := c.Query("is_active"); isActive != "" {
        query = query.Where("is_active = ?", isActive)
    }

    if err := query.Count(&total).Order("id DESC").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data": gin.H{
            "list":      users,
            "total":     total,
            "page":      page,
            "page_size": pageSize,
        },
    })
}

// AdminCreateUser 创建用户
func AdminCreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    hashed, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码处理失败"})
        return
    }

    // 根据角色ID获取角色名称，自动设置 UserType
    var role models.Role
    if err := config.GetDB().First(&role, req.RoleID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的角色ID"})
        return
    }

    // 根据角色名称设置 UserType
    userType := ""
    switch role.RoleName {
    case "admin":
        userType = "admin"
    case "teacher":
        userType = "teacher"
    case "student":
        userType = "student"
    case "parent":
        userType = "parent"
    }

    user := models.User{
        Username: req.Username,
        Password: hashed,
        RoleID:   req.RoleID,
        UserType: userType,
    }
    if req.IsActive != nil {
        user.IsActive = *req.IsActive
    }

    if err := config.GetDB().Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": user})
}

// AdminUpdateUser 更新用户
func AdminUpdateUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
        return
    }
    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    db := config.GetDB()
    var user models.User
    if err := db.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
        return
    }

    if req.Password != nil {
        hashed, err := utils.HashPassword(*req.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码处理失败"})
            return
        }
        user.Password = hashed
    }
    if req.RoleID != nil {
        user.RoleID = *req.RoleID
        // 更新角色时，同时更新 UserType
        var role models.Role
        if err := db.First(&role, *req.RoleID).Error; err == nil {
            switch role.RoleName {
            case "admin":
                user.UserType = "admin"
            case "teacher":
                user.UserType = "teacher"
            case "student":
                user.UserType = "student"
            case "parent":
                user.UserType = "parent"
            }
        }
    }
    if req.IsActive != nil {
        user.IsActive = *req.IsActive
    }

    if err := db.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": user})
}

// AdminDeleteUser 删除用户
func AdminDeleteUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
        return
    }
    if err := config.GetDB().Delete(&models.User{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


