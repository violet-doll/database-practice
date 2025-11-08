<template>
  <div class="enrollments">
    <el-card>
      <div class="toolbar">
        <el-button type="primary" @click="dialogVisible = true">添加选课</el-button>
      </div>

      <!-- 筛选区域 -->
      <div class="search-box" style="margin-top: 16px; margin-bottom: 16px">
        <el-input
          v-model="searchForm.student_id"
          placeholder="学生ID"
          style="width: 180px; margin-right: 10px"
          clearable
        />
        <el-input
          v-model="searchForm.course_id"
          placeholder="课程ID"
          style="width: 180px; margin-right: 10px"
          clearable
        />
        <el-button type="primary" @click="fetchEnrollments">查询</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </div>

      <!-- 选课列表表格 -->
      <el-table :data="enrollments" style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column label="学生ID" width="120">
          <template #default="{ row }">{{ row.student_id }}</template>
        </el-table-column>
        <el-table-column label="学生姓名" min-width="160">
          <template #default="{ row }">{{ row.Student?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="课程ID" width="120">
          <template #default="{ row }">{{ row.course_id }}</template>
        </el-table-column>
        <el-table-column label="课程名称" min-width="180">
          <template #default="{ row }">{{ row.Course?.course_name || '-' }}</template>
        </el-table-column>
        <el-table-column label="选课时间" width="180">
          <template #default="{ row }">
            {{ row.CreatedAt ? new Date(row.CreatedAt).toLocaleString('zh-CN') : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
              :loading="deletingId === row.ID"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="page.page"
          v-model:page-size="page.page_size"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="page.total"
          @size-change="() => fetchEnrollments()"
          @current-change="() => fetchEnrollments()"
        />
      </div>
    </el-card>

    <!-- 添加选课弹窗 -->
    <el-dialog v-model="dialogVisible" title="添加选课" width="520px">
      <el-form ref="enrollmentFormRef" :model="enrollmentForm" :rules="enrollmentRules" label-width="100px">
        <el-form-item label="学生ID" prop="student_id">
          <el-input v-model.number="enrollmentForm.student_id" type="number" />
        </el-form-item>
        <el-form-item label="课程ID" prop="course_id">
          <el-input v-model.number="enrollmentForm.course_id" type="number" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEnrollment" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getEnrollments, createEnrollment, deleteEnrollment } from '@/api/enrollment'

// 数据
const enrollments = ref([])
const loading = ref(false)
const page = reactive({ page: 1, page_size: 10, total: 0 })
const searchForm = reactive({ student_id: '', course_id: '' })

// 获取选课列表
const fetchEnrollments = async () => {
  loading.value = true
  try {
    const params = {
      page: page.page,
      page_size: page.page_size,
    }
    if (searchForm.student_id) {
      params.student_id = searchForm.student_id
    }
    if (searchForm.course_id) {
      params.course_id = searchForm.course_id
    }
    const res = await getEnrollments(params)
    enrollments.value = res.data.list || []
    page.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('获取选课列表失败')
  } finally {
    loading.value = false
  }
}

// 重置搜索
const resetSearch = () => {
  searchForm.student_id = ''
  searchForm.course_id = ''
  page.page = 1
  fetchEnrollments()
}

// 添加选课弹窗
const dialogVisible = ref(false)
const enrollmentFormRef = ref(null)
const submitting = ref(false)
const enrollmentForm = reactive({ student_id: '', course_id: '' })
const enrollmentRules = {
  student_id: [{ required: true, message: '请输入学生ID', trigger: 'blur' }],
  course_id: [{ required: true, message: '请输入课程ID', trigger: 'blur' }],
}

// 提交选课
const submitEnrollment = async () => {
  if (!enrollmentFormRef.value) return
  await enrollmentFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await createEnrollment({
        student_id: Number(enrollmentForm.student_id),
        course_id: Number(enrollmentForm.course_id),
      })
      ElMessage.success('添加选课成功')
      dialogVisible.value = false
      // 重置表单
      enrollmentForm.student_id = ''
      enrollmentForm.course_id = ''
      // 刷新列表
      await fetchEnrollments()
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '添加选课失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除选课
const deletingId = ref(null)
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条选课记录吗？删除后关联的成绩记录也会被删除。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    deletingId.value = row.ID
    try {
      await deleteEnrollment(row.ID)
      ElMessage.success('删除成功')
      await fetchEnrollments()
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '删除失败')
    } finally {
      deletingId.value = null
    }
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  fetchEnrollments()
})
</script>

<style scoped>
.enrollments {
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
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

