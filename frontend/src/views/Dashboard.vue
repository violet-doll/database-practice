<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#409eff"><User /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.studentCount }}</div>
              <div class="stat-label">学生总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#67c23a"><School /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.classCount }}</div>
              <div class="stat-label">班级总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#e6a23c"><Reading /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.courseCount }}</div>
              <div class="stat-label">课程总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#f56c6c"><UserFilled /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.teacherCount }}</div>
              <div class="stat-label">教师总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>欢迎使用学生管理系统</span>
            </div>
          </template>
          <div class="welcome-content">
            <h2>系统功能模块</h2>
            <el-row :gutter="20" style="margin-top: 20px">
              <el-col :span="8" v-for="module in modules" :key="module.name">
                <div class="module-card">
                  <el-icon :size="40" :color="module.color">
                    <component :is="module.icon" />
                  </el-icon>
                  <h3>{{ module.name }}</h3>
                  <p>{{ module.desc }}</p>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { User, School, Reading, UserFilled, Document, Calendar, Medal, Bell } from '@element-plus/icons-vue'
import { fetchAdminOverview } from '@/api/admin'

const stats = ref({
  studentCount: 0,
  classCount: 0,
  courseCount: 0,
  teacherCount: 0,
})

const modules = [
  { name: '学生管理', icon: User, color: '#409eff', desc: '学生信息的增删改查' },
  { name: '班级管理', icon: School, color: '#67c23a', desc: '班级信息管理' },
  { name: '课程管理', icon: Reading, color: '#e6a23c', desc: '课程信息管理' },
  { name: '成绩管理', icon: Document, color: '#f56c6c', desc: '学生成绩管理' },
  { name: '考勤管理', icon: Calendar, color: '#909399', desc: '学生考勤记录' },
  { name: '奖惩管理', icon: Medal, color: '#ff9800', desc: '学生奖惩记录' },
]

onMounted(async () => {
  try {
    const res = await fetchAdminOverview()
    if (res.code === 200) {
      const data = res.data
      stats.value = {
        studentCount: data.students_total,
        classCount: data.classes_total,
        courseCount: data.courses_total,
        teacherCount: data.teachers_total,
      }
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  font-size: 48px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}

.welcome-content {
  padding: 20px;
}

.welcome-content h2 {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
}

.module-card {
  text-align: center;
  padding: 30px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
  transition: all 0.3s;
  cursor: pointer;
}

.module-card:hover {
  background: #ecf5ff;
  transform: translateY(-5px);
}

.module-card h3 {
  font-size: 18px;
  color: #333;
  margin: 15px 0 10px;
}

.module-card p {
  font-size: 14px;
  color: #666;
}
</style>
