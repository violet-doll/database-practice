package v1

import (
    "net/http"
    "strconv"

    "student-management-system/config"
    "student-management-system/internal/models"

    "github.com/gin-gonic/gin"
)

// GetParents 家长联系方式列表（可按学生筛选，分页）
func GetParents(c *gin.Context) {
    db := config.GetDB()

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize

    studentNo := c.Query("student_no") // 仅使用学生学号筛选

    query := db.Model(&models.Parent{})
    if studentNo != "" {
        // 先用学号模糊匹配出学生主键ID集合，再用其筛选家长
        var matchedStudentIDs []uint
        var students []models.Student
        like := "%" + studentNo + "%"
        db.Where("student_id LIKE ?", like).Find(&students)
        for _, s := range students {
            matchedStudentIDs = append(matchedStudentIDs, s.ID)
        }
        if len(matchedStudentIDs) > 0 {
            query = query.Where("student_id IN ?", matchedStudentIDs)
        } else {
            query = query.Where("1 = 0")
        }
    }

    var total int64
    query.Count(&total)

    var parents []models.Parent
    query.Limit(pageSize).Offset(offset).Find(&parents)

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data": gin.H{
            "list":      parents,
            "total":     total,
            "page":      page,
            "page_size": pageSize,
        },
    })
}

type CreateParentRequest struct {
    // 优先使用 student_no；为兼容前端或历史请求，也接受 student_id 字段作为学号
    StudentNo string `json:"student_no"`
    StudentID string `json:"student_id"`
    Name      string `json:"name" binding:"required"`
    Phone     string `json:"phone" binding:"required"`
    Relation  string `json:"relation" binding:"required"`
}

// CreateParent 新增家长联系方式
func CreateParent(c *gin.Context) {
    var req CreateParentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    db := config.GetDB()

    // 仅按学号字段严格匹配（从 student_no 获取；若为空则回退使用 student_id 作为学号字符串）
    studentNo := req.StudentNo
    if studentNo == "" {
        studentNo = req.StudentID
    }
    if studentNo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": "student_no 不能为空"})
        return
    }

    var student models.Student
    if err := db.Where("student_id = ?", studentNo).First(&student).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "学生不存在"})
        return
    }

    parent := models.Parent{
        StudentID: student.ID, // 以数据库主键ID作为外键
        Name:      req.Name,
        Phone:     req.Phone,
        Relation:  req.Relation,
    }

    if err := db.Create(&parent).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": parent})
}

type UpdateParentRequest struct {
    Name     string `json:"name"`
    Phone    string `json:"phone"`
    Relation string `json:"relation"`
}

// UpdateParent 更新家长联系方式
func UpdateParent(c *gin.Context) {
    id := c.Param("id")
    var req UpdateParentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
        return
    }

    db := config.GetDB()
    var parent models.Parent
    if err := db.First(&parent, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "家长记录不存在"})
        return
    }

    if req.Name != "" {
        parent.Name = req.Name
    }
    if req.Phone != "" {
        parent.Phone = req.Phone
    }
    if req.Relation != "" {
        parent.Relation = req.Relation
    }

    if err := db.Save(&parent).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": parent})
}

// DeleteParent 删除家长联系方式
func DeleteParent(c *gin.Context) {
    id := c.Param("id")
    db := config.GetDB()

    if err := db.Delete(&models.Parent{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}


