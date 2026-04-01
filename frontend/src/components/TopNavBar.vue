<template>
  <div class="top-nav">
    <div class="nav-brand" @click="goToHome">博客系统</div>
    <div class="nav-user" v-if="isLoggedIn">
      <el-dropdown trigger="click">
        <div class="user-info">
          <el-avatar :size="40" :src="getAvatarUrl(userInfo?.avatar)">
            {{ getUserInitial(userInfo?.nickname) }}
          </el-avatar>
          <span class="username">{{ userInfo?.nickname || userInfo?.username }}</span>
          <el-icon class="el-icon--right"><arrow-down /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="goToProfile">
              <el-icon><user /></el-icon>
              个人资料
            </el-dropdown-item>
            <el-dropdown-item divided @click="handleLogout">
              <el-icon><switch-button /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <div class="nav-actions" v-else>
      <el-dropdown trigger="click">
        <div class="user-info">
          <el-avatar :size="40">
            <el-icon><user /></el-icon>
          </el-avatar>
          <span class="username">未登录</span>
          <el-icon class="el-icon--right"><arrow-down /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="goToLogin">
              <el-icon><key /></el-icon>
              登录
            </el-dropdown-item>
            <el-dropdown-item divided @click="goToRegister">
              <el-icon><user-filled /></el-icon>
              注册
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { ArrowDown, User, Key, UserFilled, SwitchButton } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.token && !userStore.isTokenExpired())
const userInfo = computed(() => userStore.userInfo)

// 获取用户首字母作为头像文字
const getUserInitial = (nickname) => {
  if (!nickname) return '用'
  return nickname.charAt(0).toUpperCase()
}

// 根据昵称生成默认头像颜色
const getAvatarFromNickname = (nickname) => {
  if (!nickname) return ''
  const colors = [
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', 
    '#909399', '#A0CFFF', '#B3D8FF', '#C6E2FF'
  ]
  const index = nickname.charCodeAt(0) % colors.length
  return colors[index]
}

// 获取头像 URL，如果没有则返回颜色
const getAvatarUrl = (avatar) => {
  if (avatar && avatar.startsWith('http')) {
    return avatar
  }
  return getAvatarFromNickname(userInfo?.nickname)
}

// 加载用户信息
onMounted(async () => {
  // 如果已登录但没有用户信息，尝试获取
  if (isLoggedIn.value && !userInfo.value) {
    try {
      await userStore.fetchUserInfo()
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }
})

// 跳转首页
const goToHome = () => {
  router.push('/')
}

// 跳转登录页面
const goToLogin = () => {
  router.push('/login')
}

// 跳转注册页面
const goToRegister = () => {
  router.push('/register')
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  setTimeout(() => {
    router.push('/')
  }, 500)
}

// 跳转个人资料页面
const goToProfile = () => {
  router.push('/profile')
}
</script>

<style lang="scss" scoped>
.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 40px;
  margin-bottom: 30px;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  
  .nav-brand {
    color: #fff;
    font-size: 24px;
    font-weight: 600;
    cursor: pointer;
    letter-spacing: 2px;
    
    &:hover {
      opacity: 0.8;
    }
  }
  
  .nav-user,
  .nav-actions {
    .user-info {
      display: flex;
      align-items: center;
      gap: 10px;
      cursor: pointer;
      padding: 8px 12px;
      border-radius: 6px;
      transition: all 0.3s ease;
      
      &:hover {
        background: rgba(255, 255, 255, 0.08);
      }
      
      .username {
        color: #fff;
        font-size: 15px;
        font-weight: 500;
      }
      
      .el-icon {
        color: rgba(255, 255, 255, 0.6);
        transition: transform 0.3s ease;
      }
    }
    
    // 下拉菜单打开时旋转箭头
    :deep(.el-dropdown-menu__item) {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 10px 16px;
      
      .el-icon {
        font-size: 16px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .top-nav {
    padding: 15px 20px;
    
    .nav-brand {
      font-size: 20px;
    }
    
    .nav-user,
    .nav-actions {
      gap: 8px;
      
      .username {
        display: none; // 移动端隐藏用户名
      }
    }
  }
}
</style>
