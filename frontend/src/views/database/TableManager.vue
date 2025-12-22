<template>
  <div class="table-manager">
    <el-container style="height: 100%;">
      <!-- 左侧表列表 -->
      <el-aside width="250px" class="table-list-aside">
        <div class="aside-header">
          <h3>数据库表</h3>
          <el-button @click="refreshTables" :icon="Refresh" circle size="small" />
        </div>
        
        <el-input
          v-model="searchTable"
          placeholder="搜索表名"
          :prefix-icon="Search"
          clearable
          class="search-input"
        />
        
        <el-menu
          :default-active="activeTable"
          @select="handleTableSelect"
          class="table-menu"
        >
          <el-menu-item
            v-for="table in filteredTables"
            :key="table.name"
            :index="table.name"
          >
            <el-icon><Document /></el-icon>
            <span>{{ table.label }}</span>
            <el-tag size="small" class="table-count">{{ table.count || 0 }}</el-tag>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <!-- 右侧表数据 -->
      <el-main class="table-content">
        <div v-if="!activeTable" class="empty-state">
          <el-icon :size="80" color="#909399"><Coin /></el-icon>
          <p>请从左侧选择一个表查看数据</p>
        </div>
        
        <div v-else class="table-data-container">
          <!-- 工具栏 -->
          <div class="toolbar">
            <div class="toolbar-left">
              <h2>{{ currentTableLabel }}</h2>
              <el-tag>{{ activeTable }}</el-tag>
            </div>
            <div class="toolbar-right">
              <el-button @click="handleRefresh" :icon="Refresh">刷新</el-button>
              <el-button @click="handleExport" :icon="Download">导出</el-button>
            </div>
          </div>
          
          <!-- 数据表格 -->
          <el-table
            :data="tableData"
            v-loading="loading"
            border
            stripe
            height="calc(100vh - 280px)"
          >
            <el-table-column type="index" label="#" width="60" align="center" />
            
            <el-table-column
              v-for="column in tableColumns"
              :key="column.prop"
              :prop="column.prop"
              :label="column.label"
              :width="column.width"
              :min-width="column.minWidth || 120"
              :formatter="column.formatter"
              show-overflow-tooltip
            >
              <template #default="{ row }" v-if="column.type === 'tag'">
                <el-tag :type="getTagType(row[column.prop])">
                  {{ row[column.prop] }}
                </el-tag>
              </template>
              <template #default="{ row }" v-else-if="column.type === 'datetime'">
                {{ formatDateTime(row[column.prop]) }}
              </template>
              <template #default="{ row }" v-else-if="column.type === 'boolean'">
                <el-tag :type="row[column.prop] ? 'success' : 'danger'">
                  {{ row[column.prop] ? '是' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="100" align="center" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="handleView(row)" :icon="View">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 分页 -->
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.page_size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="fetchTableData"
            @size-change="fetchTableData"
            class="pagination"
          />
        </div>
      </el-main>
    </el-container>
    
    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="数据详情"
      width="800px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item
          v-for="field in formFields"
          :key="field.prop"
          :label="field.label"
        >
          {{ formatValue(viewData[field.prop], field.type) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Coin, View, Refresh, Download, Search } from '@element-plus/icons-vue'
import { 
  getTableList, 
  getTableData,
  exportTableData 
} from '@/api/database'

// 状态管理
const loading = ref(false)
const viewDialogVisible = ref(false)
const searchTable = ref('')
const activeTable = ref('')
const tableData = ref([])
const viewData = ref({})
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 数据库表配置
const tables = ref([
  { name: 'users', label: '用户表', count: 0 },
  { name: 'roles', label: '角色表', count: 0 },
  { name: 'permissions', label: '权限表', count: 0 },
  { name: 'students', label: '学生表', count: 0 },
  { name: 'teachers', label: '教师表', count: 0 },
  { name: 'parents', label: '家长表', count: 0 },
  { name: 'classes', label: '班级表', count: 0 },
  { name: 'courses', label: '课程表', count: 0 },
  { name: 'enrollments', label: '选课表', count: 0 },
  { name: 'grades', label: '成绩表', count: 0 },
  { name: 'attendances', label: '考勤表', count: 0 },
  { name: 'reward_punishments', label: '奖惩表', count: 0 },
  { name: 'notifications', label: '通知表', count: 0 },
  { name: 'schedules', label: '课程表(排课)', count: 0 },
  { name: 'grade_audit_logs', label: '成绩审计日志', count: 0 },
  { name: 'vw_class_performance', label: '班级成绩视图', count: 0, isView: true },
  { name: 'vw_student_full_profile', label: '学生档案视图', count: 0, isView: true }
])

// 过滤后的表列表
const filteredTables = computed(() => {
  if (!searchTable.value) return tables.value
  return tables.value.filter(t => 
    t.name.toLowerCase().includes(searchTable.value.toLowerCase()) ||
    t.label.includes(searchTable.value)
  )
})

// 当前表标签
const currentTableLabel = computed(() => {
  const table = tables.value.find(t => t.name === activeTable.value)
  return table ? table.label : ''
})

// 表列配置 - 根据不同表动态生成
const tableColumns = computed(() => {
  const columnsMap = {
    users: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'username', label: '用户名', minWidth: 120 },
      { prop: 'role_id', label: '角色ID', width: 100 },
      { prop: 'user_type', label: '用户类型', width: 100, type: 'tag' },
      { prop: 'is_active', label: '状态', width: 100, type: 'boolean' },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    roles: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'role_name', label: '角色名称', minWidth: 120 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' },
      { prop: 'updated_at', label: '更新时间', width: 180, type: 'datetime' }
    ],
    permissions: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'name', label: '权限名称', minWidth: 150 },
      { prop: 'permission', label: '权限标识', minWidth: 150 },
      { prop: 'group', label: '权限分组', minWidth: 120 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    students: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'student_id', label: '学号', minWidth: 120 },
      { prop: 'name', label: '姓名', minWidth: 100 },
      { prop: 'gender', label: '性别', width: 80 },
      { prop: 'age', label: '年龄', width: 80 },
      { prop: 'email', label: '邮箱', minWidth: 150 },
      { prop: 'phone', label: '电话', minWidth: 120 },
      { prop: 'class_id', label: '班级ID', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    teachers: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'teacher_id', label: '工号', minWidth: 120 },
      { prop: 'name', label: '姓名', minWidth: 100 },
      { prop: 'email', label: '邮箱', minWidth: 150 },
      { prop: 'phone', label: '电话', minWidth: 120 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    parents: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'name', label: '姓名', minWidth: 100 },
      { prop: 'relation', label: '关系', minWidth: 100 },
      { prop: 'phone', label: '电话', minWidth: 120 },
      { prop: 'student_id', label: '学生ID', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    classes: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'class_name', label: '班级名称', minWidth: 150 },
      { prop: 'teacher_id', label: '班主任ID', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    courses: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'course_name', label: '课程名称', minWidth: 150 },
      { prop: 'teacher_id', label: '教师ID', width: 100 },
      { prop: 'credits', label: '学分', width: 80 },
      { prop: 'capacity', label: '容量', width: 80 },
      { prop: 'enrolled_count', label: '已选人数', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    enrollments: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'student_id', label: '学生ID', width: 100 },
      { prop: 'course_id', label: '课程ID', width: 100 },
      { prop: 'enrollment_date', label: '选课日期', width: 120, type: 'date' },
      { prop: 'status', label: '状态', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    grades: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'enrollment_id', label: '选课ID', width: 100 },
      { prop: 'score_type', label: '成绩类型', minWidth: 100 },
      { prop: 'score', label: '分数', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    attendances: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'student_id', label: '学生ID', width: 100 },
      { prop: 'date', label: '考勤日期', width: 120 },
      { prop: 'status', label: '状态', width: 100 },
      { prop: 'reason', label: '备注', minWidth: 150 },
      { prop: 'teacher_id', label: '记录人ID', width: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    reward_punishments: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'student_id', label: '学生ID', width: 100 },
      { prop: 'type', label: '类型', width: 100 },
      { prop: 'description', label: '描述', minWidth: 200 },
      { prop: 'date', label: '日期', width: 120 },
      { prop: 'issuer', label: '发布人', minWidth: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    notifications: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'title', label: '标题', minWidth: 150 },
      { prop: 'content', label: '内容', minWidth: 200 },
      { prop: 'sender_id', label: '发送人ID', width: 100 },
      { prop: 'target', label: '目标', width: 120 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    schedules: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'course_id', label: '课程ID', width: 100 },
      { prop: 'class_id', label: '班级ID', width: 100 },
      { prop: 'teacher_id', label: '教师ID', width: 100 },
      { prop: 'day_of_week', label: '星期', width: 80 },
      { prop: 'start_time', label: '开始时间', width: 100 },
      { prop: 'end_time', label: '结束时间', width: 100 },
      { prop: 'location', label: '上课地点', minWidth: 100 },
      { prop: 'semester', label: '学期', minWidth: 100 },
      { prop: 'created_at', label: '创建时间', width: 180, type: 'datetime' }
    ],
    grade_audit_logs: [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'grade_id', label: '成绩ID', width: 100 },
      { prop: 'old_score', label: '原分数', width: 100 },
      { prop: 'new_score', label: '新分数', width: 100 },
      { prop: 'modified_by', label: '修改人', minWidth: 100 },
      { prop: 'reason', label: '原因', minWidth: 200 },
      { prop: 'created_at', label: '修改时间', width: 180, type: 'datetime' }
    ]
  }
  
  return columnsMap[activeTable.value] || generateDefaultColumns()
})

// 表单字段配置
const formFields = computed(() => {
  const fieldsMap = {
    users: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'username', label: '用户名', type: 'text', span: 12 },
      { prop: 'password', label: '密码', type: 'text', span: 12 },
      { prop: 'role_id', label: '角色ID', type: 'number', span: 12 },
      { prop: 'user_type', label: '用户类型', type: 'select', span: 12, options: [
        { label: '管理员', value: 'admin' },
        { label: '教师', value: 'teacher' },
        { label: '学生', value: 'student' },
        { label: '家长', value: 'parent' }
      ]},
      { prop: 'is_active', label: '是否启用', type: 'boolean', span: 12 }
    ],
    roles: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'role_name', label: '角色名称', type: 'text', span: 12 }
    ],
    permissions: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'name', label: '权限名称', type: 'text', span: 12 },
      { prop: 'permission', label: '权限标识', type: 'text', span: 12 },
      { prop: 'group', label: '权限分组', type: 'text', span: 12 }
    ],
    students: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'student_id', label: '学号', type: 'text', span: 12 },
      { prop: 'name', label: '姓名', type: 'text', span: 12 },
      { prop: 'gender', label: '性别', type: 'select', span: 12, options: [
        { label: '男', value: '男' },
        { label: '女', value: '女' }
      ]},
      { prop: 'age', label: '年龄', type: 'number', span: 12 },
      { prop: 'email', label: '邮箱', type: 'text', span: 12 },
      { prop: 'phone', label: '电话', type: 'text', span: 12 },
      { prop: 'class_id', label: '班级ID', type: 'number', span: 12 },
      { prop: 'address', label: '地址', type: 'textarea', span: 24 }
    ],
    teachers: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'teacher_id', label: '工号', type: 'text', span: 12 },
      { prop: 'name', label: '姓名', type: 'text', span: 12 },
      { prop: 'email', label: '邮箱', type: 'text', span: 12 },
      { prop: 'phone', label: '电话', type: 'text', span: 12 },
      { prop: 'department', label: '部门', type: 'text', span: 12 }
    ],
    parents: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'name', label: '姓名', type: 'text', span: 12 },
      { prop: 'relation', label: '关系', type: 'select', span: 12, options: [
        { label: '父亲', value: '父亲' },
        { label: '母亲', value: '母亲' },
        { label: '监护人', value: '监护人' }
      ]},
      { prop: 'phone', label: '电话', type: 'text', span: 12 },
      { prop: 'student_id', label: '学生ID', type: 'number', span: 12 }
    ],
    classes: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'class_name', label: '班级名称', type: 'text', span: 12 },
      { prop: 'teacher_id', label: '班主任ID', type: 'number', span: 12 },
      { prop: 'grade_level', label: '年级', type: 'number', span: 12 }
    ],
    courses: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'course_name', label: '课程名称', type: 'text', span: 12 },
      { prop: 'teacher_id', label: '教师ID', type: 'number', span: 12 },
      { prop: 'credits', label: '学分', type: 'number', span: 12, precision: 1 },
      { prop: 'capacity', label: '容量', type: 'number', span: 12 },
      { prop: 'description', label: '课程描述', type: 'textarea', span: 24 }
    ],
    enrollments: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'student_id', label: '学生ID', type: 'number', span: 12 },
      { prop: 'course_id', label: '课程ID', type: 'number', span: 12 },
      { prop: 'enrollment_date', label: '选课日期', type: 'date', span: 12 },
      { prop: 'status', label: '状态', type: 'select', span: 12, options: [
        { label: '正常', value: '正常' },
        { label: '退课', value: '退课' }
      ]}
    ],
    grades: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'enrollment_id', label: '选课ID', type: 'number', span: 12 },
      { prop: 'score_type', label: '成绩类型', type: 'select', span: 12, options: [
        { label: '平时成绩', value: '平时成绩' },
        { label: '期中成绩', value: '期中成绩' },
        { label: '期末成绩', value: '期末成绩' }
      ]},
      { prop: 'score', label: '分数', type: 'number', span: 12, min: 0, max: 100, precision: 1 }
    ],
    attendances: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'student_id', label: '学生ID', type: 'number', span: 12 },
      { prop: 'date', label: '考勤日期', type: 'text', span: 12 },
      { prop: 'status', label: '状态', type: 'select', span: 12, options: [
        { label: '出勤', value: '出勤' },
        { label: '迟到', value: '迟到' },
        { label: '早退', value: '早退' },
        { label: '缺勤', value: '缺勤' },
        { label: '请假', value: '请假' }
      ]},
      { prop: 'reason', label: '备注', type: 'textarea', span: 24 },
      { prop: 'teacher_id', label: '记录人ID', type: 'number', span: 12 }
    ],
    reward_punishments: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'student_id', label: '学生ID', type: 'number', span: 12 },
      { prop: 'type', label: '类型', type: 'select', span: 12, options: [
        { label: '奖励', value: '奖励' },
        { label: '惩罚', value: '处分' }
      ]},
      { prop: 'description', label: '描述', type: 'textarea', span: 24 },
      { prop: 'date', label: '日期', type: 'text', span: 12 },
      { prop: 'issuer', label: '发布人', type: 'text', span: 12 }
    ],
    notifications: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'title', label: '标题', type: 'text', span: 12 },
      { prop: 'content', label: '内容', type: 'textarea', span: 24 },
      { prop: 'sender_id', label: '发送人ID', type: 'number', span: 12 },
      { prop: 'target', label: '目标', type: 'text', span: 12 }
    ],
    schedules: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'course_id', label: '课程ID', type: 'number', span: 12 },
      { prop: 'class_id', label: '班级ID', type: 'number', span: 12 },
      { prop: 'teacher_id', label: '教师ID', type: 'number', span: 12 },
      { prop: 'day_of_week', label: '星期', type: 'select', span: 12, options: [
        { label: '星期一', value: 1 },
        { label: '星期二', value: 2 },
        { label: '星期三', value: 3 },
        { label: '星期四', value: 4 },
        { label: '星期五', value: 5 }
      ]},
      { prop: 'start_time', label: '开始时间', type: 'text', span: 12 },
      { prop: 'end_time', label: '结束时间', type: 'text', span: 12 },
      { prop: 'location', label: '上课地点', type: 'text', span: 12 },
      { prop: 'semester', label: '学期', type: 'text', span: 12 }
    ],
    grade_audit_logs: [
      { prop: 'id', label: 'ID', type: 'number', disabled: true, span: 12 },
      { prop: 'grade_id', label: '成绩ID', type: 'number', span: 12 },
      { prop: 'old_score', label: '原分数', type: 'number', span: 12, precision: 1 },
      { prop: 'new_score', label: '新分数', type: 'number', span: 12, precision: 1 },
      { prop: 'modified_by', label: '修改人', type: 'text', span: 12 },
      { prop: 'reason', label: '原因', type: 'textarea', span: 24 }
    ]
  }
  
  return fieldsMap[activeTable.value] || generateDefaultFields()
})

// 生成默认字段配置
const generateDefaultFields = () => {
  if (tableData.value.length === 0) return []
  
  const firstRow = tableData.value[0]
  return Object.keys(firstRow).map(key => ({
    prop: key,
    label: key,
    type: 'text',
    span: 12,
    disabled: key === 'ID' || key === 'id' || key.toLowerCase().includes('createdat') || key.toLowerCase().includes('updatedat')
  }))
}

// 生成默认列配置
const generateDefaultColumns = () => {
  if (tableData.value.length === 0) return []
  
  const firstRow = tableData.value[0]
  return Object.keys(firstRow).map(key => ({
    prop: key,
    label: key,
    minWidth: 120
  }))
}

// 格式化日期时间
const formatDateTime = (value) => {
  if (!value) return '-'
  const date = new Date(value)
  return date.toLocaleString('zh-CN')
}

// 格式化值
const formatValue = (value, type) => {
  if (value === null || value === undefined) return '-'
  if (type === 'datetime') return formatDateTime(value)
  if (type === 'boolean') return value ? '是' : '否'
  if (typeof value === 'object') return JSON.stringify(value, null, 2)
  return value
}

// 获取标签类型
const getTagType = (value) => {
  const typeMap = {
    'admin': 'danger',
    'teacher': 'warning',
    'student': 'success',
    'parent': 'info'
  }
  return typeMap[value] || 'info'
}

// 刷新表列表
const refreshTables = async () => {
  try {
    const response = await getTableList()
    if (response.data?.tables) {
      tables.value = response.data.tables
    }
    ElMessage.success('刷新成功')
  } catch (error) {
    ElMessage.error('刷新失败: ' + error.message)
  }
}

// 选择表
const handleTableSelect = (tableName) => {
  activeTable.value = tableName
  pagination.page = 1
  fetchTableData()
}

// 获取表数据
const fetchTableData = async () => {
  if (!activeTable.value) return
  
  loading.value = true
  try {
    const response = await getTableData(activeTable.value, {
      page: pagination.page,
      page_size: pagination.page_size
    })
    
    tableData.value = response.data?.list || []
    pagination.total = response.data?.total || 0
    
    ElMessage.success('数据加载成功')
  } catch (error) {
    ElMessage.error('获取数据失败: ' + error.message)
    tableData.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

// 刷新当前表
const handleRefresh = () => {
  fetchTableData()
}

// 查看
const handleView = (row) => {
  viewData.value = { ...row }
  viewDialogVisible.value = true
}

// 导出数据
const handleExport = async () => {
  try {
    const response = await exportTableData(activeTable.value, {
      page: pagination.page,
      page_size: pagination.page_size
    })
    
    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `${activeTable.value}_${Date.now()}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message)
  }
}

onMounted(() => {
  refreshTables()
})
</script>

<style scoped>
.table-manager {
  height: calc(100vh - 120px);
  background: white;
  border-radius: 4px;
  overflow: hidden;
}

.table-list-aside {
  border-right: 1px solid #e6e6e6;
  background: #fafafa;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.aside-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e6e6e6;
  flex-shrink: 0;
}

.aside-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.search-input {
  margin: 12px;
  width: calc(100% - 24px);
  flex-shrink: 0;
}

.table-menu {
  border: none;
  background: transparent;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 美化滚动条 */
.table-menu::-webkit-scrollbar {
  width: 6px;
}

.table-menu::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.table-menu::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.table-menu::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.table-menu :deep(.el-menu-item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

.table-count {
  margin-left: auto;
}

.table-content {
  padding: 0;
  background: white;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.empty-state p {
  margin-top: 16px;
  font-size: 14px;
}

.table-data-container {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e6e6e6;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-left h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.data-form {
  max-height: 60vh;
  overflow-y: auto;
}
</style>
