import request from '@/utils/request'

export const getSchedules = (params) => {
    return request({ url: '/api/v1/schedules', method: 'get', params })
}

export const getSchedule = (id) => {
    return request({ url: `/api/v1/schedules/${id}`, method: 'get' })
}

export const createSchedule = (data) => {
    return request({ url: '/api/v1/schedules', method: 'post', data })
}

export const updateSchedule = (id, data) => {
    return request({ url: `/api/v1/schedules/${id}`, method: 'put', data })
}

export const deleteSchedule = (id) => {
    return request({ url: `/api/v1/schedules/${id}`, method: 'delete' })
}
