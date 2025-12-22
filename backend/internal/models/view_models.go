package models

// ClassPerformanceView 班级成绩统计视图模型
// 对应数据库视图: vw_class_performance
// 用途: 供教师端仪表盘使用，自动计算平均分、最高分、最低分、及格率
type ClassPerformanceView struct {
	ClassID      uint    `json:"class_id" gorm:"column:class_id"`
	ClassName    string  `json:"class_name" gorm:"column:class_name"`
	CourseName   string  `json:"course_name" gorm:"column:course_name"`
	TeacherName  string  `json:"teacher_name" gorm:"column:teacher_name"`
	StudentCount int     `json:"student_count" gorm:"column:student_count"`
	AvgScore     float64 `json:"avg_score" gorm:"column:avg_score"`
	MaxScore     float64 `json:"max_score" gorm:"column:max_score"`
	MinScore     float64 `json:"min_score" gorm:"column:min_score"`
	PassRate     string  `json:"pass_rate" gorm:"column:pass_rate"`
}

// TableName 设定表名为视图名
func (ClassPerformanceView) TableName() string {
	return "vw_class_performance"
}

// StudentFullProfileView 学生完整档案视图模型
// 对应数据库视图: vw_student_full_profile
// 用途: 将分散在 User, Student, Class, Parent 表的信息聚合，便于导出和查询
type StudentFullProfileView struct {
	StudentID    uint   `json:"student_id" gorm:"column:student_id"`
	Code         string `json:"code" gorm:"column:code"`                         // 学号
	StudentName  string `json:"student_name" gorm:"column:student_name"`
	Gender       string `json:"gender" gorm:"column:gender"`
	Email        string `json:"email" gorm:"column:email"`
	Phone        string `json:"phone" gorm:"column:phone"`
	Address      string `json:"address" gorm:"column:address"`
	ClassName    string `json:"class_name" gorm:"column:class_name"`
	LoginAccount string `json:"login_account" gorm:"column:login_account"`
	ParentsInfo  string `json:"parents_info" gorm:"column:parents_info"`         // 家长信息（聚合）
}

// TableName 设定表名为视图名
func (StudentFullProfileView) TableName() string {
	return "vw_student_full_profile"
}

// ==================== 使用示例 ====================
//
// 1. 获取某位老师的班级统计
//    var stats []ClassPerformanceView
//    db.Where("teacher_name = ?", teacherName).Find(&stats)
//
// 2. 获取某个班级的课程成绩统计
//    var stats []ClassPerformanceView
//    db.Where("class_name = ?", className).Find(&stats)
//
// 3. 获取学生完整档案
//    var profile StudentFullProfileView
//    db.Where("student_id = ?", studentID).First(&profile)
//
// 4. 获取所有学生完整档案列表（分页）
//    var profiles []StudentFullProfileView
//    db.Limit(10).Offset(0).Find(&profiles)
//
// 5. 按学号搜索学生档案
//    var profile StudentFullProfileView
//    db.Where("code = ?", studentCode).First(&profile)
//
