import request from '@/utils/request'

export function fetchNotifications(params) {
  return request({
    url: '/notifications',
    method: 'get',
    params,
  })
}

export function createNotification(data) {
  return request({
    url: '/notifications',
    method: 'post',
    data,
  })
}


