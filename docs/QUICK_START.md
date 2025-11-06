# 学生管理系统 - 快速启动指南

本指南将帮助您快速启动和运行学生管理系统。

## 第一步：准备环境

### 1.1 安装必要软件

确保您的系统已安装以下软件：

- **Go 1.21+**: [下载地址](https://golang.org/dl/)
- **Node.js 18+**: [下载地址](https://nodejs.org/)
- **MySQL 8.0+**: [下载地址](https://dev.mysql.com/downloads/mysql/)

### 1.2 验证安装

打开终端（PowerShell），运行以下命令验证：

```powershell
go version
node --version
npm --version
mysql --version
```

## 第二步：配置数据库

### 2.1 启动 MySQL 服务

确保 MySQL 服务正在运行。

### 2.2 创建数据库

打开 MySQL 命令行或使用图形化工具（如 Navicat、phpMyAdmin），执行：

```sql
CREATE DATABASE student_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2.3 配置数据库连接

编辑后端配置文件：

**方式一：使用 .env 文件**

1. 复制 `.env.example` 为 `.env`：
   ```powershell
   cd backend
   copy .env.example .env
   ```

2. 编辑 `.env` 文件，修改数据库配置：
   ```
   DB_USER=root
   DB_PASSWORD=你的MySQL密码
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=student_db
   ```

**方式二：直接修改 config.ini**

编辑 `backend/config/config.ini`：

```ini
[Database]
Type = mysql
Host = 127.0.0.1
Port = 3306
User = root
Password = 你的MySQL密码
DBName = student_db
```

## 第三步：启动后端服务

### 3.1 进入后端目录

```powershell
cd backend
```

### 3.2 安装依赖

```powershell
go mod tidy
```

### 3.3 启动服务

```powershell
go run ./cmd/main.go
```

看到以下输出表示启动成功：

```
数据库连接成功
数据库表迁移成功
默认角色初始化完成
服务器启动成功，监听端口 :8080
```

**注意**：首次启动时，Gorm 会自动创建所有数据表。

## 第四步：创建管理员账号

保持后端服务运行，打开**新的**终端窗口：

### 4.1 运行创建管理员脚本

```powershell
cd backend
go run ./scripts/create_admin/main.go
```

成功后会显示：

```
✓ 管理员账号创建成功!
  用户名: admin
  密码: admin123
```

## 第五步：启动前端应用

### 5.1 进入前端目录

打开**新的**终端窗口：

```powershell
cd frontend
```

### 5.2 安装依赖（如果还没安装）

```powershell
npm install
```

### 5.3 启动开发服务器

```powershell
npm run dev
```

看到以下输出表示启动成功：

```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

## 第六步：访问系统

### 6.1 打开浏览器

访问：http://localhost:5173

### 6.2 登录系统

使用以下账号登录：

- **用户名**: `admin`
- **密码**: `admin123`

### 6.3 开始使用

登录成功后，您可以：

1. 查看数据看板
2. 管理学生信息
3. 浏览其他功能模块

## 常见问题

### Q1: 后端启动失败，提示"数据库连接失败"

**解决方案**：
1. 检查 MySQL 服务是否启动
2. 确认数据库配置信息是否正确
3. 确认数据库 `student_db` 已创建
4. 检查用户名和密码是否正确

### Q2: 前端无法连接后端

**解决方案**：
1. 确认后端服务已启动（http://localhost:8080）
2. 检查浏览器控制台是否有 CORS 错误
3. 确认 `vite.config.js` 中的代理配置正确

### Q3: 创建管理员账号失败

**解决方案**：
1. 确保后端服务已经启动过一次（创建表结构）
2. 确保数据库连接配置正确
3. 检查是否已经创建过管理员账号

### Q4: npm install 失败

**解决方案**：
1. 检查网络连接
2. 尝试使用国内镜像源：
   ```powershell
   npm config set registry https://registry.npmmirror.com
   ```
3. 删除 `node_modules` 和 `package-lock.json` 后重试

### Q5: PowerShell 执行策略错误

如果遇到"无法加载文件，因为在此系统上禁止运行脚本"错误：

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 停止服务

### 停止后端

在后端终端按 `Ctrl + C`

### 停止前端

在前端终端按 `Ctrl + C`

## 下一步

1. 熟悉系统界面和功能
2. 添加测试数据
3. 开发其他功能模块
4. 查看开发文档

## 获取帮助

如有问题，请查看：

1. 项目 README.md
2. 后端 README: `backend/README.md`
3. 前端 README: `frontend/README.md`
4. 提交 Issue

---

祝您使用愉快！🎉
