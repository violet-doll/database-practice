import request from '@/utils/request'

export const getParents = (params) => {
    return request({ url: '/api/v1/parents', method: 'get', params })
}

export const createParent = (data) => {
    return request({ url: '/api/v1/parents', method: 'post', data })
}

export const updateParent = (id, data) => {
    return request({ url: `/api/v1/parents/${id}`, method: 'put', data })
}

export const deleteParent = (id) => {
    return request({ url: `/api/v1/parents/${id}`, method: 'delete' })
}
