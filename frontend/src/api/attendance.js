import request from '@/utils/request'

// 考勤列表（分页+筛选）
export const getAttendance = (params) => {
  return request({
    url: '/attendance',
    method: 'get',
    params,
  })
}

// 新增考勤记录
export const createAttendance = (data) => {
  return request({
    url: '/attendance',
    method: 'post',
    data,
  })
}

// 按学生查询
export const getAttendanceByStudent = (studentId, params) => {
  return request({
    url: `/attendance/student/${studentId}`,
    method: 'get',
    params,
  })
}

// 考勤统计（按学生聚合）
export const getAttendanceStats = (params) => {
  return request({
    url: '/attendance/stats',
    method: 'get',
    params,
  })
}

// 删除考勤记录
export const deleteAttendance = (id) => {
  return request({
    url: `/attendance/${id}`,
    method: 'delete',
  })
}


