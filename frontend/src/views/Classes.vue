<template>
  <div class="classes">
    <el-card>
      <!-- 工具栏：搜索 + 新增 -->
      <div class="toolbar">
        <div class="search-box">
          <el-input
            v-model="searchForm.class_name"
            placeholder="搜索班级名称"
            style="width: 220px; margin-right: 10px"
            clearable
          />
          <el-input
            v-model="searchForm.teacher_id"
            placeholder="按班主任ID过滤"
            style="width: 180px; margin-right: 10px"
            clearable
          />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <el-button type="primary" @click="handleAdd">新增班级</el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="classList" style="width: 100%; margin-top: 20px" v-loading="loading">
        <el-table-column prop="ID" label="ID" width="90" />
        <el-table-column prop="class_name" label="班级名称" min-width="200" />
        <el-table-column prop="teacher_id" label="班主任ID" width="120" />
        <el-table-column label="班主任" min-width="160">
          <template #default="{ row }">
            <span>{{ row.teacher?.name || '-' }}</span>
          </template>
        </el-table-column>
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

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="classFormRef" :model="classForm" :rules="classRules" label-width="96px">
        <el-form-item label="班级名称" prop="class_name">
          <el-input v-model="classForm.class_name" placeholder="请输入班级名称" />
        </el-form-item>
        <el-form-item label="班主任ID" prop="teacher_id">
          <el-input v-model.number="classForm.teacher_id" placeholder="请输入班主任ID（可为空）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handleSubmit">保 存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
  
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getClasses, createClass, updateClass, deleteClass } from '@/api/class'

const loading = ref(false)
const classList = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增班级')
const classFormRef = ref(null)

const searchForm = reactive({
  class_name: '',
  teacher_id: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0,
})

const classForm = reactive({
  ID: null,
  class_name: '',
  teacher_id: null,
})

const classRules = {
  class_name: [{ required: true, message: '请输入班级名称', trigger: 'blur' }],
}

const fetchClasses = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm,
    }
    const res = await getClasses(params)
    classList.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    console.error('获取班级列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchClasses()
}

const handleReset = () => {
  searchForm.class_name = ''
  searchForm.teacher_id = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增班级'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑班级'
  Object.assign(classForm, {
    ID: row.ID,
    class_name: row.class_name,
    teacher_id: row.teacher_id ?? null,
  })
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定删除班级【${row.class_name}】吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      await deleteClass(row.ID)
      ElMessage.success('删除成功')
      fetchClasses()
    })
    .catch(() => {})
}

const resetForm = () => {
  Object.assign(classForm, {
    ID: null,
    class_name: '',
    teacher_id: null,
  })
}

const handleSubmit = () => {
  classFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (classForm.ID) {
        await updateClass(classForm.ID, {
          class_name: classForm.class_name,
          teacher_id: classForm.teacher_id,
        })
        ElMessage.success('更新成功')
      } else {
        await createClass({
          class_name: classForm.class_name,
          teacher_id: classForm.teacher_id,
        })
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchClasses()
    } catch (error) {
      console.error('保存失败:', error)
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchClasses()
}

const handleCurrentChange = () => {
  fetchClasses()
}

onMounted(() => {
  fetchClasses()
})
</script>

<style scoped>
.classes {
  padding: 20px;
}
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
