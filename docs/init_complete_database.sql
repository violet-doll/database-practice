-- ============================================
-- 学生管理系统 - 完整数据库初始化脚本
-- ============================================
-- 用途：从空数据库完整构建到当前数据库状态
-- 包含：数据库创建、所有表结构、索引、视图、触发器、存储过程
-- 
-- 使用方法：
-- mysql -u root -p < docs/init_complete_database.sql
-- ============================================

-- ============================================
-- 第一部分：数据库创建
-- ============================================

-- 1. 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS student_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE student_db;

-- ============================================
-- 第二部分：基础表结构创建
-- ============================================
-- 注意：这些表通常由 GORM 自动创建，但为了完整性，这里手动定义

-- 2. 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    role_name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    KEY idx_roles_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 3. 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    permission VARCHAR(100) NOT NULL UNIQUE COMMENT '权限标识',
    `group` VARCHAR(50) COMMENT '权限分组',
    KEY idx_permissions_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 4. 角色-权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    UNIQUE KEY idx_role_perm (role_id, permission_id),
    KEY idx_role_permissions_deleted_at (deleted_at),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 5. 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    user_id BIGINT UNSIGNED COMMENT '关联的实体ID（学生/教师/家长）',
    user_type VARCHAR(20) COMMENT '用户类型：student/teacher/admin/parent',
    KEY idx_users_deleted_at (deleted_at),
    FOREIGN KEY (role_id) REFERENCES roles(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 6. 教师表
CREATE TABLE IF NOT EXISTS teachers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    name VARCHAR(100) NOT NULL COMMENT '教师姓名',
    teacher_id VARCHAR(50) NOT NULL UNIQUE COMMENT '教师工号',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '电话',
    user_id BIGINT UNSIGNED COMMENT '关联用户ID',
    KEY idx_teachers_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师表';

-- 7. 班级表
CREATE TABLE IF NOT EXISTS classes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    class_name VARCHAR(100) NOT NULL COMMENT '班级名称',
    teacher_id BIGINT UNSIGNED COMMENT '班主任ID',
    KEY idx_classes_deleted_at (deleted_at),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='班级表';

-- 8. 学生表
CREATE TABLE IF NOT EXISTS students (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    name VARCHAR(100) NOT NULL COMMENT '学生姓名',
    student_id VARCHAR(50) NOT NULL UNIQUE COMMENT '学号',
    gender VARCHAR(10) COMMENT '性别',
    age INT COMMENT '年龄',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '电话',
    address VARCHAR(255) COMMENT '地址',
    class_id BIGINT UNSIGNED COMMENT '班级ID',
    user_id BIGINT UNSIGNED COMMENT '关联用户ID',
    KEY idx_students_deleted_at (deleted_at),
    FOREIGN KEY (class_id) REFERENCES classes(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生表';

-- 9. 家长表
CREATE TABLE IF NOT EXISTS parents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    name VARCHAR(100) COMMENT '家长姓名',
    phone VARCHAR(20) NOT NULL COMMENT '电话',
    relation VARCHAR(20) COMMENT '关系（父亲/母亲等）',
    user_id BIGINT UNSIGNED COMMENT '关联用户ID',
    KEY idx_parents_deleted_at (deleted_at),
    KEY idx_parents_student_id (student_id),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家长表';

-- 10. 课程表
CREATE TABLE IF NOT EXISTS courses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    course_name VARCHAR(100) NOT NULL COMMENT '课程名称',
    teacher_id BIGINT UNSIGNED COMMENT '授课教师ID',
    credits DECIMAL(3,1) COMMENT '学分',
    capacity INT DEFAULT 50 COMMENT '课程最大容量',
    enrolled_count INT DEFAULT 0 COMMENT '当前已选人数',
    KEY idx_courses_deleted_at (deleted_at),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- 11. 选课表
CREATE TABLE IF NOT EXISTS enrollments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    course_id BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
    UNIQUE KEY idx_student_course (student_id, course_id),
    KEY idx_enrollments_deleted_at (deleted_at),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='选课表';

-- 12. 成绩表
CREATE TABLE IF NOT EXISTS grades (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    enrollment_id BIGINT UNSIGNED NOT NULL COMMENT '选课记录ID',
    score_type VARCHAR(50) COMMENT '成绩类型（平时/期末/总评）',
    score DECIMAL(5,2) COMMENT '分数',
    KEY idx_grades_deleted_at (deleted_at),
    KEY idx_grades_enrollment_id (enrollment_id),
    FOREIGN KEY (enrollment_id) REFERENCES enrollments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成绩表';

-- 13. 考勤表
CREATE TABLE IF NOT EXISTS attendances (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    date VARCHAR(20) COMMENT '日期',
    status VARCHAR(20) COMMENT '状态（出勤/缺席/请假/迟到）',
    reason TEXT COMMENT '原因或备注',
    teacher_id BIGINT UNSIGNED COMMENT '记录人',
    KEY idx_attendances_deleted_at (deleted_at),
    KEY idx_attendances_student_id (student_id),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='考勤表';

-- 14. 奖惩表
CREATE TABLE IF NOT EXISTS reward_punishments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    type VARCHAR(20) COMMENT '类型（奖励/处分）',
    description TEXT COMMENT '描述',
    date VARCHAR(20) COMMENT '日期',
    issuer VARCHAR(100) COMMENT '发布人',
    KEY idx_reward_punishments_deleted_at (deleted_at),
    KEY idx_reward_punishments_student_id (student_id),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='奖惩表';

-- 15. 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    title VARCHAR(200) COMMENT '标题',
    content TEXT COMMENT '内容',
    sender_id BIGINT UNSIGNED COMMENT '发送人ID',
    target VARCHAR(100) COMMENT '目标（all/class:5等）',
    KEY idx_notifications_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';

-- 16. 课程表（排课）
CREATE TABLE IF NOT EXISTS schedules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    course_id BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
    class_id BIGINT UNSIGNED NOT NULL COMMENT '班级ID',
    teacher_id BIGINT UNSIGNED NOT NULL COMMENT '教师ID',
    day_of_week INT COMMENT '星期几（1-7）',
    start_time VARCHAR(20) COMMENT '开始时间',
    end_time VARCHAR(20) COMMENT '结束时间',
    location VARCHAR(100) COMMENT '上课地点',
    semester VARCHAR(50) COMMENT '学期',
    KEY idx_schedules_deleted_at (deleted_at),
    KEY idx_schedules_course_id (course_id),
    KEY idx_schedules_class_id (class_id),
    KEY idx_schedules_teacher_id (teacher_id),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表（排课）';

-- ============================================
-- 第三部分：高级功能表
-- ============================================

-- 17. 课程先修关系表（多对多自引用）
CREATE TABLE IF NOT EXISTS course_prerequisites (
    course_id BIGINT UNSIGNED NOT NULL COMMENT '当前课程ID',
    prereq_id BIGINT UNSIGNED NOT NULL COMMENT '先修课程ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (course_id, prereq_id),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (prereq_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程先修关系表';

-- 为先修关系表添加索引以提高查询性能
CREATE INDEX idx_prereq_course ON course_prerequisites(prereq_id);

-- 18. 成绩修改审计日志表
CREATE TABLE IF NOT EXISTS grade_audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    grade_id BIGINT UNSIGNED NOT NULL COMMENT '被修改的成绩记录ID',
    old_score DECIMAL(5, 2) COMMENT '修改前分数',
    new_score DECIMAL(5, 2) COMMENT '修改后分数',
    created_at DATETIME(3) NULL DEFAULT NULL,
    updated_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    updated_at_audit DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    FOREIGN KEY (grade_id) REFERENCES grades(id) ON DELETE CASCADE,
    KEY idx_grade_audit_logs_deleted_at (deleted_at),
    KEY idx_grade_id (grade_id),
    KEY idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成绩修改审计日志表';

-- ============================================
-- 第四部分：创建触发器
-- ============================================

-- 19. 创建成绩修改触发器（数据看门狗）
-- 功能：在更新 grades 表之前，自动记录旧值
-- 亮点：完全由数据库自动完成，后端代码无需额外逻辑

-- 先删除已存在的触发器（如果存在）
DROP TRIGGER IF EXISTS trg_audit_grade_update;

DELIMITER //

CREATE TRIGGER trg_audit_grade_update
BEFORE UPDATE ON grades
FOR EACH ROW
BEGIN
    -- 只有当分数确实发生变化时才记录
    IF OLD.score != NEW.score THEN
        INSERT INTO grade_audit_logs (grade_id, old_score, new_score, created_at, updated_at, updated_at_audit)
        VALUES (OLD.id, OLD.score, NEW.score, NOW(3), NOW(3), NOW());
    END IF;
END //

DELIMITER ;

-- ============================================
-- 第五部分：创建存储过程
-- ============================================

-- 20. 创建选课存储过程
-- 功能：智能选课，检查先修课程、课程容量等
-- 优势：事务安全、业务逻辑下沉数据库层

-- 先删除已存在的存储过程（如果存在）
DROP PROCEDURE IF EXISTS sp_enroll_student;

DELIMITER //

CREATE PROCEDURE sp_enroll_student(
    IN p_student_id BIGINT UNSIGNED,
    IN p_course_id BIGINT UNSIGNED,
    OUT p_status INT,       -- 0: 成功, 1: 失败
    OUT p_message VARCHAR(255)
)
BEGIN
    DECLARE v_capacity INT;
    DECLARE v_enrolled INT;
    DECLARE v_already_enrolled INT;
    DECLARE v_prereq_count INT;
    DECLARE v_prereq_met INT;

    -- 开始事务
    START TRANSACTION;

    -- 1. 检查是否已经选课
    SELECT COUNT(*) INTO v_already_enrolled 
    FROM enrollments 
    WHERE student_id = p_student_id AND course_id = p_course_id AND deleted_at IS NULL;

    IF v_already_enrolled > 0 THEN
        SET p_status = 1;
        SET p_message = '已经选过该课程';
        ROLLBACK;
    ELSE
        -- 2. 检查先修课程要求
        -- 统计该课程有多少先修课程
        SELECT COUNT(*) INTO v_prereq_count 
        FROM course_prerequisites 
        WHERE course_id = p_course_id;

        -- 统计学生已经完成且及格的先修课程数量（分数 >= 60）
        SELECT COUNT(*) INTO v_prereq_met
        FROM course_prerequisites cp
        JOIN enrollments e ON cp.prereq_id = e.course_id
        JOIN grades g ON e.id = g.enrollment_id
        WHERE cp.course_id = p_course_id 
          AND e.student_id = p_student_id
          AND g.score >= 60
          AND e.deleted_at IS NULL; 

        IF v_prereq_met < v_prereq_count THEN
            SET p_status = 1;
            SET p_message = '未完成先修课程要求';
            ROLLBACK;
        ELSE
            -- 3. 检查并锁定课程容量
            SELECT capacity, enrolled_count INTO v_capacity, v_enrolled
            FROM courses
            WHERE id = p_course_id
            FOR UPDATE;  -- 行锁，防止并发选课超额

            IF v_enrolled >= v_capacity THEN
                SET p_status = 1;
                SET p_message = '课程已满';
                ROLLBACK;
            ELSE
                -- 4. 执行选课操作
                INSERT INTO enrollments (created_at, updated_at, student_id, course_id)
                VALUES (NOW(3), NOW(3), p_student_id, p_course_id);

                -- 更新已选人数
                UPDATE courses 
                SET enrolled_count = enrolled_count + 1 
                WHERE id = p_course_id;

                SET p_status = 0;
                SET p_message = '选课成功';
                COMMIT;
            END IF;
        END IF;
    END IF;
END //

DELIMITER ;

-- ============================================
-- 第六部分：创建智能报表视图
-- ============================================

-- 21. 创建班级成绩统计视图（供教师端仪表盘使用）
-- 功能：自动计算平均分、最高分、最低分、及格率
-- 优势：简化前端查询逻辑，提供数据逻辑独立性

DROP VIEW IF EXISTS vw_class_performance;

CREATE VIEW vw_class_performance AS
SELECT 
    c.id AS class_id,
    c.class_name,
    co.course_name,
    t.name AS teacher_name,
    COUNT(g.score) AS student_count,
    ROUND(AVG(g.score), 2) AS avg_score,
    MAX(g.score) AS max_score,
    MIN(g.score) AS min_score,
    -- 计算及格率 (假设 >= 60 及格)
    CONCAT(ROUND(SUM(CASE WHEN g.score >= 60 THEN 1 ELSE 0 END) / COUNT(g.score) * 100, 1), '%') AS pass_rate
FROM classes c
JOIN students s ON c.id = s.class_id
JOIN enrollments e ON s.id = e.student_id
JOIN courses co ON e.course_id = co.id
JOIN teachers t ON co.teacher_id = t.id
JOIN grades g ON e.id = g.enrollment_id
WHERE e.deleted_at IS NULL AND g.deleted_at IS NULL
GROUP BY c.id, c.class_name, co.course_name, t.name;

-- 22. 创建学生完整档案视图（供学生端/导出使用）
-- 功能：将分散在 User, Student, Class, Parent 表的信息聚合
-- 优势：一次查询获取完整学生信息

DROP VIEW IF EXISTS vw_student_full_profile;

CREATE VIEW vw_student_full_profile AS
SELECT 
    s.id AS student_id,
    s.student_id AS code,   -- 学号
    s.name AS student_name,
    s.gender,
    s.email,
    s.phone,
    s.address,
    c.class_name,
    u.username AS login_account,
    -- 聚合家长信息 (如果有多个家长，用逗号连接)
    GROUP_CONCAT(CONCAT(p.relation, ':', p.name, '(', p.phone, ')') SEPARATOR '; ') AS parents_info
FROM students s
LEFT JOIN classes c ON s.class_id = c.id
LEFT JOIN users u ON s.user_id = u.id
LEFT JOIN parents p ON s.id = p.student_id AND p.deleted_at IS NULL
WHERE s.deleted_at IS NULL
GROUP BY s.id, s.student_id, s.name, s.gender, s.email, s.phone, s.address, c.class_name, u.username;

-- ============================================
-- 第七部分：初始数据插入
-- ============================================

-- 23. 插入默认角色
INSERT INTO roles (role_name, created_at, updated_at) VALUES
('admin', NOW(3), NOW(3)),
('teacher', NOW(3), NOW(3)),
('student', NOW(3), NOW(3)),
('parent', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- ============================================
-- 验证和使用说明
-- ============================================

-- 测试触发器：
-- UPDATE grades SET score = 95.0 WHERE id = 1;
-- SELECT * FROM grade_audit_logs WHERE grade_id = 1;

-- 测试存储过程：
-- CALL sp_enroll_student(1, 2, @status, @message);
-- SELECT @status AS status, @message AS message;

-- 查询班级成绩统计：
-- SELECT * FROM vw_class_performance WHERE class_name = '计算机1班';

-- 查询学生完整档案：
-- SELECT * FROM vw_student_full_profile WHERE code = 'S2021001';

-- 查看所有视图：
-- SHOW FULL TABLES WHERE Table_type = 'VIEW';

-- 查看所有触发器：
-- SHOW TRIGGERS;

-- 查看所有存储过程：
-- SHOW PROCEDURE STATUS WHERE Db = 'student_db';

-- ============================================
-- 脚本执行完成
-- ============================================
