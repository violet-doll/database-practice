package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

type CreateRoleRequest struct {
    RoleName string `json:"role_name" binding:"required"`
}

// AdminListRoles 角色列表（分页）
func AdminListRoles(c *gin.Context) {
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
    var roles []models.Role
    q := db.Model(&models.Role{})
    if name := c.Query("role_name"); name != "" {
        q = q.Where("role_name LIKE ?", "%"+name+"%")
    }

    if err := q.Count(&total).Order("id ASC").Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data": gin.H{
            "list":      roles,
            "total":     total,
            "page":      page,
            "page_size": pageSize,
        },
    })
}

// AdminCreateRole 新增角色
func AdminCreateRole(c *gin.Context) {
    var req CreateRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    role := models.Role{RoleName: req.RoleName}
    if err := config.GetDB().Create(&role).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": role})
}

// AdminUpdateRole 更新角色名
func AdminUpdateRole(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
        return
    }
    var req CreateRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }
    db := config.GetDB()
    var role models.Role
    if err := db.First(&role, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
        return
    }
    role.RoleName = req.RoleName
    if err := db.Save(&role).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": role})
}

// AdminDeleteRole 删除角色
func AdminDeleteRole(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
        return
    }
    if err := config.GetDB().Delete(&models.Role{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


