package v1

import (
	"net/http"
	"strconv"

	"student-management-system/config"
	"student-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// GetGrades 成绩列表（全部，支持分页与可选筛选 student_id, course_id）
func GetGrades(c *gin.Context) {
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

type CreateGradeRequest struct {
	StudentID uint    `json:"student_id" binding:"required"`
	CourseID  uint    `json:"course_id" binding:"required"`
	ScoreType string  `json:"score_type" binding:"required"`
	Score     float64 `json:"score" binding:"required"`
}

// CreateGrade 成绩录入：若不存在 Enrollment 则创建，再插入成绩
func CreateGrade(c *gin.Context) {
	var req CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	db := config.GetDB()

	// 查找或创建 Enrollment
	var enrollment models.Enrollment
	if err := db.Where("student_id = ? AND course_id = ?", req.StudentID, req.CourseID).First(&enrollment).Error; err != nil {
		enrollment = models.Enrollment{StudentID: req.StudentID, CourseID: req.CourseID}
		if err := db.Create(&enrollment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建选课记录失败", "error": err.Error()})
			return
		}
	}

	grade := models.Grade{
		EnrollmentID: enrollment.ID,
		ScoreType:    req.ScoreType,
		Score:        req.Score,
	}

	if err := db.Create(&grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "录入成绩失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "录入成功", "data": grade})
}

// GetGradesByStudent 按学生查询所有课程成绩（含课程信息、成绩明细）
func GetGradesByStudent(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 查询 Enrollment + 关联（需要课程、成绩，也补充学生信息以便前端显示）
	var enrollments []models.Enrollment
	query := db.Where("student_id = ?", id).
		Preload("Course").
		Preload("Student").
		Preload("Grades")

	var total int64
	query.Model(&models.Enrollment{}).Count(&total)

	query.Limit(pageSize).Offset(offset).Find(&enrollments)

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

// GetGradesByCourse 按课程查询所有学生成绩（含学生信息、成绩明细）
func GetGradesByCourse(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var enrollments []models.Enrollment
	query := db.Where("course_id = ?", id).
		Preload("Course").
		Preload("Student").
		Preload("Grades")

	var total int64
	query.Model(&models.Enrollment{}).Count(&total)

	query.Limit(pageSize).Offset(offset).Find(&enrollments)

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

// GetGradeAuditLogs 获取成绩修改审计日志
func GetGradeAuditLogs(c *gin.Context) {
	db := config.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	gradeID := c.Query("grade_id")

	query := db.Model(&models.GradeAuditLog{})
	if gradeID != "" {
		query = query.Where("grade_id = ?", gradeID)
	}

	var total int64
	query.Count(&total)

	var logs []models.GradeAuditLog
	// Preload Grade and its related Enrollment/Student/Course for context if needed
	// Note: Preloading deep relationships might be heavy, but useful for audit context
	query.Preload("Grade").
		Preload("Grade.Enrollment").
		Preload("Grade.Enrollment.Student").
		Preload("Grade.Enrollment.Course").
		Order("created_at desc").
		Limit(pageSize).Offset(offset).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
