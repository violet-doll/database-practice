import request from '@/utils/request'

/**
 * 获取学生列表
 */
export const getStudents = (params) => {
    return request({
        url: '/api/v1/students',
        method: 'get',
        params
    })
}

/**
 * 获取学生详情
 */
export const getStudent = (id) => {
    return request({
        url: `/api/v1/students/${id}`,
        method: 'get'
    })
}

/**
 * 创建学生
 */
export const createStudent = (data) => {
    return request({
        url: '/api/v1/students',
        method: 'post',
        data
    })
}

/**
 * 更新学生
 */
export const updateStudent = (id, data) => {
    return request({
        url: `/api/v1/students/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除学生
 */
export const deleteStudent = (id) => {
    return request({
        url: `/api/v1/students/${id}`,
        method: 'delete'
    })
}
