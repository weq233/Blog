import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'

// 创建 axios 实例 (普通 API，使用 /api 前缀)
const request = axios.create({
  baseURL: '/api', // API 基础路径
  timeout: 10000, // 请求超时时间
})

// 创建管理后台专用的 axios 实例 (不使用 baseURL，直接使用完整路径)
export const adminRequest = axios.create({
  baseURL: '', // 不使用前缀
  timeout: 10000,
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 直接从 localStorage 获取 token，避免 store 未初始化的问题
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    
    // 如果返回的状态码不是 0，说明接口有错误
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      
      // 401: 未授权，需要重新登录
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        router.push('/login')
      }
      
      // 创建错误对象并附加响应信息
      const error = new Error(res.message || '请求失败')
      error.response = response
      error.data = res
      return Promise.reject(error)
    }
    
    return res
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('未授权，请重新登录')
          break
        case 403:
          ElMessage.error('拒绝访问')
          break
        case 404:
          ElMessage.error('请求资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(error.response.data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    
    return Promise.reject(error)
  }
)

export default request
