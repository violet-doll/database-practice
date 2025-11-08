import request from '@/utils/request'

// 获取班级列表
export const getClasses = (params) => {
    return request({
        url: '/classes',
        method: 'get',
        params,
    })
}

// 获取班级详情
export const getClass = (id) => {
    return request({
        url: `/classes/${id}`,
        method: 'get',
    })
}

// 新增班级
export const createClass = (data) => {
    return request({
        url: '/classes',
        method: 'post',
        data,
    })
}

// 更新班级
export const updateClass = (id, data) => {
    return request({
        url: `/classes/${id}`,
        method: 'put',
        data,
    })
}

// 删除班级
export const deleteClass = (id) => {
    return request({
        url: `/classes/${id}`,
        method: 'delete',
    })
}


