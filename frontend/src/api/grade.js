import request from '@/utils/request'

export const getGrades = (params) => {
    return request({ url: '/api/v1/grades', method: 'get', params })
}

export const getGradesByStudent = (id) => {
    return request({ url: `/api/v1/grades/student/${id}`, method: 'get' })
}

export const getGradesByCourse = (id) => {
    return request({ url: `/api/v1/grades/course/${id}`, method: 'get' })
}

export const createGrade = (data) => {
    return request({ url: '/api/v1/grades', method: 'post', data })
}
