package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

// GetClasses 获取班级列表（分页、搜索 class_name、teacher_id）
func GetClasses(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    className := c.Query("class_name")
    teacherID := c.Query("teacher_id")

    query := db.Model(&models.Class{})
    if className != "" {
        query = query.Where("class_name LIKE ?", "%"+className+"%")
    }
    if teacherID != "" {
        query = query.Where("teacher_id = ?", teacherID)
    }

    var total int64
    query.Count(&total)

    var list []models.Class
    query.Preload("Teacher").Limit(pageSize).Offset(offset).Find(&list)

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data": gin.H{
            "list":      list,
            "total":     total,
            "page":      page,
            "page_size": pageSize,
        },
    })
}

// GetClass 获取班级详情
func GetClass(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    var cls models.Class
    if err := db.Preload("Teacher").Preload("Students").First(&cls, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "班级不存在"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": cls})
}

// CreateClass 新增班级
func CreateClass(c *gin.Context) {
    var cls models.Class
    if err := c.ShouldBindJSON(&cls); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    db := config.GetDB()
    if err := db.Create(&cls).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": cls})
}

// UpdateClass 更新班级
func UpdateClass(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    var cls models.Class
    if err := db.First(&cls, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "班级不存在"})
        return
    }

    if err := c.ShouldBindJSON(&cls); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    if err := db.Save(&cls).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": cls})
}

// DeleteClass 删除班级
func DeleteClass(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    var cls models.Class
    if err := db.First(&cls, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "班级不存在"})
        return
    }

    if err := db.Delete(&cls).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


