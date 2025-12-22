<template>
  <div class="page-container">
    <el-form :model="searchForm" :inline="true" class="search-form">
      <el-form-item label="姓名">
        <el-input v-model="searchForm.name" placeholder="请输入学生姓名" clearable/>
      </el-form-item>
      <el-form-item label="学号">
        <el-input v-model="searchForm.student_id" placeholder="请输入学号" clearable/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch" icon="Search">搜索</el-button>
        <el-button @click="handleReset" icon="Refresh">重置</el-button>
      </el-form-item>
    </el-form>
    
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd" icon="Plus">新增学生</el-button>
      <el-button @click="fetchData" icon="Refresh">刷新</el-button>
    </div>
    
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column type="index" label="序号" width="60" align="center"/>
      <el-table-column prop="student_id" label="学号" min-width="120"/>
      <el-table-column prop="name" label="姓名" min-width="100"/>
      <el-table-column prop="gender" label="性别" width="80" align="center"/>
      <el-table-column prop="age" label="年龄" width="80" align="center"/>
      <el-table-column prop="email" label="邮箱" min-width="180"/>
      <el-table-column prop="phone" label="电话" min-width="120"/>
      <el-table-column prop="class.class_name" label="班级" min-width="100"/>
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
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="fetchData"
      @current-change="fetchData"
      class="mt-20"
    />
    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="学号" prop="student_id">
          <el-input v-model="form.student_id" placeholder="请输入学号"/>
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名"/>
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio label="男">男</el-radio>
            <el-radio label="女">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input-number v-model="form.age" :min="1" :max="100"/>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱"/>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入电话"/>
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="form.address" type="textarea" placeholder="请输入地址"/>
        </el-form-item>
        <el-form-item label="班级" prop="class_id">
          <el-select v-model="form.class_id" placeholder="请选择班级" style="width: 100%">
            <el-option
              v-for="item in classList"
              :key="item.ID"
              :label="item.class_name"
              :value="item.ID"
            />
          </el-select>
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
import { getStudents, createStudent, updateStudent, deleteStudent } from '@/api/student'
import { getClasses } from '@/api/class'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()

const searchForm = reactive({
  name: '',
  student_id: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const tableData = ref([])
const classList = ref([])

const form = reactive({
  name: '',
  student_id: '',
  gender: '男',
  age: 18,
  email: '',
  phone: '',
  address: '',
  class_id: null
})

const formRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  student_id: [{ required: true, message: '请输入学号', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    }
    const response = await getStudents(params)
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const fetchClasses = async () => {
  try {
    const response = await getClasses({ page: 1, page_size: 100 })
    classList.value = response.data || []
  } catch (error) {
    console.error('获取班级列表失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.student_id = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增学生'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑学生'
  Object.assign(form, row)
  dialogVisible.value = true
}

const resetForm = () => {
  Object.assign(form, {
    ID: undefined,
    name: '',
    student_id: '',
    gender: '男',
    age: 18,
    email: '',
    phone: '',
    address: '',
    class_id: null
  })
  formRef.value?.clearValidate()
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    
    if (form.ID) {
      await updateStudent(form.ID, form)
      ElMessage.success('更新成功')
    } else {
      await createStudent(form)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data?.error || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该学生吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteStudent(row.ID)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchData()
  fetchClasses()
})
</script>
