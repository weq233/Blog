<template>
  <div class="articles-page">
    <TopNavBar />
    
    <div class="articles-header">
      <h1>文章列表</h1>
      <div class="search-box">
        <el-input
          v-model="searchQuery"
          placeholder="搜索文章标题..."
          clearable
          @keyup.enter="handleSearch"
        >
          <template #suffix>
            <el-icon 
              class="search-icon"
              style="cursor: pointer; color: rgba(255, 255, 255, 0.6);"
              @click="handleSearch"
            >
              <Search />
            </el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <div class="articles-filters">
      <div class="filter-group">
        <span class="filter-label">
          <el-icon><Folder /></el-icon>
          分类：
        </span>
        <el-select
          v-model="selectedCategory"
          placeholder="全部分类"
          clearable
          filterable
          @change="handleFilterChange"
          style="width: 180px;"
        >
          <el-option
            v-for="cat in categories"
            :key="cat.id"
            :label="cat.name"
            :value="cat.slug || cat.id?.toString()"
          />
        </el-select>
      </div>

      <div class="filter-group">
        <span class="filter-label">
          <el-icon><PriceTag /></el-icon>
          标签：
        </span>
        <el-select
          v-model="selectedTag"
          placeholder="全部标签"
          clearable
          filterable
          @change="handleFilterChange"
          style="width: 180px;"
        >
          <el-option
            v-for="tag in tags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.slug || tag.id?.toString()"
          />
        </el-select>
      </div>
    </div>

    <div class="articles-list" v-loading="loading">
      <el-card
        v-for="article in articles"
        :key="article.id"
        class="article-card"
        shadow="hover"
        @click.native="goToArticle(article.id)"
      >
        <div class="article-content">
          <div class="article-cover" v-if="article.cover_image">
            <img :src="article.cover_image" :alt="article.title" />
          </div>
          <div class="article-info">
            <h2 class="article-title">{{ article.title }}</h2>
            <p class="article-summary">{{ article.summary }}</p>
            <div class="article-meta">
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(article.created_at) }}
              </span>
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ article.view_count }}
              </span>
              <span class="meta-item" v-if="article.category">
                <el-icon><Folder /></el-icon>
                {{ article.category.name }}
              </span>
            </div>
            <div class="article-tags" v-if="article.tags && article.tags.length">
              <el-tag
                v-for="tag in article.tags"
                :key="tag.id"
                size="small"
                type="info"
                effect="plain"
              >
                {{ tag.name }}
              </el-tag>
            </div>
          </div>
        </div>
      </el-card>

      <el-empty v-if="!loading && articles.length === 0" description="暂无文章" />
    </div>

    <div class="articles-pagination" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Calendar, View, Folder, PriceTag } from '@element-plus/icons-vue'
import { getArticles } from '@/api/auth'
import TopNavBar from '@/components/TopNavBar.vue'

const router = useRouter()

// 加载文章列表
onMounted(() => {
  loadArticles()
})

const loading = ref(false)
const articles = ref([])
const categories = ref([])
const tags = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchQuery = ref('')
const selectedCategory = ref()
const selectedTag = ref()

// 加载文章列表
const loadArticles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
    }

    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    if (selectedTag.value) {
      params.tag = selectedTag.value
    }

    if (searchQuery.value) {
      params.search = searchQuery.value
    }

    const res = await getArticles(params)
    articles.value = res.data.articles || []
    total.value = res.data.total || 0
    categories.value = res.data.categories || []
    tags.value = res.data.tags || []
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    loading.value = false
  }
}

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  loadArticles()
}

// 筛选变化
const handleFilterChange = () => {
  currentPage.value = 1
  loadArticles()
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  loadArticles()
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 跳转文章详情
const goToArticle = (id) => {
  router.push(`/article/${id}`)
}

</script>

<style lang="scss" scoped>
.articles-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e2a4a 0%, #2d1f3d 50%, #1a3a5c 100%);
  
  .articles-header {
    max-width: 1200px;
    margin: 30px auto;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 20px;
    
    h1 {
      color: #fff;
      font-size: 32px;
      margin: 0;
      font-weight: 300;
      letter-spacing: 2px;
    }
    
    .search-box {
      width: 300px;
      
      :deep(.el-input__wrapper) {
        background: rgba(255, 255, 255, 0.05) !important;
        background-image: none !important;
        border: 1px solid rgba(255, 255, 255, 0.15) !important;
        backdrop-filter: blur(10px);
        box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.15) inset !important;
        transition: all 0.3s ease;
        
        &:hover {
          background-color: rgba(255, 255, 255, 0.08) !important;
          border-color: rgba(103, 194, 58, 0.5) !important;
          box-shadow: 0 0 0 1px rgba(103, 194, 58, 0.5) inset !important;
        }
      }
      
      :deep(.el-input__inner) {
        color: rgba(255, 255, 255, 0.9) !important;
        
        &::placeholder {
          color: rgba(255, 255, 255, 0.4) !important;
        }
      }
      
      .search-icon {
        transition: all 0.3s ease;
        
        &:hover {
          color: #67c23a !important;
        }
      }
    }
  }
  
  .articles-filters {
    max-width: 1200px;
    margin: 0 auto 30px;
    padding: 0 20px;
    display: flex;
    gap: 30px;
    flex-wrap: wrap;
    
    .filter-group {
      display: flex;
      align-items: center;
      gap: 10px;
      
      .filter-label {
        display: flex;
        align-items: center;
        gap: 6px;
        color: rgba(255, 255, 255, 0.8);
        font-size: 14px;
        font-weight: 500;
        
        .el-icon {
          font-size: 16px;
          color: #67c23a;
        }
      }
    }
    
    :deep(.el-select) {
      // 直接覆盖 el-select__wrapper
      & .el-select__wrapper,
      & .el-input__wrapper {
        background-color: rgba(255, 255, 255, 0.05) !important;
        background-image: none !important;
        background: rgba(255, 255, 255, 0.05) !important;
        border: 1px solid rgba(255, 255, 255, 0.15) !important;
        backdrop-filter: blur(10px);
        box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.15) inset !important;
        transition: all 0.3s ease;
        
        &:hover {
          background-color: rgba(255, 255, 255, 0.08) !important;
          background-image: none !important;
          background: rgba(255, 255, 255, 0.08) !important;
          border-color: rgba(103, 194, 58, 0.5) !important;
          box-shadow: 0 0 0 1px rgba(103, 194, 58, 0.5) inset !important;
        }
        
        &:focus-within {
          background-color: rgba(255, 255, 255, 0.08) !important;
          background-image: none !important;
          background: rgba(255, 255, 255, 0.08) !important;
          border-color: rgba(103, 194, 58, 0.8) !important;
          box-shadow: 0 0 0 1px rgba(103, 194, 58, 0.8) inset !important;
        }
      }
      
      // 输入框 - 确保文字颜色
      .el-input__inner {
        color: rgba(255, 255, 255, 0.9) !important;
        
        &::placeholder {
          color: rgba(255, 255, 255, 0.4) !important;
        }
      }
      
      // 选中项 - 文字颜色
      .el-select__selected-item {
        color: rgba(255, 255, 255, 0.9) !important;
        background-color: transparent !important;
      }
      
      // 箭头图标
      .el-select__caret {
        color: rgba(255, 255, 255, 0.6) !important;
      }
      
      // 下拉菜单容器 - 使用全局选择器覆盖
      .el-select__popper {
        background-color: rgba(30, 30, 50, 0.98) !important;
        background-image: none !important;
        backdrop-filter: blur(20px);
        border: 1px solid rgba(255, 255, 255, 0.08) !important;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3) !important;
        
        // 下拉选项
        .el-select-dropdown__item {
          color: rgba(255, 255, 255, 0.85) !important;
          background-color: transparent !important;
          
          &.selected {
            color: #67c23a !important;
            background-color: rgba(103, 194, 58, 0.15) !important;
          }
          
          &:hover {
            background-color: rgba(103, 194, 58, 0.08) !important;
          }
        }
      }
    }
    
    // 全局覆盖下拉菜单样式（穿透 scoped）
    :deep(.el-select-dropdown) {
      background-color: rgba(30, 30, 50, 0.98) !important;
      
      .el-select-dropdown__item {
        color: rgba(255, 255, 255, 0.85) !important;
        
        &.selected {
          color: #67c23a !important;
          background-color: rgba(103, 194, 58, 0.15) !important;
        }
        
        &:hover {
          background-color: rgba(103, 194, 58, 0.08) !important;
        }
      }
    }
  }
  
  .articles-list {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    
    .article-card {
      cursor: pointer;
      transition: all 0.3s ease;
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(20px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
        border-color: rgba(255, 255, 255, 0.2);
      }
      
      :deep(.el-card__body) {
        padding: 0;
      }
      
      .article-content {
        display: flex;
        gap: 25px;
        padding: 25px;
        
        .article-cover {
          flex-shrink: 0;
          width: 200px;
          height: 140px;
          border-radius: 8px;
          overflow: hidden;
          background: rgba(255, 255, 255, 0.1);
          
          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: transform 0.3s ease;
          }
          
          &:hover img {
            transform: scale(1.1);
          }
        }
        
        .article-info {
          flex: 1;
          display: flex;
          flex-direction: column;
          gap: 12px;
          
          .article-title {
            color: #fff;
            font-size: 20px;
            margin: 0;
            font-weight: 400;
            letter-spacing: 1px;
            line-height: 1.4;
          }
          
          .article-summary {
            color: rgba(255, 255, 255, 0.7);
            font-size: 14px;
            margin: 0;
            line-height: 1.6;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
          }
          
          .article-meta {
            display: flex;
            gap: 20px;
            flex-wrap: wrap;
            
            .meta-item {
              display: flex;
              align-items: center;
              gap: 6px;
              color: rgba(255, 255, 255, 0.5);
              font-size: 13px;
              
              .el-icon {
                font-size: 16px;
              }
            }
          }
          
          .article-tags {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
            
            .el-tag {
              background: rgba(255, 255, 255, 0.1) !important;
              border-color: rgba(255, 255, 255, 0.2);
              color: rgba(255, 255, 255, 0.8) !important;
            }
          }
        }
      }
    }
  }
  
  .articles-pagination {
    max-width: 1200px;
    margin: 40px auto 0;
    padding: 0 20px;
    display: flex;
    justify-content: center;
    
    :deep(.el-pagination) {
      .el-pager li {
        background: rgba(255, 255, 255, 0.1) !important;
        color: #fff !important;
        border: 1px solid rgba(255, 255, 255, 0.2);
        
        &.is-active {
          background: #409eff !important;
          border-color: #409eff !important;
        }
        
        &:hover {
          background: rgba(255, 255, 255, 0.2) !important;
        }
      }
      
      .btn-prev,
      .btn-next {
        background: rgba(255, 255, 255, 0.1) !important;
        color: #fff !important;
        border: 1px solid rgba(255, 255, 255, 0.2);
        
        &:hover {
          background: rgba(255, 255, 255, 0.2) !important;
        }
      }
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .articles-page {
    .articles-header {
      margin: 20px auto;
      padding: 0 15px;
      flex-direction: column;
      align-items: stretch;
      
      h1 {
        font-size: 24px;
        text-align: center;
        margin-bottom: 15px;
      }
      
      .search-box {
        width: 100%;
      }
    }
    
    .articles-filters {
      margin: 0 auto 20px;
      padding: 0 15px;
      flex-direction: column;
      gap: 20px;
      
      .filter-group {
        width: 100%;
        
        .filter-label {
          min-width: 70px;
        }
      }
      
      :deep(.el-select) {
        width: 100% !important;
      }
    }
    
    .articles-list {
      padding: 0 15px;
      gap: 15px;
      
      .article-card {
        .article-content {
          flex-direction: column;
          padding: 15px;
          gap: 15px;
          
          .article-cover {
            width: 100%;
            height: 180px;
          }
        }
      }
    }
    
    .articles-pagination {
      margin: 30px auto 0;
      padding: 0 15px;
    }
  }
}
</style>
