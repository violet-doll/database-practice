package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

// GetCourses 获取课程列表（分页、搜索course_name、teacher_id）
func GetCourses(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    courseName := c.Query("course_name")
    teacherID := c.Query("teacher_id")

    query := db.Model(&models.Course{})
    if courseName != "" {
        query = query.Where("course_name LIKE ?", "%"+courseName+"%")
    }
    if teacherID != "" {
        query = query.Where("teacher_id = ?", teacherID)
    }

    var total int64
    query.Count(&total)

    var list []models.Course
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

// GetCourse 获取课程详情
func GetCourse(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    var course models.Course
    if err := db.Preload("Teacher").First(&course, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课程不存在"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": course})
}

// CreateCourse 新增课程
func CreateCourse(c *gin.Context) {
    var course models.Course
    if err := c.ShouldBindJSON(&course); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }
    db := config.GetDB()
    if err := db.Create(&course).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": course})
}

// UpdateCourse 更新课程
func UpdateCourse(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    var course models.Course
    if err := db.First(&course, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课程不存在"})
        return
    }
    if err := c.ShouldBindJSON(&course); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }
    if err := db.Save(&course).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": course})
}

// DeleteCourse 删除课程
func DeleteCourse(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    var course models.Course
    if err := db.First(&course, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课程不存在"})
        return
    }
    if err := db.Delete(&course).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


