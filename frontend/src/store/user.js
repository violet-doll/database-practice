import { defineStore } from 'pinia'
import { ref } from 'vue'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
    // 状态
    const token = ref(localStorage.getItem('token') || '')
    const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || 'null'))
    const permissions = ref(JSON.parse(localStorage.getItem('permissions') || '[]'))

    // 设置Token
    const setToken = (newToken) => {
        token.value = newToken
        localStorage.setItem('token', newToken)
    }

    // 设置用户信息
    const setUserInfo = (info) => {
        userInfo.value = info
        localStorage.setItem('userInfo', JSON.stringify(info))
    }

    // 设置权限列表
    const setPermissions = (perms) => {
        permissions.value = perms
        localStorage.setItem('permissions', JSON.stringify(perms))
    }

    // 登录
    const login = (loginData) => {
        setToken(loginData.token)
        setUserInfo(loginData.user)
        setPermissions(loginData.permissions || [])
    }

    // 登出
    const logout = () => {
        token.value = ''
        userInfo.value = null
        permissions.value = []
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        localStorage.removeItem('permissions')
        router.push('/login')
    }

    // 检查是否有某个权限
    const hasPermission = (permission) => {
        return permissions.value.includes(permission)
    }

    return {
        token,
        userInfo,
        permissions,
        setToken,
        setUserInfo,
        setPermissions,
        login,
        logout,
        hasPermission
    }
})
