<template>
  <div class="creator-center-page">
    <TopNavBar />
    
    <div class="creator-container">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">
            <el-icon><Edit /></el-icon>
            创作中心
          </h1>
          <p class="page-subtitle">管理您的文章，分享知识与见解</p>
        </div>
        <el-button type="primary" @click="handleCreate" class="create-btn">
          <el-icon><Plus /></el-icon>
          创建文章
        </el-button>
      </div>

      <!-- 统计卡片 -->
      <div class="stats-cards">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">总文章数</div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.published }}</div>
            <div class="stat-label">已发布</div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.draft }}</div>
            <div class="stat-label">草稿</div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.views }}</div>
            <div class="stat-label">总阅读量</div>
          </div>
        </el-card>
      </div>

      <!-- 筛选和操作栏 -->
      <div class="filter-bar">
        <div class="filter-tabs">
          <el-radio-group v-model="filterStatus" @change="loadArticles">
            <el-radio-button value="all">全部</el-radio-button>
            <el-radio-button value="published">已发布</el-radio-button>
            <el-radio-button value="draft">草稿</el-radio-button>
          </el-radio-group>
        </div>
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章标题..."
            clearable
            @keyup.enter="loadArticles"
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="loadArticles">搜索</el-button>
        </div>
      </div>

      <!-- 文章列表 -->
      <el-card shadow="hover" class="articles-card">
        <div class="article-list" v-loading="loading">
          <div v-if="articles.length === 0 && !loading" class="empty-state">
            <el-empty description="暂无文章">
              <el-button type="primary" @click="handleCreate">创建第一篇文章</el-button>
            </el-empty>
          </div>

          <div 
            v-for="article in articles" 
            :key="article.id" 
            class="article-item"
          >
            <div class="article-main">
              <div class="article-cover-wrapper">
                <img 
                  v-if="article.cover_image" 
                  :src="article.cover_image" 
                  :alt="article.title"
                  class="article-cover"
                  @error="handleImageError"
                />
                <div v-else class="article-cover-placeholder">
                  <el-icon><Document /></el-icon>
                </div>
              </div>
              
              <div class="article-info">
                <div class="article-header">
                  <h3 class="article-title">{{ article.title }}</h3>
                  <el-tag 
                    :type="article.status === 1 ? 'success' : 'info'" 
                    size="small"
                    class="status-tag"
                    :class="article.status === 1 ? 'published' : 'draft'"
                  >
                    <el-icon class="status-icon">
                      <component :is="article.status === 1 ? 'Check' : 'Clock'" />
                    </el-icon>
                    {{ article.status === 1 ? '已发布' : '草稿' }}
                  </el-tag>
                </div>
                <p class="article-summary">{{ article.summary || '暂无摘要' }}</p>
                <div class="article-meta">
                  <span class="meta-item">
                    <el-icon><Calendar /></el-icon>
                    {{ formatDate(article.created_at) }}
                  </span>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ article.view_count || 0 }} 次阅读
                  </span>
                  <span class="meta-item" v-if="article.category">
                    <el-icon><Folder /></el-icon>
                    {{ article.category.name }}
                  </span>
                </div>
              </div>
            </div>
            
            <div class="article-actions">
              <el-button 
                type="primary" 
                size="small" 
                @click="handleEdit(article)"
                class="action-btn edit-btn"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                @click="handleDelete(article)"
                class="action-btn delete-btn"
              >
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-wrapper" v-if="total > pageSize">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="total"
            layout="total, prev, pager, next"
            @current-change="loadArticles"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  Edit, Plus, Document, Calendar, View, 
  Folder, Search, Delete, Check, Clock
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMyArticles, deleteArticle } from '@/api/auth'
import TopNavBar from '@/components/TopNavBar.vue'

const router = useRouter()
const userStore = useUserStore()

// 统计数据
const stats = reactive({
  total: 0,
  published: 0,
  draft: 0,
  views: 0
})

// 文章列表
const articles = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const filterStatus = ref('all')
const searchKeyword = ref('')

// 加载文章列表
const loadArticles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (filterStatus.value === 'published') {
      params.status = 1
    } else if (filterStatus.value === 'draft') {
      params.status = 0
    }
    
    if (searchKeyword.value) {
      params.search = searchKeyword.value
    }
    
    const res = await getMyArticles(params)
    if (res.code === 0 && res.data) {
      articles.value = res.data.articles || []
      total.value = res.data.total || 0
      
      // 更新统计
      updateStats(articles.value)
    }
  } catch (error) {
    console.error('加载文章列表失败:', error)
    ElMessage.error('加载文章列表失败')
  } finally {
    loading.value = false
  }
}

// 更新统计数据
const updateStats = (articleList) => {
  if (!articleList || !Array.isArray(articleList)) {
    return
  }
  stats.total = total.value
  stats.published = articleList.filter(a => a.status === 1).length
  stats.draft = articleList.filter(a => a.status === 0).length
  stats.views = articleList.reduce((sum, a) => sum + (a.view_count || 0), 0)
}

// 处理图片加载错误
const handleImageError = (e) => {
  const img = e.target
  img.style.display = 'none'
  const placeholder = img.parentElement?.querySelector('.article-cover-placeholder')
  if (placeholder) {
    placeholder.style.display = 'flex'
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '未知'
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return '未知'
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    console.error('formatDate 错误:', error)
    return '未知'
  }
}

// 创建文章
const handleCreate = () => {
  router.push('/creator/edit')
}

// 编辑文章
const handleEdit = (article) => {
  router.push(`/creator/edit/${article.id}`)
}

// 删除文章
const handleDelete = (article) => {
  ElMessageBox.confirm(
    `确定要删除文章"${article.title}"吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteArticle(article.id)
      ElMessage.success('删除成功')
      loadArticles()
    } catch (error) {
      console.error('删除文章失败:', error)
      ElMessage.error('删除文章失败')
    }
  }).catch(() => {})
}

// 初始化
onMounted(() => {
  if (!userStore.userInfo) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  loadArticles()
})
</script>

<style lang="scss" scoped>
.creator-center-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e2a4a 0%, #2d1f3d 50%, #1a3a5c 100%);
  padding: 24px;
  
  .creator-container {
    max-width: 1200px;
    margin: 0 auto;
    
    // 页面头部
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;
      
      .header-content {
        .page-title {
          font-size: 32px;
          font-weight: 700;
          color: rgba(255, 255, 255, 0.95);
          margin: 0 0 8px;
          display: flex;
          align-items: center;
          gap: 12px;
          
          .el-icon {
            color: rgba(102, 126, 234, 0.8);
          }
        }
        
        .page-subtitle {
          font-size: 14px;
          color: rgba(255, 255, 255, 0.6);
          margin: 0;
        }
      }
      
      .create-btn {
        height: 44px;
        padding: 0 24px;
        font-size: 15px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.3) 100%);
        border: 1px solid rgba(102, 126, 234, 0.6);
        backdrop-filter: blur(10px);
        
        &:hover {
          background: linear-gradient(135deg, rgba(102, 126, 234, 0.4) 0%, rgba(118, 75, 162, 0.4) 100%);
        }
      }
    }
    
    // 统计卡片
    .stats-cards {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
      margin-bottom: 24px;
      
      .stat-card {
        background: rgba(255, 255, 255, 0.05);
        backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.1);
        
        :deep(.el-card__body) {
          padding: 24px;
          color: rgba(255, 255, 255, 0.95) !important;
        }
        
        .stat-item {
          text-align: center;
          
          .stat-value {
            font-size: 32px;
            font-weight: 700;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            margin-bottom: 8px;
          }
          
          .stat-label {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.7);
          }
        }
      }
    }
    
    // 筛选栏
    .filter-bar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;
      
      .filter-tabs {
        :deep(.el-radio-button__inner) {
          background: rgba(255, 255, 255, 0.05);
          border-color: rgba(255, 255, 255, 0.1);
          color: rgba(255, 255, 255, 0.8);
          
          &:hover {
            color: rgba(255, 255, 255, 0.95);
          }
        }
        
        :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
          background: linear-gradient(135deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.3) 100%);
          border-color: rgba(102, 126, 234, 0.6);
          color: rgba(255, 255, 255, 0.95);
        }
      }
      
      .search-box {
        display: flex;
        gap: 12px;
        
        .search-input {
          width: 300px;
          
          :deep(.el-input__wrapper) {
            background: rgba(255, 255, 255, 0.05);
            border-color: rgba(255, 255, 255, 0.1);
            
            .el-input__inner {
              color: rgba(255, 255, 255, 0.95);
            }
          }
        }
      }
    }
    
    // 文章卡片
    .articles-card {
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      
      :deep(.el-card__body) {
        padding: 20px;
        color: rgba(255, 255, 255, 0.95) !important;
      }
      
      .article-list {
        display: flex;
        flex-direction: column;
        gap: 16px;
        
        .empty-state {
          padding: 60px 20px;
          text-align: center;
        }
        
        .article-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 20px 24px;
          background: rgba(255, 255, 255, 0.03);
          border: 1px solid rgba(255, 255, 255, 0.08);
          border-radius: 12px;
          transition: all 0.3s ease;
          
          &:hover {
            background: rgba(255, 255, 255, 0.06);
            border-color: rgba(102, 126, 234, 0.4);
            transform: translateX(4px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
          }
          
          .article-main {
            display: flex;
            gap: 20px;
            flex: 1;
            
            @media (max-width: 768px) {
              flex-direction: column;
              gap: 16px;
            }
            
            .article-cover-wrapper {
              flex-shrink: 0;
              width: 200px;
              height: 120px;
              position: relative;
              
              @media (max-width: 768px) {
                width: 100%;
                height: 200px;
              }
              
              .article-cover {
                width: 100%;
                height: 100%;
                object-fit: cover;
                border-radius: 6px;
              }
              
              .article-cover-placeholder {
                width: 100%;
                height: 100%;
                background: rgba(102, 126, 234, 0.2);
                border-radius: 6px;
                display: none;
                align-items: center;
                justify-content: center;
                position: absolute;
                top: 0;
                left: 0;
                
                .el-icon {
                  font-size: 48px;
                  color: rgba(102, 126, 234, 0.6);
                }
              }
            }
            
            .article-info {
              flex: 1;
              min-width: 0;
              
              .article-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 12px;
                
                .article-title {
                  font-size: 18px;
                  font-weight: 600;
                  color: rgba(255, 255, 255, 0.95);
                  margin: 0;
                  overflow: hidden;
                  text-overflow: ellipsis;
                  white-space: nowrap;
                }
                
                .status-tag {
                  border: none;
                  font-weight: 600;
                  font-size: 12px;
                  padding: 6px 14px;
                  border-radius: 20px;
                  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
                  position: relative;
                  overflow: hidden;
                  display: inline-flex;
                  align-items: center;
                  gap: 4px;
                  
                  .status-icon {
                    font-size: 14px;
                    animation: pulse 2s infinite;
                  }
                  
                  &::before {
                    content: '';
                    position: absolute;
                    top: 0;
                    left: -100%;
                    width: 100%;
                    height: 100%;
                    background: linear-gradient(
                      90deg,
                      transparent,
                      rgba(255, 255, 255, 0.3),
                      transparent
                    );
                    transition: left 0.5s;
                  }
                  
                  &:hover {
                    transform: translateY(-2px) scale(1.05);
                    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.3);
                    
                    &::before {
                      left: 100%;
                    }
                  }
                  
                  &.published {
                    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
                    color: white;
                    
                    &:hover {
                      background: linear-gradient(135deg, #38ef7d 0%, #11998e 100%);
                      box-shadow: 0 6px 20px rgba(56, 239, 125, 0.4);
                    }
                  }
                  
                  &.draft {
                    background: linear-gradient(135deg, #868f96 0%, #596167 100%);
                    color: white;
                    
                    &:hover {
                      background: linear-gradient(135deg, #596167 0%, #868f96 100%);
                      box-shadow: 0 6px 20px rgba(134, 143, 150, 0.4);
                    }
                  }
                }
              }
              
              .article-summary {
                font-size: 14px;
                color: rgba(255, 255, 255, 0.6);
                margin: 0 0 12px;
                overflow: hidden;
                text-overflow: ellipsis;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                line-clamp: 2;
                -webkit-box-orient: vertical;
                line-height: 1.5;
              }
              
              .article-meta {
                display: flex;
                gap: 16px;
                font-size: 13px;
                color: rgba(255, 255, 255, 0.5);
                
                .meta-item {
                  display: flex;
                  align-items: center;
                  gap: 4px;
                  
                  .el-icon {
                    font-size: 14px;
                  }
                }
              }
            }
          }
          
          .article-actions {
            display: flex;
            flex-direction: column;
            gap: 12px;
            margin-left: 24px;
            flex-shrink: 0;
            align-items: flex-end;
            justify-content: flex-start;
            padding: 8px 0;
            
            @media (max-width: 768px) {
              flex-direction: row;
              margin-left: 0;
              margin-top: 12px;
              width: 100%;
              justify-content: flex-end;
              gap: 12px;
              align-items: center;
            }
            
            .action-btn {
              min-width: 85px;
              font-weight: 600;
              border: none;
              transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
              position: relative;
              overflow: hidden;
              
              &::before {
                content: '';
                position: absolute;
                top: 0;
                left: -100%;
                width: 100%;
                height: 100%;
                background: linear-gradient(
                  90deg,
                  transparent,
                  rgba(255, 255, 255, 0.3),
                  transparent
                );
                transition: left 0.5s;
              }
              
              &:hover {
                transform: translateY(-2px);
                box-shadow: 0 6px 16px rgba(0, 0, 0, 0.3);
                
                &::before {
                  left: 100%;
                }
              }
              
              &:active {
                transform: translateY(0);
              }
            }
            
            .edit-btn {
              background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
              color: white;
              
              &:hover {
                background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);
                box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
              }
            }
            
            .delete-btn {
              background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
              color: white;
              
              &:hover {
                background: linear-gradient(135deg, #f5576c 0%, #f093fb 100%);
                box-shadow: 0 6px 20px rgba(245, 87, 108, 0.4);
              }
            }
          }
        }
      }
      
      // 分页
      .pagination-wrapper {
        margin-top: 24px;
        display: flex;
        justify-content: center;
        
        :deep(.el-pagination) {
          .el-pager li {
            background: rgba(255, 255, 255, 0.05);
            color: rgba(255, 255, 255, 0.8);
            
            &.is-active {
              background: linear-gradient(135deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.3) 100%);
              color: rgba(255, 255, 255, 0.95);
            }
            
            &:hover {
              color: rgba(255, 255, 255, 0.95);
            }
          }
          
          .btn-prev, .btn-next {
            background: rgba(255, 255, 255, 0.05);
            color: rgba(255, 255, 255, 0.8);
            
            &:hover {
              color: rgba(255, 255, 255, 0.95);
            }
          }
        }
      }
    }
  }
  
  // 响应式设计
  @media (max-width: 768px) {
    padding: 16px;
    
    .creator-container {
      .page-header {
        flex-direction: column;
        gap: 16px;
        text-align: center;
        
        .header-content {
          .page-title {
            font-size: 24px;
          }
        }
      }
      
      .stats-cards {
        grid-template-columns: repeat(2, 1fr);
      }
      
      .filter-bar {
        flex-direction: column;
        gap: 16px;
        
        .search-box {
          width: 100%;
          
          .search-input {
            flex: 1;
          }
        }
      }
      
      .article-item {
        flex-direction: column;
        align-items: flex-start;
        padding: 16px;
        
        .article-main {
          flex-direction: column;
          width: 100%;
        }
      }
    }
  }
}

// 动画效果
@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}

@keyframes shimmer {
  0% {
    background-position: -1000px 0;
  }
  100% {
    background-position: 1000px 0;
  }
}
</style>
