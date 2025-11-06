import request from '@/utils/request'

// 奖惩列表（分页+筛选）
export const getRewards = (params) => {
  return request({
    url: '/rewards',
    method: 'get',
    params,
  })
}

// 新增奖惩记录
export const createReward = (data) => {
  return request({
    url: '/rewards',
    method: 'post',
    data,
  })
}

// 删除奖惩记录
export const deleteReward = (id) => {
  return request({
    url: `/rewards/${id}`,
    method: 'delete',
  })
}


