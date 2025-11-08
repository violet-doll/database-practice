<template>
  <div class="schedule">
    <el-card>
      <div class="toolbar">
        <div class="search-box">
          <el-input
            v-model="searchForm.class_id"
            placeholder="按班级ID过滤"
            style="width: 180px; margin-right: 10px"
            clearable
          />
          <el-input
            v-model="searchForm.teacher_id"
            placeholder="按教师ID过滤"
            style="width: 180px; margin-right: 10px"
            clearable
          />
          <el-input
            v-model="searchForm.semester"
            placeholder="按学期过滤"
            style="width: 180px; margin-right: 10px"
            clearable
          />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <el-button type="primary" @click="handleAdd">新增排课</el-button>
      </div>

      <el-table :data="scheduleList" style="width: 100%; margin-top: 20px" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column label="课程" min-width="150">
          <template #default="{ row }">
            <div v-if="row.Course">
              <div>{{ row.Course.course_name }}</div>
              <div style="font-size: 12px; color: #909399">ID: {{ row.course_id }}</div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="班级" min-width="120">
          <template #default="{ row }">
            <div v-if="row.Class">
              <div>{{ row.Class.class_name }}</div>
              <div style="font-size: 12px; color: #909399">ID: {{ row.class_id }}</div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="教师" min-width="120">
          <template #default="{ row }">
            <div v-if="row.Teacher">
              <div>{{ row.Teacher.name }}</div>
              <div style="font-size: 12px; color: #909399">ID: {{ row.teacher_id }}</div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="星期" width="100">
          <template #default="{ row }">
            {{ getDayOfWeekText(row.day_of_week) }}
          </template>
        </el-table-column>
        <el-table-column label="时间" width="150">
          <template #default="{ row }">
            <div>{{ row.start_time }} - {{ row.end_time }}</div>
          </template>
        </el-table-column>
        <el-table-column label="地点" width="120">
          <template #default="{ row }">{{ row.location || '-' }}</template>
        </el-table-column>
        <el-table-column label="学期" width="120">
          <template #default="{ row }">{{ row.semester || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="180">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)" :loading="deletingId === row.ID">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="scheduleFormRef" :model="scheduleForm" :rules="scheduleRules" label-width="100px">
        <el-form-item label="课程ID" prop="course_id">
          <el-input-number v-model="scheduleForm.course_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="班级ID" prop="class_id">
          <el-input-number v-model="scheduleForm.class_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="教师ID" prop="teacher_id">
          <el-input-number v-model="scheduleForm.teacher_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="星期" prop="day_of_week">
          <el-select v-model="scheduleForm.day_of_week" style="width: 100%">
            <el-option label="周一" :value="1" />
            <el-option label="周二" :value="2" />
            <el-option label="周三" :value="3" />
            <el-option label="周四" :value="4" />
            <el-option label="周五" :value="5" />
            <el-option label="周六" :value="6" />
            <el-option label="周日" :value="7" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间" prop="start_time">
          <el-input v-model="scheduleForm.start_time" placeholder="例如: 08:00 或 1-2节" />
        </el-form-item>
        <el-form-item label="结束时间" prop="end_time">
          <el-input v-model="scheduleForm.end_time" placeholder="例如: 09:40" />
        </el-form-item>
        <el-form-item label="上课地点" prop="location">
          <el-input v-model="scheduleForm.location" placeholder="例如: 教5-101" />
        </el-form-item>
        <el-form-item label="学期" prop="semester">
          <el-input v-model="scheduleForm.semester" placeholder="例如: 2025-Fall" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSchedules, createSchedule, updateSchedule, deleteSchedule } from '@/api/schedule'

const loading = ref(false)
const scheduleList = ref([])
const deletingId = ref(null)
const submitting = ref(false)

const searchForm = reactive({
  class_id: '',
  teacher_id: '',
  semester: '',
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0,
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增排课')
const scheduleFormRef = ref(null)
const scheduleForm = reactive({
  ID: null,
  course_id: null,
  class_id: null,
  teacher_id: null,
  day_of_week: 1,
  start_time: '',
  end_time: '',
  location: '',
  semester: '',
})

const scheduleRules = {
  course_id: [{ required: true, message: '请输入课程ID', trigger: 'blur' }],
  class_id: [{ required: true, message: '请输入班级ID', trigger: 'blur' }],
  teacher_id: [{ required: true, message: '请输入教师ID', trigger: 'blur' }],
  day_of_week: [{ required: true, message: '请选择星期', trigger: 'change' }],
  start_time: [{ required: true, message: '请输入开始时间', trigger: 'blur' }],
  end_time: [{ required: true, message: '请输入结束时间', trigger: 'blur' }],
}

const getDayOfWeekText = (day) => {
  const days = ['', '周一', '周二', '周三', '周四', '周五', '周六', '周日']
  return days[day] || '-'
}

const fetchSchedules = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
    }
    if (searchForm.class_id) {
      params.class_id = searchForm.class_id
    }
    if (searchForm.teacher_id) {
      params.teacher_id = searchForm.teacher_id
    }
    if (searchForm.semester) {
      params.semester = searchForm.semester
    }
    const res = await getSchedules(params)
    scheduleList.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('获取排课列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchSchedules()
}

const handleReset = () => {
  searchForm.class_id = ''
  searchForm.teacher_id = ''
  searchForm.semester = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增排课'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑排课'
  Object.assign(scheduleForm, {
    ID: row.ID,
    course_id: row.course_id,
    class_id: row.class_id,
    teacher_id: row.teacher_id,
    day_of_week: row.day_of_week,
    start_time: row.start_time,
    end_time: row.end_time,
    location: row.location || '',
    semester: row.semester || '',
  })
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该排课记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      deletingId.value = row.ID
      try {
        await deleteSchedule(row.ID)
        ElMessage.success('删除成功')
        fetchSchedules()
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '删除失败')
      } finally {
        deletingId.value = null
      }
    })
    .catch(() => {
      // 用户取消
    })
}

const handleSubmit = async () => {
  if (!scheduleFormRef.value) return
  await scheduleFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const payload = {
        course_id: Number(scheduleForm.course_id),
        class_id: Number(scheduleForm.class_id),
        teacher_id: Number(scheduleForm.teacher_id),
        day_of_week: Number(scheduleForm.day_of_week),
        start_time: scheduleForm.start_time,
        end_time: scheduleForm.end_time,
        location: scheduleForm.location,
        semester: scheduleForm.semester,
      }
      if (scheduleForm.ID) {
        await updateSchedule(scheduleForm.ID, payload)
        ElMessage.success('更新成功')
      } else {
        await createSchedule(payload)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchSchedules()
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const resetForm = () => {
  scheduleForm.ID = null
  scheduleForm.course_id = null
  scheduleForm.class_id = null
  scheduleForm.teacher_id = null
  scheduleForm.day_of_week = 1
  scheduleForm.start_time = ''
  scheduleForm.end_time = ''
  scheduleForm.location = ''
  scheduleForm.semester = ''
}

const handleSizeChange = (val) => {
  pagination.page_size = val
  fetchSchedules()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchSchedules()
}

onMounted(() => {
  fetchSchedules()
})
</script>

<style scoped>
.schedule {
  padding: 20px;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.search-box {
  display: flex;
  align-items: center;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

