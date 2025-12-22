<template>
  <el-container class="main-layout">
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h3>学生管理系统</h3>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><House /></el-icon>
          <span>仪表板</span>
        </el-menu-item>
        
        <el-sub-menu index="student">
          <template #title>
            <el-icon><User /></el-icon>
            <span>学籍管理</span>
          </template>
          <el-menu-item index="/students">学生管理</el-menu-item>
          <el-menu-item index="/classes">班级管理</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="academic">
          <template #title>
            <el-icon><Reading /></el-icon>
            <span>教务管理</span>
          </template>
          <el-menu-item index="/courses">课程管理</el-menu-item>
          <el-menu-item index="/schedules">排课管理</el-menu-item>
          <el-menu-item index="/enrollments">选课管理</el-menu-item>
          <el-menu-item index="/grades">成绩管理</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="daily">
          <template #title>
            <el-icon><Calendar /></el-icon>
            <span>日常管理</span>
          </template>
          <el-menu-item index="/attendance">考勤管理</el-menu-item>
          <el-menu-item index="/rewards">奖惩管理</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="communication">
          <template #title>
            <el-icon><ChatDotRound /></el-icon>
            <span>家校互通</span>
          </template>
          <el-menu-item index="/parents">家长管理</el-menu-item>
          <el-menu-item index="/notifications">通知管理</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="system" v-if="hasAdminPermission">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin/users">用户管理</el-menu-item>
          <el-menu-item index="/admin/roles">角色管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header class="header">
        <div class="header-right">
          <span class="username">{{ userInfo?.username || '未登录' }}</span>
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-icon><Avatar /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="password">修改密码</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const userInfo = computed(() => userStore.userInfo)

const hasAdminPermission = computed(() => {
  return userStore.hasPermission('admin:user:read') || 
         userStore.hasPermission('admin:role:read')
})

const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      ElMessage.info('个人信息功能待开发')
      break
    case 'password':
      ElMessage.info('修改密码功能待开发')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        userStore.logout()
      } catch {
        // 取消操作
      }
      break
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background: #304156;
  color: white;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #263445;
}

.logo h3 {
  margin: 0;
  color: white;
  font-size: 16px;
}

.sidebar-menu {
  border: none;
  background: #304156;
}

.sidebar-menu :deep(.el-menu-item),
.sidebar-menu :deep(.el-sub-menu__title) {
  color: #bfcbd9;
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background: #263445 !important;
  color: #ffffff;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: #409eff !important;
  color: #ffffff;
}

.header {
  background: white;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.username {
  font-size: 14px;
  color: #333;
}

.user-dropdown {
  cursor: pointer;
  font-size: 20px;
  color: #409eff;
}

.main-content {
  background: #f0f2f5;
  padding: 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
