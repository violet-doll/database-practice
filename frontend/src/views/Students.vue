<template>
  <div class="students">
    <el-card>
      <!-- 搜索和操作栏 -->
      <div class="toolbar">
        <div class="search-box">
          <el-input
            v-model="searchForm.name"
            placeholder="搜索学生姓名"
            style="width: 200px; margin-right: 10px"
            clearable
          />
          <el-input
            v-model="searchForm.student_id"
            placeholder="搜索学号"
            style="width: 200px; margin-right: 10px"
            clearable
          />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <el-button type="primary" @click="handleAdd">添加学生</el-button>
      </div>

      <!-- 学生列表表格 -->
      <el-table
        :data="studentList"
        style="width: 100%; margin-top: 20px"
        v-loading="loading"
      >
        <el-table-column prop="student_id" label="学号" width="120" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="gender" label="性别" width="80" />
        <el-table-column prop="age" label="年龄" width="80" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="class.class_name" label="班级" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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

    <!-- 添加/编辑学生对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form
        ref="studentFormRef"
        :model="studentForm"
        :rules="studentRules"
        label-width="100px"
      >
        <el-form-item label="学号" prop="student_id">
          <el-input v-model="studentForm.student_id" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="studentForm.name" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-select v-model="studentForm.gender" style="width: 100%">
            <el-option label="男" value="男" />
            <el-option label="女" value="女" />
          </el-select>
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input-number v-model="studentForm.age" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="studentForm.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="studentForm.email" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="studentForm.address" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getStudents, createStudent, updateStudent, deleteStudent } from '@/api/student'

const loading = ref(false)
const studentList = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('添加学生')
const studentFormRef = ref(null)

const searchForm = reactive({
  name: '',
  student_id: '',
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0,
})

const studentForm = reactive({
  id: null,
  student_id: '',
  name: '',
  gender: '男',
  age: 18,
  phone: '',
  email: '',
  address: '',
})

const studentRules = {
  student_id: [{ required: true, message: '请输入学号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
}

const fetchStudents = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm,
    }
    const res = await getStudents(params)
    studentList.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    console.error('获取学生列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchStudents()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.student_id = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '添加学生'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑学生'
  Object.assign(studentForm, row)
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该学生吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await deleteStudent(row.ID)
      ElMessage.success('删除成功')
      fetchStudents()
    } catch (error) {
      console.error('删除失败:', error)
    }
  })
}

const handleSubmit = async () => {
  if (!studentFormRef.value) return
  
  await studentFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (studentForm.id) {
          await updateStudent(studentForm.ID, studentForm)
          ElMessage.success('更新成功')
        } else {
          await createStudent(studentForm)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchStudents()
      } catch (error) {
        console.error('提交失败:', error)
      }
    }
  })
}

const resetForm = () => {
  studentForm.id = null
  studentForm.student_id = ''
  studentForm.name = ''
  studentForm.gender = '男'
  studentForm.age = 18
  studentForm.phone = ''
  studentForm.email = ''
  studentForm.address = ''
}

const handleSizeChange = (val) => {
  pagination.page_size = val
  fetchStudents()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchStudents()
}

onMounted(() => {
  fetchStudents()
})
</script>

<style scoped>
.students {
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
