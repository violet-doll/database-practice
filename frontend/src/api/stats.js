import request from '@/utils/request'

// 获取数据看板统计信息
export function fetchDashboardStats() {
    return request({ url: '/stats/dashboard', method: 'get' })
}

