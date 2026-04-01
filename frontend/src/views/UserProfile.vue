<template>
  <div class="user-profile-page">
    <TopNavBar />
    
    <div class="user-back" v-if="!loading">
      <el-button @click="goBack" text>
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
    </div>

    <div class="user-container" v-loading="loading">
      <!-- 用户基本信息 -->
      <el-card v-if="user" class="user-card">
        <div class="user-header">
          <el-avatar :size="100" :src="user.avatar || defaultAvatar" />
          <h2 class="user-nickname">{{ user.nickname || user.username }}</h2>
          <p class="user-username">@{{ user.username }}</p>
        </div>
        
        <div class="user-info">
          <div class="info-item">
            <el-icon><User /></el-icon>
            <span>角色：{{ roleText }}</span>
          </div>
          <div class="info-item">
            <el-icon><Calendar /></el-icon>
            <span>加入时间：{{ formatDate(user.created_at) }}</span>
          </div>
          <div class="info-item" v-if="stats">
            <el-icon><Document /></el-icon>
            <span>文章数：{{ stats.articleCount }}</span>
          </div>
          <div class="info-item" v-if="stats">
            <el-icon><ChatDotRound /></el-icon>
            <span>评论数：{{ stats.commentCount }}</span>
          </div>
          <div class="info-item" v-if="stats">
            <el-icon><Star /></el-icon>
            <span>获赞数：{{ stats.likeCount }}</span>
          </div>
        </div>
      </el-card>

      <!-- 用户文章列表 -->
      <el-card v-if="user" class="user-articles">
        <template #header>
          <div class="card-header">
            <span>TA 的文章</span>
          </div>
        </template>
        
        <div v-if="articles && articles.length" class="article-list">
          <div
            v-for="article in articles"
            :key="article.id"
            class="article-item"
            @click="goToArticle(article.id)"
          >
            <h3 class="article-title">{{ article.title }}</h3>
            <p class="article-summary">{{ article.summary }}</p>
            <div class="article-meta">
              <span><el-icon><View /></el-icon> {{ article.view_count }} 阅读</span>
              <span>{{ formatDate(article.created_at) }}</span>
            </div>
          </div>
        </div>
        
        <el-empty v-else description="暂无文章" />
        
        <div class="pagination" v-if="articles && articles.length">
          <el-pagination
            v-model:current-page="articlePage"
            :page-size="pageSize"
            :total="articleTotal"
            layout="prev, pager, next"
            @current-change="loadUserArticles"
          />
        </div>
      </el-card>

      <!-- 用户评论列表 -->
      <el-card v-if="user" class="user-comments">
        <template #header>
          <div class="card-header">
            <span>TA 的评论</span>
          </div>
        </template>
        
        <div v-if="comments && comments.length" class="comment-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-meta">
              <span>发表于 {{ comment.article?.title || '文章' }}</span>
              <span>{{ formatDate(comment.created_at) }}</span>
            </div>
          </div>
        </div>
        
        <el-empty v-else description="暂无评论" />
        
        <div class="pagination" v-if="comments && comments.length">
          <el-pagination
            v-model:current-page="commentPage"
            :page-size="pageSize"
            :total="commentTotal"
            layout="prev, pager, next"
            @current-change="loadUserComments"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, User, Calendar, Document, ChatDotRound, Star, View } from '@element-plus/icons-vue'
import { ElMessage, ElCard, ElAvatar, ElEmpty, ElPagination, ElButton, ElIcon } from 'element-plus'
import TopNavBar from '@/components/TopNavBar.vue'
import { getUserProfile, getUserArticles, getUserComments, getUserStatsById } from '@/api/auth'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const user = ref(null)
const articles = ref([])
const comments = ref([])
const articlePage = ref(1)
const commentPage = ref(1)
const pageSize = ref(10)
const articleTotal = ref(0)
const commentTotal = ref(0)
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b9d3989743988cee5471b90a8672e.png'
const isUnmounted = ref(false)

const roleText = computed(() => {
  if (!user.value) return ''
  switch (user.value.role) {
    case 0: return '管理员'
    case 1: return '普通用户'
    default: return '未知'
  }
})

const stats = ref(null)

const goBack = () => {
  router.back()
}

const goToArticle = (articleId) => {
  router.push(`/article/${articleId}`)
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const loadUserProfile = async () => {
  try {
    const userId = route.params.id
    const res = await getUserProfile(userId)
    if (!isUnmounted.value) {
      user.value = res.data
      loadUserStats()
    }
  } catch (error) {
    if (!isUnmounted.value) {
      console.error('加载用户信息失败:', error)
      ElMessage.error('加载用户信息失败')
    }
  }
}

const loadUserStats = async () => {
  try {
    const userId = route.params.id
    const res = await getUserStatsById(userId)
    if (!isUnmounted.value) {
      stats.value = res.data
    }
  } catch (error) {
    if (!isUnmounted.value) {
      console.error('加载用户统计失败:', error)
    }
  }
}

const loadUserArticles = async () => {
  try {
    const userId = route.params.id
    const res = await getUserArticles(userId, {
      page: articlePage.value,
      pageSize: pageSize.value
    })
    if (!isUnmounted.value) {
      articles.value = res.data.articles || []
      articleTotal.value = res.data.total || 0
    }
  } catch (error) {
    if (!isUnmounted.value) {
      console.error('加载用户文章失败:', error)
    }
  }
}

const loadUserComments = async () => {
  try {
    const userId = route.params.id
    const res = await getUserComments(userId, {
      page: commentPage.value,
      pageSize: pageSize.value
    })
    if (!isUnmounted.value) {
      comments.value = res.data.comments || []
      commentTotal.value = res.data.total || 0
    }
  } catch (error) {
    if (!isUnmounted.value) {
      console.error('加载用户评论失败:', error)
    }
  }
}

onMounted(() => {
  loadUserProfile()
  loadUserArticles()
  loadUserComments()
  loading.value = false
})

onBeforeUnmount(() => {
  isUnmounted.value = true
})

</script>

<style lang="scss" scoped>
.user-profile-page {
  min-height: 100vh;
  background: #1a1a2e;
  padding: 20px;
}

.user-back {
  max-width: 1200px;
  margin: 0 auto 20px;
}

.user-container {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  gap: 20px;
}

.user-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  
  :deep(.el-card__body) {
    padding: 0;
  }
}

.user-header {
  text-align: center;
  padding: 40px 20px 20px;
  
  .user-nickname {
    margin: 20px 0 10px;
    font-size: 28px;
    font-weight: bold;
  }
  
  .user-username {
    color: rgba(255, 255, 255, 0.8);
    font-size: 16px;
    margin: 0;
  }
}

.user-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
  padding: 20px;
  background: rgba(0, 0, 0, 0.2);
  
  .info-item {
    display: flex;
    align-items: center;
    gap: 10px;
    
    .el-icon {
      font-size: 20px;
    }
  }
}

.user-articles,
.user-comments {
  background: #16213e;
  color: #e4e4e4;
  
  :deep(.el-card__header) {
    background: #0f3460;
    border-color: #0f3460;
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      span {
        color: #e4e4e4;
        font-weight: bold;
        font-size: 16px;
      }
    }
  }
}

.article-list {
  display: grid;
  gap: 15px;
}

.article-item {
  padding: 20px;
  background: #0f3460;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
  
  &:hover {
    transform: translateY(-2px);
  }
  
  .article-title {
    margin: 0 0 10px;
    color: #e4e4e4;
    font-size: 18px;
  }
  
  .article-summary {
    color: #a0a0a0;
    margin: 0 0 15px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .article-meta {
    display: flex;
    gap: 20px;
    color: #888;
    font-size: 14px;
    
    span {
      display: flex;
      align-items: center;
      gap: 5px;
    }
  }
}

.comment-list {
  display: grid;
  gap: 15px;
}

.comment-item {
  padding: 20px;
  background: #0f3460;
  border-radius: 8px;
  
  .comment-content {
    color: #e4e4e4;
    margin: 0 0 15px;
    line-height: 1.6;
  }
  
  .comment-meta {
    display: flex;
    justify-content: space-between;
    color: #888;
    font-size: 14px;
  }
}

.pagination {
  margin-top: 20px;
  text-align: center;
  
  :deep(.el-pagination) {
    .el-pager li {
      background: #0f3460;
      color: #e4e4e4;
      
      &.is-active {
        background: #667eea;
      }
      
      &:hover {
        background: #1a1a2e;
      }
    }
  }
}
</style>
