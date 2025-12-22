<template>
  <div class="page-container">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增用户</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="username" label="用户名" min-width="120"/>
      <el-table-column prop="role.role_name" label="角色" min-width="100"/>
      <el-table-column prop="user_type" label="用户类型" width="100"/>
      <el-table-column prop="is_active" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'danger'">
            {{ row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" align="center">
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
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username"/>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!form.ID">
          <el-input v-model="form.password" type="password"/>
        </el-form-item>
        <el-form-item label="角色ID" prop="role_id">
          <el-input-number v-model="form.role_id" :min="1" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="用户类型" prop="user_type">
          <el-select v-model="form.user_type" style="width: 100%">
            <el-option value="admin" label="管理员"/>
            <el-option value="teacher" label="教师"/>
            <el-option value="student" label="学生"/>
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="is_active">
          <el-switch v-model="form.is_active"/>
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
import { getUsers, createUser, updateUser, deleteUser } from '@/api/admin'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const tableData = ref([])
const form = reactive({ username: '', password: '', role_id: null, user_type: 'admin', is_active: true })
const formRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role_id: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getUsers({ page: pagination.page, page_size: pagination.page_size })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增用户'
  Object.assign(form, { ID: undefined, username: '', password: '', role_id: null, user_type: 'admin', is_active: true })
  formRef.value?.clearValidate()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑用户'
  Object.assign(form, row)
  delete form.password
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    if (form.ID) {
      await updateUser(form.ID, form)
      ElMessage.success('更新成功')
    } else {
      await createUser(form)
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
    await ElMessageBox.confirm('确定要删除该用户吗？', '提示', { type: 'warning' })
    await deleteUser(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchData())
</script>
