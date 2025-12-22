<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增考勤记录</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="student.name" label="学生" min-width="100"/>
      <el-table-column prop="date" label="日期" min-width="120"/>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === '出勤' ? 'success' : row.status === '请假' ? 'warning' : 'danger'">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="备注" min-width="150"/>
      <el-table-column label="操作" width="100" align="center">
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
    
    <el-dialog v-model="dialogVisible" title="新增考勤记录" width="500px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="学生ID" prop="student_id">
          <el-input-number v-model="form.student_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="日期" prop="date">
          <el-date-picker v-model="form.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" style="width: 100%">
            <el-option value="出勤"/>
            <el-option value="缺席"/>
            <el-option value="请假"/>
            <el-option value="迟到"/>
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="reason">
          <el-input v-model="form.reason" type="textarea" rows="3"/>
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
import { getAttendance, createAttendance, deleteAttendance } from '@/api/attendance'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const formRef = ref()
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({ student_id: null, date: '', status: '出勤', reason: '' })
const formRules = {
  student_id: [{ required: true, message: '请输入学生ID', trigger: 'change' }],
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getAttendance({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  Object.assign(form, { student_id: null, date: '', status: '出勤', reason: '' })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    await createAttendance(form)
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
    await ElMessageBox.confirm('确定要删除该记录吗？', '提示', { type: 'warning' })
    await deleteAttendance(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchData())
</script>
