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
        redirect: '/database',
        children: [
            // 数据库表管理
            {
                path: 'database',
                name: 'TableManager',
                component: () => import('@/views/database/TableManager.vue'),
                meta: { title: '数据库表管理' }
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
