# 学生管理系统

这是一个基于 Go + Gin + Gorm + Vue 3 + Element Plus 的现代化学生管理系统，是《软件工程》课程的实验项目。

## 项目概述

本系统旨在为学校或教育机构提供一个高效、全面的学生信息管理平台，涵盖学籍管理、课程管理、成绩管理、考勤管理、奖惩管理等多个功能模块。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: Gorm
- **数据库**: MySQL 9.2
- **认证**: JWT (JSON Web Tokens)
- **密码加密**: bcrypt

### 前端
- **框架**: Vue 3
- **构建工具**: Vite
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **HTTP客户端**: Axios
- **UI组件库**: Element Plus

## 项目结构

```
student-management-system/
├── backend/                # 后端服务
│   ├── cmd/               # 程序入口
│   ├── config/            # 配置文件
│   ├── internal/          # 内部代码
│   │   ├── api/          # API处理器
│   │   ├── middleware/   # 中间件
│   │   ├── models/       # 数据模型
│   │   ├── repository/   # 数据仓库层
│   │   ├── service/      # 业务逻辑层
│   │   └── utils/        # 工具函数
│   └── go.mod
├── frontend/              # 前端应用
│   ├── src/
│   │   ├── api/          # API接口
│   │   ├── components/   # 公共组件
│   │   ├── layouts/      # 布局组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理
│   │   ├── utils/        # 工具函数
│   │   └── views/        # 页面视图
│   └── package.json
├── docs/                  # 软件工程文档
└── README.md
```

## 功能模块

### ✅ 已实现
- [x] 用户认证（登录/登出）
- [x] 学生信息管理（CRUD）
- [x] 课程管理（CRUD）
- [x] 成绩管理（录入、按学生/课程查询）
- [x] 数据看板
- [x] 基于角色的权限控制（RBAC）

### 🚧 开发中
- [ ] 班级管理
- [ ] 教师管理
- [ ] 考勤管理
- [ ] 奖惩管理
- [ ] 通知管理
- [ ] 家长联系管理
- [ ] 数据统计和分析

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

### 3. 启动后端

```bash
cd backend

# 复制配置文件
copy .env.example .env

# 修改 .env 或 config/config.ini 中的数据库配置

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

首次运行需要手动创建管理员账号。可以通过以下方式：

1. 直接在数据库中插入记录
2. 使用API创建（需要先临时关闭认证中间件）
3. 使用数据库初始化脚本

示例SQL：

```sql
-- 插入管理员用户（密码: admin123，已使用bcrypt加密）
INSERT INTO users (username, password, role_id, is_active, user_type, created_at, updated_at)
VALUES ('admin', '$2a$10$...(bcrypt hash)...', 1, 1, 'admin', NOW(), NOW());
```

## 默认账号

- **用户名**: admin
- **密码**: admin123

（需要先创建该账号）

## API文档

后端API接口文档：

- 基础路径: `/api/v1`
- 认证方式: Bearer Token (JWT)

主要接口（节选）：

```
/api/v1
  /courses
    GET    /            # 课程列表（分页、搜索）
    GET    /:id         # 课程详情
    POST   /            # 新增课程
    PUT    /:id         # 更新课程
    DELETE /:id         # 删除课程

  /grades
    POST   /            # 录入成绩（自动创建选课关系）
    GET    /student/:id # 按学生查询成绩（返回选课+成绩明细）
    GET    /course/:id  # 按课程查询成绩（返回选课+成绩明细）
```

详细接口文档请参考 `backend/internal/api/v1` 源码注释。

## 开发指南

### 后端开发

1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/api/v1/` 中创建API处理器
3. 在 `internal/api/router.go` 中注册路由
4. 使用中间件进行权限控制

### 前端开发

1. 在 `src/api/` 中定义API接口
2. 在 `src/views/` 中创建页面组件
3. 在 `src/router/` 中注册路由
4. 使用Pinia进行状态管理

## 软件工程文档

所有软件工程相关文档存放在 `docs/` 目录下：

- 需求分析（含用例图）
- 概要设计（含系统架构图）
- 详细设计（含序列图、活动图、类图）
- 测试报告
- 项目总结
- 原型系统截图

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

## 待优化

- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 完善错误处理
- [ ] 添加日志系统
- [ ] 性能优化
- [ ] 代码文档
- [ ] Docker支持
- [ ] CI/CD流程

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

**项目状态**: 🚧 开发中

**最后更新**: 2025年11月6日
