import { defineStore } from 'pinia'
import { login as loginApi, getCurrentUser } from '@/api/auth'

export const useUserStore = defineStore('user', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        userInfo: null,
    }),

    getters: {
        isLoggedIn: (state) => !!state.token,
    },

    actions: {
        // 登录
        async login(credentials) {
            try {
                const res = await loginApi(credentials)
                this.token = res.data.token
                this.userInfo = res.data.user
                localStorage.setItem('token', res.data.token)
                return res
            } catch (error) {
                throw error
            }
        },

        // 获取用户信息
        async fetchUserInfo() {
            try {
                const res = await getCurrentUser()
                this.userInfo = res.data
                return res
            } catch (error) {
                throw error
            }
        },

        // 登出
        logout() {
            this.token = ''
            this.userInfo = null
            localStorage.removeItem('token')
        },
    },
})
