import request from '@/utils/request'

// 用户管理
export const getUsers = (params) => {
    return request({ url: '/api/v1/admin/users', method: 'get', params })
}

export const createUser = (data) => {
    return request({ url: '/api/v1/admin/users', method: 'post', data })
}

export const updateUser = (id, data) => {
    return request({ url: `/api/v1/admin/users/${id}`, method: 'put', data })
}

export const deleteUser = (id) => {
    return request({ url: `/api/v1/admin/users/${id}`, method: 'delete' })
}

// 角色管理
export const getRoles = (params) => {
    return request({ url: '/api/v1/admin/roles', method: 'get', params })
}

export const createRole = (data) => {
    return request({ url: '/api/v1/admin/roles', method: 'post', data })
}

export const updateRole = (id, data) => {
    return request({ url: `/api/v1/admin/roles/${id}`, method: 'put', data })
}

export const deleteRole = (id) => {
    return request({ url: `/api/v1/admin/roles/${id}`, method: 'delete' })
}

// 权限管理
export const getPermissions = () => {
    return request({ url: '/api/v1/admin/permissions', method: 'get' })
}

export const getRolePermissions = (roleId) => {
    return request({ url: `/api/v1/admin/roles/${roleId}/permissions`, method: 'get' })
}

export const updateRolePermissions = (roleId, data) => {
    return request({ url: `/api/v1/admin/roles/${roleId}/permissions`, method: 'post', data })
}
