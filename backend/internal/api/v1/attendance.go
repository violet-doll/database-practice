package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type CreateAttendanceRequest struct {
    StudentID uint   `json:"student_id" binding:"required"`
    Date      string `json:"date" binding:"required"`   // YYYY-MM-DD
    Status    string `json:"status" binding:"required"` // 出勤/缺席/请假/迟到
    Reason    string `json:"reason"`
    TeacherID uint   `json:"teacher_id"`
}

// CreateAttendance 新增考勤记录
func CreateAttendance(c *gin.Context) {
    var req CreateAttendanceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    // 简单校验状态
    switch req.Status {
    case "出勤", "缺席", "请假", "迟到":
    default:
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的考勤状态"})
        return
    }

    db := config.GetDB()
    record := models.Attendance{
        StudentID: req.StudentID,
        Date:      req.Date,
        Status:    req.Status,
        Reason:    req.Reason,
        TeacherID: req.TeacherID,
    }
    if err := db.Create(&record).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建考勤记录失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": record})
}

// GetAttendance 考勤列表（分页+筛选: student_id, status, date_from, date_to）
func GetAttendance(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    studentID := c.Query("student_id")
    status := c.Query("status")
    dateFrom := c.Query("date_from")
    dateTo := c.Query("date_to")

    query := db.Model(&models.Attendance{}).Preload("Student")
    if studentID != "" {
        query = query.Where("student_id = ?", studentID)
    }
    if status != "" {
        query = query.Where("status = ?", status)
    }
    if dateFrom != "" {
        query = query.Where("date >= ?", dateFrom)
    }
    if dateTo != "" {
        query = query.Where("date <= ?", dateTo)
    }

    var total int64
    query.Count(&total)

    var list []models.Attendance
    if err := query.Order("date DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
        return
    }

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

// GetAttendanceByStudent 按学生查询
func GetAttendanceByStudent(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    var total int64
    db.Model(&models.Attendance{}).Where("student_id = ?", id).Count(&total)

    var list []models.Attendance
    if err := db.Preload("Student").Where("student_id = ?", id).
        Order("date DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
        return
    }

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

type AttendanceStatsItem struct {
    StudentID uint                `json:"student_id"`
    Student   *models.Student     `json:"student"`
    Present   int                 `json:"present"`
    Absent    int                 `json:"absent"`
    Leave     int                 `json:"leave"`
    Late      int                 `json:"late"`
    Total     int                 `json:"total"`
}

// GetAttendanceStats 考勤统计（支持按日期范围、可选 student_id；聚合到学生粒度）
func GetAttendanceStats(c *gin.Context) {
    db := config.GetDB()

    studentID := c.Query("student_id")
    dateFrom := c.Query("date_from")
    dateTo := c.Query("date_to")

    var records []models.Attendance
    query := db.Preload("Student").Model(&models.Attendance{})
    if studentID != "" {
        query = query.Where("student_id = ?", studentID)
    }
    if dateFrom != "" {
        query = query.Where("date >= ?", dateFrom)
    }
    if dateTo != "" {
        query = query.Where("date <= ?", dateTo)
    }
    if err := query.Find(&records).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
        return
    }

    statsMap := map[uint]*AttendanceStatsItem{}
    for i := range records {
        rec := records[i]
        item, ok := statsMap[rec.StudentID]
        if !ok {
            item = &AttendanceStatsItem{StudentID: rec.StudentID}
            // 如果已预载 Student，带上
            if rec.Student.ID != 0 {
                stu := rec.Student
                item.Student = &stu
            }
            statsMap[rec.StudentID] = item
        }
        item.Total++
        switch rec.Status {
        case "出勤":
            item.Present++
        case "缺席":
            item.Absent++
        case "请假":
            item.Leave++
        case "迟到":
            item.Late++
        }
    }

    // 输出为切片
    result := make([]AttendanceStatsItem, 0, len(statsMap))
    for _, v := range statsMap {
        // 确保 Student 填充（若未预载，补查一次）
        if v.Student == nil {
            var stu models.Student
            _ = db.First(&stu, v.StudentID).Error
            if stu.ID != 0 {
                s := stu
                v.Student = &s
            }
        }
        result = append(result, *v)
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data": gin.H{
            "list": result,
        },
    })
}

// helper to avoid unused import if certain build tags change
var _ *gorm.DB

// DeleteAttendance 删除考勤记录
func DeleteAttendance(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    // 先检查是否存在
    var rec models.Attendance
    if err := db.First(&rec, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
        return
    }

    if err := db.Delete(&models.Attendance{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


