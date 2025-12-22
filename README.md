# 数据库表管理系统

这是一个基于 Go + Gin + Gorm + Vue 3 + Element Plus 的通用数据库表管理系统。项目最初是一个学生管理系统，现已重构为专注于数据库表结构的通用管理工具。

## 项目概述

本系统提供了一个通用的后台管理界面，用于管理数据库中的各种表数据。它保留了完善的用户认证和权限管理系统，并提供了对数据库表的增删改查（CRUD）、数据导出、SQL执行等通用功能。

### 功能亮点

- 🔐 **完善的权限系统**：基于 RBAC 的细粒度权限控制，支持角色和权限自定义
- 🗃️ **通用表管理**：支持对数据库中任意注册表的 CRUD 操作
- 📊 **数据可视化**：提供基础的数据统计和概览
- 🛠️ **SQL 执行**：支持直接执行 SQL 查询（开发模式）
- 📤 **数据导出**：支持将表数据导出为 JSON 格式
- 🎨 **现代化界面**：基于 Element Plus 的响应式设计，操作简单直观

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: Gorm
- **数据库**: MySQL 8.0+
- **认证**: JWT (JSON Web Tokens)
- **密码加密**: bcrypt

### 前端
- **框架**: Vue 3.3+
- **构建工具**: Vite 5
- **路由**: Vue Router 4.2+
- **状态管理**: Pinia 2.1+
- **HTTP客户端**: Axios 1.6+
- **UI组件库**: Element Plus 2.4+

## 项目结构

```
database-management-system/
├── backend/                # 后端服务
│   ├── cmd/               # 程序入口
│   │   ├── main.go       # 主程序
│   │   └── create_admin.go # 创建管理员工具
│   ├── config/            # 配置文件和数据库初始化
│   │   └── database.go   # 数据库连接与表迁移
│   ├── internal/          # 内部代码
│   │   ├── api/          # API路由和处理器
│   │   │   ├── router.go # 路由配置
│   │   │   └── v1/       # v1版本API处理器
│   │   │       ├── database.go # 通用数据库表接口
│   │   │       ├── auth.go     # 认证接口
│   │   │       └── admin_*.go  # 管理员接口
│   │   ├── middleware/   # 中间件
│   │   └── models/       # 数据模型
│   └── go.mod
├── frontend/              # 前端应用
│   ├── src/
│   │   ├── api/          # API接口封装
│   │   ├── views/        # 页面视图组件
│   │   └── ...
│   └── package.json
└── README.md
```

## 功能模块

### ✅ 已实现

#### 核心功能
- [x] **用户认证**：登录/登出、JWT令牌认证
- [x] **通用表管理**：
    - 获取表列表及统计信息
    - 获取表结构（Schema）
    - 表数据分页查询
    - 表数据新增、修改、删除
    - 表数据导出
- [x] **SQL 执行工具**：直接运行 SQL 语句并查看结果

#### 系统管理
- [x] **用户管理**：用户 CRUD、角色分配
- [x] **角色管理**：角色 CRUD、权限分配
- [x] **权限管理**：查看系统权限列表

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 1. 克隆项目

```bash
git clone <repository-url>
cd database-management-system
```

### 2. 配置数据库

创建MySQL数据库：

```sql
CREATE DATABASE student_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 启动后端

```bash
cd backend

# 安装依赖
go mod tidy

# 运行服务
go run ./cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 4. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端应用将在 `http://localhost:5173` 启动。

### 5. 创建管理员账号

首次运行需要创建管理员账号：

```bash
cd backend
go run ./cmd/create_admin.go
```

## API文档

后端API接口文档：

- **基础路径**: `/api/v1`
- **认证方式**: Bearer Token (JWT)

### 主要接口

#### 认证相关
```
POST   /api/v1/auth/login      # 用户登录
GET    /api/v1/auth/me         # 获取当前用户信息
```

#### 数据库管理 (核心)
```
GET    /api/v1/database/tables                 # 获取所有表列表及统计
GET    /api/v1/database/tables/:table          # 获取指定表数据（分页）
GET    /api/v1/database/tables/:table/schema   # 获取表结构
POST   /api/v1/database/tables/:table          # 新增数据
PUT    /api/v1/database/tables/:table/:id      # 更新数据
DELETE /api/v1/database/tables/:table/:id      # 删除数据
GET    /api/v1/database/tables/:table/export   # 导出数据
POST   /api/v1/database/execute                # 执行 SQL
```

#### 管理员模块
```
GET    /api/v1/admin/users     # 用户列表
GET    /api/v1/admin/roles     # 角色列表
GET    /api/v1/admin/permissions # 权限列表
```

## 注意事项

1.  **安全性**: 本系统包含直接操作数据库的功能（如 SQL 执行），请务必在生产环境中严格限制该功能的访问权限，或将其禁用。
2.  **数据备份**: 在执行批量删除或修改操作前，建议备份数据库。

## 许可证

MIT License
