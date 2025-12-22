import request from '@/utils/request'

export const getNotifications = (params) => {
    return request({ url: '/api/v1/notifications', method: 'get', params })
}

export const createNotification = (data) => {
    return request({ url: '/api/v1/notifications', method: 'post', data })
}
