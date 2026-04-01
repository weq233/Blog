<template>
  <div class="home">
    <TopNavBar />

    <div class="home-content">
      <h1>欢迎来到博客系统</h1>
      <p class="subtitle">探索精彩文章，分享知识与见解</p>
      <div class="action-buttons">
        <el-button type="primary" size="large" @click="goToArticles">
          浏览文章
        </el-button>
        <el-button size="large" @click="goToLogin" v-if="!isLoggedIn">
          登录
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { computed } from 'vue'
import TopNavBar from '@/components/TopNavBar.vue'

const router = useRouter()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.token && !userStore.isTokenExpired())

const goToArticles = () => {
  router.push('/articles')
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<style lang="scss" scoped>
.home {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.85) 0%, rgba(50, 30, 60, 0.8) 100%),
              url('https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1920&q=80') center/cover no-repeat;

  .home-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    animation: fadeInUp 0.8s ease-out;
    
    @keyframes fadeInUp {
      from {
        opacity: 0;
        transform: translateY(30px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
    
    h1 {
      color: #fff;
      font-size: 48px;
      margin-bottom: 20px;
      font-weight: 300;
      letter-spacing: 3px;
    }
    
    .subtitle {
      color: rgba(255, 255, 255, 0.7);
      font-size: 18px;
      margin-bottom: 40px;
      font-weight: 300;
    }
    
    .action-buttons {
      display: flex;
      gap: 20px;
      justify-content: center;
      
      .el-button {
        min-width: 120px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .home {
    .home-content {
      h1 {
        font-size: 32px;
      }
      
      .subtitle {
        font-size: 16px;
      }
    }
  }
}
</style>
