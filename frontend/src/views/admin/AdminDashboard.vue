<template>
  <div class="admin-dashboard">
    <!-- 顶部导航栏 -->
    <TopNavBar />
    
    <div class="admin-container">
      <!-- 侧边栏 -->
      <aside class="admin-sidebar">
        <div class="sidebar-header">
          <h2>📊 管理后台</h2>
        </div>
        <nav class="sidebar-nav">
          <router-link to="/admin" class="nav-item active">
            <span class="nav-icon">🏠</span>
            <span>数据概览</span>
          </router-link>
          <router-link to="/admin/articles" class="nav-item">
            <span class="nav-icon">📝</span>
            <span>文章管理</span>
          </router-link>
          <router-link to="/admin/categories" class="nav-item">
            <span class="nav-icon">📁</span>
            <span>分类管理</span>
          </router-link>
          <router-link to="/admin/tags" class="nav-item">
            <span class="nav-icon">🏷️</span>
            <span>标签管理</span>
          </router-link>
          <router-link to="/admin/users" class="nav-item">
            <span class="nav-icon">👥</span>
            <span>用户管理</span>
          </router-link>
        </nav>
      </aside>
      
      <!-- 主内容区 -->
      <main class="admin-content">
        <div class="content-header">
          <h1>数据概览</h1>
          <p>欢迎回来，管理员！</p>
        </div>
        
        <!-- 统计卡片 -->
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon article">📝</div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.articleCount }}</div>
              <div class="stat-label">文章总数</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon user">👥</div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon category">📁</div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.categoryCount }}</div>
              <div class="stat-label">分类总数</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon tag">🏷️</div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.tagCount }}</div>
              <div class="stat-label">标签总数</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon comment">💬</div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.commentCount }}</div>
              <div class="stat-label">评论总数</div>
            </div>
          </div>
        </div>
        
        <!-- 最新文章 -->
        <div class="content-section">
          <h2>最新文章</h2>
          <el-table :data="recentArticles" style="width: 100%" class="dark-table">
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="author.username" label="作者" width="120" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'">
                  {{ row.status === 1 ? '已发布' : '草稿' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="发布时间" width="180" />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="viewArticle(row)">查看</el-button>
                <el-button size="small" type="warning" @click="editArticle(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 活跃用户 -->
        <div class="content-section">
          <h2>活跃用户</h2>
          <el-table :data="activeUsers" style="width: 100%" class="dark-table">
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column prop="created_at" label="注册时间" width="180" />
          </el-table>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import TopNavBar from '@/components/TopNavBar.vue'
import { adminRequest } from '@/utils/request'

const router = useRouter()

const stats = ref({
  articleCount: 0,
  userCount: 0,
  categoryCount: 0,
  tagCount: 0,
  commentCount: 0
})

const recentArticles = ref([])
const activeUsers = ref([])

// 加载统计数据
const loadStats = async () => {
  try {
    const token = localStorage.getItem('token') // 修改为 'token'
    if (!token) {
      ElMessage.error('请先登录')
      router.push('/login')
      return
    }
    
    const response = await adminRequest.get('/admin/stats', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    
    if (response.data.code === 0) {
      stats.value = response.data.data.stats || {}
      
      // 如果有返回数据，设置最新文章和活跃用户
      if (response.data.data.recentArticles) {
        recentArticles.value = response.data.data.recentArticles
      }
      if (response.data.data.activeUsers) {
        activeUsers.value = response.data.data.activeUsers
      }
    } else {
      ElMessage.error(response.data.message || '加载失败')
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    
    if (error.response?.status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      router.push('/login')
    } else {
      ElMessage.error('加载统计数据失败')
    }
  }
}

// 查看文章
const viewArticle = (article) => {
  router.push(`/article/${article.id}`)
}

// 编辑文章
const editArticle = (article) => {
  router.push(`/admin/article/edit/${article.id}`)
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped lang="scss">
.admin-dashboard {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e1e2f 0%, #2d2d44 100%);
}

.admin-container {
  display: flex;
  min-height: calc(100vh - 60px);
}

.admin-sidebar {
  width: 250px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 0;
  
  .sidebar-header {
    padding: 0 20px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    
    h2 {
      color: rgba(255, 255, 255, 0.95);
      font-size: 20px;
      margin: 0;
    }
  }
  
  .sidebar-nav {
    margin-top: 20px;
    
    .nav-item {
      display: flex;
      align-items: center;
      padding: 12px 20px;
      color: rgba(255, 255, 255, 0.7);
      text-decoration: none;
      transition: all 0.3s;
      
      &:hover {
        background: rgba(255, 255, 255, 0.05);
        color: rgba(255, 255, 255, 0.95);
      }
      
      &.active {
        background: rgba(102, 126, 234, 0.2);
        color: rgba(102, 126, 234, 1);
        border-right: 3px solid rgba(102, 126, 234, 1);
      }
      
      .nav-icon {
        margin-right: 10px;
        font-size: 18px;
      }
    }
  }
}

.admin-content {
  flex: 1;
  padding: 30px;
  overflow-y: auto;
  
  .content-header {
    margin-bottom: 30px;
    
    h1 {
      color: rgba(255, 255, 255, 0.95);
      font-size: 28px;
      margin: 0 0 10px 0;
    }
    
    p {
      color: rgba(255, 255, 255, 0.6);
      margin: 0;
    }
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
    margin-bottom: 30px;
    
    .stat-card {
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 12px;
      padding: 20px;
      display: flex;
      align-items: center;
      gap: 15px;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
        border-color: rgba(102, 126, 234, 0.5);
      }
      
      .stat-icon {
        width: 50px;
        height: 50px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 24px;
        
        &.article {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        
        &.user {
          background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }
        
        &.category {
          background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }
        
        &.tag {
          background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        }
        
        &.comment {
          background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
        }
      }
      
      .stat-info {
        flex: 1;
        
        .stat-value {
          color: rgba(255, 255, 255, 0.95);
          font-size: 28px;
          font-weight: bold;
          margin-bottom: 5px;
        }
        
        .stat-label {
          color: rgba(255, 255, 255, 0.6);
          font-size: 14px;
        }
      }
    }
  }
  
  .content-section {
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
    
    h2 {
      color: rgba(255, 255, 255, 0.95);
      font-size: 20px;
      margin: 0 0 20px 0;
      padding-bottom: 10px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }
  }
}
</style>
