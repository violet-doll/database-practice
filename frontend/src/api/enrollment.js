import request from '@/utils/request'

export const getEnrollments = (params) => {
    return request({ url: '/api/v1/enrollments', method: 'get', params })
}

export const createEnrollment = (data) => {
    return request({ url: '/api/v1/enrollments', method: 'post', data })
}

export const deleteEnrollment = (id) => {
    return request({ url: `/api/v1/enrollments/${id}`, method: 'delete' })
}
