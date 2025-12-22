package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"` // 存储哈希后的密码
	RoleID   uint   `json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role"` // Gorm 自动关联
	IsActive bool   `gorm:"default:true" json:"is_active"`
	UserID   uint   `json:"user_id"`                           // 关联学生、教师或家长 (多态关联)
	UserType string `gorm:"type:varchar(20)" json:"user_type"` // "student", "teacher", "admin", "parent"
}

// 2. 角色表 (RBAC)
type Role struct {
	gorm.Model
	RoleName    string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"role_name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"` // 角色拥有的权限
}

// 权限表 (Permission)
type Permission struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100);not null" json:"name"`                   // 权限名称 (e.g., "创建学生")
	Permission string `gorm:"type:varchar(100);uniqueIndex;not null" json:"permission"` // 权限标识 (e.g., "student:create")
	Group      string `gorm:"type:varchar(50)" json:"group"`                            // 分组 (e.g., "student")
}

// 角色-权限 关联表 (RolePermission)
type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"index:idx_role_perm,unique" json:"role_id"`
	PermissionID uint `gorm:"index:idx_role_perm,unique" json:"permission_id"`
}

// 3. 学生表 (对应要求 1)
type Student struct {
	gorm.Model
	Name      string   `gorm:"type:varchar(100);not null" json:"name"`
	StudentID string   `gorm:"type:varchar(50);uniqueIndex;not null" json:"student_id"` // 学号
	Gender    string   `gorm:"type:varchar(10)" json:"gender"`                          // "男", "女"
	Age       int      `json:"age"`
	Email     string   `gorm:"type:varchar(100)" json:"email"`
	Phone     string   `gorm:"type:varchar(20)" json:"phone"`
	Address   string   `gorm:"type:varchar(255)" json:"address"`
	ClassID   uint     `json:"class_id"` // 关联班级
	Class     Class    `gorm:"foreignKey:ClassID" json:"class"`
	UserID    uint     `json:"user_id"`                             // 关联登录用户 User (如果学生可以登录)
	Parents   []Parent `gorm:"foreignKey:StudentID" json:"parents"` // 一对多关联家长
}

// 4. 家长表 (对应要求 1, 5)
type Parent struct {
	gorm.Model
	StudentID uint    `gorm:"index" json:"student_id"`                       // 关联学生
	Student   Student `gorm:"foreignKey:StudentID" json:"student,omitempty"` // 关联学生信息
	Name      string  `gorm:"type:varchar(100)" json:"name"`
	Phone     string  `gorm:"type:varchar(20);not null" json:"phone"`
	Relation  string  `gorm:"type:varchar(20)" json:"relation"` // e.g., "父亲", "母亲"
	UserID    uint    `json:"user_id"`                          // 关联登录用户 User (如果家长可以登录)
}

// 5. 教师表
type Teacher struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	TeacherID string `gorm:"type:varchar(50);uniqueIndex;not null" json:"teacher_id"` // 教师工号
	Email     string `gorm:"type:varchar(100)" json:"email"`
	Phone     string `gorm:"type:varchar(20)" json:"phone"`
	UserID    uint   `json:"user_id"` // 关联登录用户 User
}

// 6. 班级表
type Class struct {
	gorm.Model
	ClassName string    `gorm:"not null" json:"class_name"`
	TeacherID uint      `json:"teacher_id"` // 班主任 (关联 Teacher)
	Teacher   Teacher   `gorm:"foreignKey:TeacherID" json:"teacher"`
	Students  []Student `gorm:"foreignKey:ClassID" json:"students,omitempty"` // 一个班级有多个学生
}

// 7. 课程表 (对应要求 2)
type Course struct {
	gorm.Model
	CourseName    string  `gorm:"not null" json:"course_name"`
	TeacherID     uint    `json:"teacher_id"` // 授课教师 (关联 Teacher)
	Teacher       Teacher `gorm:"foreignKey:TeacherID" json:"teacher"`
	Credits       float64 `json:"credits"`                         // 学分
	Capacity      int     `gorm:"default:50" json:"capacity"`      // 课程容量
	EnrolledCount int     `gorm:"default:0" json:"enrolled_count"` // 已选人数
}

// 8. 选课表 (学生和课程的中间表) (对应要求 2)
type Enrollment struct {
	gorm.Model
	StudentID uint    `gorm:"index:idx_student_course,unique" json:"student_id"`
	Student   Student `gorm:"foreignKey:StudentID" json:"Student"` // 明确指定JSON字段名
	CourseID  uint    `gorm:"index:idx_student_course,unique" json:"course_id"`
	Course    Course  `gorm:"foreignKey:CourseID" json:"Course"`               // 明确指定JSON字段名
	Grades    []Grade `gorm:"foreignKey:EnrollmentID" json:"Grades,omitempty"` // 明确指定JSON字段名
}

// 课程先修关系表 (对应课程依赖关系)
type CoursePrerequisite struct {
	CourseID  uint   `gorm:"primaryKey" json:"course_id"` // 当前课程ID
	PrereqID  uint   `gorm:"primaryKey" json:"prereq_id"` // 先修课程ID
	CreatedAt string `json:"created_at"`
	Course    Course `gorm:"foreignKey:CourseID" json:"course,omitempty"` // 当前课程
	Prereq    Course `gorm:"foreignKey:PrereqID" json:"prereq,omitempty"` // 先修课程
}

// 9. 成绩表 (对应要求 2)
type Grade struct {
	gorm.Model
	EnrollmentID uint    `gorm:"index" json:"enrollment_id"` // 关联选课记录
	ScoreType    string  `json:"score_type"`                 // e.g., "平时成绩", "期末成绩", "总评"
	Score        float64 `json:"score"`
}

// 10. 考勤表 (对应要求 3)
type Attendance struct {
	gorm.Model
	StudentID uint    `gorm:"index" json:"student_id"` // 关联学生
	Student   Student `gorm:"foreignKey:StudentID" json:"student"`
	Date      string  `json:"date"`       // 日期 (e.g., "2023-10-27")
	Status    string  `json:"status"`     // "出勤", "缺席", "请假", "迟到"
	Reason    string  `json:"reason"`     // (可选) 备注或请假理由
	TeacherID uint    `json:"teacher_id"` // (可选) 记录人
}

// 11. 奖惩表 (对应要求 4)
type RewardPunishment struct {
	gorm.Model
	StudentID   uint    `gorm:"index" json:"student_id"` // 关联学生
	Student     Student `gorm:"foreignKey:StudentID" json:"student"`
	Type        string  `json:"type"`        // "奖励", "处分"
	Description string  `json:"description"` // 事由
	Date        string  `json:"date"`        // 日期
	Issuer      string  `json:"issuer"`      // (可选) 发布人/机构
}

// 12. 通知表 (对应要求 5)
type Notification struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	SenderID uint   `json:"sender_id"` // (可选) 发送人 (关联 User)
	Target   string `json:"target"`    // (可选) 发送目标 e.g., "all", "class:5"
}

// 13. 课程表 (排课)
type Schedule struct {
	gorm.Model
	CourseID  uint    `gorm:"index" json:"course_id"` // 关联课程
	Course    Course  `gorm:"foreignKey:CourseID" json:"course"`
	ClassID   uint    `gorm:"index" json:"class_id"` // 关联班级
	Class     Class   `gorm:"foreignKey:ClassID" json:"class"`
	TeacherID uint    `gorm:"index" json:"teacher_id"` // 关联教师 (可从Course获取, 但显式存储更灵活)
	Teacher   Teacher `gorm:"foreignKey:TeacherID" json:"teacher"`
	DayOfWeek int     `json:"day_of_week"` // 星期几 (例如 1=周一, 7=周日)
	StartTime string  `json:"start_time"`  // 节次或时间 (e.g., "08:00" 或 "1-2节")
	EndTime   string  `json:"end_time"`    // (e.g., "09:40")
	Location  string  `json:"location"`    // 上课地点 (e.g., "教5-101")
	Semester  string  `json:"semester"`    // (可选) 学期 (e.g., "2025-Fall")
}

// 14. 成绩审计日志表 (对应修改要求 - 数据看门狗)
// 说明：该表配合触发器使用，自动记录成绩修改历史
// 触发器在 docs/update_schema.sql 中定义
type GradeAuditLog struct {
	gorm.Model
	GradeID  uint    `gorm:"index;not null" json:"grade_id"`            // 被修改的成绩记录ID
	OldScore float64 `gorm:"type:decimal(5,2)" json:"old_score"`        // 修改前分数
	NewScore float64 `gorm:"type:decimal(5,2)" json:"new_score"`        // 修改后分数
	Grade    Grade   `gorm:"foreignKey:GradeID" json:"grade,omitempty"` // 关联成绩记录
}
