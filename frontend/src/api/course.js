import request from '@/utils/request'

// 获取课程列表
export const getCourses = (params) => {
    return request({
        url: '/courses',
        method: 'get',
        params,
    })
}

// 获取课程详情
export const getCourse = (id) => {
    return request({
        url: `/courses/${id}`,
        method: 'get',
    })
}

// 创建课程
export const createCourse = (data) => {
    return request({
        url: '/courses',
        method: 'post',
        data,
    })
}

// 更新课程
export const updateCourse = (id, data) => {
    return request({
        url: `/courses/${id}`,
        method: 'put',
        data,
    })
}

// 删除课程
export const deleteCourse = (id) => {
    return request({
        url: `/courses/${id}`,
        method: 'delete',
    })
}


