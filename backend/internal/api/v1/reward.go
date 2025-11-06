package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

type CreateRewardRequest struct {
    StudentID   uint   `json:"student_id" binding:"required"`
    Type        string `json:"type" binding:"required"`        // 奖励/处分
    Description string `json:"description" binding:"required"` // 事由
    Date        string `json:"date" binding:"required"`        // YYYY-MM-DD
    Issuer      string `json:"issuer"`
}

// CreateReward 奖惩记录录入
func CreateReward(c *gin.Context) {
    var req CreateRewardRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    if req.Type != "奖励" && req.Type != "处分" {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "type 仅支持 '奖励' 或 '处分'"})
        return
    }

    db := config.GetDB()
    record := models.RewardPunishment{
        StudentID:   req.StudentID,
        Type:        req.Type,
        Description: req.Description,
        Date:        req.Date,
        Issuer:      req.Issuer,
    }
    if err := db.Create(&record).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建奖惩记录失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": record})
}

// GetRewards 奖惩记录查询（分页+筛选: student_id, type, date_from, date_to）
func GetRewards(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    studentID := c.Query("student_id")
    rtype := c.Query("type")
    dateFrom := c.Query("date_from")
    dateTo := c.Query("date_to")

    query := db.Model(&models.RewardPunishment{}).Preload("Student")
    if studentID != "" {
        query = query.Where("student_id = ?", studentID)
    }
    if rtype != "" {
        query = query.Where("type = ?", rtype)
    }
    if dateFrom != "" {
        query = query.Where("date >= ?", dateFrom)
    }
    if dateTo != "" {
        query = query.Where("date <= ?", dateTo)
    }

    var total int64
    query.Count(&total)

    var list []models.RewardPunishment
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

// GetRewardsByStudent 按学生查询
func GetRewardsByStudent(c *gin.Context) {
    id := c.Param("id")
    c.Request.URL.Query()
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    db := config.GetDB()
    query := db.Where("student_id = ?", id).Model(&models.RewardPunishment{}).Preload("Student")

    var total int64
    query.Count(&total)

    var list []models.RewardPunishment
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

// DeleteReward 删除奖惩记录
func DeleteReward(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()
    if err := db.Delete(&models.RewardPunishment{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


