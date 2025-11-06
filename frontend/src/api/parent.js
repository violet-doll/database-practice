import request from '@/utils/request'

export function fetchParents(params) {
  return request({
    url: '/parents',
    method: 'get',
    params,
  })
}

export function createParent(data) {
  return request({
    url: '/parents',
    method: 'post',
    data,
  })
}

export function updateParent(id, data) {
  return request({
    url: `/parents/${id}`,
    method: 'put',
    data,
  })
}

export function deleteParent(id) {
  return request({
    url: `/parents/${id}`,
    method: 'delete',
  })
}


