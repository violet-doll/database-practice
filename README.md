# 学生管理系统

这是一个基于 Go + Gin + Gorm + Vue 3 + Element Plus 的现代化学生管理系统，是《软件工程》课程的实验项目。

## 项目概述

本系统旨在为学校或教育机构提供一个高效、全面的学生信息管理平台，涵盖学籍管理、课程管理、成绩管理、考勤管理、奖惩管理等多个功能模块。

### 功能亮点

- 🔐 **完善的权限系统**：基于 RBAC 的细粒度权限控制，支持角色和权限自定义
- 📊 **数据可视化**：使用 ECharts 展示统计数据，直观了解系统运行状态
- 🎓 **全面的学生管理**：学生信息、班级、课程、成绩、考勤、奖惩一体化管理
- 👨‍🏫 **灵活的排课系统**：支持课程安排、教室分配、时间管理
- 📱 **家校互通**：家长联系管理和通知系统，便于家校沟通
- 🔍 **高效查询**：支持分页、搜索、筛选等多种查询方式
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
- **数据可视化**: ECharts 6.0+

## 项目结构

```
student-management-system/
├── backend/                # 后端服务
│   ├── cmd/               # 程序入口
│   │   ├── main.go       # 主程序
│   │   └── create_admin.go # 创建管理员工具
│   ├── config/            # 配置文件和数据库初始化
│   │   ├── config.ini    # 配置文件
│   │   └── database.go   # 数据库连接
│   ├── internal/          # 内部代码
│   │   ├── api/          # API路由和处理器
│   │   │   ├── router.go # 路由配置
│   │   │   └── v1/       # v1版本API处理器
│   │   ├── middleware/   # 中间件
│   │   │   ├── auth.go   # JWT认证和权限中间件
│   │   │   └── cors.go   # CORS中间件
│   │   ├── models/       # 数据模型
│   │   │   └── models.go # 所有数据模型定义
│   │   └── utils/        # 工具函数
│   │       ├── jwt.go    # JWT工具
│   │       └── password.go # 密码加密
│   ├── scripts/           # 脚本工具
│   └── go.mod
├── frontend/              # 前端应用
│   ├── src/
│   │   ├── api/          # API接口封装
│   │   ├── layouts/      # 布局组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # Pinia状态管理
│   │   ├── utils/        # 工具函数
│   │   └── views/        # 页面视图组件
│   └── package.json
├── docs/                  # 项目文档
│   └── init.sql          # 数据库初始化脚本
└── README.md
```

## 功能模块

### ✅ 已实现

#### 核心功能
- [x] **用户认证**：登录/登出、JWT令牌认证、路由守卫、自动登出
- [x] **学生信息管理**：完整的CRUD操作、分页、搜索（姓名、学号、班级）
- [x] **班级管理**：班级CRUD、关联班主任、查看班级学生
- [x] **课程管理**：课程CRUD、关联授课教师、学分管理
- [x] **排课管理**：课程时间安排、教室分配
- [x] **选课管理**：学生选课、批量选课、选课记录查询
- [x] **成绩管理**：成绩录入、按学生查询、按课程查询、成绩列表（分页、筛选）
- [x] **考勤管理**：考勤记录（出勤/缺席/请假/迟到）、按学生查询、考勤统计
- [x] **奖惩管理**：奖惩记录录入、按学生查询、奖惩列表（分页、筛选）
- [x] **家长联系管理**：家长联系方式CRUD、关联学生
- [x] **通知管理**：通知发布、通知列表、模拟短信/邮件发送
- [x] **数据看板**：系统概览统计、关键数据展示

#### 系统管理
- [x] **基于角色的权限控制（RBAC）**：角色管理、权限分配、细粒度权限控制
- [x] **用户管理**：用户CRUD、用户状态管理、角色分配
- [x] **角色管理**：角色CRUD、权限关联
- [x] **权限管理**：权限列表、角色权限分配
- [x] **统计概览**：系统数据统计（用户、学生、课程、成绩、考勤等）

### 🚧 待优化功能
- [ ] 数据导入导出（Excel）
- [ ] 批量操作
- [ ] 高级搜索和筛选
- [ ] 操作日志记录
- [ ] 数据备份和恢复
- [ ] 系统设置

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 1. 克隆项目

```bash
git clone <repository-url>
cd student-management-system
```

### 2. 配置数据库

创建MySQL数据库：

```sql
CREATE DATABASE student_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

可选：使用 `docs/init.sql` 初始化数据库表结构和基础数据。

### 3. 启动后端

```bash
cd backend

# 修改 .env 或 config/config.ini 中的数据库配置
# .env 示例：
# DB_HOST=localhost
# DB_PORT=3306
# DB_USER=root
# DB_PASSWORD=your_password
# DB_NAME=student_db

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

首次运行需要创建管理员账号。推荐使用以下方式：

**方式一：使用创建管理员工具**
```bash
cd backend
go run ./cmd/create_admin.go
```

**方式二：使用初始化脚本**
```bash
mysql -u root -p student_db < docs/init.sql
```

**方式三：手动在数据库中插入**
```sql
-- 先创建角色（如果不存在）
INSERT INTO roles (role_name, created_at, updated_at) VALUES ('管理员', NOW(), NOW());

-- 插入管理员用户（密码: admin123）
-- 注意：需要使用bcrypt加密后的密码哈希值
INSERT INTO users (username, password, role_id, is_active, user_type, created_at, updated_at)
VALUES ('admin', '$2a$10$YourBcryptHashHere', 1, 1, 'admin', NOW(), NOW());
```

## 默认账号

- **用户名**: admin
- **密码**: admin123

（需要先创建该账号）

## API文档

后端API接口文档：

- **基础路径**: `/api/v1`
- **认证方式**: Bearer Token (JWT)
- **响应格式**: JSON

### 主要接口

#### 认证相关
```
POST   /api/v1/auth/login      # 用户登录
POST   /api/v1/auth/logout     # 用户登出
GET    /api/v1/auth/me         # 获取当前用户信息
```

#### 学生管理
```
GET    /api/v1/students        # 学生列表（分页、搜索：name, student_id, class_id）
GET    /api/v1/students/:id    # 学生详情
POST   /api/v1/students        # 新增学生
PUT    /api/v1/students/:id    # 更新学生信息
DELETE /api/v1/students/:id    # 删除学生
```

#### 班级管理
```
GET    /api/v1/classes         # 班级列表
GET    /api/v1/classes/:id     # 班级详情
POST   /api/v1/classes         # 新增班级
PUT    /api/v1/classes/:id     # 更新班级
DELETE /api/v1/classes/:id    # 删除班级
```

#### 课程管理
```
GET    /api/v1/courses         # 课程列表（分页、搜索：course_name, teacher_id）
GET    /api/v1/courses/:id     # 课程详情
POST   /api/v1/courses         # 新增课程
PUT    /api/v1/courses/:id     # 更新课程
DELETE /api/v1/courses/:id     # 删除课程
```

#### 排课管理
```
GET    /api/v1/schedules       # 排课列表
GET    /api/v1/schedules/:id   # 排课详情
POST   /api/v1/schedules       # 新增排课
PUT    /api/v1/schedules/:id   # 更新排课
DELETE /api/v1/schedules/:id  # 删除排课
```

#### 选课管理
```
GET    /api/v1/enrollments     # 选课列表（分页、筛选）
POST   /api/v1/enrollments     # 新增选课
DELETE /api/v1/enrollments/:id # 删除选课
```

#### 成绩管理
```
GET    /api/v1/grades          # 成绩列表（分页、筛选）
GET    /api/v1/grades/student/:id  # 按学生查询成绩
GET    /api/v1/grades/course/:id   # 按课程查询成绩
POST   /api/v1/grades          # 录入成绩
```

#### 考勤管理
```
GET    /api/v1/attendance      # 考勤列表（分页、筛选）
GET    /api/v1/attendance/student/:id  # 按学生查询考勤
GET    /api/v1/attendance/stats # 考勤统计
POST   /api/v1/attendance      # 新增考勤记录
DELETE /api/v1/attendance/:id  # 删除考勤记录
```

#### 奖惩管理
```
GET    /api/v1/rewards         # 奖惩列表（分页、筛选）
GET    /api/v1/rewards/student/:id  # 按学生查询奖惩
POST   /api/v1/rewards         # 新增奖惩记录
DELETE /api/v1/rewards/:id     # 删除奖惩记录
```

#### 家长联系管理
```
GET    /api/v1/parents         # 家长联系方式列表（分页，支持 student_id 筛选）
POST   /api/v1/parents         # 新增家长联系方式
PUT    /api/v1/parents/:id     # 更新家长联系方式
DELETE /api/v1/parents/:id     # 删除家长联系方式
```

#### 通知管理
```
GET    /api/v1/notifications   # 通知列表（分页，按 target 筛选）
POST   /api/v1/notifications   # 发布通知（模拟发送短信/邮件）
```

#### 管理员模块
```
# 用户管理
GET    /api/v1/admin/users    # 用户列表
POST   /api/v1/admin/users     # 新增用户
PUT    /api/v1/admin/users/:id # 更新用户
DELETE /api/v1/admin/users/:id # 删除用户

# 角色管理
GET    /api/v1/admin/roles     # 角色列表
POST   /api/v1/admin/roles     # 新增角色
PUT    /api/v1/admin/roles/:id  # 更新角色
DELETE /api/v1/admin/roles/:id # 删除角色

# 权限管理
GET    /api/v1/admin/permissions                    # 所有权限列表
GET    /api/v1/admin/roles/:id/permissions          # 获取角色权限
POST   /api/v1/admin/roles/:id/permissions          # 更新角色权限

# 统计概览
GET    /api/v1/admin/stats/overview                 # 系统概览统计
```

详细接口文档请参考 `backend/internal/api/v1` 源码注释。

## 开发指南

### 后端开发

项目采用简洁的分层架构，所有业务逻辑直接在 API 处理器中实现：

1. **定义数据模型**：在 `internal/models/models.go` 中添加新的数据结构
2. **创建API处理器**：在 `internal/api/v1/` 中创建对应的处理器文件
3. **注册路由**：在 `internal/api/router.go` 中注册新的路由和权限
4. **权限控制**：使用 `AuthMiddleware()` 和 `PermissionMiddleware()` 进行认证和授权

**示例**：添加新功能模块
```go
// 1. 在 models.go 中定义模型
type NewFeature struct {
    gorm.Model
    Name string `json:"name"`
}

// 2. 在 v1/ 中创建 new_feature.go
func GetNewFeatures(c *gin.Context) {
    // 实现业务逻辑
}

// 3. 在 router.go 中注册路由
newFeatures := apiV1.Group("/new-features")
newFeatures.Use(middleware.AuthMiddleware())
{
    newFeatures.GET("", middleware.PermissionMiddleware("feature:read"), v1.GetNewFeatures)
}
```

### 前端开发

前端采用 Vue 3 Composition API + Element Plus 组件库：

1. **定义API接口**：在 `src/api/` 中创建对应的 API 封装文件
2. **创建页面组件**：在 `src/views/` 中创建 Vue 组件
3. **注册路由**：在 `src/router/index.js` 中添加路由配置
4. **状态管理**：如需全局状态，在 `src/store/` 中使用 Pinia

**示例**：添加新页面
```javascript
// 1. 在 api/ 中定义接口
export const getNewFeatures = (params) => {
    return request.get('/api/v1/new-features', { params })
}

// 2. 在 views/ 中创建组件
<template>
    <div>
        <!-- 使用 Element Plus 组件 -->
    </div>
</template>

// 3. 在 router/index.js 中注册
{
    path: '/new-features',
    component: () => import('@/views/NewFeatures.vue'),
    meta: { requiresAuth: true }
}
```

### 数据库迁移

项目使用 GORM 的自动迁移功能：

```go
// 在 config/database.go 中添加新模型
db.AutoMigrate(&models.NewFeature{})
```

## 项目特点

### 架构设计
- **前后端分离**：前端 Vue 3 + 后端 Go Gin，通过 RESTful API 通信
- **模块化设计**：清晰的代码组织结构，易于维护和扩展
- **数据库设计**：使用 GORM 进行 ORM 映射，支持自动迁移

### 安全性
- **JWT 认证**：基于 Token 的无状态认证机制
- **密码加密**：使用 bcrypt 算法加密存储密码
- **RBAC 权限控制**：细粒度的角色权限管理
- **权限中间件**：API 级别的权限验证

### 用户体验
- **响应式设计**：基于 Element Plus 的现代化 UI
- **数据可视化**：使用 ECharts 展示统计数据
- **分页和搜索**：支持大数据量的高效查询
- **表单验证**：完善的前后端数据验证

## 注意事项

1. **安全性**: 
   - 修改JWT密钥（config/config.ini）
   - 使用强密码
   - 生产环境启用HTTPS

2. **数据库**:
   - 定期备份数据
   - 合理设置索引
   - 监控性能

3. **部署**:
   - 使用环境变量管理配置
   - 配置反向代理（Nginx）
   - 设置日志轮转

## 待优化与改进

### 功能增强
- [ ] **数据导入导出**：支持 Excel 格式的批量导入导出
- [ ] **批量操作**：支持批量删除、批量修改等操作
- [ ] **高级搜索**：多条件组合搜索、模糊查询优化
- [ ] **操作日志**：记录用户操作历史，便于审计
- [ ] **数据备份**：自动化数据备份和恢复机制
- [ ] **系统设置**：可配置的系统参数管理

### 技术优化
- [ ] **单元测试**：为核心业务逻辑添加单元测试
- [ ] **集成测试**：API 接口的自动化测试
- [ ] **错误处理**：统一的错误处理和日志记录机制
- [ ] **日志系统**：结构化日志，支持日志级别和轮转
- [ ] **性能优化**：数据库查询优化、缓存机制
- [ ] **API文档**：使用 Swagger 生成 API 文档

### 部署与运维
- [ ] **Docker支持**：容器化部署方案
- [ ] **CI/CD流程**：自动化构建、测试和部署
- [ ] **监控告警**：系统性能监控和异常告警
- [ ] **负载均衡**：支持水平扩展

### 安全增强
- [ ] **API限流**：防止接口滥用
- [ ] **SQL注入防护**：参数化查询验证
- [ ] **XSS防护**：前端输入过滤和转义
- [ ] **CSRF防护**：跨站请求伪造防护

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。

---

**项目状态**: ✅ 核心功能已完成，持续优化中

**最后更新**: 2025年12月

## 技术支持

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发起 Pull Request
- 查看项目文档
