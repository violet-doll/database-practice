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
        ],
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()

    if (to.meta.requiresAuth && !userStore.isLoggedIn) {
        // 需要认证但未登录，跳转到登录页
        next('/login')
    } else if (to.path === '/login' && userStore.isLoggedIn) {
        // 已登录但访问登录页，跳转到首页
        next('/')
    } else {
        next()
    }
})

export default router
