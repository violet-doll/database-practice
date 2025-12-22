# 数据库完整初始化指南

## 概述

本指南说明如何使用 `init_complete_database.sql` 脚本从零开始构建完整的学生管理系统数据库。

## 脚本功能

`init_complete_database.sql` 是一个**一站式**数据库初始化脚本，包含：

### 1. 数据库创建
- 创建 `student_db` 数据库（UTF-8编码）

### 2. 基础表结构（18张表）
- `roles` - 角色表
- `permissions` - 权限表
- `role_permissions` - 角色权限关联表
- `users` - 用户表
- `teachers` - 教师表
- `classes` - 班级表
- `students` - 学生表
- `parents` - 家长表
- `courses` - 课程表（含容量控制字段）
- `enrollments` - 选课表
- `grades` - 成绩表
- `attendances` - 考勤表
- `reward_punishments` - 奖惩表
- `notifications` - 通知表
- `schedules` - 课程表（排课）
- `course_prerequisites` - 课程先修关系表
- `grade_audit_logs` - 成绩审计日志表

### 3. 索引优化
- 所有外键索引
- 软删除索引（deleted_at）
- 业务查询索引（如 student_id, course_id 等）

### 4. 触发器
- `trg_audit_grade_update` - 成绩修改自动审计触发器

### 5. 存储过程
- `sp_enroll_student` - 智能选课存储过程（含先修课程检查、容量控制）

### 6. 视图
- `vw_class_performance` - 班级成绩统计视图
- `vw_student_full_profile` - 学生完整档案视图

### 7. 初始数据
- 4个默认角色：admin, teacher, student, parent

## 使用场景

### 场景一：全新安装（推荐）

适用于：
- 第一次部署系统
- 开发环境初始化
- 测试环境重置

```powershell
# 直接执行完整初始化脚本（PowerShell）
mysql -u root -p  -e "source docs/init_complete_database.sql"

# 或者使用 Linux/Mac 的 bash
# mysql -u root -p < docs/init_complete_database.sql
```

执行后输入MySQL密码，脚本将自动完成所有初始化工作。

### 场景二：现有数据库升级

如果您已有 `student_db` 数据库，但需要添加新功能：

**选项A：完全重建（会丢失数据）**
```powershell
# 1. 备份现有数据
mysqldump -u root -p student_db > "backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').sql"

# 2. 删除旧数据库
mysql -u root -p -e "DROP DATABASE IF EXISTS student_db;"

# 3. 执行完整初始化
mysql -u root -p  -e "source docs/init_complete_database.sql"

# 4. 如需恢复数据，可选择性导入备份
```

**选项B：增量更新（保留数据）**
```powershell
# 仅执行 complete_schema.sql（不包含基础表创建）
mysql -u root -p -e "source docs/complete_schema.sql"
```

### 场景三：验证数据库结构

执行初始化后，验证数据库是否正确创建：

```powershell
# 执行验证脚本
mysql -u root -p -e "source docs/verify_schema.sql"
```

## 执行步骤详解

### 步骤1：准备环境

确保满足以下条件：
- MySQL 8.0+ 已安装并运行
- 有数据库管理权限（root 或具有 CREATE DATABASE 权限的用户）
- 如有旧数据库，已完成备份

### 步骤2：执行初始化脚本

**Windows PowerShell:**
```powershell
# 进入项目目录
cd c:\Users\violet\Desktop\student-management-system

# 执行初始化脚本
mysql -u root -p  -e "source docs/init_complete_database.sql"
```

**Linux/Mac:**
```bash
# 进入项目目录
cd /path/to/student-management-system

# 执行初始化脚本
mysql -u root -p < docs/init_complete_database.sql
```

### 步骤3：验证安装

连接到数据库并检查：

```sql
-- 连接数据库
mysql -u root -p student_db

-- 查看所有表
SHOW TABLES;

-- 查看所有视图
SHOW FULL TABLES WHERE Table_type = 'VIEW';

-- 查看所有触发器
SHOW TRIGGERS;

-- 查看所有存储过程
SHOW PROCEDURE STATUS WHERE Db = 'student_db';

-- 查看角色数据
SELECT * FROM roles;
```

预期结果：
- 18张表
- 2个视图
- 1个触发器
- 1个存储过程
- 4条角色记录

### 步骤4：启动后端服务

数据库初始化完成后，启动后端服务：

```bash
cd backend
go run cmd/server/main.go
```

GORM 会自动同步模型与数据库，但不会删除已有的表和数据。

## 与其他脚本的关系

项目中有多个SQL脚本，它们的关系如下：

| 脚本文件 | 用途 | 使用场景 |
|---------|------|---------|
| `init_complete_database.sql` | **完整初始化** | 从零开始构建数据库 |
| `init.sql` | 简单初始化（已过时） | 仅创建数据库，依赖GORM创建表 |
| `complete_schema.sql` | 高级功能更新 | 在已有数据库上添加视图、触发器等 |
| `update_schema.sql` | 部分功能更新 | 仅添加审计功能 |
| `verify_schema.sql` | 验证脚本 | 检查数据库结构是否正确 |

**推荐使用：**
- **新项目**：使用 `init_complete_database.sql`
- **已有项目升级**：使用 `complete_schema.sql`
- **验证**：使用 `verify_schema.sql`

## 功能测试

### 测试触发器（成绩审计）

```sql
-- 1. 插入测试数据（需要先有学生、课程、选课记录）
-- 假设已有 enrollment_id = 1
INSERT INTO grades (enrollment_id, score_type, score, created_at, updated_at)
VALUES (1, '期末', 85.0, NOW(3), NOW(3));

-- 2. 修改成绩，触发器应自动记录
UPDATE grades SET score = 90.0 WHERE id = 1;

-- 3. 查看审计日志
SELECT * FROM grade_audit_logs WHERE grade_id = 1;
-- 应该看到一条记录：old_score = 85.0, new_score = 90.0
```

### 测试存储过程（智能选课）

```sql
-- 调用选课存储过程
CALL sp_enroll_student(1, 2, @status, @message);

-- 查看结果
SELECT @status AS 状态码, @message AS 消息;
-- 状态码 0 = 成功, 1 = 失败
```

### 测试视图（班级成绩统计）

```sql
-- 查询班级成绩统计
SELECT * FROM vw_class_performance LIMIT 5;

-- 查询学生完整档案
SELECT * FROM vw_student_full_profile LIMIT 5;
```

## 常见问题

### Q1: 执行脚本时报错 "Database exists"
**A:** 脚本使用了 `CREATE DATABASE IF NOT EXISTS`，不会报错。如果仍报错，可能是权限问题。

### Q2: 触发器创建失败
**A:** 检查MySQL版本是否 >= 5.7，并确保有 TRIGGER 权限：
```sql
GRANT TRIGGER ON student_db.* TO 'your_user'@'localhost';
```

### Q3: 视图查询为空
**A:** 视图依赖基础表的数据，初始化后表是空的。需要通过后端API或手动插入测试数据。

### Q4: GORM 启动时报错 "Table already exists"
**A:** 这是正常的，GORM 会尝试创建表，但因为表已存在会跳过。不影响使用。

### Q5: 如何重置数据库到初始状态？
**A:** 
```powershell
# 删除数据库
mysql -u root -p -e "DROP DATABASE IF EXISTS student_db;"

# 重新执行初始化
mysql -u root -p  -e "source docs/init_complete_database.sql"
```

## 数据备份与恢复

### 备份数据库

```powershell
# 完整备份（包含结构和数据）
mysqldump -u root -p student_db > backup_full.sql

# 仅备份数据（不含结构）
mysqldump -u root -p --no-create-info student_db > backup_data.sql

# 仅备份结构（不含数据）
mysqldump -u root -p --no-data student_db > backup_schema.sql
```

### 恢复数据库

```powershell
# 恢复完整备份
mysql -u root -p -e "source backup_full.sql"

# 恢复数据（需要先有表结构）
mysql -u root -p -e "source backup_data.sql"
```

## 性能优化建议

1. **定期分析表**
```sql
ANALYZE TABLE students, courses, enrollments, grades;
```

2. **监控慢查询**
```sql
-- 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
```

3. **视图性能**
- 视图 `vw_class_performance` 和 `vw_student_full_profile` 是实时计算的
- 如果数据量大，考虑创建物化视图或定时任务缓存结果

## 安全建议

1. **修改默认密码**
   - 脚本创建的角色没有关联用户，需要通过后端API创建管理员账号

2. **限制数据库权限**
```sql
-- 创建应用专用用户（不使用root）
CREATE USER 'student_app'@'localhost' IDENTIFIED BY 'strong_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON student_db.* TO 'student_app'@'localhost';
FLUSH PRIVILEGES;
```

3. **定期备份**
   - 建议每天自动备份数据库
   - 保留至少7天的备份历史

## 下一步

数据库初始化完成后：

1. ✅ 启动后端服务
2. ✅ 通过API创建管理员账号
3. ✅ 登录系统并开始使用
4. ✅ 根据需要添加教师、学生、课程等数据

## 技术支持

如有问题，请查看：
- 项目 README.md
- MySQL 错误日志
- 后端应用日志

---

**最后更新：** 2025-12-22
