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
    
    <!-- 数据可视化图表区域 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 学生性别分布饼图 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>学生性别分布</span>
            </div>
          </template>
          <div ref="genderChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>
      
      <!-- 考勤状态统计柱状图 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>考勤状态统计</span>
            </div>
          </template>
          <div ref="attendanceChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 成绩分布统计柱状图 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>成绩分布统计</span>
            </div>
          </template>
          <div ref="gradeDistChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <!-- 热门课程统计柱状图 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>热门课程 Top 5</span>
            </div>
          </template>
          <div ref="coursePopChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 成绩统计信息卡片 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>成绩统计</span>
            </div>
          </template>
          <div class="grade-stats">
            <div class="grade-stat-item">
              <div class="grade-stat-label">平均分</div>
              <div class="grade-stat-value">{{ dashboardStats.grade_stats?.average_score?.toFixed(2) || '0.00' }}</div>
            </div>
            <div class="grade-stat-item">
              <div class="grade-stat-label">成绩记录总数</div>
              <div class="grade-stat-value">{{ dashboardStats.grade_stats?.total_grades || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { User, School, Reading, UserFilled, Document, Calendar, Medal, Bell } from '@element-plus/icons-vue'
import { fetchAdminOverview } from '@/api/admin'
import { fetchDashboardStats } from '@/api/stats'
import * as echarts from 'echarts'

const stats = ref({
  studentCount: 0,
  classCount: 0,
  courseCount: 0,
  teacherCount: 0,
})

const dashboardStats = ref({
  student_gender: {
    male_count: 0,
    female_count: 0,
  },
  attendance_status: {
    present: 0,
    absent: 0,
    leave: 0,
    late: 0,
  },
  grade_stats: {
    average_score: 0,
    total_grades: 0,
  },
  class_stats: {
    total_classes: 0,
  },
  course_stats: {
    total_courses: 0,
  },
  grade_distribution: {},
  course_popularity: [],
})

const genderChartRef = ref(null)
const attendanceChartRef = ref(null)
const gradeDistChartRef = ref(null)
const coursePopChartRef = ref(null)
let genderChart = null
let attendanceChart = null
let gradeDistChart = null
let coursePopChart = null



// 初始化学生性别分布饼图
const initGenderChart = () => {
  if (!genderChartRef.value) return
  
  genderChart = echarts.init(genderChartRef.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '学生性别',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: '{b}: {c}\n({d}%)'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        data: [
          { value: dashboardStats.value.student_gender.male_count, name: '男', itemStyle: { color: '#5470c6' } },
          { value: dashboardStats.value.student_gender.female_count, name: '女', itemStyle: { color: '#ee6666' } }
        ]
      }
    ]
  }
  genderChart.setOption(option)
}

// 初始化考勤状态统计柱状图
const initAttendanceChart = () => {
  if (!attendanceChartRef.value) return
  
  attendanceChart = echarts.init(attendanceChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['出勤', '缺席', '请假', '迟到'],
      axisTick: {
        alignWithLabel: true
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '考勤数量',
        type: 'bar',
        barWidth: '60%',
        data: [
          { value: dashboardStats.value.attendance_status.present, itemStyle: { color: '#91cc75' } },
          { value: dashboardStats.value.attendance_status.absent, itemStyle: { color: '#ff7875' } },
          { value: dashboardStats.value.attendance_status.leave, itemStyle: { color: '#fac858' } },
          { value: dashboardStats.value.attendance_status.late, itemStyle: { color: '#73c0de' } }
        ]
      }
    ]
  }
  attendanceChart.setOption(option)
}

// 初始化成绩分布图表
const initGradeDistChart = () => {
  if (!gradeDistChartRef.value) return

  gradeDistChart = echarts.init(gradeDistChartRef.value)
  const data = dashboardStats.value.grade_distribution || {}
  const categories = ['<60', '60-69', '70-79', '80-89', '>=90']
  const values = categories.map(key => data[key] || 0)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: categories,
      axisTick: { alignWithLabel: true }
    },
    yAxis: { type: 'value' },
    series: [
      {
        name: '人数',
        type: 'bar',
        barWidth: '60%',
        itemStyle: {
          color: function(params) {
            const colors = ['#ff7875', '#ff9c6e', '#fac858', '#91cc75', '#3ba272'];
            return colors[params.dataIndex] || '#5470c6';
          }
        },
        data: values
      }
    ]
  }
  gradeDistChart.setOption(option)
}

// 初始化热门课程图表
const initCoursePopChart = () => {
  if (!coursePopChartRef.value) return

  coursePopChart = echarts.init(coursePopChartRef.value)
  const data = dashboardStats.value.course_popularity || []
  const categories = data.map(item => item.course_name)
  const values = data.map(item => item.count)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'value'
    },
    yAxis: {
      type: 'category',
      data: categories,
      axisLabel: { interval: 0, rotate: 0 }
    },
    series: [
      {
        name: '选课人数',
        type: 'bar',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#83bff6' },
            { offset: 0.5, color: '#188df0' },
            { offset: 1, color: '#188df0' }
          ])
        },
        data: values
      }
    ]
  }
  coursePopChart.setOption(option)
}

onMounted(async () => {
  try {
    // 获取基础统计数据
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
    
    // 获取详细统计数据用于图表
    const statsRes = await fetchDashboardStats()
    if (statsRes.code === 200) {
      dashboardStats.value = statsRes.data
      
      // 等待 DOM 更新后初始化图表
      await nextTick()
      initGenderChart()
      initAttendanceChart()
      initGradeDistChart()
      initCoursePopChart()
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

.card-header {
  font-weight: bold;
  font-size: 16px;
}

.grade-stats {
  display: flex;
  justify-content: space-around;
  padding: 20px 0;
}

.grade-stat-item {
  text-align: center;
}

.grade-stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.grade-stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
}
</style>
