<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增班级</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="class_name" label="班级名称" min-width="150"/>
      <el-table-column prop="teacher.name" label="班主任" min-width="100"/>
      <el-table-column prop="teacher.phone" label="班主任电话" min-width="120"/>
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
        <el-form-item label="班级名称" prop="class_name">
          <el-input v-model="form.class_name" placeholder="请输入班级名称"/>
        </el-form-item>
        <el-form-item label="班主任" prop="teacher_id">
          <el-input-number v-model="form.teacher_id" :min="1" placeholder="请输入教师ID" style="width: 100%"/>
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
import { getClasses, createClass, updateClass, deleteClass } from '@/api/class'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()

const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({ class_name: '', teacher_id: null })
const formRules = {
  class_name: [{ required: true, message: '请输入班级名称', trigger: 'blur' }],
  teacher_id: [{ required: true, message: '请输入教师ID', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getClasses({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增班级'
  Object.assign(form, { ID: undefined, class_name: '', teacher_id: null })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑班级'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    if (form.ID) {
      await updateClass(form.ID, form)
      ElMessage.success('更新成功')
    } else {
      await createClass(form)
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
    await ElMessageBox.confirm('确定要删除该班级吗？', '提示', { type: 'warning' })
    await deleteClass(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchData())
</script>
