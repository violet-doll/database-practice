package v1

import (
	"net/http"
	"strconv"

	"student-management-system/config"
	"student-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// GetStudents 获取学生列表
func GetStudents(c *gin.Context) {
	db := config.GetDB()

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 获取搜索参数
	name := c.Query("name")
	studentID := c.Query("student_id")
	classID := c.Query("class_id")

	// 构建查询
	query := db.Model(&models.Student{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if studentID != "" {
		query = query.Where("student_id LIKE ?", "%"+studentID+"%")
	}
	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var students []models.Student
	query.Preload("Class").Preload("Parents").
		Limit(pageSize).Offset(offset).
		Find(&students)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      students,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetStudent 获取单个学生详情
func GetStudent(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var student models.Student
	if err := db.Preload("Class").Preload("Parents").First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    student,
	})
}

// CreateStudent 创建学生
func CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	db := config.GetDB()

	var existingStudent models.Student
	if err := db.Unscoped().Where("student_id = ?", student.StudentID).First(&existingStudent).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "创建失败：该学号已存在",
		})
		return
	}

	if err := db.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    student,
	})
}

// UpdateStudent 更新学生信息
func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var student models.Student
	if err := db.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生不存在",
		})
		return
	}

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := db.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    student,
	})
}

// DeleteStudent 删除学生
func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var student models.Student
	if err := db.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生不存在",
		})
		return
	}

	if err := db.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
