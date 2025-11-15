import { defineStore } from 'pinia'
import { login as loginApi, getCurrentUser } from '@/api/auth'

export const useUserStore = defineStore('user', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        userInfo: null,
        permissions: JSON.parse(localStorage.getItem('permissions') || '[]'), // 权限列表
    }),

    getters: {
        isLoggedIn: (state) => !!state.token,
        // 检查是否有某个权限
        hasPermission: (state) => (permission) => {
            return state.permissions.includes(permission)
        },
    },

    actions: {
        // 登录
        async login(credentials) {
            try {
                const res = await loginApi(credentials)
                this.token = res.data.token
                this.userInfo = res.data.user
                this.permissions = res.data.permissions || []
                localStorage.setItem('token', res.data.token)
                localStorage.setItem('permissions', JSON.stringify(this.permissions))
                return res
            } catch (error) {
                throw error
            }
        },

        // 获取用户信息
        async fetchUserInfo() {
            try {
                const res = await getCurrentUser()
                this.userInfo = res.data.user || res.data
                this.permissions = res.data.permissions || []
                localStorage.setItem('permissions', JSON.stringify(this.permissions))
                return res
            } catch (error) {
                throw error
            }
        },

        // 登出
        logout() {
            this.token = ''
            this.userInfo = null
            this.permissions = []
            localStorage.removeItem('token')
            localStorage.removeItem('permissions')
        },
    },
})
