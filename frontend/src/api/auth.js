import request from '@/utils/request'

// 用户登录
export const login = (data) => {
    return request({
        url: '/auth/login',
        method: 'post',
        data,
    })
}

// 用户登出
export const logout = () => {
    return request({
        url: '/auth/logout',
        method: 'post',
    })
}

// 获取当前用户信息
export const getCurrentUser = () => {
    return request({
        url: '/auth/me',
        method: 'get',
    })
}
