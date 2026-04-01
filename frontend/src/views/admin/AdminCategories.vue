<template>
  <div class="admin-categories">
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
          <router-link to="/admin/articles" class="nav-item">
            <span class="nav-icon">📝</span>
            <span>文章管理</span>
          </router-link>
          <router-link to="/admin/categories" class="nav-item active">
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
          <h1>分类管理</h1>
          <p>管理系统分类（标签）</p>
        </div>
        
        <!-- 操作栏 -->
        <div class="action-bar">
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新建分类
          </el-button>
        </div>
        
        <!-- 分类列表 -->
        <el-table :data="categories" style="width: 100%" class="dark-table" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="分类名称" />
          <el-table-column prop="slug" label="分类标识" />
          <el-table-column prop="article_count" label="文章数" width="100" />
          <el-table-column prop="created_at" label="创建时间" width="180" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button size="small" type="warning" @click="editCategory(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteCategory(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 创建/编辑对话框 -->
        <el-dialog
          v-model="dialogVisible"
          :title="isEdit ? '编辑分类' : '创建分类'"
          width="400px"
          class="dark-dialog"
          @close="resetForm"
        >
          <el-form :model="form" label-width="80px">
            <el-form-item label="分类名称" required>
              <el-input
                v-model="form.name"
                placeholder="请输入分类名称"
                class="dark-input"
              />
            </el-form-item>
            <el-form-item label="分类标识" required>
              <el-input
                v-model="form.slug"
                placeholder="请输入分类标识（英文）"
                class="dark-input"
              />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="dialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitForm" :loading="saving">
              {{ isEdit ? '保存' : '创建' }}
            </el-button>
          </template>
        </el-dialog>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import TopNavBar from '@/components/TopNavBar.vue'
import { adminRequest } from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)

const categories = ref([])
const form = ref({
  id: null,
  name: '',
  slug: ''
})

// 加载分类列表
const loadCategories = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await adminRequest.get('/admin/categories', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    
    categories.value = response.data.data.categories || []
  } catch (error) {
    ElMessage.error('加载分类列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  form.value = { name: '', slug: '' }
  dialogVisible.value = true
}

// 编辑分类
const editCategory = (category) => {
  isEdit.value = true
  form.value = {
    id: category.id,
    name: category.name,
    slug: category.slug
  }
  dialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  form.value = { id: null, name: '', slug: '' }
}

// 提交表单
const submitForm = async () => {
  if (!form.value.name || !form.value.slug) {
    ElMessage.warning('请填写完整信息')
    return
  }
  
  saving.value = true
  try {
    const url = isEdit.value 
      ? `/admin/category/update/${form.value.id}`
      : '/admin/category/create'
    
    const method = isEdit.value ? 'put' : 'post'
    
    await adminRequest[method](url, form.value)
    
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    dialogVisible.value = false
    loadCategories()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    console.error(error)
  } finally {
    saving.value = false
  }
}

// 删除分类
const deleteCategory = async (category) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类"${category.name}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await adminRequest.delete(`/admin/category/delete/${category.id}`)
    
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
      console.error(error)
    }
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped lang="scss">
.admin-categories {
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
  
  .action-bar {
    margin-bottom: 20px;
  }
}
</style>
