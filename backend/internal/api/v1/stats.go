package v1

import (
	"net/http"

	"student-management-system/config"
	"student-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// StatsResponse 数据统计响应结构体
type StatsResponse struct {
	// 学生性别统计
	StudentGender struct {
		MaleCount   int64 `json:"male_count"`
		FemaleCount int64 `json:"female_count"`
	} `json:"student_gender"`

	// 考勤状态统计
	AttendanceStatus struct {
		Present int64 `json:"present"` // 出勤
		Absent  int64 `json:"absent"`  // 缺席
		Leave   int64 `json:"leave"`   // 请假
		Late    int64 `json:"late"`    // 迟到
	} `json:"attendance_status"`

	// 成绩统计
	GradeStats struct {
		AverageScore float64 `json:"average_score"` // 平均分
		TotalGrades  int64   `json:"total_grades"`  // 成绩记录总数
	} `json:"grade_stats"`

	// 班级统计
	ClassStats struct {
		TotalClasses int64 `json:"total_classes"`
	} `json:"class_stats"`

	// 课程统计
	CourseStats struct {
		TotalCourses int64 `json:"total_courses"`
	} `json:"course_stats"`

	// 成绩分布 (新增)
	GradeDistribution map[string]int64 `json:"grade_distribution"`

	// 热门课程 (新增)
	CoursePopularity []struct {
		CourseName string `json:"course_name"`
		Count      int64  `json:"count"`
	} `json:"course_popularity"`
}

// GetDashboardStats 获取数据看板统计信息
// 使用 SQL 聚合函数在数据库层面完成统计计算，只返回聚合后的数据
func GetDashboardStats(c *gin.Context) {
	db := config.GetDB()

	var stats StatsResponse

	// 1. 学生性别统计（按性别分组）
	var genderStats []struct {
		Gender string
		Count  int64
	}
	db.Model(&models.Student{}).
		Select("gender, COUNT(*) as count").
		Group("gender").
		Scan(&genderStats)

	// 初始化性别统计
	stats.StudentGender.MaleCount = 0
	stats.StudentGender.FemaleCount = 0
	for _, g := range genderStats {
		if g.Gender == "男" {
			stats.StudentGender.MaleCount = g.Count
		} else if g.Gender == "女" {
			stats.StudentGender.FemaleCount = g.Count
		}
	}

	// 2. 考勤状态统计（按状态分组）
	var attendanceStats []struct {
		Status string
		Count  int64
	}
	db.Model(&models.Attendance{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&attendanceStats)

	// 初始化考勤统计
	stats.AttendanceStatus.Present = 0
	stats.AttendanceStatus.Absent = 0
	stats.AttendanceStatus.Leave = 0
	stats.AttendanceStatus.Late = 0
	for _, a := range attendanceStats {
		switch a.Status {
		case "出勤":
			stats.AttendanceStatus.Present = a.Count
		case "缺席":
			stats.AttendanceStatus.Absent = a.Count
		case "请假":
			stats.AttendanceStatus.Leave = a.Count
		case "迟到":
			stats.AttendanceStatus.Late = a.Count
		}
	}

	// 3. 成绩统计（计算平均分）
	var avgResult struct {
		AverageScore float64
		TotalCount   int64
	}
	err := db.Model(&models.Grade{}).
		Select("AVG(score) as average_score, COUNT(*) as total_count").
		Scan(&avgResult).Error

	if err == nil {
		stats.GradeStats.AverageScore = avgResult.AverageScore
		stats.GradeStats.TotalGrades = avgResult.TotalCount
	} else {
		// 如果没有成绩记录，设置为0
		stats.GradeStats.AverageScore = 0
		stats.GradeStats.TotalGrades = 0
	}

	// 4. 班级统计
	db.Model(&models.Class{}).Count(&stats.ClassStats.TotalClasses)

	// 5. 课程统计
	db.Model(&models.Course{}).Count(&stats.CourseStats.TotalCourses)

	// 6. 成绩分布统计 (新增)
	// 统计各分数段人数: <60, 60-69, 70-79, 80-89, >=90
	stats.GradeDistribution = make(map[string]int64)
	var grades []float64
	db.Model(&models.Grade{}).Pluck("score", &grades)

	for _, score := range grades {
		if score < 60 {
			stats.GradeDistribution["<60"]++
		} else if score < 70 {
			stats.GradeDistribution["60-69"]++
		} else if score < 80 {
			stats.GradeDistribution["70-79"]++
		} else if score < 90 {
			stats.GradeDistribution["80-89"]++
		} else {
			stats.GradeDistribution[">=90"]++
		}
	}
	// 确保所有区间都有值，方便前端展示
	if _, ok := stats.GradeDistribution["<60"]; !ok {
		stats.GradeDistribution["<60"] = 0
	}
	if _, ok := stats.GradeDistribution["60-69"]; !ok {
		stats.GradeDistribution["60-69"] = 0
	}
	if _, ok := stats.GradeDistribution["70-79"]; !ok {
		stats.GradeDistribution["70-79"] = 0
	}
	if _, ok := stats.GradeDistribution["80-89"]; !ok {
		stats.GradeDistribution["80-89"] = 0
	}
	if _, ok := stats.GradeDistribution[">=90"]; !ok {
		stats.GradeDistribution[">=90"] = 0
	}

	// 7. 热门课程统计 (新增)
	// 统计选课人数最多的前5门课程
	type CoursePop struct {
		CourseName string `json:"course_name"`
		Count      int64  `json:"count"`
	}
	var coursePops []CoursePop
	db.Table("enrollments").
		Select("courses.course_name, count(*) as count").
		Joins("left join courses on enrollments.course_id = courses.id").
		Where("enrollments.deleted_at IS NULL AND courses.deleted_at IS NULL").
		Group("courses.course_name").
		Order("count desc").
		Limit(5).
		Scan(&coursePops)

	stats.CoursePopularity = make([]struct {
		CourseName string `json:"course_name"`
		Count      int64  `json:"count"`
	}, len(coursePops))
	for i, cp := range coursePops {
		stats.CoursePopularity[i] = struct {
			CourseName string `json:"course_name"`
			Count      int64  `json:"count"`
		}{
			CourseName: cp.CourseName,
			Count:      cp.Count,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    stats,
	})
}
