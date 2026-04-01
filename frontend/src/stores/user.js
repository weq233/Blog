import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCurrentUser } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)
  const tokenExpiry = ref(localStorage.getItem('tokenExpiry') || '')

  // 初始化时立即恢复用户信息
  function initUser() {
    const savedUserInfo = localStorage.getItem('userInfo')
    if (savedUserInfo && token.value && !isTokenExpired()) {
      try {
        userInfo.value = JSON.parse(savedUserInfo)
        return true
      } catch (e) {
        localStorage.removeItem('userInfo')
      }
    }
    return false
  }

  // 立即执行初始化
  initUser()

  // 检查 token 是否过期
  function isTokenExpired() {
    if (!tokenExpiry.value) return true
    const expiryTime = new Date(tokenExpiry.value).getTime()
    const currentTime = new Date().getTime()
    // 提前 5 分钟判定为过期
    return currentTime >= (expiryTime - 5 * 60 * 1000)
  }

  function setToken(newToken, expiresIn = 24 * 60 * 60 * 1000) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
      // 计算过期时间（默认 24 小时）
      const expiryTime = new Date(new Date().getTime() + expiresIn)
      tokenExpiry.value = expiryTime.toISOString()
      localStorage.setItem('tokenExpiry', tokenExpiry.value)
    } else {
      localStorage.removeItem('token')
      localStorage.removeItem('tokenExpiry')
      tokenExpiry.value = ''
    }
  }

  function setUserInfo(info) {
    userInfo.value = info
    // 将用户信息也持久化到 localStorage
    if (info) {
      localStorage.setItem('userInfo', JSON.stringify(info))
    } else {
      localStorage.removeItem('userInfo')
    }
  }

  // 获取当前用户信息
  async function fetchUserInfo() {
    if (!token.value) {
      return null
    }
    
    // 如果 token 已过期，直接清除
    if (isTokenExpired()) {
      logout()
      return null
    }

    try {
      const res = await getCurrentUser()
      if (res.code === 0 && res.data) {
        setUserInfo(res.data)
        return res.data
      }
      // 如果 token 无效，清除登录状态
      logout()
      return null
    } catch (error) {
      // 网络错误时不清除用户信息，保持登录状态
      return null
    }
  }

  // 自动登录（页面加载时调用）
  async function autoLogin() {
    // 先尝试从本地存储恢复
    const restored = initUser()
    if (restored) {
      // 后台验证 token 有效性
      await fetchUserInfo()
      return !!userInfo.value
    }
    
    // 如果没有本地数据，但有 token，尝试获取用户信息
    if (token.value && !isTokenExpired()) {
      await fetchUserInfo()
      return !!userInfo.value
    }
    
    return false
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    tokenExpiry.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('tokenExpiry')
    localStorage.removeItem('userInfo')
  }

  return {
    token,
    userInfo,
    tokenExpiry,
    isTokenExpired,
    setToken,
    setUserInfo,
    initUser,
    fetchUserInfo,
    autoLogin,
    logout,
  }
})