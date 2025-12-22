<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="学生总数" :value="stats.studentCount">
            <template #prefix>
              <el-icon style="vertical-align: middle"><User /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="班级总数" :value="stats.classCount">
            <template #prefix>
              <el-icon style="vertical-align: middle"><School /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="课程总数" :value="stats.courseCount">
            <template #prefix>
              <el-icon style="vertical-align: middle"><Reading /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="教师总数" :value="stats.teacherCount">
            <template #prefix>
              <el-icon style="vertical-align: middle"><Avatar /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="mt-20">
      <template #header>
        <div class="card-header">
          <span>欢迎使用学生管理系统</span>
        </div>
      </template>
      <el-empty description="请从左侧菜单选择功能模块" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const stats = ref({
  studentCount: 0,
  classCount: 0,
  courseCount: 0,
  teacherCount: 0
})

const fetchStats = async () => {
  try {
    const response = await request.get('/api/v1/admin/stats/overview')
    stats.value = response.data || {}
  } catch (error) {
    console.error('获取统计数据失败:', error)
    // 如果没有权限，不显示错误消息
    if (error.response?.status !== 403) {
      ElMessage.error('获取统计数据失败')
    }
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
