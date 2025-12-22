<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增选课</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="Student.name" label="学生姓名" min-width="100"/>
      <el-table-column prop="Student.student_id" label="学号" min-width="120"/>
      <el-table-column prop="Course.course_name" label="课程" min-width="150"/>
      <el-table-column prop="CreatedAt" label="选课时间" min-width="180">
        <template #default="{ row }">
          {{ new Date(row.CreatedAt).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="100" align="center">
        <template #default="{ row }">
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
    
    <el-dialog v-model="dialogVisible" title="新增选课" width="500px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="学生ID" prop="student_id">
          <el-input-number v-model="form.student_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="课程ID" prop="course_id">
          <el-input-number v-model="form.course_id" :min="1" style="width: 100%"/>
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
import { getEnrollments, createEnrollment, deleteEnrollment } from '@/api/enrollment'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const formRef = ref()
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({ student_id: null, course_id: null })
const formRules = {
  student_id: [{ required: true, message: '请输入学生ID', trigger: 'change' }],
  course_id: [{ required: true, message: '请输入课程ID', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getEnrollments({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  Object.assign(form, { student_id: null, course_id: null })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    await createEnrollment(form)
    ElMessage.success('创建成功')
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
    await ElMessageBox.confirm('确定要删除该选课记录吗？', '提示', { type: 'warning' })
    await deleteEnrollment(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchData())
</script>
