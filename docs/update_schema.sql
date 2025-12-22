-- 学生管理系统 - 数据库更新脚本
-- 第二阶段：成绩修改审计功能（Trigger）

-- ============================================
-- 1. 创建成绩修改审计表
-- ============================================
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
    KEY idx_grade_audit_logs_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成绩修改审计日志表';

-- 为审计表添加索引以提高查询性能
CREATE INDEX idx_grade_id ON grade_audit_logs(grade_id);
CREATE INDEX idx_updated_at ON grade_audit_logs(updated_at);

-- ============================================
-- 2. 创建成绩修改触发器
-- ============================================
-- 功能：在更新 grades 表之前，自动记录旧值
-- 亮点：完全由数据库自动完成，后端代码无需额外逻辑

DELIMITER //

CREATE TRIGGER trg_audit_grade_update
BEFORE UPDATE ON grades
FOR EACH ROW
BEGIN
    -- 只有当分数确实发生变化时才记录
    IF OLD.score != NEW.score THEN
        INSERT INTO grade_audit_logs (grade_id, old_score, new_score, created_at, updated_at, updated_at_audit)
        VALUES (OLD.id, OLD.score, NEW.score, NOW(), NOW(), NOW());
    END IF;
END //

DELIMITER ;

-- ============================================
-- 验证说明
-- ============================================
-- 测试触发器：
-- 1. 更新一条成绩：UPDATE grades SET score = 95.0 WHERE id = 1;
-- 2. 查看审计日志：SELECT * FROM grade_audit_logs WHERE grade_id = 1;
-- 3. 验证只记录分数变化：UPDATE grades SET score_type = '期末' WHERE id = 1;
--    (不应该产生新的审计记录)
