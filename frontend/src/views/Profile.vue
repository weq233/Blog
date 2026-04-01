
<template>
  <div class="profile-page">
    <TopNavBar />
    
    <div class="profile-container">
      <!-- 用户信息卡片 -->
      <el-card class="profile-header-card" shadow="hover">
        <div class="profile-header">
          <div class="avatar-info">
            <el-avatar 
              :size="100" 
              :src="getAvatarUrl(userStore.userInfo?.avatar)"
              class="profile-avatar"
            >
              {{ getUserInitial(userStore.userInfo?.nickname) }}
            </el-avatar>
            <div class="user-basic-info">
              <div class="user-name-row">
                <h1 class="user-nickname">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</h1>
                <el-tag :type="getRoleType(userStore.userInfo?.role)" effect="dark" size="large">
                  <el-icon><Star /></el-icon>
                  {{ getRoleText(userStore.userInfo?.role) }}
                </el-tag>
              </div>
              <div class="user-meta">
                <span class="username">@{{ userStore.userInfo?.username }}</span>
                <span class="divider">·</span>
                <span class="join-date">
                  <el-icon><Calendar /></el-icon>
                  加入于 {{ formatDate(userStore.userInfo?.createdAt) }}
                </span>
              </div>
              <p class="user-bio">这个人很懒，什么都没写~</p>
            </div>
          </div>
          
          <div class="profile-actions">
            <el-button type="primary" @click="toggleEdit" class="edit-btn">
              <el-icon><Edit /></el-icon>
              编辑资料
            </el-button>
          </div>
        </div>
      </el-card>

      <!-- 统计卡片 -->
      <div class="stats-cards">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.articles }}</div>
            <div class="stat-label">文章</div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.favorites }}</div>
            <div class="stat-label">收藏</div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.following }}</div>
            <div class="stat-label">关注</div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-value">{{ stats.followers }}</div>
            <div class="stat-label">粉丝</div>
          </div>
        </el-card>
      </div>

      <!-- 主要内容区 -->
      <div class="profile-content">
        <!-- 左侧近期阅读 -->
        <div class="main-content">
          <el-card shadow="hover" class="recent-views-card">
            <template #header>
              <div class="card-header">
                <el-icon><Document /></el-icon>
                <span>近期阅读</span>
              </div>
            </template>
            
            <div class="article-list" v-loading="loadingRecent">
              <div v-if="recentArticles.length === 0 && !loadingRecent" class="empty-state">
                <el-empty description="暂无阅读记录" :image-size="80" />
              </div>
              
              <div 
                v-for="article in recentArticles" 
                :key="article.id" 
                class="article-item"
                @click="$router.push(`/article/${article.id}`)"
              >
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
                  <h3 class="article-title">{{ article.title }}</h3>
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
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </div>

        <!-- 右侧快捷操作 -->
        <div class="sidebar">
          <el-card shadow="hover" class="quick-actions-card">
            <template #header>
              <div class="card-header">
                <el-icon><Grid /></el-icon>
                <span>快捷操作</span>
              </div>
            </template>
            
            <div class="quick-actions">
              <el-button class="action-item" @click="$router.push('/articles')">
                <el-icon><Document /></el-icon>
                浏览文章
              </el-button>
              <el-button class="action-item" @click="$router.push('/creator')">
                <el-icon><Edit /></el-icon>
                创作中心
              </el-button>
              <el-button 
                class="action-item" 
                v-if="userStore.userInfo?.role === 2" 
                @click="$router.push('/admin')"
              >
                <el-icon><Setting /></el-icon>
                后台管理
              </el-button>
            </div>
          </el-card>
        </div>
      </div>
    </div>

    <!-- 编辑资料对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑个人资料"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="用户名">
          <el-input v-model="editForm.username" placeholder="用户名不可修改" disabled />
        </el-form-item>
        
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        
        <el-form-item label="邮箱">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        
        <el-form-item label="头像 URL">
          <el-input v-model="editForm.avatar" placeholder="请输入头像链接" />
        </el-form-item>
        
        <el-form-item label="角色">
          <el-tag :type="getRoleType(editForm.role)" effect="dark">
            <el-icon><Star /></el-icon>
            {{ getRoleText(editForm.role) }}
          </el-tag>
          <div class="form-tip"></div>
        </el-form-item>
        
        <el-form-item label="注册时间">
          <span>{{ formatDate(editForm.createdAt) }}</span>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSaveProfile" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="修改密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="当前密码">
          <el-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="handleChangePassword" :loading="changingPassword">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  Edit, User, Lock, Message, 
  Calendar, Star, Check, Close, 
  DataAnalysis, Document, Grid, Setting, View
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getArticles, getUserStats } from '@/api/auth'
import TopNavBar from '@/components/TopNavBar.vue'

const router = useRouter()
const userStore = useUserStore()

// 统计数据
const stats = reactive({
  articles: 0,
  favorites: 0,
  following: 0,
  followers: 0
})

// 加载用户统计
const loadUserStats = async () => {
  try {
    const res = await getUserStats()
    if (res.code === 0 && res.data) {
      stats.articles = res.data.articles || 0
      stats.favorites = res.data.favorites || 0
      stats.following = res.data.following || 0
      stats.followers = res.data.followers || 0
    }
  } catch (error) {
    console.error('加载用户统计失败:', error)
  }
}

// 近期阅读文章
const recentArticles = ref([])
const loadingRecent = ref(false)

// 加载近期阅读
const loadRecentArticles = async () => {
  loadingRecent.value = true
  try {
    const res = await getArticles({ page: 1, page_size: 5 })
    recentArticles.value = res.data.articles || []
  } catch (error) {
    console.error('加载近期阅读失败:', error)
  } finally {
    loadingRecent.value = false
  }
}

// 编辑表单
const showEditDialog = ref(false)
const saving = ref(false)
const editForm = reactive({
  username: '',
  nickname: '',
  email: '',
  avatar: '',
  role: '',
  createdAt: ''
})

// 修改密码相关
const showPasswordDialog = ref(false)
const changingPassword = ref(false)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

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

// 获取头像 URL
const getAvatarUrl = (avatar) => {
  if (avatar && avatar.startsWith('http')) {
    return avatar
  }
  return getAvatarFromNickname(userStore.userInfo?.nickname)
}

// 获取用户首字母
const getUserInitial = (nickname) => {
  if (!nickname) return 'U'
  const firstChar = nickname.charAt(0).toUpperCase()
  return /[A-Za-z]/.test(firstChar) ? firstChar : 'U'
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

// 获取角色类型
const getRoleType = (role) => {
  if (role === 'admin') return 'danger'
  if (role === 'editor') return 'warning'
  return 'info'
}

// 获取角色文字
const getRoleText = (role) => {
  if (role === 'admin') return '管理员'
  if (role === 'editor') return '编辑'
  return '普通用户'
}

// 切换编辑状态
const toggleEdit = () => {
  const userInfo = userStore.userInfo
  if (userInfo) {
    editForm.username = userInfo.username || ''
    editForm.nickname = userInfo.nickname || ''
    editForm.email = userInfo.email || ''
    editForm.avatar = userInfo.avatar || ''
    editForm.role = userInfo.role || ''
    editForm.createdAt = userInfo.createdAt || ''
  }
  showEditDialog.value = true
}

// 保存个人资料
const handleSaveProfile = async () => {
  if (!editForm.nickname.trim()) {
    ElMessage.warning('请输入昵称')
    return
  }
  
  saving.value = true
  
  try {
    // TODO: 调用后端 API 更新用户信息
    // await updateUserProfile(editForm)
    
    // 更新本地存储
    await userStore.fetchUserInfo()
    
    ElMessage.success('保存成功')
    showEditDialog.value = false
  } catch (error) {
    ElMessage.error('保存失败')
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  if (!passwordForm.oldPassword) {
    ElMessage.warning('请输入当前密码')
    return
  }
  
  if (!passwordForm.newPassword) {
    ElMessage.warning('请输入新密码')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的新密码不一致')
    return
  }
  
  changingPassword.value = true
  
  try {
    // TODO: 调用后端 API 修改密码
    // await changePassword(passwordForm)
    
    ElMessage.success('密码修改成功，请重新登录')
    showPasswordDialog.value = false
    userStore.logout()
    router.push('/login')
  } catch (error) {
    ElMessage.error('密码修改失败')
    console.error('密码修改失败:', error)
  } finally {
    changingPassword.value = false
  }
}

// 初始化
onMounted(() => {
  if (!userStore.userInfo) {
    userStore.autoLogin()
  }
  // 加载近期阅读文章
  loadRecentArticles()
  // 加载用户统计
  loadUserStats()
})
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e2a4a 0%, #2d1f3d 50%, #1a3a5c 100%);
  padding: 24px;
  
  .profile-container {
    max-width: 1200px;
    margin: 0 auto;
    
    // 用户信息卡片
    .profile-header-card {
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      margin-bottom: 24px;
      
      :deep(.el-card__body) {
        padding: 0;
      }
      
      .profile-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 32px;
        
        .avatar-info {
          display: flex;
          align-items: center;
          gap: 24px;
          
          .profile-avatar {
            border: 3px solid rgba(255, 255, 255, 0.2);
            border-radius: 16px;
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
          }
          
          .user-basic-info {
            .user-name-row {
              display: flex;
              align-items: center;
              gap: 16px;
              margin-bottom: 12px;
              
              .user-nickname {
                font-size: 26px;
                font-weight: 700;
                color: rgba(255, 255, 255, 0.95);
                margin: 0;
              }
              
              .el-tag {
                font-size: 13px;
                padding: 4px 12px;
              }
            }
            
            .user-meta {
              display: flex;
              align-items: center;
              gap: 8px;
              margin-bottom: 12px;
              color: rgba(255, 255, 255, 0.6);
              font-size: 14px;
              
              .username {
                color: rgba(255, 255, 255, 0.5);
              }
              
              .divider {
                color: rgba(255, 255, 255, 0.3);
              }
              
              .join-date {
                display: flex;
                align-items: center;
                gap: 4px;
              }
            }
            
            .user-bio {
              font-size: 14px;
              color: rgba(255, 255, 255, 0.7);
              line-height: 1.6;
              margin: 0;
            }
          }
        }
        
        .profile-actions {
          .edit-btn {
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
    
    // 主要内容区
    .profile-content {
      display: grid;
      grid-template-columns: 1fr 320px;
      gap: 24px;
      
      .main-content {
        display: flex;
        flex-direction: column;
        gap: 24px;
        
        .recent-views-card {
          background: rgba(255, 255, 255, 0.05);
          backdrop-filter: blur(10px);
          border: 1px solid rgba(255, 255, 255, 0.1);
          
          :deep(.el-card__header) {
            background: rgba(255, 255, 255, 0.03);
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            padding: 16px 20px;
            
            .card-header {
              display: flex;
              align-items: center;
              gap: 10px;
              color: rgba(255, 255, 255, 0.95) !important;
              font-weight: 600;
              font-size: 16px;
              
              .el-icon {
                color: rgba(102, 126, 234, 0.8);
              }
            }
          }
          
          :deep(.el-card__body) {
            padding: 20px;
            color: rgba(255, 255, 255, 0.95) !important;
          }
          
          .article-list {
            display: flex;
            flex-direction: column;
            gap: 16px;
            
            .empty-state {
              padding: 40px 20px;
              text-align: center;
            }
            
            .article-item {
              display: flex;
              gap: 16px;
              padding: 16px;
              background: rgba(255, 255, 255, 0.03);
              border: 1px solid rgba(255, 255, 255, 0.08);
              border-radius: 8px;
              cursor: pointer;
              transition: all 0.3s ease;
              
              &:hover {
                background: rgba(255, 255, 255, 0.06);
                border-color: rgba(102, 126, 234, 0.4);
                transform: translateX(4px);
              }
              
              .article-cover-wrapper {
                flex-shrink: 0;
                width: 160px;
                height: 100px;
                position: relative;
                
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
                    font-size: 40px;
                    color: rgba(102, 126, 234, 0.6);
                  }
                }
              }
              
              .article-info {
                flex: 1;
                min-width: 0;
                
                .article-title {
                  font-size: 16px;
                  font-weight: 600;
                  color: rgba(255, 255, 255, 0.95);
                  margin: 0 0 8px;
                  overflow: hidden;
                  text-overflow: ellipsis;
                  display: -webkit-box;
                  -webkit-line-clamp: 2;
                  line-clamp: 2;
                  -webkit-box-orient: vertical;
                  line-height: 1.4;
                }
                
                .article-summary {
                  font-size: 14px;
                  color: rgba(255, 255, 255, 0.6);
                  margin: 0 0 10px;
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
          }
        }
        
        // 安全卡片 - 修改密码入口
        .security-card {
          .security-items {
            .security-item {
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding: 16px 0;
              
              .security-info {
                .security-label {
                  font-size: 15px;
                  color: rgba(255, 255, 255, 0.95);
                  margin-bottom: 6px;
                  font-weight: 500;
                }
                
                .security-desc {
                  font-size: 13px;
                  color: rgba(255, 255, 255, 0.6);
                }
              }
              
              .el-button {
                font-weight: 600;
                color: rgba(102, 126, 234, 0.9);
              }
            }
          }
        }
      }
      
      // 右侧边栏
      .sidebar {
        .quick-actions-card {
          background: rgba(255, 255, 255, 0.05);
          backdrop-filter: blur(10px);
          border: 1px solid rgba(255, 255, 255, 0.1);
          
          :deep(.el-card__header) {
            background: rgba(255, 255, 255, 0.03);
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            padding: 16px 20px;
            
            .card-header {
              display: flex;
              align-items: center;
              gap: 10px;
              color: rgba(255, 255, 255, 0.95) !important;
              font-weight: 600;
              font-size: 16px;
              
              .el-icon {
                color: rgba(102, 126, 234, 0.8);
              }
            }
          }
          
          :deep(.el-card__body) {
            padding: 20px;
            color: rgba(255, 255, 255, 0.95) !important;
          }
          
          .quick-actions {
            display: flex;
            flex-direction: column;
            gap: 12px;
            width: 100%;
            
            .action-item {
              width: 100%;
              height: 48px;
              font-size: 15px;
              display: flex;
              align-items: center;
              justify-content: flex-start;
              gap: 0;
              background: rgba(255, 255, 255, 0.05);
              border: 1px solid rgba(255, 255, 255, 0.1);
              color: rgba(255, 255, 255, 0.95) !important;
              padding: 0 16px;
              text-align: left;
              line-height: 1;
              box-sizing: border-box;
              position: relative;
              
              // 强制重置所有可能的 margin
              margin-left: 0 !important;
              margin-right: 0 !important;
              
              // 覆盖 Element Plus 按钮默认的 flex 布局
              :deep(.el-button__content) {
                display: flex;
                align-items: center;
                gap: 10px;
                width: 100%;
                line-height: 1;
                justify-content: flex-start;
                margin-left: 0 !important;
              }
              
              &:hover {
                background: rgba(102, 126, 234, 0.15);
                border-color: rgba(102, 126, 234, 0.4);
                color: rgba(255, 255, 255, 0.95) !important;
              }
              
              .el-icon {
                width: 20px;
                min-width: 20px;
                font-size: 18px;
                color: rgba(102, 126, 234, 0.8);
                flex-shrink: 0;
                display: inline-flex;
                align-items: center;
                justify-content: center;
                vertical-align: middle;
                line-height: 1;
                margin-left: 0 !important;
                margin-right: 10px;
              }
            }
          }
        }
      }
    }
  }
  
  // 编辑资料对话框样式
  :deep(.el-dialog) {
    background: rgba(30, 42, 74, 0.95);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    
    .el-dialog__title {
      color: rgba(255, 255, 255, 0.95);
      font-weight: 600;
    }
    
    .el-dialog__body {
      color: rgba(255, 255, 255, 0.9);
      padding-top: 20px;
    }
    
    .el-form-item__label {
      color: rgba(255, 255, 255, 0.8) !important;
      font-weight: 500;
    }
    
    .el-input__wrapper {
      background: rgba(255, 255, 255, 0.05);
      border-color: rgba(255, 255, 255, 0.1);
      
      .el-input__inner {
        color: rgba(255, 255, 255, 0.95);
      }
    }
    
    .form-tip {
      font-size: 12px;
      color: rgba(255, 255, 255, 0.5);
      margin-top: 4px;
    }
  }
  
  // 修改密码对话框样式
  :deep(.el-dialog) {
    background: rgba(30, 42, 74, 0.95);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    
    .el-dialog__title {
      color: rgba(255, 255, 255, 0.95);
      font-weight: 600;
    }
    
    .el-dialog__body {
      color: rgba(255, 255, 255, 0.9);
      padding-top: 20px;
    }
    
    .el-form-item__label {
      color: rgba(255, 255, 255, 0.8) !important;
      font-weight: 500;
    }
    
    .el-input__wrapper {
      background: rgba(255, 255, 255, 0.05);
      border-color: rgba(255, 255, 255, 0.1);
      
      .el-input__inner {
        color: rgba(255, 255, 255, 0.95);
      }
    }
  }
  
  // 响应式设计
  @media (max-width: 768px) {
    padding: 16px;
    
    .profile-container {
      .profile-header-card {
        .profile-header {
          flex-direction: column;
          text-align: center;
          gap: 20px;
          
          .avatar-info {
            flex-direction: column;
            text-align: center;
          }
        }
      }
      
      .stats-cards {
        grid-template-columns: repeat(2, 1fr);
      }
      
      .profile-content {
        grid-template-columns: 1fr;
      }
    }
  }
}
</style>