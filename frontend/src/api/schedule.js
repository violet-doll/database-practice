import request from '@/utils/request'

// 获取排课列表（分页，可按 class_id、teacher_id、semester 过滤）
export const getSchedules = (params) => {
    return request({
        url: '/schedules',
        method: 'get',
        params,
    })
}

// 获取排课详情
export const getSchedule = (id) => {
    return request({
        url: `/schedules/${id}`,
        method: 'get',
    })
}

// 创建排课记录
export const createSchedule = (data) => {
    return request({
        url: '/schedules',
        method: 'post',
        data,
    })
}

// 更新排课记录
export const updateSchedule = (id, data) => {
    return request({
        url: `/schedules/${id}`,
        method: 'put',
        data,
    })
}

// 删除排课记录
export const deleteSchedule = (id) => {
    return request({
        url: `/schedules/${id}`,
        method: 'delete',
    })
}

