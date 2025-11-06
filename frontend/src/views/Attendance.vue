<template>
  <div class="attendance">
    <el-card>
      <div class="toolbar">
        <div class="search-box">
          <el-input v-model="query.student_id" placeholder="学生ID" style="width: 160px; margin-right: 10px" />
          <el-select v-model="query.status" placeholder="状态" clearable style="width: 140px; margin-right: 10px">
            <el-option label="出勤" value="出勤" />
            <el-option label="缺席" value="缺席" />
            <el-option label="请假" value="请假" />
            <el-option label="迟到" value="迟到" />
          </el-select>
          <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="margin-right: 10px" />
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </div>
        <div>
          <el-button type="primary" @click="dialogVisible = true">新增考勤</el-button>
          <router-link to="/attendance-stats"><el-button type="success" style="margin-left: 8px">查看统计</el-button></router-link>
        </div>
      </div>

      <el-table :data="list" style="width: 100%; margin-top: 16px" v-loading="loading">
        <el-table-column label="学生ID" width="120">
          <template #default="{ row }">{{ row.student_id }}</template>
        </el-table-column>
        <el-table-column label="学生姓名" min-width="160">
          <template #default="{ row }">{{ row.student && row.student.name }}</template>
        </el-table-column>
        <el-table-column prop="date" label="日期" width="140" />
        <el-table-column prop="status" label="状态" width="120" />
        <el-table-column prop="reason" label="备注/原因" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-popconfirm title="确定删除这条考勤记录吗？" @confirm="() => handleDelete(row)">
              <template #reference>
                <el-button type="danger" size="small" plain>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page.page"
          v-model:page-size="page.page_size"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="page.total"
          @size-change="() => fetchList()"
          @current-change="() => fetchList()"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增考勤" width="520px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="96px">
        <el-form-item label="学生ID" prop="student_id">
          <el-input v-model="form.student_id" />
        </el-form-item>
        <el-form-item label="日期" prop="date">
          <el-date-picker v-model="form.date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status">
            <el-option label="出勤" value="出勤" />
            <el-option label="缺席" value="缺席" />
            <el-option label="请假" value="请假" />
            <el-option label="迟到" value="迟到" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注/原因" prop="reason">
          <el-input v-model="form.reason" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { getAttendance, createAttendance, deleteAttendance } from '@/api/attendance'

const loading = ref(false)
const list = ref([])
const page = reactive({ page: 1, page_size: 10, total: 0 })

const query = reactive({ student_id: '', status: '', date_from: '', date_to: '' })
const dateRange = ref('')

watch(dateRange, (v) => {
  if (Array.isArray(v) && v.length === 2) {
    query.date_from = dayjs(v[0]).format('YYYY-MM-DD')
    query.date_to = dayjs(v[1]).format('YYYY-MM-DD')
  } else {
    query.date_from = ''
    query.date_to = ''
  }
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getAttendance({
      page: page.page,
      page_size: page.page_size,
      student_id: query.student_id || undefined,
      status: query.status || undefined,
      date_from: query.date_from || undefined,
      date_to: query.date_to || undefined,
    })
    list.value = res.data.list || []
    page.total = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  query.student_id = ''
  query.status = ''
  dateRange.value = ''
  fetchList()
}

const dialogVisible = ref(false)
const formRef = ref(null)
const form = reactive({ student_id: '', date: '', status: '出勤', reason: '' })
const rules = {
  student_id: [{ required: true, message: '请输入学生ID', trigger: 'blur' }],
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    await createAttendance({
      student_id: Number(form.student_id),
      date: dayjs(form.date).format('YYYY-MM-DD'),
      status: form.status,
      reason: form.reason,
    })
    ElMessage.success('新增成功')
    dialogVisible.value = false
    fetchList()
  })
}

onMounted(() => {
  fetchList()
})

const handleDelete = async (row) => {
  await deleteAttendance(row.ID || row.id)
  ElMessage.success('删除成功')
  fetchList()
}
</script>

<style scoped>
.attendance { padding: 20px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.search-box { display: flex; align-items: center; }
.pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>