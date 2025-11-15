<template>
  <div class="main-layout">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside width="240px" class="sidebar">
        <div class="logo">
          <h2>学生管理系统</h2>
        </div>
        <el-menu
          :default-active="activeMenu"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <!-- 数据看板：所有已登录用户可见 -->
          <el-menu-item index="/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据看板</span>
          </el-menu-item>
          
          <!-- 学生管理：需要 student:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('student:read')" index="/students">
            <el-icon><User /></el-icon>
            <span>学生管理</span>
          </el-menu-item>
          
          <!-- 班级管理：需要 class:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('class:read')" index="/classes">
            <el-icon><School /></el-icon>
            <span>班级管理</span>
          </el-menu-item>
          
          <!-- 课程管理：需要 course:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('course:read')" index="/courses">
            <el-icon><Reading /></el-icon>
            <span>课程管理</span>
          </el-menu-item>
          
          <!-- 选课管理：需要 enrollment:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('enrollment:read')" index="/enrollments">
            <el-icon><DocumentAdd /></el-icon>
            <span>选课管理</span>
          </el-menu-item>
          
          <!-- 排课管理：需要 schedule:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('schedule:read')" index="/schedule">
            <el-icon><Grid /></el-icon>
            <span>排课管理</span>
          </el-menu-item>
          
          <!-- 成绩管理：需要 grade:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('grade:read')" index="/grades">
            <el-icon><Document /></el-icon>
            <span>成绩管理</span>
          </el-menu-item>
          
          <!-- 考勤管理：需要 attendance:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('attendance:read')" index="/attendance">
            <el-icon><Calendar /></el-icon>
            <span>考勤管理</span>
          </el-menu-item>
          
          <!-- 奖惩管理：需要 reward:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('reward:read')" index="/rewards">
            <el-icon><Medal /></el-icon>
            <span>奖惩管理</span>
          </el-menu-item>
          
          <!-- 通知管理：需要 notification:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('notification:read')" index="/notifications">
            <el-icon><Bell /></el-icon>
            <span>通知管理</span>
          </el-menu-item>
          
          <!-- 家长联系方式：需要 parent:read 权限 -->
          <el-menu-item v-if="userStore.hasPermission('parent:read')" index="/parents">
            <el-icon><User /></el-icon>
            <span>家长联系方式</span>
          </el-menu-item>
          
          <!-- 管理员：系统设置菜单，需要 admin:user:read 或 admin:role:read 权限 -->
          <template v-if="userStore.hasPermission('admin:user:read') || userStore.hasPermission('admin:role:read')">
            <el-sub-menu index="/admin">
              <template #title>
                <el-icon><Setting /></el-icon>
                <span>系统设置</span>
              </template>
              <!-- 统计概览：需要 admin:stats:read 权限 -->
              <el-menu-item v-if="userStore.hasPermission('admin:stats:read')" index="/admin/overview">
                <el-icon><DataAnalysis /></el-icon>
                <span>统计概览</span>
              </el-menu-item>
              <!-- 用户与权限：需要 admin:user:read 权限 -->
              <el-menu-item v-if="userStore.hasPermission('admin:user:read')" index="/admin/users">
                <el-icon><User /></el-icon>
                <span>用户与权限</span>
              </el-menu-item>
              <!-- 角色管理：需要 admin:role:read 权限 -->
              <el-menu-item v-if="userStore.hasPermission('admin:role:read')" index="/admin/roles">
                <el-icon><User /></el-icon>
                <span>角色管理</span>
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
      </el-aside>

      <!-- 主体内容 -->
      <el-container>
        <!-- 顶部导航栏 -->
        <el-header class="header">
          <div class="header-left">
            <h3>{{ currentTitle }}</h3>
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-icon><Avatar /></el-icon>
                {{ userStore.userInfo?.username || '用户' }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <!-- 内容区域 -->
        <el-main class="main-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  DataAnalysis,
  User,
  School,
  Reading,
  Document,
  DocumentAdd,
  Grid,
  Calendar,
  Medal,
  Bell,
  Avatar,
  ArrowDown,
  Setting,
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title || '')

const handleCommand = (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }).then(() => {
      userStore.logout()
      router.push('/login')
      ElMessage.success('退出登录成功')
    })
  } else if (command === 'profile') {
    ElMessage.info('个人信息页面开发中...')
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.sidebar {
  background-color: #304156;
  overflow-x: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3947;
  color: #fff;
}

.logo h2 {
  font-size: 18px;
  font-weight: 600;
}

.el-menu {
  border-right: none;
}

.header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left h3 {
  font-size: 18px;
  color: #333;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
}

.user-info:hover {
  color: #409eff;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}
</style>
