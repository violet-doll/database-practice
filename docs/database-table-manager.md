# 数据库表管理系统

## 概述

前端已从"学生管理系统"模式改为通用的"数据库表管理"模式,允许用户直接查看和管理数据库中的所有表。

## 主要功能

### 1. 数据库表浏览
- 左侧菜单显示所有数据库表
- 支持搜索表名
- 显示每个表的记录数量
- 包含普通表和视图(如班级成绩视图、学生档案视图)

### 2. 表数据管理
- **查看**: 以表格形式展示表数据
- **新增**: 添加新记录
- **编辑**: 修改现有记录
- **删除**: 删除记录
- **导出**: 导出表数据为 Excel 文件

### 3. 分页和过滤
- 支持分页浏览
- 可调整每页显示数量(10/20/50/100)
- 显示总记录数

### 4. 智能表单
- 根据不同表自动生成表单字段
- 支持多种字段类型:
  - 文本输入
  - 数字输入
  - 下拉选择
  - 日期/日期时间选择
  - 开关(布尔值)
  - 文本域

## 支持的表

1. **用户管理**
   - users (用户表)
   - roles (角色表)
   - permissions (权限表)

2. **学籍管理**
   - students (学生表)
   - teachers (教师表)
   - parents (家长表)
   - classes (班级表)

3. **教务管理**
   - courses (课程表)
   - enrollments (选课表)
   - grades (成绩表)
   - schedules (排课表)

4. **日常管理**
   - attendances (考勤表)
   - reward_punishments (奖惩表)

5. **通知管理**
   - notifications (通知表)

6. **审计日志**
   - grade_audit_logs (成绩审计日志)

7. **数据库视图**
   - vw_class_performance (班级成绩统计视图)
   - vw_student_full_profile (学生完整档案视图)

## API 接口

前端已创建通用的数据库 API (`/api/database.js`),包含以下接口:

```javascript
// 获取所有表列表
getTableList()

// 获取表数据
getTableData(tableName, { page, page_size })

// 获取表结构
getTableSchema(tableName)

// 创建记录
createTableData(tableName, data)

// 更新记录
updateTableData(tableName, id, data)

// 删除记录
deleteTableData(tableName, id)

// 批量删除
batchDeleteTableData(tableName, ids)

// 导出数据
exportTableData(tableName, params)

// 执行 SQL
executeSQL(sql)
```

## 后端需要实现的接口

为了让前端正常工作,后端需要实现以下 RESTful API:

### 1. 获取表列表
```
GET /api/database/tables
Response: {
  tables: [
    { name: "users", label: "用户表", count: 10 },
    { name: "students", label: "学生表", count: 50 },
    ...
  ]
}
```

### 2. 获取表数据
```
GET /api/database/tables/:tableName?page=1&page_size=20
Response: {
  list: [...],
  total: 100
}
```

### 3. 获取表结构
```
GET /api/database/tables/:tableName/schema
Response: {
  columns: [
    { name: "id", type: "int", nullable: false },
    { name: "name", type: "varchar", nullable: false },
    ...
  ]
}
```

### 4. 创建记录
```
POST /api/database/tables/:tableName
Body: { field1: value1, field2: value2, ... }
```

### 5. 更新记录
```
PUT /api/database/tables/:tableName/:id
Body: { field1: value1, field2: value2, ... }
```

### 6. 删除记录
```
DELETE /api/database/tables/:tableName/:id
```

### 7. 导出数据
```
GET /api/database/tables/:tableName/export
Response: Excel 文件流
```

## 使用方法

1. 登录系统后,默认进入"数据库表管理"页面
2. 从左侧菜单选择要查看的表
3. 使用工具栏按钮进行新增、刷新、导出操作
4. 点击表格中的"查看"、"编辑"、"删除"按钮管理记录

## 注意事项

1. 当前 API 接口尚未在后端实现,需要后端开发配合
2. 表单字段配置目前只针对部分表进行了定制,其他表会使用默认配置
3. 导出功能需要后端支持生成 Excel 文件
4. 视图(如 vw_class_performance)应该是只读的,不支持新增、编辑、删除操作

## 下一步工作

1. 后端实现通用的数据库表管理 API
2. 完善所有表的字段配置
3. 添加数据验证规则
4. 实现高级搜索和过滤功能
5. 添加批量操作功能
6. 实现 SQL 查询编辑器
