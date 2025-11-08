package v1

import (
	"net/http"
	"strconv"

	"student-management-system/config"
	"student-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// GetSchedules 获取排课列表（分页、支持按班级、教师、学期筛选）
func GetSchedules(c *gin.Context) {
	db := config.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	classID := c.Query("class_id")
	teacherID := c.Query("teacher_id")
	semester := c.Query("semester")

	query := db.Model(&models.Schedule{})
	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}

	var total int64
	query.Count(&total)

	var list []models.Schedule
	query.Preload("Course").Preload("Class").Preload("Teacher").
		Limit(pageSize).Offset(offset).Find(&list)

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

// GetSchedule 获取排课详情
func GetSchedule(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var schedule models.Schedule
	if err := db.Preload("Course").Preload("Class").Preload("Teacher").First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "排课记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": schedule})
}

// CreateSchedule 新增排课
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	db := config.GetDB()

	// 验证课程是否存在
	var course models.Course
	if err := db.First(&course, schedule.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "课程不存在", "error": err.Error()})
		return
	}

	// 验证班级是否存在
	var class models.Class
	if err := db.First(&class, schedule.ClassID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "班级不存在", "error": err.Error()})
		return
	}

	// 验证教师是否存在
	var teacher models.Teacher
	if err := db.First(&teacher, schedule.TeacherID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "教师不存在", "error": err.Error()})
		return
	}

	// 验证星期几范围 (1-7)
	if schedule.DayOfWeek < 1 || schedule.DayOfWeek > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "星期几必须在1-7之间"})
		return
	}

	if err := db.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
		return
	}

	// 预加载关联数据以便返回完整信息
	db.Preload("Course").Preload("Class").Preload("Teacher").First(&schedule, schedule.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": schedule})
}

// UpdateSchedule 更新排课
func UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var schedule models.Schedule
	if err := db.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "排课记录不存在"})
		return
	}

	// 绑定更新数据
	var updateData models.Schedule
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 如果更新了关联ID，验证关联是否存在
	if updateData.CourseID != 0 && updateData.CourseID != schedule.CourseID {
		var course models.Course
		if err := db.First(&course, updateData.CourseID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "课程不存在", "error": err.Error()})
			return
		}
		schedule.CourseID = updateData.CourseID
	}

	if updateData.ClassID != 0 && updateData.ClassID != schedule.ClassID {
		var class models.Class
		if err := db.First(&class, updateData.ClassID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "班级不存在", "error": err.Error()})
			return
		}
		schedule.ClassID = updateData.ClassID
	}

	if updateData.TeacherID != 0 && updateData.TeacherID != schedule.TeacherID {
		var teacher models.Teacher
		if err := db.First(&teacher, updateData.TeacherID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "教师不存在", "error": err.Error()})
			return
		}
		schedule.TeacherID = updateData.TeacherID
	}

	// 更新其他字段
	if updateData.DayOfWeek != 0 {
		if updateData.DayOfWeek < 1 || updateData.DayOfWeek > 7 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "星期几必须在1-7之间"})
			return
		}
		schedule.DayOfWeek = updateData.DayOfWeek
	}
	if updateData.StartTime != "" {
		schedule.StartTime = updateData.StartTime
	}
	if updateData.EndTime != "" {
		schedule.EndTime = updateData.EndTime
	}
	if updateData.Location != "" {
		schedule.Location = updateData.Location
	}
	if updateData.Semester != "" {
		schedule.Semester = updateData.Semester
	}

	if err := db.Save(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}

	// 预加载关联数据以便返回完整信息
	db.Preload("Course").Preload("Class").Preload("Teacher").First(&schedule, schedule.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": schedule})
}

// DeleteSchedule 删除排课
func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var schedule models.Schedule
	if err := db.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "排课记录不存在"})
		return
	}
	if err := db.Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

