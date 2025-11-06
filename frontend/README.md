# 学生管理系统 - 前端

基于 Vue 3 + Vite + Element Plus 的学生管理系统前端应用。

## 技术栈

- **Vue 3**: 渐进式 JavaScript 框架
- **Vite**: 下一代前端构建工具
- **Vue Router**: 官方路由管理器
- **Pinia**: Vue 3 状态管理库
- **Axios**: HTTP 客户端
- **Element Plus**: Vue 3 UI 组件库

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口
│   │   ├── auth.js       # 认证相关接口
│   │   └── student.js    # 学生管理接口
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   ├── layouts/          # 布局组件
│   │   └── MainLayout.vue
│   ├── router/           # 路由配置
│   │   └── index.js
│   ├── store/            # 状态管理
│   │   └── user.js
│   ├── utils/            # 工具函数
│   │   └── request.js    # Axios 封装
│   ├── views/            # 页面视图
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   ├── Students.vue
│   │   └── ...
│   ├── App.vue
│   └── main.js
├── index.html
├── package.json
└── vite.config.js
```

## 安装和运行

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

应用将在 `http://localhost:5173` 启动。

### 3. 构建生产版本

```bash
npm run build
```

构建输出将在 `dist/` 目录下。

### 4. 预览生产构建

```bash
npm run preview
```

## 功能模块

### 已实现

- ✅ 用户登录/登出
- ✅ 路由守卫
- ✅ 主布局（侧边栏 + 顶部导航）
- ✅ 数据看板（Dashboard）
- ✅ 学生管理（列表、添加、编辑、删除、搜索、分页）

### 开发中

- 🚧 班级管理
- 🚧 课程管理
- 🚧 成绩管理
- 🚧 考勤管理
- 🚧 奖惩管理
- 🚧 通知管理

## 配置说明

### API 代理

在 `vite.config.js` 中配置了开发服务器代理，将 `/api` 请求转发到后端服务器：

```javascript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

### 请求拦截

在 `src/utils/request.js` 中配置了请求和响应拦截器：

- 请求拦截：自动添加 JWT token
- 响应拦截：统一处理错误和消息提示

## 默认账号

- 用户名：`admin`
- 密码：`admin123`

（注意：需要先在后端创建该账号）

## 开发说明

### 添加新页面

1. 在 `src/views/` 下创建 Vue 组件
2. 在 `src/router/index.js` 中注册路由
3. 在 `src/layouts/MainLayout.vue` 中添加菜单项

### 添加新 API

1. 在 `src/api/` 下创建对应的 API 文件
2. 使用 `src/utils/request.js` 发起请求

### 状态管理

使用 Pinia 进行状态管理，示例：

```javascript
import { defineStore } from 'pinia'

export const useMyStore = defineStore('my-store', {
  state: () => ({ ... }),
  getters: { ... },
  actions: { ... },
})
```

## 注意事项

1. 确保后端服务已启动（默认 `http://localhost:8080`）
2. 首次运行需要安装依赖
3. 开发时浏览器可能会显示CORS错误，这是正常的，代理会处理
4. 生产环境需要配置实际的后端 API 地址

## 浏览器支持

- Chrome (推荐)
- Firefox
- Safari
- Edge

## 许可证

MIT
