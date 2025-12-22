import request from '@/utils/request'

/**
 * 获取所有表的列表和统计信息
 */
export const getTableList = () => {
    return request({
        url: '/api/v1/database/tables',
        method: 'get'
    })
}

/**
 * 获取指定表的数据
 * @param {string} tableName - 表名
 * @param {object} params - 查询参数 { page, page_size, filters }
 */
export const getTableData = (tableName, params) => {
    return request({
        url: `/api/v1/database/tables/${tableName}`,
        method: 'get',
        params
    })
}

/**
 * 获取表结构信息
 * @param {string} tableName - 表名
 */
export const getTableSchema = (tableName) => {
    return request({
        url: `/api/v1/database/tables/${tableName}/schema`,
        method: 'get'
    })
}

/**
 * 创建表数据
 * @param {string} tableName - 表名
 * @param {object} data - 数据
 */
export const createTableData = (tableName, data) => {
    return request({
        url: `/api/v1/database/tables/${tableName}`,
        method: 'post',
        data
    })
}

/**
 * 更新表数据
 * @param {string} tableName - 表名
 * @param {number} id - 记录ID
 * @param {object} data - 数据
 */
export const updateTableData = (tableName, id, data) => {
    return request({
        url: `/api/v1/database/tables/${tableName}/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除表数据
 * @param {string} tableName - 表名
 * @param {number} id - 记录ID
 */
export const deleteTableData = (tableName, id) => {
    return request({
        url: `/api/v1/database/tables/${tableName}/${id}`,
        method: 'delete'
    })
}

/**
 * 批量删除表数据
 * @param {string} tableName - 表名
 * @param {array} ids - 记录ID数组
 */
export const batchDeleteTableData = (tableName, ids) => {
    return request({
        url: `/api/v1/database/tables/${tableName}/batch`,
        method: 'delete',
        data: { ids }
    })
}

/**
 * 导出表数据
 * @param {string} tableName - 表名
 * @param {object} params - 查询参数
 */
export const exportTableData = (tableName, params) => {
    return request({
        url: `/api/v1/database/tables/${tableName}/export`,
        method: 'get',
        params,
        responseType: 'blob'
    })
}

/**
 * 执行SQL查询
 * @param {string} sql - SQL语句
 */
export const executeSQL = (sql) => {
    return request({
        url: '/api/v1/database/execute',
        method: 'post',
        data: { sql }
    })
}
