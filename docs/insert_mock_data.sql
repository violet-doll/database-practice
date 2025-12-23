-- ============================================
-- Student Management System - Mock Data Insert Script
-- ============================================
-- Purpose: Insert test data into all database tables
-- Usage: mysql -u root -p student_db < docs/insert_mock_data.sql
-- ============================================

SET NAMES utf8mb4;
SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;

USE student_db;

-- ============================================
-- 清空现有数据（可选,谨慎使用）
-- ============================================
-- 禁用外键检查,避免插入顺序问题
SET FOREIGN_KEY_CHECKS = 0;

-- 清空所有表数据
TRUNCATE TABLE grade_audit_logs;
TRUNCATE TABLE grades;
TRUNCATE TABLE enrollments;
TRUNCATE TABLE attendances;
TRUNCATE TABLE reward_punishments;
TRUNCATE TABLE notifications;
TRUNCATE TABLE schedules;
TRUNCATE TABLE course_prerequisites;
TRUNCATE TABLE parents;
TRUNCATE TABLE students;
TRUNCATE TABLE courses;
TRUNCATE TABLE classes;
TRUNCATE TABLE teachers;
TRUNCATE TABLE users;
TRUNCATE TABLE role_permissions;
TRUNCATE TABLE permissions;
TRUNCATE TABLE roles;

-- 注意: 外键检查将在脚本末尾重新启用

-- ============================================
-- 1. 插入角色数据
-- ============================================
INSERT INTO roles (role_name, created_at, updated_at) VALUES
('admin', NOW(3), NOW(3)),
('teacher', NOW(3), NOW(3)),
('student', NOW(3), NOW(3)),
('parent', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- ============================================
-- 2. Insert permissions data
-- ============================================
INSERT INTO permissions (name, permission, `group`, created_at, updated_at) VALUES
-- User management (Admin)
('View User', 'admin:user:read', 'Admin', NOW(3), NOW(3)),
('Create User', 'admin:user:create', 'Admin', NOW(3), NOW(3)),
('Edit User', 'admin:user:update', 'Admin', NOW(3), NOW(3)),
('Delete User', 'admin:user:delete', 'Admin', NOW(3), NOW(3)),

-- Role management (Admin)
('View Role', 'admin:role:read', 'Admin', NOW(3), NOW(3)),
('Create Role', 'admin:role:create', 'Admin', NOW(3), NOW(3)),
('Edit Role', 'admin:role:update', 'Admin', NOW(3), NOW(3)),
('Delete Role', 'admin:role:delete', 'Admin', NOW(3), NOW(3)),

-- Stats (Admin)
('View Stats', 'admin:stats:read', 'Admin', NOW(3), NOW(3)),

-- Student management
('View Student', 'student:read', 'Student', NOW(3), NOW(3)),
('Create Student', 'student:create', 'Student', NOW(3), NOW(3)),
('Edit Student', 'student:update', 'Student', NOW(3), NOW(3)),
('Delete Student', 'student:delete', 'Student', NOW(3), NOW(3)),

-- Class management
('View Class', 'class:read', 'Class', NOW(3), NOW(3)),
('Create Class', 'class:create', 'Class', NOW(3), NOW(3)),
('Edit Class', 'class:update', 'Class', NOW(3), NOW(3)),
('Delete Class', 'class:delete', 'Class', NOW(3), NOW(3)),

-- Course management
('View Course', 'course:read', 'Course', NOW(3), NOW(3)),
('Create Course', 'course:create', 'Course', NOW(3), NOW(3)),
('Edit Course', 'course:update', 'Course', NOW(3), NOW(3)),
('Delete Course', 'course:delete', 'Course', NOW(3), NOW(3)),

-- Schedule management
('View Schedule', 'schedule:read', 'Schedule', NOW(3), NOW(3)),
('Create Schedule', 'schedule:create', 'Schedule', NOW(3), NOW(3)),
('Edit Schedule', 'schedule:update', 'Schedule', NOW(3), NOW(3)),
('Delete Schedule', 'schedule:delete', 'Schedule', NOW(3), NOW(3)),

-- Enrollment management
('View Enrollment', 'enrollment:read', 'Enrollment', NOW(3), NOW(3)),
('Create Enrollment', 'enrollment:create', 'Enrollment', NOW(3), NOW(3)),
('Delete Enrollment', 'enrollment:delete', 'Enrollment', NOW(3), NOW(3)),

-- Grade management
('View Grade', 'grade:read', 'Grade', NOW(3), NOW(3)),
('Create Grade', 'grade:create', 'Grade', NOW(3), NOW(3)),

-- Attendance management
('View Attendance', 'attendance:read', 'Attendance', NOW(3), NOW(3)),
('Create Attendance', 'attendance:create', 'Attendance', NOW(3), NOW(3)),
('Delete Attendance', 'attendance:delete', 'Attendance', NOW(3), NOW(3)),

-- Reward management
('View Reward', 'reward:read', 'Reward', NOW(3), NOW(3)),
('Create Reward', 'reward:create', 'Reward', NOW(3), NOW(3)),
('Delete Reward', 'reward:delete', 'Reward', NOW(3), NOW(3)),

-- Parent management
('View Parent', 'parent:read', 'Parent', NOW(3), NOW(3)),
('Create Parent', 'parent:create', 'Parent', NOW(3), NOW(3)),
('Edit Parent', 'parent:update', 'Parent', NOW(3), NOW(3)),
('Delete Parent', 'parent:delete', 'Parent', NOW(3), NOW(3)),

-- Notification management
('View Notification', 'notification:read', 'Notification', NOW(3), NOW(3)),
('Create Notification', 'notification:create', 'Notification', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- ============================================
-- 3. 插入角色-权限关联数据
-- ============================================
-- 管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE role_name = 'admin'),
    p.id,
    NOW(3),
    NOW(3)
FROM permissions p
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- 教师拥有部分权限
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE role_name = 'teacher'),
    p.id,
    NOW(3),
    NOW(3)
FROM permissions p
WHERE p.permission IN ('student:read', 'course:read', 'grade:read', 'grade:create', 'attendance:read', 'attendance:create', 'schedule:read')
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- 学生只能查看
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE role_name = 'student'),
    p.id,
    NOW(3),
    NOW(3)
FROM permissions p
WHERE p.permission IN ('course:read', 'grade:read', 'schedule:read', 'notification:read')
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- 4. Insert teachers data
-- ============================================
INSERT INTO teachers (created_at, updated_at, name, teacher_id, email, phone) VALUES
(NOW(3), NOW(3), 'Zhang Wei', 2021001, 'zhangwei@tju.edu.cn', '13800138001'),
(NOW(3), NOW(3), 'Li Na', 2021002, 'lina@tju.edu.cn', '13800138002'),
(NOW(3), NOW(3), 'Wang Qiang', 2021003, 'wangqiang@tju.edu.cn', '13800138003'),
(NOW(3), NOW(3), 'Liu Fang', 2021004, 'liufang@tju.edu.cn', '13800138004'),
(NOW(3), NOW(3), 'Chen Ming', 2021005, 'chenming@tju.edu.cn', '13800138005'),
(NOW(3), NOW(3), 'Zhao Li', 2021006, 'zhaoli@tju.edu.cn', '13800138006'),
(NOW(3), NOW(3), 'Sun Tao', 2021007, 'suntao@tju.edu.cn', '13800138007'),
(NOW(3), NOW(3), 'Zhou Min', 2021008, 'zhoumin@tju.edu.cn', '13800138008');

-- ============================================
-- 5. Insert classes data
-- ============================================
-- 注意: teacher_id 引用 teachers 表的自增 id 字段
INSERT INTO classes (class_name, teacher_id, created_at, updated_at) VALUES
('CS Class 1', 1, NOW(3), NOW(3)),
('CS Class 2', 2, NOW(3), NOW(3)),
('SE Class 1', 3, NOW(3), NOW(3)),
('SE Class 2', 4, NOW(3), NOW(3)),
('DS Class 1', 5, NOW(3), NOW(3)),
('AI Class 1', 6, NOW(3), NOW(3));

-- ============================================
-- 6. Insert users data (admin, teachers, students)
-- ============================================
-- Admin account (password: admin123)
INSERT INTO users (username, password, role_id, is_active, user_type, created_at, updated_at) VALUES
('admin', '$2a$10$QCeaQ2f9NrYkhlUCrilEWOW0UBOUZgZyoBa7DrsPt0vsSsLbbB3W2',  
 (SELECT id FROM roles WHERE role_name = 'admin'), TRUE, 'admin', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

-- Teacher accounts (password: teacher123)
INSERT INTO users (username, password, role_id, is_active, user_id, user_type, created_at, updated_at) VALUES
('teacher001', '$2a$10$8SuerLKQMwTas4N5fX9DJeIAJS9EaahTJkch2G82l4ehHW7mqIR8G', 
 (SELECT id FROM roles WHERE role_name = 'teacher'), TRUE, 1, 'teacher', NOW(3), NOW(3)),
('teacher002', '$2a$10$8SuerLKQMwTas4N5fX9DJeIAJS9EaahTJkch2G82l4ehHW7mqIR8G', 
 (SELECT id FROM roles WHERE role_name = 'teacher'), TRUE, 2, 'teacher', NOW(3), NOW(3)),
('teacher003', '$2a$10$8SuerLKQMwTas4N5fX9DJeIAJS9EaahTJkch2G82l4ehHW7mqIR8G',  
 (SELECT id FROM roles WHERE role_name = 'teacher'), TRUE, 3, 'teacher', NOW(3), NOW(3));

-- 7. Insert students data
-- ============================================
INSERT INTO students (created_at, updated_at, name, student_id, gender, age, email, phone, address, class_id) VALUES
-- CS Class 1
(NOW(3), NOW(3), 'Zhang San', 2021001, 'Male', 20, 'zhangsan@tju.edu.cn', '13900139001', 'Beijing Haidian', 1),
(NOW(3), NOW(3), 'Li Si', 2021002, 'Male', 21, 'lisi@tju.edu.cn', '13900139002', 'Beijing Chaoyang', 1),
(NOW(3), NOW(3), 'Wang Wu', 2021003, 'Female', 20, 'wangwu@tju.edu.cn', '13900139003', 'Beijing Xicheng', 1),
(NOW(3), NOW(3), 'Zhao Liu', 2021004, 'Male', 22, 'zhaoliu@tju.edu.cn', '13900139004', 'Beijing Dongcheng', 1),
(NOW(3), NOW(3), 'Sun Qi', 2021005, 'Female', 19, 'sunqi@tju.edu.cn', '13900139005', 'Beijing Fengtai', 1),

-- CS Class 2
(NOW(3), NOW(3), 'Zhou Ba', 2021006, 'Male', 20, 'zhouba@tju.edu.cn', '13900139006', 'Shanghai Pudong', 2),
(NOW(3), NOW(3), 'Wu Jiu', 2021007, 'Female', 21, 'wujiu@tju.edu.cn', '13900139007', 'Shanghai Huangpu', 2),
(NOW(3), NOW(3), 'Zheng Shi', 2021008, 'Male', 20, 'zhengshi@tju.edu.cn', '13900139008', 'Shanghai Jingan', 2),
(NOW(3), NOW(3), 'Qian Yi', 2021009, 'Female', 22, 'qianyi@tju.edu.cn', '13900139009', 'Shanghai Xuhui', 2),
(NOW(3), NOW(3), 'Sun Er', 2021010, 'Male', 19, 'suner@tju.edu.cn', '13900139010', 'Shanghai Changning', 2),

-- SE Class 1
(NOW(3), NOW(3), 'Li Ming', 2021011, 'Male', 20, 'liming@tju.edu.cn', '13900139011', 'Guangzhou Tianhe', 3),
(NOW(3), NOW(3), 'Chen Hong', 2021012, 'Female', 21, 'chenhong@tju.edu.cn', '13900139012', 'Guangzhou Yuexiu', 3),
(NOW(3), NOW(3), 'Lin Feng', 2021013, 'Male', 20, 'linfeng@tju.edu.cn', '13900139013', 'Guangzhou Haizhu', 3),
(NOW(3), NOW(3), 'Huang Li', 2021014, 'Female', 22, 'huangli@tju.edu.cn', '13900139014', 'Guangzhou Baiyun', 3),
(NOW(3), NOW(3), 'Yang Gang', 2021015, 'Male', 19, 'yanggang@tju.edu.cn', '13900139015', 'Guangzhou Panyu', 3),

-- SE Class 2
(NOW(3), NOW(3), 'Liu Yang', 2021016, 'Male', 20, 'liuyang@tju.edu.cn', '13900139016', 'Shenzhen Nanshan', 4),
(NOW(3), NOW(3), 'Zhang Min', 2021017, 'Female', 21, 'zhangmin@tju.edu.cn', '13900139017', 'Shenzhen Futian', 4),
(NOW(3), NOW(3), 'Wang Tao', 2021018, 'Male', 20, 'wangtao@tju.edu.cn', '13900139018', 'Shenzhen Luohu', 4),
(NOW(3), NOW(3), 'Zhao Jing', 2021019, 'Female', 22, 'zhaojing@tju.edu.cn', '13900139019', 'Shenzhen Baoan', 4),
(NOW(3), NOW(3), 'Sun Qiang', 2021020, 'Male', 19, 'sunqiang@tju.edu.cn', '13900139020', 'Shenzhen Longgang', 4),

-- DS Class 1
(NOW(3), NOW(3), 'Zhou Jie', 2021021, 'Male', 20, 'zhoujie@tju.edu.cn', '13900139021', 'Hangzhou Xihu', 5),
(NOW(3), NOW(3), 'Wu Na', 2021022, 'Female', 21, 'wuna@tju.edu.cn', '13900139022', 'Hangzhou Gongshu', 5),
(NOW(3), NOW(3), 'Zheng Wei', 2021023, 'Male', 20, 'zhengwei@tju.edu.cn', '13900139023', 'Hangzhou Jianggan', 5),
(NOW(3), NOW(3), 'Qian Fang', 2021024, 'Female', 22, 'qianfang@tju.edu.cn', '13900139024', 'Hangzhou Xiacheng', 5),
(NOW(3), NOW(3), 'Sun Tao2', 2021025, 'Male', 19, 'suntao2@tju.edu.cn', '13900139025', 'Hangzhou Binjiang', 5),

-- AI Class 1
(NOW(3), NOW(3), 'Li Hua', 2021026, 'Male', 20, 'lihua@tju.edu.cn', '13900139026', 'Nanjing Xuanwu', 6),
(NOW(3), NOW(3), 'Chen Jing', 2021027, 'Female', 21, 'chenjing@tju.edu.cn', '13900139027', 'Nanjing Qinhuai', 6),
(NOW(3), NOW(3), 'Lin Qiang', 2021028, 'Male', 20, 'linqiang@tju.edu.cn', '13900139028', 'Nanjing Gulou', 6),
(NOW(3), NOW(3), 'Huang Min', 2021029, 'Female', 22, 'huangmin@tju.edu.cn', '13900139029', 'Nanjing Jianye', 6),
(NOW(3), NOW(3), 'Yang Feng', 2021030, 'Male', 19, 'yangfeng@tju.edu.cn', '13900139030', 'Nanjing Qixia', 6);

-- Update teachers user_id
UPDATE teachers SET user_id = 1 WHERE id = 1;
UPDATE teachers SET user_id = 2 WHERE id = 2;
UPDATE teachers SET user_id = 3 WHERE id = 3;

-- ============================================
-- 8. 为学生创建用户账号
-- ============================================
-- 学生账号（密码：student123）
INSERT INTO users (username, password, role_id, is_active, user_id, user_type, created_at, updated_at)
SELECT 
    CONCAT('stu', LPAD(id, 3, '0')),
    '$2a$10$Z5zieIBr18cAIbAvjj9BFejm3dbIKlj9eF70Q3il15rjSwdkNT3Cu',
    (SELECT id FROM roles WHERE role_name = 'student'),
    TRUE,
    id,
    'student',
    NOW(3),
    NOW(3)
FROM students;

-- 更新学生表的user_id
UPDATE students s
JOIN users u ON u.user_id = s.id AND u.user_type = 'student'
SET s.user_id = u.id;

-- ============================================
-- 9. Insert parents data
-- ============================================
INSERT INTO parents (student_id, name, phone, relation, created_at, updated_at) VALUES
-- Add parents for first 10 students
(1, 'Zhang Daming', '13700137001', 'Father', NOW(3), NOW(3)),
(1, 'Li Xiaohong', '13700137002', 'Mother', NOW(3), NOW(3)),
(2, 'Li Jianguo', '13700137003', 'Father', NOW(3), NOW(3)),
(2, 'Wang Meili', '13700137004', 'Mother', NOW(3), NOW(3)),
(3, 'Wang Tianyi', '13700137005', 'Father', NOW(3), NOW(3)),
(4, 'Zhao Guoqiang', '13700137006', 'Father', NOW(3), NOW(3)),
(5, 'Sun Jianhua', '13700137007', 'Father', NOW(3), NOW(3)),
(5, 'Liu Xiuying', '13700137008', 'Mother', NOW(3), NOW(3)),
(6, 'Zhou Dawei', '13700137009', 'Father', NOW(3), NOW(3)),
(7, 'Wu Jianjun', '13700137010', 'Father', NOW(3), NOW(3)),
(8, 'Zheng Guoqing', '13700137011', 'Father', NOW(3), NOW(3)),
(9, 'Qian Wenhua', '13700137012', 'Father', NOW(3), NOW(3)),
(10, 'Sun Zhiqiang', '13700137013', 'Father', NOW(3), NOW(3));

-- ============================================
-- 10. Insert courses data
-- ============================================
INSERT INTO courses (course_name, teacher_id, credits, capacity, enrolled_count, created_at, updated_at) VALUES
-- Basic courses
('Data Structure', 1, 4.0, 50, 0, NOW(3), NOW(3)),
('Computer Arch', 2, 3.5, 50, 0, NOW(3), NOW(3)),
('Operating System', 3, 4.0, 50, 0, NOW(3), NOW(3)),
('Computer Network', 4, 3.5, 50, 0, NOW(3), NOW(3)),
('Database System', 5, 4.0, 50, 0, NOW(3), NOW(3)),

-- Major courses
('Software Eng', 6, 3.0, 40, 0, NOW(3), NOW(3)),
('AI Introduction', 7, 3.5, 40, 0, NOW(3), NOW(3)),
('Machine Learning', 8, 4.0, 35, 0, NOW(3), NOW(3)),
('Deep Learning', 1, 4.0, 30, 0, NOW(3), NOW(3)),
('Big Data Tech', 2, 3.5, 40, 0, NOW(3), NOW(3)),

-- Elective courses
('Web Development', 3, 2.5, 45, 0, NOW(3), NOW(3)),
('Mobile Dev', 4, 2.5, 45, 0, NOW(3), NOW(3)),
('Cloud Computing', 5, 3.0, 40, 0, NOW(3), NOW(3)),
('Info Security', 6, 3.0, 40, 0, NOW(3), NOW(3)),
('Blockchain Tech', 7, 2.5, 35, 0, NOW(3), NOW(3));

-- 扩展课程：批量新增 400 门课程，为大规模选课提供唯一组合
SET @course_existing_count := (SELECT COUNT(*) FROM courses);
SET @extra_course_total := 400;

INSERT INTO courses (course_name, teacher_id, credits, capacity, enrolled_count, created_at, updated_at)
SELECT 
    CONCAT('Extra Course ', LPAD(n, 3, '0')) AS course_name,
    ((n - 1) MOD 8) + 1 AS teacher_id,   -- 循环分配给 8 名教师
    3.0 AS credits,
    60 AS capacity,
    0 AS enrolled_count,
    NOW(3),
    NOW(3)
FROM (
    SELECT (t1.i + t2.i*10 + t3.i*100) + 1 AS n
    FROM (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t1
    CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t2
    CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t3
) seq
WHERE n <= @extra_course_total;

SET @course_new_start := @course_existing_count + 1;
SET @course_new_count := (SELECT COUNT(*) FROM courses) - @course_existing_count;

-- ============================================
-- 11. 插入课程先修关系
-- ============================================
INSERT INTO course_prerequisites (course_id, prereq_id, created_at) VALUES
-- 操作系统需要先修数据结构
(3, 1, NOW()),
-- 数据库系统需要先修数据结构
(5, 1, NOW()),
-- 机器学习需要先修人工智能导论
(8, 7, NOW()),
-- 深度学习需要先修机器学习
(9, 8, NOW()),
-- 大数据技术需要先修数据库系统
(10, 5, NOW());

-- ============================================
-- 12. 插入选课数据
-- ============================================
-- 为每个学生选3-5门课程
INSERT INTO enrollments (student_id, course_id, created_at, updated_at) VALUES
-- 学生1-5（计算机科学1班）
(1, 1, NOW(3), NOW(3)), (1, 2, NOW(3), NOW(3)), (1, 3, NOW(3), NOW(3)), (1, 4, NOW(3), NOW(3)),
(2, 1, NOW(3), NOW(3)), (2, 2, NOW(3), NOW(3)), (2, 4, NOW(3), NOW(3)), (2, 5, NOW(3), NOW(3)),
(3, 1, NOW(3), NOW(3)), (3, 2, NOW(3), NOW(3)), (3, 3, NOW(3), NOW(3)), (3, 5, NOW(3), NOW(3)),
(4, 1, NOW(3), NOW(3)), (4, 3, NOW(3), NOW(3)), (4, 4, NOW(3), NOW(3)), (4, 6, NOW(3), NOW(3)),
(5, 2, NOW(3), NOW(3)), (5, 3, NOW(3), NOW(3)), (5, 4, NOW(3), NOW(3)), (5, 5, NOW(3), NOW(3)),

-- 学生6-10（计算机科学2班）
(6, 1, NOW(3), NOW(3)), (6, 2, NOW(3), NOW(3)), (6, 6, NOW(3), NOW(3)), (6, 7, NOW(3), NOW(3)),
(7, 1, NOW(3), NOW(3)), (7, 3, NOW(3), NOW(3)), (7, 5, NOW(3), NOW(3)), (7, 7, NOW(3), NOW(3)),
(8, 2, NOW(3), NOW(3)), (8, 4, NOW(3), NOW(3)), (8, 6, NOW(3), NOW(3)), (8, 11, NOW(3), NOW(3)),
(9, 1, NOW(3), NOW(3)), (9, 5, NOW(3), NOW(3)), (9, 7, NOW(3), NOW(3)), (9, 11, NOW(3), NOW(3)),
(10, 2, NOW(3), NOW(3)), (10, 3, NOW(3), NOW(3)), (10, 4, NOW(3), NOW(3)), (10, 6, NOW(3), NOW(3)),

-- 学生11-15（软件工程1班）
(11, 1, NOW(3), NOW(3)), (11, 5, NOW(3), NOW(3)), (11, 6, NOW(3), NOW(3)), (11, 11, NOW(3), NOW(3)),
(12, 1, NOW(3), NOW(3)), (12, 2, NOW(3), NOW(3)), (12, 6, NOW(3), NOW(3)), (12, 12, NOW(3), NOW(3)),
(13, 3, NOW(3), NOW(3)), (13, 5, NOW(3), NOW(3)), (13, 6, NOW(3), NOW(3)), (13, 11, NOW(3), NOW(3)),
(14, 1, NOW(3), NOW(3)), (14, 4, NOW(3), NOW(3)), (14, 6, NOW(3), NOW(3)), (14, 12, NOW(3), NOW(3)),
(15, 2, NOW(3), NOW(3)), (15, 5, NOW(3), NOW(3)), (15, 6, NOW(3), NOW(3)), (15, 11, NOW(3), NOW(3)),

-- 学生16-20（软件工程2班）
(16, 1, NOW(3), NOW(3)), (16, 6, NOW(3), NOW(3)), (16, 11, NOW(3), NOW(3)), (16, 12, NOW(3), NOW(3)),
(17, 2, NOW(3), NOW(3)), (17, 5, NOW(3), NOW(3)), (17, 6, NOW(3), NOW(3)), (17, 11, NOW(3), NOW(3)),
(18, 1, NOW(3), NOW(3)), (18, 3, NOW(3), NOW(3)), (18, 6, NOW(3), NOW(3)), (18, 12, NOW(3), NOW(3)),
(19, 4, NOW(3), NOW(3)), (19, 5, NOW(3), NOW(3)), (19, 6, NOW(3), NOW(3)), (19, 11, NOW(3), NOW(3)),
(20, 1, NOW(3), NOW(3)), (20, 2, NOW(3), NOW(3)), (20, 6, NOW(3), NOW(3)), (20, 12, NOW(3), NOW(3)),

-- 学生21-25（数据科学1班）
(21, 1, NOW(3), NOW(3)), (21, 5, NOW(3), NOW(3)), (21, 7, NOW(3), NOW(3)), (21, 10, NOW(3), NOW(3)),
(22, 1, NOW(3), NOW(3)), (22, 5, NOW(3), NOW(3)), (22, 7, NOW(3), NOW(3)), (22, 13, NOW(3), NOW(3)),
(23, 2, NOW(3), NOW(3)), (23, 5, NOW(3), NOW(3)), (23, 10, NOW(3), NOW(3)), (23, 13, NOW(3), NOW(3)),
(24, 1, NOW(3), NOW(3)), (24, 4, NOW(3), NOW(3)), (24, 7, NOW(3), NOW(3)), (24, 10, NOW(3), NOW(3)),
(25, 3, NOW(3), NOW(3)), (25, 5, NOW(3), NOW(3)), (25, 7, NOW(3), NOW(3)), (25, 13, NOW(3), NOW(3)),

-- 学生26-30（人工智能1班）
(26, 1, NOW(3), NOW(3)), (26, 7, NOW(3), NOW(3)), (26, 8, NOW(3), NOW(3)), (26, 14, NOW(3), NOW(3)),
(27, 1, NOW(3), NOW(3)), (27, 5, NOW(3), NOW(3)), (27, 7, NOW(3), NOW(3)), (27, 8, NOW(3), NOW(3)),
(28, 2, NOW(3), NOW(3)), (28, 7, NOW(3), NOW(3)), (28, 8, NOW(3), NOW(3)), (28, 14, NOW(3), NOW(3)),
(29, 1, NOW(3), NOW(3)), (29, 4, NOW(3), NOW(3)), (29, 7, NOW(3), NOW(3)), (29, 8, NOW(3), NOW(3)),
(30, 3, NOW(3), NOW(3)), (30, 5, NOW(3), NOW(3)), (30, 7, NOW(3), NOW(3)), (30, 8, NOW(3), NOW(3));

-- 兼容未开启/不支持递归 CTE 的 MySQL 版本：用数字表生成 1..10000
SET @student_count := (SELECT COUNT(*) FROM students);

INSERT INTO enrollments (student_id, course_id, created_at, updated_at)
SELECT 
    FLOOR((n - 1) / @course_new_count) + 1 AS student_id,              -- 400 门扩展课为一组，顺序铺满学生
    @course_new_start + ((n - 1) MOD @course_new_count) AS course_id,   -- 每 400 条切换到下一位学生
    NOW(3),
    NOW(3)
FROM (
    SELECT (t1.i + t2.i*10 + t3.i*100 + t4.i*1000) + 1 AS n
    FROM (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t1
    CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t2
    CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t3
    CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) t4
) seq
WHERE n <= 10000;

-- 更新课程的已选人数
UPDATE courses c
SET enrolled_count = (
    SELECT COUNT(*) 
    FROM enrollments e 
    WHERE e.course_id = c.id AND e.deleted_at IS NULL
);

-- ============================================
-- 13. Insert grades data
-- ============================================
-- Add grades for each enrollment (regular, final, overall)
INSERT INTO grades (enrollment_id, score_type, score, created_at, updated_at)
SELECT 
    e.id,
    'Regular',
    ROUND(60 + RAND() * 35, 2),  -- Random 60-95
    NOW(3),
    NOW(3)
FROM enrollments e;

INSERT INTO grades (enrollment_id, score_type, score, created_at, updated_at)
SELECT 
    e.id,
    'Final',
    ROUND(55 + RAND() * 40, 2),  -- Random 55-95
    NOW(3),
    NOW(3)
FROM enrollments e;

INSERT INTO grades (enrollment_id, score_type, score, created_at, updated_at)
SELECT 
    e.id,
    'Overall',
    ROUND(60 + RAND() * 35, 2),  -- Random 60-95
    NOW(3),
    NOW(3)
FROM enrollments e;

-- ============================================
-- 13.5. Insert grade audit logs data
-- ============================================
-- Simulate some grade changes
INSERT INTO grade_audit_logs (grade_id, old_score, new_score, created_at, updated_at, updated_at_audit) VALUES
(1, 85.00, 90.00, NOW(3), NOW(3), NOW()),
(2, 70.00, 75.50, NOW(3), NOW(3), NOW()),
(5, 60.00, 55.00, NOW(3), NOW(3), NOW()),
(10, 88.00, 92.00, NOW(3), NOW(3), NOW()),
(15, 78.00, 80.00, NOW(3), NOW(3), NOW());

-- ============================================
-- 14. Insert attendance data
-- ============================================
-- Add 30 days attendance for first 10 students
INSERT INTO attendances (student_id, date, status, reason, teacher_id, created_at, updated_at)
SELECT 
    s.id,
    DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL n.n DAY), '%Y-%m-%d'),
    CASE 
        WHEN RAND() < 0.85 THEN 'Present'
        WHEN RAND() < 0.92 THEN 'Late'
        WHEN RAND() < 0.97 THEN 'Leave'
        ELSE 'Absent'
    END,
    CASE 
        WHEN RAND() < 0.85 THEN NULL
        WHEN RAND() < 0.92 THEN 'Traffic'
        WHEN RAND() < 0.97 THEN 'Sick'
        ELSE 'Unknown'
    END,
    FLOOR(1 + RAND() * 8),  -- Random teacher ID 1-8
    NOW(3),
    NOW(3)
FROM students s
CROSS JOIN (
    SELECT 0 AS n UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 
    UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9
    UNION SELECT 10 UNION SELECT 11 UNION SELECT 12 UNION SELECT 13 UNION SELECT 14
    UNION SELECT 15 UNION SELECT 16 UNION SELECT 17 UNION SELECT 18 UNION SELECT 19
    UNION SELECT 20 UNION SELECT 21 UNION SELECT 22 UNION SELECT 23 UNION SELECT 24
    UNION SELECT 25 UNION SELECT 26 UNION SELECT 27 UNION SELECT 28 UNION SELECT 29
) n
WHERE s.id <= 10;  -- Only for first 10 students

-- ============================================
-- 15. Insert reward/punishment data
-- ============================================
INSERT INTO reward_punishments (student_id, type, description, date, issuer, created_at, updated_at) VALUES
(1, 'Reward', 'Math Competition 1st Prize', '2024-03-15', 'Academic Office', NOW(3), NOW(3)),
(2, 'Reward', 'Excellent Student Leader', '2024-04-20', 'Student Office', NOW(3), NOW(3)),
(3, 'Reward', 'National Scholarship', '2024-05-10', 'Academic Office', NOW(3), NOW(3)),
(5, 'Reward', 'Excellent Member', '2024-06-01', 'Youth League', NOW(3), NOW(3)),
(7, 'Reward', 'Innovation Award', '2024-07-15', 'Research Office', NOW(3), NOW(3)),
(8, 'Punishment', 'Warning for Absence', '2024-03-25', 'Student Office', NOW(3), NOW(3)),
(11, 'Reward', 'Social Practice Award', '2024-08-20', 'Student Office', NOW(3), NOW(3)),
(15, 'Reward', 'Excellent Volunteer', '2024-09-10', 'Youth League', NOW(3), NOW(3)),
(18, 'Punishment', 'Cheating Demerit', '2024-06-30', 'Academic Office', NOW(3), NOW(3)),
(21, 'Reward', 'Programming Contest 2nd', '2024-10-15', 'CS Department', NOW(3), NOW(3));

-- ============================================
-- 16. Insert notifications data
-- ============================================
INSERT INTO notifications (title, content, sender_id, target, created_at, updated_at) VALUES
('Final Exam Notice', 'Final exams start on 15th next month. Prepare well.', 1, 'all', NOW(3), NOW(3)),
('Library Closed', 'Library closed this Saturday for maintenance.', 1, 'all', NOW(3), NOW(3)),
('CS Class 1 Meeting', 'Class meeting Friday 3PM at A101. All attend.', 2, 'class:1', NOW(3), NOW(3)),
('Scholarship Notice', 'Scholarship application starts. Submit by Monday.', 1, 'all', NOW(3), NOW(3)),
('Course Selection', 'Course selection opens Wednesday. Plan ahead.', 3, 'class:3', NOW(3), NOW(3)),
('Campus Safety', 'Stay safe on campus. Travel in groups at night.', 1, 'all', NOW(3), NOW(3)),
('DS Class Internship', 'Internship starts next month. Details in meeting.', 5, 'class:5', NOW(3), NOW(3)),
('Sports Day', 'Sports day next month. Sign up for events.', 1, 'all', NOW(3), NOW(3));

-- ============================================
-- 17. 插入课程表（排课）数据
-- ============================================
INSERT INTO schedules (course_id, class_id, teacher_id, day_of_week, start_time, end_time, location, semester, created_at, updated_at) VALUES
-- 计算机科学1班课表
(1, 1, 1, 1, '08:00', '09:40', 'A101', '2024-2025-1', NOW(3), NOW(3)),
(2, 1, 2, 1, '10:00', '11:40', 'A102', '2024-2025-1', NOW(3), NOW(3)),
(3, 1, 3, 2, '08:00', '09:40', 'A103', '2024-2025-1', NOW(3), NOW(3)),
(4, 1, 4, 2, '14:00', '15:40', 'A104', '2024-2025-1', NOW(3), NOW(3)),
(5, 1, 5, 3, '08:00', '09:40', 'A105', '2024-2025-1', NOW(3), NOW(3)),

-- 计算机科学2班课表
(1, 2, 1, 2, '10:00', '11:40', 'B101', '2024-2025-1', NOW(3), NOW(3)),
(2, 2, 2, 3, '08:00', '09:40', 'B102', '2024-2025-1', NOW(3), NOW(3)),
(6, 2, 6, 3, '14:00', '15:40', 'B103', '2024-2025-1', NOW(3), NOW(3)),
(7, 2, 7, 4, '08:00', '09:40', 'B104', '2024-2025-1', NOW(3), NOW(3)),

-- 软件工程1班课表
(1, 3, 1, 3, '10:00', '11:40', 'C101', '2024-2025-1', NOW(3), NOW(3)),
(5, 3, 5, 4, '10:00', '11:40', 'C102', '2024-2025-1', NOW(3), NOW(3)),
(6, 3, 6, 4, '14:00', '15:40', 'C103', '2024-2025-1', NOW(3), NOW(3)),
(11, 3, 3, 5, '08:00', '09:40', 'C104', '2024-2025-1', NOW(3), NOW(3)),

-- 软件工程2班课表
(1, 4, 1, 4, '08:00', '09:40', 'D101', '2024-2025-1', NOW(3), NOW(3)),
(6, 4, 6, 5, '10:00', '11:40', 'D102', '2024-2025-1', NOW(3), NOW(3)),
(11, 4, 3, 5, '14:00', '15:40', 'D103', '2024-2025-1', NOW(3), NOW(3)),
(12, 4, 4, 1, '14:00', '15:40', 'D104', '2024-2025-1', NOW(3), NOW(3)),

-- 数据科学1班课表
(1, 5, 1, 5, '08:00', '09:40', 'E101', '2024-2025-1', NOW(3), NOW(3)),
(5, 5, 5, 1, '10:00', '11:40', 'E102', '2024-2025-1', NOW(3), NOW(3)),
(7, 5, 7, 2, '10:00', '11:40', 'E103', '2024-2025-1', NOW(3), NOW(3)),
(10, 5, 2, 3, '14:00', '15:40', 'E104', '2024-2025-1', NOW(3), NOW(3)),

-- 人工智能1班课表
(1, 6, 1, 1, '14:00', '15:40', 'F101', '2024-2025-1', NOW(3), NOW(3)),
(7, 6, 7, 3, '10:00', '11:40', 'F102', '2024-2025-1', NOW(3), NOW(3)),
(8, 6, 8, 4, '14:00', '15:40', 'F103', '2024-2025-1', NOW(3), NOW(3)),
(14, 6, 6, 5, '14:00', '15:40', 'F104', '2024-2025-1', NOW(3), NOW(3));

-- ============================================
-- 数据插入完成
-- ============================================

-- View data statistics
SELECT 'Data Insert Statistics' AS info;
SELECT 'roles' AS table_name, COUNT(*) AS record_count FROM roles
UNION ALL SELECT 'permissions', COUNT(*) FROM permissions
UNION ALL SELECT 'role_permissions', COUNT(*) FROM role_permissions
UNION ALL SELECT 'teachers', COUNT(*) FROM teachers
UNION ALL SELECT 'classes', COUNT(*) FROM classes
UNION ALL SELECT 'students', COUNT(*) FROM students
UNION ALL SELECT 'users', COUNT(*) FROM users
UNION ALL SELECT 'parents', COUNT(*) FROM parents
UNION ALL SELECT 'courses', COUNT(*) FROM courses
UNION ALL SELECT 'course_prerequisites', COUNT(*) FROM course_prerequisites
UNION ALL SELECT 'enrollments', COUNT(*) FROM enrollments
UNION ALL SELECT 'grades', COUNT(*) FROM grades
UNION ALL SELECT 'grade_audit_logs', COUNT(*) FROM grade_audit_logs
UNION ALL SELECT 'attendances', COUNT(*) FROM attendances
UNION ALL SELECT 'reward_punishments', COUNT(*) FROM reward_punishments
UNION ALL SELECT 'notifications', COUNT(*) FROM notifications
UNION ALL SELECT 'schedules', COUNT(*) FROM schedules;

-- ============================================
-- 使用说明
-- ============================================
-- 1. 默认账号密码：
--    管理员：admin / admin123
--    教师：teacher001 / teacher123
--    学生：stu001 / student123
--
-- 2. 数据说明：
--    - 8名教师
--    - 6个班级
--    - 30名学生
--    - 13个家长记录
--    - 15门课程
--    - 120条选课记录
--    - 360条成绩记录
--    - 300条考勤记录
--    - 10条奖惩记录
--    - 8条通知
--    - 24条排课记录
--
-- 3. 注意事项：
--    - 密码哈希是示例值,实际使用时需要后端生成真实的bcrypt哈希
--    - 可以根据需要调整数据量和内容
--    - 如需重新插入,请先清空相关表
-- ============================================

-- ============================================
-- 重新启用外键检查
-- ============================================
SET FOREIGN_KEY_CHECKS = 1;

-- 数据插入完成!
SELECT '✓ All data inserted successfully!' AS status;
