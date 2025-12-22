-- ============================================
-- 数据库功能验证脚本
-- ============================================
-- 用途：验证 complete_schema.sql 执行后的数据库功能
-- 使用方法：
-- mysql -u root -p -e "source docs/verify_schema.sql"
-- ============================================

USE student_db;

SELECT '============================================' AS '';
SELECT '开始验证数据库架构...' AS '';
SELECT '============================================' AS '';

-- ============================================
-- 1. 检查表是否存在
-- ============================================
SELECT '' AS '';
SELECT '1. 检查表结构...' AS '';
SELECT '-------------------' AS '';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ course_prerequisites 表存在'
        ELSE '✗ course_prerequisites 表不存在'
    END AS '检查结果'
FROM information_schema.tables 
WHERE table_schema = 'student_db' 
  AND table_name = 'course_prerequisites';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ grade_audit_logs 表存在'
        ELSE '✗ grade_audit_logs 表不存在'
    END AS '检查结果'
FROM information_schema.tables 
WHERE table_schema = 'student_db' 
  AND table_name = 'grade_audit_logs';

-- ============================================
-- 2. 检查视图是否存在
-- ============================================
SELECT '' AS '';
SELECT '2. 检查视图...' AS '';
SELECT '-------------------' AS '';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ vw_class_performance 视图存在'
        ELSE '✗ vw_class_performance 视图不存在'
    END AS '检查结果'
FROM information_schema.views 
WHERE table_schema = 'student_db' 
  AND table_name = 'vw_class_performance';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ vw_student_full_profile 视图存在'
        ELSE '✗ vw_student_full_profile 视图不存在'
    END AS '检查结果'
FROM information_schema.views 
WHERE table_schema = 'student_db' 
  AND table_name = 'vw_student_full_profile';

-- ============================================
-- 3. 检查触发器是否存在
-- ============================================
SELECT '' AS '';
SELECT '3. 检查触发器...' AS '';
SELECT '-------------------' AS '';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ trg_audit_grade_update 触发器存在'
        ELSE '✗ trg_audit_grade_update 触发器不存在'
    END AS '检查结果'
FROM information_schema.triggers 
WHERE trigger_schema = 'student_db' 
  AND trigger_name = 'trg_audit_grade_update';

-- ============================================
-- 4. 检查存储过程是否存在
-- ============================================
SELECT '' AS '';
SELECT '4. 检查存储过程...' AS '';
SELECT '-------------------' AS '';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ sp_enroll_student 存储过程存在'
        ELSE '✗ sp_enroll_student 存储过程不存在'
    END AS '检查结果'
FROM information_schema.routines 
WHERE routine_schema = 'student_db' 
  AND routine_name = 'sp_enroll_student'
  AND routine_type = 'PROCEDURE';

-- ============================================
-- 5. 检查courses表字段
-- ============================================
SELECT '' AS '';
SELECT '5. 检查课程表字段...' AS '';
SELECT '-------------------' AS '';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ courses.capacity 字段存在'
        ELSE '✗ courses.capacity 字段不存在'
    END AS '检查结果'
FROM information_schema.columns 
WHERE table_schema = 'student_db' 
  AND table_name = 'courses' 
  AND column_name = 'capacity';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ courses.enrolled_count 字段存在'
        ELSE '✗ courses.enrolled_count 字段不存在'
    END AS '检查结果'
FROM information_schema.columns 
WHERE table_schema = 'student_db' 
  AND table_name = 'courses' 
  AND column_name = 'enrolled_count';

-- ============================================
-- 6. 检查索引
-- ============================================
SELECT '' AS '';
SELECT '6. 检查索引...' AS '';
SELECT '-------------------' AS '';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ grade_audit_logs 表有 idx_grade_id 索引'
        ELSE '✗ grade_audit_logs 表缺少 idx_grade_id 索引'
    END AS '检查结果'
FROM information_schema.statistics 
WHERE table_schema = 'student_db' 
  AND table_name = 'grade_audit_logs'
  AND index_name = 'idx_grade_id';

SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN '✓ course_prerequisites 表有 idx_prereq_course 索引'
        ELSE '✗ course_prerequisites 表缺少 idx_prereq_course 索引'
    END AS '检查结果'
FROM information_schema.statistics 
WHERE table_schema = 'student_db' 
  AND table_name = 'course_prerequisites'
  AND index_name = 'idx_prereq_course';

-- ============================================
-- 7. 统计数据库对象
-- ============================================
SELECT '' AS '';
SELECT '7. 数据库对象统计...' AS '';
SELECT '-------------------' AS '';

SELECT COUNT(*) AS '表总数' 
FROM information_schema.tables 
WHERE table_schema = 'student_db' 
  AND table_type = 'BASE TABLE';

SELECT COUNT(*) AS '视图总数' 
FROM information_schema.views 
WHERE table_schema = 'student_db';

SELECT COUNT(*) AS '触发器总数' 
FROM information_schema.triggers 
WHERE trigger_schema = 'student_db';

SELECT COUNT(*) AS '存储过程总数' 
FROM information_schema.routines 
WHERE routine_schema = 'student_db' 
  AND routine_type = 'PROCEDURE';

-- ============================================
-- 验证完成
-- ============================================
SELECT '' AS '';
SELECT '============================================' AS '';
SELECT '验证完成！' AS '';
SELECT '============================================' AS '';
