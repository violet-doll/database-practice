import request from '@/utils/request'

// 获取选课列表（分页，可按 student_id、course_id 过滤）
export const getEnrollments = (params) => {
    return request({
        url: '/enrollments',
        method: 'get',
        params,
    })
}

// 创建选课记录
export const createEnrollment = (data) => {
    return request({
        url: '/enrollments',
        method: 'post',
        data,
    })
}

// 删除选课记录
export const deleteEnrollment = (id) => {
    return request({
        url: `/enrollments/${id}`,
        method: 'delete',
    })
}

