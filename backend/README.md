# 学生管理系统 - 后端

基于 Go + Gin + Gorm 的学生管理系统后端服务。

## 技术栈

- **Go**: 1.21+
- **Web框架**: Gin
- **ORM**: Gorm
- **数据库**: MySQL 9.2
- **认证**: JWT (JSON Web Tokens)
- **密码加密**: bcrypt

## 项目结构

```
backend/
├── cmd/
│   └── main.go              # 程序入口
├── config/
│   ├── config.ini           # 配置文件
│   └── database.go          # 数据库配置
├── internal/
│   ├── api/
│   │   ├── v1/              # API v1 版本
│   │   │   ├── auth.go      # 认证相关
│   │   │   └── student.go   # 学生管理
│   │   └── router.go        # 路由配置
│   ├── middleware/
│   │   ├── auth.go          # 认证中间件
│   │   └── cors.go          # CORS中间件
│   ├── models/
│   │   └── models.go        # 数据模型
│   ├── repository/          # 数据仓库层
│   ├── service/             # 业务逻辑层
│   └── utils/
│       ├── jwt.go           # JWT工具
│       └── password.go      # 密码加密工具
└── go.mod
```

## 安装和运行

### 1. 安装依赖

```bash
cd backend
go mod tidy
```

### 2. 配置数据库

复制并编辑配置文件：

```bash
copy .env.example .env
```

修改 `.env` 或 `config/config.ini` 中的数据库配置：

```ini
[Database]
Type = mysql
Host = 127.0.0.1
Port = 3306
User = root
Password = your_password
DBName = student_db
```

### 3. 创建数据库

```sql
CREATE DATABASE student_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行项目

```bash
go run ./cmd/main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 文档

### 认证相关

#### 登录
- **POST** `/api/v1/auth/login`
- 请求体:
  ```json
  {
    "username": "admin",
    "password": "password"
  }
  ```
- 响应:
  ```json
  {
    "code": 200,
    "message": "登录成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": {...}
    }
  }
  ```

#### 获取当前用户信息
- **GET** `/api/v1/auth/me`
- Headers: `Authorization: Bearer {token}`

#### 登出
- **POST** `/api/v1/auth/logout`
- Headers: `Authorization: Bearer {token}`

### 学生管理

#### 获取学生列表
- **GET** `/api/v1/students`
- 查询参数: `page`, `page_size`, `name`, `student_id`, `class_id`
- Headers: `Authorization: Bearer {token}`

#### 获取学生详情
- **GET** `/api/v1/students/:id`
- Headers: `Authorization: Bearer {token}`

#### 创建学生
- **POST** `/api/v1/students`
- Headers: `Authorization: Bearer {token}`
- 请求体:
  ```json
  {
    "name": "张三",
    "student_id": "2023001",
    "gender": "男",
    "age": 18,
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "class_id": 1
  }
  ```

#### 更新学生信息
- **PUT** `/api/v1/students/:id`
- Headers: `Authorization: Bearer {token}`

#### 删除学生
- **DELETE** `/api/v1/students/:id`
- Headers: `Authorization: Bearer {token}`

## 默认数据

系统启动时会自动创建以下默认角色：

- `admin` - 管理员
- `teacher` - 教师
- `student` - 学生
- `parent` - 家长

## 开发说明

### 添加新的API接口

1. 在 `internal/api/v1/` 下创建对应的处理器文件
2. 在 `internal/api/router.go` 中注册路由
3. 如需权限控制，使用 `middleware.AuthMiddleware()` 和 `middleware.RoleMiddleware()`

### 数据库迁移

Gorm 会在启动时自动迁移数据表结构。如果修改了模型，重启服务即可自动更新表结构。

## 待实现功能

- [ ] 班级管理 API
- [ ] 课程管理 API
- [ ] 成绩管理 API
- [ ] 考勤管理 API
- [ ] 奖惩管理 API
- [ ] 通知管理 API
- [ ] 用户管理 API (Admin)
- [ ] 数据统计和分析 API

## 注意事项

1. 修改 JWT 密钥（在 `config/config.ini` 和 `cmd/main.go` 中）
2. 使用环境变量管理敏感配置
3. 生产环境建议使用 HTTPS
4. 定期备份数据库
