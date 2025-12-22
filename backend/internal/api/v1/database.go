package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"student-management-system/config"

	"github.com/gin-gonic/gin"
)

// TableInfo 表信息
type TableInfo struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	Count  int64  `json:"count"`
	IsView bool   `json:"isView,omitempty"`
}

// GetTableList 获取所有表的列表和统计信息
func GetTableList(c *gin.Context) {
	db := config.DB

	// 定义所有表的配置
	tables := []TableInfo{
		{Name: "users", Label: "用户表"},
		{Name: "roles", Label: "角色表"},
		{Name: "permissions", Label: "权限表"},
		{Name: "students", Label: "学生表"},
		{Name: "teachers", Label: "教师表"},
		{Name: "parents", Label: "家长表"},
		{Name: "classes", Label: "班级表"},
		{Name: "courses", Label: "课程表"},
		{Name: "enrollments", Label: "选课表"},
		{Name: "grades", Label: "成绩表"},
		{Name: "attendances", Label: "考勤表"},
		{Name: "reward_punishments", Label: "奖惩表"},
		{Name: "notifications", Label: "通知表"},
		{Name: "schedules", Label: "课程表(排课)"},
		{Name: "grade_audit_logs", Label: "成绩审计日志"},
		{Name: "vw_class_performance", Label: "班级成绩视图", IsView: true},
		{Name: "vw_student_full_profile", Label: "学生档案视图", IsView: true},
	}

	// 获取每个表的记录数
	for i := range tables {
		if !tables[i].IsView {
			var count int64
			db.Table(tables[i].Name).Count(&count)
			tables[i].Count = count
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"tables": tables,
		},
	})
}

// GetTableData 获取指定表的数据
func GetTableData(c *gin.Context) {
	tableName := c.Param("table")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := config.DB
	offset := (page - 1) * pageSize

	// 验证表名是否合法（防止SQL注入）
	validTables := map[string]bool{
		"users": true, "roles": true, "permissions": true,
		"students": true, "teachers": true, "parents": true,
		"classes": true, "courses": true, "enrollments": true,
		"grades": true, "attendances": true, "reward_punishments": true,
		"notifications": true, "schedules": true, "grade_audit_logs": true,
		"vw_class_performance": true, "vw_student_full_profile": true,
	}

	if !validTables[tableName] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid table name",
		})
		return
	}

	// 获取总数
	var total int64
	db.Table(tableName).Count(&total)

	// 获取数据
	var results []map[string]interface{}
	db.Table(tableName).Limit(pageSize).Offset(offset).Find(&results)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      results,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateTableData 创建表数据
func CreateTableData(c *gin.Context) {
	tableName := c.Param("table")

	// 验证表名
	validTables := map[string]bool{
		"users": true, "roles": true, "permissions": true,
		"students": true, "teachers": true, "parents": true,
		"classes": true, "courses": true, "enrollments": true,
		"grades": true, "attendances": true, "reward_punishments": true,
		"notifications": true, "schedules": true,
	}

	if !validTables[tableName] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid table name",
		})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	db := config.DB
	if err := db.Table(tableName).Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    data,
	})
}

// UpdateTableData 更新表数据
func UpdateTableData(c *gin.Context) {
	tableName := c.Param("table")
	id := c.Param("id")

	// 验证ID是否为有效数字
	if id == "" || id == "undefined" || id == "null" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	// 尝试将ID转换为整数以验证其有效性
	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID必须是有效的数字",
		})
		return
	}

	// 验证表名
	validTables := map[string]bool{
		"users": true, "roles": true, "permissions": true,
		"students": true, "teachers": true, "parents": true,
		"classes": true, "courses": true, "enrollments": true,
		"grades": true, "attendances": true, "reward_punishments": true,
		"notifications": true, "schedules": true,
	}

	if !validTables[tableName] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid table name",
		})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	db := config.DB
	if err := db.Table(tableName).Where("id = ?", id).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteTableData 删除表数据
func DeleteTableData(c *gin.Context) {
	tableName := c.Param("table")
	id := c.Param("id")

	// 验证ID是否为有效数字
	if id == "" || id == "undefined" || id == "null" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	// 尝试将ID转换为整数以验证其有效性
	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID必须是有效的数字",
		})
		return
	}

	// 验证表名
	validTables := map[string]bool{
		"users": true, "roles": true, "permissions": true,
		"students": true, "teachers": true, "parents": true,
		"classes": true, "courses": true, "enrollments": true,
		"grades": true, "attendances": true, "reward_punishments": true,
		"notifications": true, "schedules": true,
	}

	if !validTables[tableName] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid table name",
		})
		return
	}

	db := config.DB

	// 使用Exec执行删除操作
	result := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": result.Error.Error(),
		})
		return
	}

	// 检查是否有记录被删除
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ExportTableData 导出表数据
func ExportTableData(c *gin.Context) {
	tableName := c.Param("table")

	// 验证表名
	validTables := map[string]bool{
		"users": true, "roles": true, "permissions": true,
		"students": true, "teachers": true, "parents": true,
		"classes": true, "courses": true, "enrollments": true,
		"grades": true, "attendances": true, "reward_punishments": true,
		"notifications": true, "schedules": true, "grade_audit_logs": true,
		"vw_class_performance": true, "vw_student_full_profile": true,
	}

	if !validTables[tableName] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid table name",
		})
		return
	}

	db := config.DB
	var results []map[string]interface{}
	db.Table(tableName).Find(&results)

	// 简单返回 JSON 格式（实际应用中应该生成 Excel 文件）
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", tableName))
	c.JSON(http.StatusOK, results)
}

// ExecuteSQL 执行 SQL 查询（仅用于开发环境，生产环境应禁用）
func ExecuteSQL(c *gin.Context) {
	var req struct {
		SQL string `json:"sql" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	db := config.DB
	var results []map[string]interface{}

	if err := db.Raw(req.SQL).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    results,
	})
}

// GetTableSchema 获取表结构信息
func GetTableSchema(c *gin.Context) {
	tableName := c.Param("table")

	// 验证表名
	validTables := map[string]bool{
		"users": true, "roles": true, "permissions": true,
		"students": true, "teachers": true, "parents": true,
		"classes": true, "courses": true, "enrollments": true,
		"grades": true, "attendances": true, "reward_punishments": true,
		"notifications": true, "schedules": true, "grade_audit_logs": true,
		"vw_class_performance": true, "vw_student_full_profile": true,
	}

	if !validTables[tableName] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid table name",
		})
		return
	}

	db := config.DB
	var columns []map[string]interface{}

	// 获取表结构（MySQL）
	query := fmt.Sprintf("DESCRIBE %s", tableName)
	if err := db.Raw(query).Scan(&columns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"columns": columns,
		},
	})
}
