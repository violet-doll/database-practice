<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">录入成绩</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="score_type" label="成绩类型" min-width="100"/>
      <el-table-column prop="score" label="分数" width="100" align="center"/>
      <el-table-column prop="CreatedAt" label="录入时间" min-width="180">
        <template #default="{ row }">
          {{ new Date(row.CreatedAt).toLocaleString() }}
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
    
    <el-dialog v-model="dialogVisible" title="录入成绩" width="500px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="选课记录ID" prop="enrollment_id">
          <el-input-number v-model="form.enrollment_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="成绩类型" prop="score_type">
          <el-select v-model="form.score_type" style="width: 100%">
            <el-option value="平时成绩"/>
            <el-option value="期中成绩"/>
            <el-option value="期末成绩"/>
            <el-option value="总评"/>
          </el-select>
        </el-form-item>
        <el-form-item label="分数" prop="score">
          <el-input-number v-model="form.score" :min="0" :max="100" :precision="1" style="width: 100%"/>
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
import { ElMessage } from 'element-plus'
import { getGrades, createGrade } from '@/api/grade'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const formRef = ref()
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({ enrollment_id: null, score_type: '总评', score: 0 })
const formRules = {
  enrollment_id: [{ required: true, message: '请输入选课记录ID', trigger: 'change' }],
  score_type: [{ required: true, message: '请选择成绩类型', trigger: 'change' }],
  score: [{ required: true, message: '请输入分数', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getGrades({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  Object.assign(form, { enrollment_id: null, score_type: '总评', score: 0 })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    await createGrade(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    if (error.response) ElMessage.error(error.response.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => fetchData())
</script>
