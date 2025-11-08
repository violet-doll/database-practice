import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { requiresAuth: false },
    },
    {
        path: '/',
        name: 'Layout',
        component: () => import('@/layouts/MainLayout.vue'),
        redirect: '/dashboard',
        meta: { requiresAuth: true },
        children: [
            {
                path: '/dashboard',
                name: 'Dashboard',
                component: () => import('@/views/Dashboard.vue'),
                meta: { title: '数据看板' },
            },
            {
                path: '/students',
                name: 'Students',
                component: () => import('@/views/Students.vue'),
                meta: { title: '学生管理' },
            },
            {
                path: '/classes',
                name: 'Classes',
                component: () => import('@/views/Classes.vue'),
                meta: { title: '班级管理' },
            },
            {
                path: '/courses',
                name: 'Courses',
                component: () => import('@/views/Courses.vue'),
                meta: { title: '课程管理' },
            },
            {
                path: '/enrollments',
                name: 'Enrollments',
                component: () => import('@/views/Enrollments.vue'),
                meta: { title: '选课管理' },
            },
            {
                path: '/schedule',
                name: 'Schedule',
                component: () => import('@/views/Schedule.vue'),
                meta: { title: '排课管理' },
            },
            {
                path: '/grades',
                name: 'Grades',
                component: () => import('@/views/Grades.vue'),
                meta: { title: '成绩管理' },
            },
            {
                path: '/attendance',
                name: 'Attendance',
                component: () => import('@/views/Attendance.vue'),
                meta: { title: '考勤管理' },
            },
            {
                path: '/attendance-stats',
                name: 'AttendanceStats',
                component: () => import('@/views/AttendanceStats.vue'),
                meta: { title: '考勤统计' },
            },
            {
                path: '/rewards',
                name: 'Rewards',
                component: () => import('@/views/Rewards.vue'),
                meta: { title: '奖惩管理' },
            },
            {
                path: '/notifications',
                name: 'Notifications',
                component: () => import('@/views/Notifications.vue'),
                meta: { title: '通知管理' },
            },
            {
                path: '/parents',
                name: 'Parents',
                component: () => import('@/views/Parents.vue'),
                meta: { title: '家长联系方式' },
            },
            // 管理员：系统设置
            {
                path: '/admin/overview',
                name: 'AdminOverview',
                component: () => import('@/views/AdminOverview.vue'),
                meta: { title: '统计概览（管理员）' },
            },
            {
                path: '/admin/users',
                name: 'AdminUsers',
                component: () => import('@/views/AdminUsers.vue'),
                meta: { title: '用户与权限（管理员）' },
            },
            {
                path: '/admin/roles',
                name: 'AdminRoles',
                component: () => import('@/views/AdminRoles.vue'),
                meta: { title: '角色管理（管理员）' },
            },
        ],
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

// 路由守卫（刷新后恢复用户信息，保证管理员菜单渲染）
router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()

    // 若已登录但 userInfo 为空，则尝试拉取当前用户信息（刷新场景）
    if (userStore.isLoggedIn && !userStore.userInfo && to.path !== '/login') {
        try {
            await userStore.fetchUserInfo()
        } catch (e) {
            // token 失效，清理并跳登录
            userStore.logout()
            return next('/login')
        }
    }

    if (to.meta.requiresAuth && !userStore.isLoggedIn) {
        return next('/login')
    }
    if (to.path === '/login' && userStore.isLoggedIn) {
        return next('/')
    }
    return next()
})

export default router
