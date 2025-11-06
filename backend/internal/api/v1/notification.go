package v1

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

// GetNotifications 通知列表（分页，可按 target、keyword(标题/内容) 筛选）
func GetNotifications(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    target := c.Query("target")
    keyword := c.Query("keyword")

    query := db.Model(&models.Notification{})
    if target != "" {
        query = query.Where("target LIKE ?", "%"+target+"%")
    }
    if keyword != "" {
        like := "%" + keyword + "%"
        query = query.Where("title LIKE ? OR content LIKE ?", like, like)
    }

    var total int64
    query.Count(&total)

    var list []models.Notification
    query.Order("id desc").Limit(pageSize).Offset(offset).Find(&list)

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

type CreateNotificationRequest struct {
    Title    string   `json:"title" binding:"required"`
    Content  string   `json:"content" binding:"required"`
    Target   string   `json:"target" binding:"required"` // e.g. "all", "student:12", "class:5", "parent:9"
    Channels []string `json:"channels"`                   // e.g. ["sms","email"], 可为空
}

// CreateNotification 发布通知（模拟发送短信/邮件）
func CreateNotification(c *gin.Context) {
    var req CreateNotificationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    db := config.GetDB()

    // 记录通知
    notification := models.Notification{
        Title:   req.Title,
        Content: req.Content,
        Target:  req.Target,
        // SenderID: 可结合登录用户，从上下文设置，这里简化忽略
    }

    if err := db.Create(&notification).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存通知失败", "error": err.Error()})
        return
    }

    // 模拟发送：根据 Target 解析范围，查询收件人联系方式（家长Phone或学生Phone/Email），此处做最小可行：
    // - 若 target 为 parent:{id}，尝试读取该家长Phone
    // - 若 target 为 student:{id}，尝试读取学生的所有家长Phone
    // - 若 target 为 all/class:{id} 等，当前仅记录为已广播，不逐一枚举
    sentCount := 0
    details := []string{}

    lowerTarget := strings.ToLower(req.Target)
    if strings.HasPrefix(lowerTarget, "parent:") {
        pid := strings.TrimPrefix(lowerTarget, "parent:")
        var parent models.Parent
        if err := db.First(&parent, pid).Error; err == nil {
            if contains(req.Channels, "sms") && parent.Phone != "" {
                details = append(details, fmt.Sprintf("SMS->%s", parent.Phone))
                sentCount++
            }
            // email 通道：Parent 模型暂无 email 字段，略过
        }
    } else if strings.HasPrefix(lowerTarget, "student:") {
        sid := strings.TrimPrefix(lowerTarget, "student:")
        var parents []models.Parent
        db.Where("student_id = ?", sid).Find(&parents)
        for _, p := range parents {
            if contains(req.Channels, "sms") && p.Phone != "" {
                details = append(details, fmt.Sprintf("SMS->%s", p.Phone))
                sentCount++
            }
        }
    } else {
        // 对 all / class:{id} 等情况：不逐一枚举，记录一个广播描述
        details = append(details, "broadcast")
    }

    c.JSON(http.StatusOK, gin.H{
        "code":       200,
        "message":    "发布成功(模拟发送)",
        "data":       notification,
        "sent_count": sentCount,
        "channels":   req.Channels,
        "details":    details,
    })
}

func contains(arr []string, needle string) bool {
    for _, v := range arr {
        if strings.EqualFold(v, needle) {
            return true
        }
    }
    return false
}


