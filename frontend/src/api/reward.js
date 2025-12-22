import request from '@/utils/request'

export const getRewards = (params) => {
    return request({ url: '/api/v1/rewards', method: 'get', params })
}

export const getRewardsByStudent = (id) => {
    return request({ url: `/api/v1/rewards/student/${id}`, method: 'get' })
}

export const createReward = (data) => {
    return request({ url: '/api/v1/rewards', method: 'post', data })
}

export const deleteReward = (id) => {
    return request({ url: `/api/v1/rewards/${id}`, method: 'delete' })
}
