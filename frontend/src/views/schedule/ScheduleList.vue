<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增排课</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="course.course_name" label="课程" min-width="120"/>
      <el-table-column prop="class.class_name" label="班级" min-width="100"/>
      <el-table-column prop="teacher.name" label="教师" min-width="100"/>
      <el-table-column prop="day_of_week" label="星期" width="80" align="center">
        <template #default="{ row }">
          {{ ['', '周一', '周二', '周三', '周四', '周五', '周六', '周日'][row.day_of_week] }}
        </template>
      </el-table-column>
      <el-table-column prop="start_time" label="开始时间" width="100"/>
      <el-table-column prop="end_time" label="结束时间" width="100"/>
      <el-table-column prop="location" label="地点" min-width="100"/>
      <el-table-column prop="semester" label="学期" min-width="100"/>
      <el-table-column label="操作" fixed="right" width="180" align="center">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)" icon="Edit">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)" icon="Delete">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.page_size"
      :total="pagination.total"
      layout="total, prev, pager, next"
      @current-change="fetchData"
      class="mt-20"
    />
    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="课程ID" prop="course_id">
          <el-input-number v-model="form.course_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="班级ID" prop="class_id">
          <el-input-number v-model="form.class_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="教师ID" prop="teacher_id">
          <el-input-number v-model="form.teacher_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="星期" prop="day_of_week">
          <el-select v-model="form.day_of_week" style="width: 100%">
            <el-option :value="1" label="周一"/>
            <el-option :value="2" label="周二"/>
            <el-option :value="3" label="周三"/>
            <el-option :value="4" label="周四"/>
            <el-option :value="5" label="周五"/>
            <el-option :value="6" label="周六"/>
            <el-option :value="7" label="周日"/>
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间" prop="start_time">
          <el-input v-model="form.start_time" placeholder="如: 08:00"/>
        </el-form-item>
        <el-form-item label="结束时间" prop="end_time">
          <el-input v-model="form.end_time" placeholder="如: 09:40"/>
        </el-form-item>
        <el-form-item label="地点" prop="location">
          <el-input v-model="form.location" placeholder="如: 教5-101"/>
        </el-form-item>
        <el-form-item label="学期" prop="semester">
          <el-input v-model="form.semester" placeholder="如: 2025-Fall"/>
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
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({
  course_id: null,
  class_id: null,
  teacher_id: null,
  day_of_week: 1,
  start_time: '',
  end_time: '',
  location: '',
  semester: ''
})
const formRules = {
  course_id: [{ required: true, message: '请选择课程', trigger: 'change' }],
  class_id: [{ required: true, message: '请选择班级', trigger: 'change' }],
  teacher_id: [{ required: true, message: '请选择教师', trigger: 'change' }],
  day_of_week: [{ required: true, message: '请选择星期', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getSchedules({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增排课'
  Object.assign(form, { ID: undefined, course_id: null, class_id: null, teacher_id: null, day_of_week: 1, start_time: '', end_time: '', location: '', semester: '' })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑排课'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    if (form.ID) {
      await updateSchedule(form.ID, form)
      ElMessage.success('更新成功')
    } else {
      await createSchedule(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    if (error.response) ElMessage.error(error.response.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该排课吗？', '提示', { type: 'warning' })
    await deleteSchedule(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchData())
</script>
