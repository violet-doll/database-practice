import request from '@/utils/request'

// 用户管理
export function fetchUsers(params) {
    return request({ url: '/admin/users', method: 'get', params })
}

export function createUser(data) {
    return request({ url: '/admin/users', method: 'post', data })
}

export function updateUser(id, data) {
    return request({ url: `/admin/users/${id}`, method: 'put', data })
}

export function deleteUser(id) {
    return request({ url: `/admin/users/${id}`, method: 'delete' })
}

// 角色管理
export function fetchRoles(params) {
    return request({ url: '/admin/roles', method: 'get', params })
}

export function createRole(data) {
    return request({ url: '/admin/roles', method: 'post', data })
}

export function updateRole(id, data) {
    return request({ url: `/admin/roles/${id}`, method: 'put', data })
}

export function deleteRole(id) {
    return request({ url: `/admin/roles/${id}`, method: 'delete' })
}

// 统计概览
export function fetchAdminOverview() {
    return request({ url: '/admin/stats/overview', method: 'get' })
}

// 权限管理
export function fetchPermissions() {
    return request({ url: '/admin/permissions', method: 'get' })
}

export function fetchRolePermissions(roleId) {
    return request({ url: `/admin/roles/${roleId}/permissions`, method: 'get' })
}

export function updateRolePermissions(roleId, permissions) {
    return request({ url: `/admin/roles/${roleId}/permissions`, method: 'post', data: { permissions } })
}


