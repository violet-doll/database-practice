<template>
  <div class="admin-overview">
    <el-card>
      <div class="cards" style="display:grid; grid-template-columns: repeat(4, 1fr); gap:12px;">
        <el-statistic title="用户总数" :value="stats.users_total" />
        <el-statistic title="活跃用户" :value="stats.users_active" />
        <el-statistic title="学生总数" :value="stats.students_total" />
        <el-statistic title="教师总数" :value="stats.teachers_total" />
        <el-statistic title="班级总数" :value="stats.classes_total" />
        <el-statistic title="课程总数" :value="stats.courses_total" />
        <el-statistic title="选课记录" :value="stats.enrollments_total" />
        <el-statistic title="成绩记录" :value="stats.grades_total" />
        <el-statistic title="考勤记录" :value="stats.attendance_total" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { fetchAdminOverview } from '@/api/admin'

const stats = reactive({
  users_total: 0,
  users_active: 0,
  students_total: 0,
  teachers_total: 0,
  classes_total: 0,
  courses_total: 0,
  enrollments_total: 0,
  grades_total: 0,
  attendance_total: 0,
})

async function load() {
  const res = await fetchAdminOverview()
  Object.assign(stats, res.data)
}

onMounted(load)
</script>

<style scoped>
</style>


