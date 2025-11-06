<template>
  <div class="attendance-stats">
    <el-card>
      <div class="toolbar">
        <div class="search-box">
          <el-input v-model="query.student_id" placeholder="学生ID（可选）" style="width: 180px; margin-right: 10px" />
          <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="margin-right: 10px" />
          <el-button type="primary" @click="fetchStats">统计</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </div>
      </div>

      <el-table :data="list" style="width: 100%; margin-top: 16px" v-loading="loading">
        <el-table-column label="学生ID" width="120">
          <template #default="{ row }">{{ row.student_id }}</template>
        </el-table-column>
        <el-table-column label="学生姓名" min-width="160">
          <template #default="{ row }">{{ row.student && row.student.name }}</template>
        </el-table-column>
        <el-table-column prop="present" label="出勤" width="100" />
        <el-table-column prop="absent" label="缺席" width="100" />
        <el-table-column prop="leave" label="请假" width="100" />
        <el-table-column prop="late" label="迟到" width="100" />
        <el-table-column prop="total" label="总计" width="120" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import dayjs from 'dayjs'
import { getAttendanceStats } from '@/api/attendance'

const loading = ref(false)
const list = ref([])

const query = reactive({ student_id: '', date_from: '', date_to: '' })
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

const fetchStats = async () => {
  loading.value = true
  try {
    const res = await getAttendanceStats({
      student_id: query.student_id || undefined,
      date_from: query.date_from || undefined,
      date_to: query.date_to || undefined,
    })
    list.value = res.data.list || []
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  query.student_id = ''
  dateRange.value = ''
  fetchStats()
}
</script>

<style scoped>
.attendance-stats { padding: 20px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.search-box { display: flex; align-items: center; }
</style>


