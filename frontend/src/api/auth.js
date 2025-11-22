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

// 修改密码
export const updatePassword = (data) => {
    return request({
        url: '/auth/password',
        method: 'put',
        data,
    })
}
