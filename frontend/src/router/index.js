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
                meta: { title: '数据看板', permission: 'admin:stats:read' },
            },
            {
                path: '/students',
                name: 'Students',
                component: () => import('@/views/Students.vue'),
                meta: { title: '学生管理', permission: 'student:read' },
            },
            {
                path: '/classes',
                name: 'Classes',
                component: () => import('@/views/Classes.vue'),
                meta: { title: '班级管理', permission: 'class:read' },
            },
            {
                path: '/courses',
                name: 'Courses',
                component: () => import('@/views/Courses.vue'),
                meta: { title: '课程管理', permission: 'course:read' },
            },
            {
                path: '/enrollments',
                name: 'Enrollments',
                component: () => import('@/views/Enrollments.vue'),
                meta: { title: '选课管理', permission: 'enrollment:read' },
            },
            {
                path: '/schedule',
                name: 'Schedule',
                component: () => import('@/views/Schedule.vue'),
                meta: { title: '排课管理', permission: 'schedule:read' },
            },
            {
                path: '/grades',
                name: 'Grades',
                component: () => import('@/views/Grades.vue'),
                meta: { title: '成绩管理', permission: 'grade:read' },
            },
            {
                path: '/attendance',
                name: 'Attendance',
                component: () => import('@/views/Attendance.vue'),
                meta: { title: '考勤管理', permission: 'attendance:read' },
            },
            {
                path: '/attendance-stats',
                name: 'AttendanceStats',
                component: () => import('@/views/AttendanceStats.vue'),
                meta: { title: '考勤统计', permission: 'attendance:read' },
            },
            {
                path: '/rewards',
                name: 'Rewards',
                component: () => import('@/views/Rewards.vue'),
                meta: { title: '奖惩管理', permission: 'reward:read' },
            },
            {
                path: '/notifications',
                name: 'Notifications',
                component: () => import('@/views/Notifications.vue'),
                meta: { title: '通知管理', permission: 'notification:read' },
            },
            {
                path: '/parents',
                name: 'Parents',
                component: () => import('@/views/Parents.vue'),
                meta: { title: '家长联系方式', permission: 'parent:read' },
            },
            // 管理员：系统设置
            {
                path: '/admin/overview',
                name: 'AdminOverview',
                component: () => import('@/views/AdminOverview.vue'),
                meta: { title: '统计概览（管理员）', permission: 'admin:stats:read' },
            },
            {
                path: '/admin/users',
                name: 'AdminUsers',
                component: () => import('@/views/AdminUsers.vue'),
                meta: { title: '用户与权限（管理员）', permission: 'admin:user:read' },
            },
            {
                path: '/admin/roles',
                name: 'AdminRoles',
                component: () => import('@/views/AdminRoles.vue'),
                meta: { title: '角色管理（管理员）', permission: 'admin:role:read' },
            },
            {
                path: '/profile',
                name: 'Profile',
                component: () => import('@/views/Profile.vue'),
                meta: { title: '个人信息' },
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

    // 权限检查
    if (userStore.isLoggedIn && to.meta.permission) {
        if (!userStore.hasPermission(to.meta.permission)) {
            // 无权限，寻找第一个有权限的路由进行跳转
            // 避免无限循环：如果当前就是尝试跳转的目标，则停止（虽然逻辑上不应该发生，因为我们是在找有权限的）

            // 获取 Layout 的子路由
            const layoutRoute = routes.find(r => r.name === 'Layout')
            if (layoutRoute && layoutRoute.children) {
                for (const child of layoutRoute.children) {
                    // 检查是否有权限 (如果没有 permission 字段，或者是 profile，则视为有权限)
                    if (!child.meta.permission || userStore.hasPermission(child.meta.permission)) {
                        // 构建完整路径 (这里假设都是一级子路由)
                        const targetPath = child.path
                        // 如果目标路径就是当前路径，说明当前路径其实是有权限的？不对，我们进这个if是因为 !hasPermission
                        // 所以 targetPath 肯定不是 to.path
                        return next(targetPath)
                    }
                }
            }
            // 如果找不到任何有权限的路由，跳转到 Profile (保底)
            return next('/profile')
        }
    }

    // 特殊处理根路径跳转：如果去 /dashboard 但没权限，也会被上面的逻辑拦截并重定向
    // 但如果是直接访问 /，redirect 是 /dashboard，也会触发上面的逻辑

    return next()
})

export default router
