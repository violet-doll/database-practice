package v1

import (
	"net/http"
	"strconv"

	"student-management-system/config"
	"student-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// GetEnrollments 选课列表（全部，支持分页与可选筛选 student_id, course_id）
func GetEnrollments(c *gin.Context) {
	db := config.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	studentID := c.Query("student_id")
	courseID := c.Query("course_id")

	query := db.Model(&models.Enrollment{})
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}

	var total int64
	query.Count(&total)

	var enrollments []models.Enrollment
	query.Preload("Course").Preload("Student").Preload("Grades").
		Limit(pageSize).Offset(offset).Find(&enrollments)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      enrollments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

type CreateEnrollmentRequest struct {
	StudentNumber string `json:"student_number" binding:"required"` // 学号
	CourseID      uint   `json:"course_id" binding:"required"`
}

// CreateEnrollment 创建选课记录
func CreateEnrollment(c *gin.Context) {
	var req CreateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	db := config.GetDB()

	// 根据学号查找学生
	var student models.Student
	if err := db.Where("student_id = ?", req.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "学生不存在（学号不存在）", "error": err.Error()})
		return
	}

	// 检查该 Enrollment 是否已存在
	var existingEnrollment models.Enrollment
	if err := db.Where("student_id = ? AND course_id = ?", student.ID, req.CourseID).First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该选课记录已存在"})
		return
	}

	// 验证课程是否存在
	var course models.Course
	if err := db.First(&course, req.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "课程不存在", "error": err.Error()})
		return
	}

	// 创建新的 Enrollment 记录
	enrollment := models.Enrollment{
		StudentID: student.ID,
		CourseID:  req.CourseID,
	}

	if err := db.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建选课记录失败", "error": err.Error()})
		return
	}

	// 预加载关联数据以便返回完整信息
	db.Preload("Course").Preload("Student").First(&enrollment, enrollment.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": enrollment})
}

// DeleteEnrollment 删除选课记录
func DeleteEnrollment(c *gin.Context) {
	id := c.Param("id")
	enrollmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的选课记录ID"})
		return
	}

	db := config.GetDB()

	// 查找选课记录
	var enrollment models.Enrollment
	if err := db.First(&enrollment, uint(enrollmentID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "选课记录不存在"})
		return
	}

	// 检查是否有关联的成绩记录
	var gradeCount int64
	db.Model(&models.Grade{}).Where("enrollment_id = ?", enrollment.ID).Count(&gradeCount)

	// 如果有成绩记录，可以选择级联删除或阻止删除
	// 这里我们选择级联删除：先删除成绩记录，再删除选课记录
	if gradeCount > 0 {
		if err := db.Where("enrollment_id = ?", enrollment.ID).Delete(&models.Grade{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除关联成绩记录失败", "error": err.Error()})
			return
		}
	}

	// 删除选课记录
	if err := db.Delete(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除选课记录失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

