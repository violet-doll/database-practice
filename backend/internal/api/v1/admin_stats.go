package v1

import (
    "net/http"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

// AdminOverviewStats 系统概览统计
func AdminOverviewStats(c *gin.Context) {
    db := config.GetDB()

    var totalUsers int64
    var activeUsers int64
    var totalStudents int64
    var totalTeachers int64
    var totalClasses int64
    var totalCourses int64
    var totalEnrollments int64
    var totalGrades int64
    var totalAttendance int64

    db.Model(&models.User{}).Count(&totalUsers)
    db.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)
    db.Model(&models.Student{}).Count(&totalStudents)
    db.Model(&models.Teacher{}).Count(&totalTeachers)
    db.Model(&models.Class{}).Count(&totalClasses)
    db.Model(&models.Course{}).Count(&totalCourses)
    db.Model(&models.Enrollment{}).Count(&totalEnrollments)
    db.Model(&models.Grade{}).Count(&totalGrades)
    db.Model(&models.Attendance{}).Count(&totalAttendance)

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data": gin.H{
            "users_total":        totalUsers,
            "users_active":       activeUsers,
            "students_total":     totalStudents,
            "teachers_total":     totalTeachers,
            "classes_total":      totalClasses,
            "courses_total":      totalCourses,
            "enrollments_total":  totalEnrollments,
            "grades_total":       totalGrades,
            "attendance_total":   totalAttendance,
        },
    })
}


