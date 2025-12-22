import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        meta: { requiresAuth: true },
        redirect: '/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/Dashboard.vue'),
                meta: { title: '仪表板' }
            },
            // 学生管理
            {
                path: 'students',
                name: 'StudentList',
                component: () => import('@/views/student/StudentList.vue'),
                meta: { title: '学生管理', permission: 'student:read' }
            },
            // 班级管理
            {
                path: 'classes',
                name: 'ClassList',
                component: () => import('@/views/class/ClassList.vue'),
                meta: { title: '班级管理', permission: 'class:read' }
            },
            // 课程管理
            {
                path: 'courses',
                name: 'CourseList',
                component: () => import('@/views/course/CourseList.vue'),
                meta: { title: '课程管理', permission: 'course:read' }
            },
            // 排课管理
            {
                path: 'schedules',
                name: 'ScheduleList',
                component: () => import('@/views/schedule/ScheduleList.vue'),
                meta: { title: '排课管理', permission: 'schedule:read' }
            },
            // 选课管理
            {
                path: 'enrollments',
                name: 'EnrollmentList',
                component: () => import('@/views/enrollment/EnrollmentList.vue'),
                meta: { title: '选课管理', permission: 'enrollment:read' }
            },
            // 成绩管理
            {
                path: 'grades',
                name: 'GradeList',
                component: () => import('@/views/grade/GradeList.vue'),
                meta: { title: '成绩管理', permission: 'grade:read' }
            },
            // 考勤管理
            {
                path: 'attendance',
                name: 'AttendanceList',
                component: () => import('@/views/attendance/AttendanceList.vue'),
                meta: { title: '考勤管理', permission: 'attendance:read' }
            },
            // 奖惩管理
            {
                path: 'rewards',
                name: 'RewardList',
                component: () => import('@/views/reward/RewardList.vue'),
                meta: { title: '奖惩管理', permission: 'reward:read' }
            },
            // 家长管理
            {
                path: 'parents',
                name: 'ParentList',
                component: () => import('@/views/parent/ParentList.vue'),
                meta: { title: '家长管理', permission: 'parent:read' }
            },
            // 通知管理
            {
                path: 'notifications',
                name: 'NotificationList',
                component: () => import('@/views/notification/NotificationList.vue'),
                meta: { title: '通知管理', permission: 'notification:read' }
            },
            // 系统管理
            {
                path: 'admin/users',
                name: 'UserList',
                component: () => import('@/views/admin/UserList.vue'),
                meta: { title: '用户管理', permission: 'admin:user:read' }
            },
            {
                path: 'admin/roles',
                name: 'RoleList',
                component: () => import('@/views/admin/RoleList.vue'),
                meta: { title: '角色管理', permission: 'admin:role:read' }
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()

    if (to.meta.requiresAuth !== false) {
        // 需要登录
        if (!userStore.token) {
            ElMessage.warning('请先登录')
            next('/login')
            return
        }

        // 检查权限
        if (to.meta.permission && !userStore.hasPermission(to.meta.permission)) {
            ElMessage.error('权限不足')
            next(from.path)
            return
        }
    }

    next()
})

export default router
