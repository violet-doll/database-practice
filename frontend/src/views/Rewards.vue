<template>
  <div class="page-container">
    <div class="filter-bar">
      <el-input v-model="filters.student_id" placeholder="学生ID" style="width: 160px; margin-right: 8px;" />
      <el-select v-model="filters.type" placeholder="类型" clearable style="width: 140px; margin-right: 8px;">
        <el-option label="奖励" value="奖励" />
        <el-option label="处分" value="处分" />
      </el-select>
      <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="margin-right: 8px;" />
      <el-button type="primary" @click="loadData(1)">查询</el-button>
      <el-button @click="resetFilters">重置</el-button>
      <el-button type="success" @click="openCreate">录入奖惩</el-button>
    </div>

    <el-table :data="list" border style="width: 100%;" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="student_id" label="学生ID" width="100" />
      <el-table-column prop="student.name" label="学生姓名" width="140" />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column prop="description" label="事由" />
      <el-table-column prop="date" label="日期" width="140" />
      <el-table-column prop="issuer" label="发布人" width="140" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-popconfirm title="确认删除该记录？" @confirm="onDelete(row.id)">
            <template #reference>
              <el-button type="danger" size="small">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-bar">
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="loadData"
        @size-change="onSizeChange"
      />
    </div>

    <el-dialog v-model="createVisible" title="录入奖惩">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px">
        <el-form-item label="学生ID" prop="student_id">
          <el-input v-model.number="form.student_id" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择">
            <el-option label="奖励" value="奖励" />
            <el-option label="处分" value="处分" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期" prop="date">
          <el-date-picker v-model="form.date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="事由" prop="description">
          <el-input v-model="form.description" type="textarea" rows="3" />
        </el-form-item>
        <el-form-item label="发布人" prop="issuer">
          <el-input v-model="form.issuer" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="onCreate">提交</el-button>
      </template>
    </el-dialog>
  </div>
  
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getRewards, createReward, deleteReward } from '@/api/reward'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = reactive({
  student_id: '',
  type: '',
  date_from: '',
  date_to: '',
})
const dateRange = ref([])

const createVisible = ref(false)
const formRef = ref()
const form = reactive({
  student_id: '',
  type: '',
  date: '',
  description: '',
  issuer: '',
})

const rules = {
  student_id: [{ required: true, message: '请输入学生ID', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  description: [{ required: true, message: '请输入事由', trigger: 'blur' }],
}

function resetFilters() {
  filters.student_id = ''
  filters.type = ''
  filters.date_from = ''
  filters.date_to = ''
  dateRange.value = []
  loadData(1)
}

function onSizeChange(size) {
  pageSize.value = size
  loadData(1)
}

async function loadData(p) {
  if (Array.isArray(dateRange.value) && dateRange.value.length === 2) {
    filters.date_from = dateRange.value[0]
    filters.date_to = dateRange.value[1]
  } else {
    filters.date_from = ''
    filters.date_to = ''
  }

  loading.value = true
  try {
    page.value = typeof p === 'number' ? p : page.value
    const res = await getRewards({
      ...filters,
      page: page.value,
      page_size: pageSize.value,
    })
    list.value = res?.data?.list || []
    total.value = res?.data?.total || 0
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { student_id: '', type: '', date: '', description: '', issuer: '' })
  createVisible.value = true
}

function onCreate() {
  formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      await createReward({
        student_id: Number(form.student_id),
        type: form.type,
        date: form.date,
        description: form.description,
        issuer: form.issuer,
      })
      ElMessage.success('创建成功')
      createVisible.value = false
      loadData(1)
    } catch (e) {
      // ignore
    }
  })
}

async function onDelete(id) {
  try {
    await deleteReward(id)
    ElMessage.success('删除成功')
    loadData(1)
  } catch (e) {
    // ignore
  }
}

onMounted(() => loadData(1))
</script>

<style scoped>
.page-container { padding: 16px; }
.filter-bar { margin-bottom: 12px; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.pagination-bar { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>