import request from '@/utils/request'

export const getAttendance = (params) => {
    return request({ url: '/api/v1/attendance', method: 'get', params })
}

export const getAttendanceByStudent = (id) => {
    return request({ url: `/api/v1/attendance/student/${id}`, method: 'get' })
}

export const createAttendance = (data) => {
    return request({ url: '/api/v1/attendance', method: 'post', data })
}

export const deleteAttendance = (id) => {
    return request({ url: `/api/v1/attendance/${id}`, method: 'delete' })
}
