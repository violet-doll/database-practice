import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import router from '@/router'

// 创建axios实例
const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
    timeout: 15000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        const userStore = useUserStore()
        if (userStore.token) {
            config.headers['Authorization'] = `Bearer ${userStore.token}`
        }
        return config
    },
    error => {
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const res = response.data

        // 如果返回的状态码不是200，说明接口有问题
        if (response.status !== 200) {
            ElMessage.error(res.message || '请求失败')
            return Promise.reject(new Error(res.message || '请求失败'))
        }

        return res
    },
    error => {
        console.error('响应错误:', error)

        if (error.response) {
            const { status, data } = error.response

            switch (status) {
                case 401:
                    ElMessage.error('未授权，请重新登录')
                    const userStore = useUserStore()
                    userStore.logout()
                    router.push('/login')
                    break
                case 403:
                    ElMessage.error('拒绝访问，权限不足')
                    break
                case 404:
                    ElMessage.error('请求的资源不存在')
                    break
                case 500:
                    ElMessage.error(data.error || '服务器内部错误')
                    break
                default:
                    ElMessage.error(data.error || data.message || '请求失败')
            }
        } else if (error.request) {
            ElMessage.error('网络错误，请检查您的网络连接')
        } else {
            ElMessage.error('请求配置出错')
        }

        return Promise.reject(error)
    }
)

export default request
