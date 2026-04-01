<template>
  <div class="admin-articles">
    <!-- 顶部导航栏 -->
    <TopNavBar />
    
    <div class="admin-container">
      <!-- 侧边栏 -->
      <aside class="admin-sidebar">
        <div class="sidebar-header">
          <h2>📊 管理后台</h2>
        </div>
        <nav class="sidebar-nav">
          <router-link to="/admin" class="nav-item">
            <span class="nav-icon">🏠</span>
            <span>数据概览</span>
          </router-link>
          <router-link to="/admin/articles" class="nav-item active">
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
          <h1>文章管理</h1>
          <p>管理系统中的所有文章</p>
        </div>
        
        <!-- 搜索栏 -->
        <div class="search-bar">
          <el-input
            v-model="searchTitle"
            placeholder="搜索文章标题"
            style="width: 300px"
            class="dark-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
          
          <el-select
            v-model="filterStatus"
            placeholder="状态筛选"
            style="width: 150px"
            class="dark-select"
            clearable
            @change="handleSearch"
          >
            <el-option label="已发布" :value="1" />
            <el-option label="草稿" :value="0" />
          </el-select>
        </div>
        
        <!-- 文章列表 -->
        <el-table
          :data="articles"
          style="width: 100%"
          class="dark-table"
          v-loading="loading"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="title" label="标题" />
          <el-table-column prop="author.username" label="作者" width="120" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '已发布' : '草稿' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="view_count" label="浏览量" width="100" />
          <el-table-column prop="created_at" label="发布时间" width="180" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="viewArticle(row)">查看</el-button>
              <el-button size="small" type="warning" @click="editArticle(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteArticle(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSearch"
            @current-change="handleSearch"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import TopNavBar from '@/components/TopNavBar.vue'
import { adminRequest } from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const searchTitle = ref('')
const filterStatus = ref(null)

const articles = ref([])
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 加载文章列表
const loadArticles = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    
    if (searchTitle.value) {
      params.title = searchTitle.value
    }
    
    if (filterStatus.value !== null) {
      params.status = filterStatus.value
    }
    
    const response = await adminRequest.get('/admin/articles', {
      params,
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    // 检查响应数据结构
    if (!response.data) {
      ElMessage.error('服务器返回数据为空')
      return
    }
    
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '加载失败')
      return
    }
    
    if (!response.data.data) {
      ElMessage.error('数据格式错误')
      return
    }
    
    articles.value = response.data.data.articles || []
    pagination.value.total = response.data.data.total || 0
  } catch (error) {
    if (error.response) {
      if (error.response.status === 401) {
        ElMessage.error('请先登录')
        router.push('/login')
      } else if (error.response.status === 403) {
        ElMessage.error('权限不足')
      } else {
        ElMessage.error(`加载失败：${error.response.status}`)
      }
    } else {
      ElMessage.error('网络错误或服务器无响应')
    }
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.value.page = 1
  loadArticles()
}

// 查看文章
const viewArticle = (article) => {
  router.push(`/article/${article.id}`)
}

// 编辑文章
const editArticle = (article) => {
  router.push(`/admin/article/edit/${article.id}`)
}

// 删除文章
const deleteArticle = async (article) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文章《${article.title}》吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const token = localStorage.getItem('token')
    await axios.delete(`/admin/article/delete/${article.id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    ElMessage.success('删除成功')
    loadArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped lang="scss">
.admin-articles {
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
  
  .search-bar {
    display: flex;
    gap: 15px;
    margin-bottom: 20px;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
