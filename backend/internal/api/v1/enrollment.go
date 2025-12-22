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

	// 使用存储过程进行选课
	// GORM 调用存储过程
	// 注意：这里假设 req.StudentNumber 是学号字符串，但存储过程需要 student_id (INT)
	// 我们需要先查询 student_id
	var student models.Student
	if err := db.Where("student_id = ?", req.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "学生不存在"})
		return
	}

	// 执行存储过程
	// CALL sp_enroll_student(p_student_id, p_course_id, @p_status, @p_message)
	err := db.Raw("CALL sp_enroll_student(?, ?, @status, @message)", student.ID, req.CourseID).Scan(&struct{}{}).Error
	if err != nil {
		// 存储过程执行错误（可能是SQL错误）
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "选课系统异常", "error": err.Error()})
		return
	}

	// 获取 OUT 参数
	type Result struct {
		Status  int
		Message string
	}
	var res Result
	db.Raw("SELECT @status as status, @message as message").Scan(&res)

	if res.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": res.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "选课成功"})
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
