-- 学生管理系统 - 数据库初始化脚本

-- 1. 创建数据库
CREATE DATABASE IF NOT EXISTS student_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE student_db;

-- 注意: Gorm 会自动创建表结构，这里只需要插入初始数据

-- 2. 插入默认角色（如果表已存在）
-- 注意: 运行此脚本前，请先启动后端服务让 Gorm 创建表结构

-- 插入管理员角色
-- INSERT INTO roles (role_name, created_at, updated_at) 
-- VALUES ('admin', NOW(), NOW());

-- 插入教师角色
-- INSERT INTO roles (role_name, created_at, updated_at) 
-- VALUES ('teacher', NOW(), NOW());

-- 插入学生角色
-- INSERT INTO roles (role_name, created_at, updated_at) 
-- VALUES ('student', NOW(), NOW());

-- 插入家长角色
-- INSERT INTO roles (role_name, created_at, updated_at) 
-- VALUES ('parent', NOW(), NOW());

-- 3. 创建管理员账号
-- 用户名: admin
-- 密码: admin123 (已使用 bcrypt 加密)
-- 注意: 下面的密码hash是示例，需要实际生成

-- INSERT INTO users (username, password, role_id, is_active, user_type, created_at, updated_at)
-- VALUES ('admin', '$2a$10$5vJZIB7h5Y7HvXN4pXqYOeGF9xX.K5JoqD0J3lKW0hH.FZJ9Y3GJm', 1, 1, 'admin', NOW(), NOW());

-- 提示: 
-- 1. 先启动后端服务，让 Gorm 自动创建表结构
-- 2. 角色数据会在后端启动时自动创建
-- 3. 管理员账号需要手动创建或通过后端程序创建
