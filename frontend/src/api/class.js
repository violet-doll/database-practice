import request from '@/utils/request'

export const getClasses = (params) => {
    return request({
        url: '/api/v1/classes',
        method: 'get',
        params
    })
}

export const getClass = (id) => {
    return request({
        url: `/api/v1/classes/${id}`,
        method: 'get'
    })
}

export const createClass = (data) => {
    return request({
        url: '/api/v1/classes',
        method: 'post',
        data
    })
}

export const updateClass = (id, data) => {
    return request({
        url: `/api/v1/classes/${id}`,
        method: 'put',
        data
    })
}

export const deleteClass = (id) => {
    return request({
        url: `/api/v1/classes/${id}`,
        method: 'delete'
    })
}
