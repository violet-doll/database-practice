<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增课程</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="course_name" label="课程名称" min-width="150"/>
      <el-table-column prop="teacher.name" label="授课教师" min-width="100"/>
      <el-table-column prop="credits" label="学分" width="80" align="center"/>
      <el-table-column prop="capacity" label="容量" width="80" align="center"/>
      <el-table-column prop="enrolled_count" label="已选人数" width="100" align="center"/>
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
    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="课程名称" prop="course_name">
          <el-input v-model="form.course_name" placeholder="请输入课程名称"/>
        </el-form-item>
        <el-form-item label="教师ID" prop="teacher_id">
          <el-input-number v-model="form.teacher_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="学分" prop="credits">
          <el-input-number v-model="form.credits" :min="0" :step="0.5" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="容量" prop="capacity">
          <el-input-number v-model="form.capacity" :min="1" style="width: 100%"/>
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
import { getCourses, createCourse, updateCourse, deleteCourse } from '@/api/course'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({ course_name: '', teacher_id: null, credits: 1, capacity: 50 })
const formRules = {
  course_name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
  teacher_id: [{ required: true, message: '请输入教师ID', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getCourses({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增课程'
  Object.assign(form, { ID: undefined, course_name: '', teacher_id: null, credits: 1, capacity: 50 })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑课程'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    if (form.ID) {
      await updateCourse(form.ID, form)
      ElMessage.success('更新成功')
    } else {
      await createCourse(form)
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
    await ElMessageBox.confirm('确定要删除该课程吗？', '提示', { type: 'warning' })
    await deleteCourse(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchData())
</script>
