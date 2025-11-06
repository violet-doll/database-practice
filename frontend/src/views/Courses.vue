<template>
  <div class="courses">
    <el-card>
      <div class="toolbar">
        <div class="search-box">
          <el-input
            v-model="searchForm.course_name"
            placeholder="搜索课程名称"
            style="width: 220px; margin-right: 10px"
            clearable
          />
          <el-input
            v-model="searchForm.teacher_id"
            placeholder="按教师ID过滤"
            style="width: 180px; margin-right: 10px"
            clearable
          />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <el-button type="primary" @click="handleAdd">新增课程</el-button>
      </div>

      <el-table :data="courseList" style="width: 100%; margin-top: 20px" v-loading="loading">
        <el-table-column prop="ID" label="ID" width="90" />
        <el-table-column prop="course_name" label="课程名称" min-width="180" />
        <el-table-column prop="credits" label="学分" width="100" />
        <el-table-column prop="teacher_id" label="教师ID" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="courseFormRef" :model="courseForm" :rules="courseRules" label-width="100px">
        <el-form-item label="课程名称" prop="course_name">
          <el-input v-model="courseForm.course_name" />
        </el-form-item>
        <el-form-item label="学分" prop="credits">
          <el-input-number v-model="courseForm.credits" :min="0" :max="20" :step="0.5" />
        </el-form-item>
        <el-form-item label="教师ID" prop="teacher_id">
          <el-input v-model="courseForm.teacher_id" />
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
import { getCourses, createCourse, updateCourse, deleteCourse } from '@/api/course'

const loading = ref(false)
const courseList = ref([])

const searchForm = reactive({
  course_name: '',
  teacher_id: '',
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0,
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增课程')
const courseFormRef = ref(null)
const courseForm = reactive({
  ID: null,
  course_name: '',
  credits: 0,
  teacher_id: ''
})

const courseRules = {
  course_name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
}

const fetchCourses = async () => {
  loading.value = true
  try {
    const params = { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    const res = await getCourses(params)
    courseList.value = res.data.list || []
    pagination.total = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchCourses()
}

const handleReset = () => {
  searchForm.course_name = ''
  searchForm.teacher_id = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增课程'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑课程'
  Object.assign(courseForm, row)
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该课程吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    await deleteCourse(row.ID)
    ElMessage.success('删除成功')
    fetchCourses()
  })
}

const handleSubmit = async () => {
  if (!courseFormRef.value) return
  await courseFormRef.value.validate(async (valid) => {
    if (!valid) return
    if (courseForm.ID) {
      const payload = {
        course_name: courseForm.course_name,
        credits: Number(courseForm.credits),
        teacher_id: courseForm.teacher_id === '' ? 0 : Number(courseForm.teacher_id),
      }
      await updateCourse(courseForm.ID, payload)
      ElMessage.success('更新成功')
    } else {
      const payload = {
        course_name: courseForm.course_name,
        credits: Number(courseForm.credits),
        teacher_id: courseForm.teacher_id === '' ? 0 : Number(courseForm.teacher_id),
      }
      await createCourse(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchCourses()
  })
}

const resetForm = () => {
  courseForm.ID = null
  courseForm.course_name = ''
  courseForm.credits = 0
  courseForm.teacher_id = ''
}

const handleSizeChange = (val) => {
  pagination.page_size = val
  fetchCourses()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchCourses()
}

onMounted(() => {
  fetchCourses()
})
</script>

<style scoped>
.courses {
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
