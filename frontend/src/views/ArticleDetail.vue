<template>
  <div class="article-detail-page">
    <TopNavBar />
    
    <div class="article-back" v-if="!loading">
      <el-button @click="goBack" text>
        <el-icon><ArrowLeft /></el-icon>
        返回文章列表
      </el-button>
    </div>

    <div class="article-container" v-loading="loading">
      <article v-if="article" class="article-content">
        <header class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>
          
          <!-- 作者信息区域 -->
          <div class="author-section">
            <div class="author-info">
              <el-avatar :size="50" :src="article.author?.avatar || defaultAvatar" />
              <div class="author-details">
                <div class="author-name">{{ article.author?.nickname || article.author?.username || '匿名用户' }}</div>
                <div class="author-meta">
                  <span class="meta-label">发布于 {{ formatDate(article.created_at) }}</span>
                  <span class="meta-divider">•</span>
                  <span class="meta-label">{{ article.view_count }} 次阅读</span>
                </div>
              </div>
            </div>
            <div class="author-actions">
              <el-button 
                v-if="showFollowButton" 
                @click="toggleFollow" 
                :type="isFollowing ? 'success' : 'primary'"
                :icon="isFollowing ? 'Check' : 'Plus'"
                size="default"
              >
                {{ isFollowing ? '已关注' : '关注' }}
              </el-button>
              <div class="follower-count">{{ followerCount }} 粉丝</div>
            </div>
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
        </header>

        <div class="article-cover" v-if="article.cover_image">
          <img :src="article.cover_image" :alt="article.title" />
        </div>

        <!-- 文章正文内容 - 使用简单的 Markdown 渲染方法 -->
        <div class="article-body" v-html="renderedContent"></div>

        <!-- 点赞和互动区域 -->
        <div class="article-interactions">
          <div class="interaction-item">
            <el-button 
              @click="toggleLike" 
              :type="isLiked ? 'danger' : 'default'"
              :icon="isLiked ? 'StarFilled' : 'Star'"
              size="large"
              circle
              class="like-button"
            />
            <span class="interaction-count">{{ likeCount }}</span>
            <span class="interaction-label">点赞</span>
          </div>
        </div>

        <!-- 评论区域 -->
        <div class="article-comments">
          <h2 class="comments-title">评论 ({{ comments.length }})</h2>
          
          <!-- 发表评论 -->
          <div class="comment-form" v-if="isLoggedIn">
            <el-input
              v-model="newComment"
              type="textarea"
              :rows="3"
              placeholder="写下你的评论..."
              maxlength="500"
              show-word-limit
              class="comment-input"
            />
            <el-button 
              @click="submitComment" 
              type="primary" 
              size="default"
              :loading="submitting"
              class="submit-comment-btn"
            >
              发表评论
            </el-button>
          </div>
          
          <div class="no-login-tip" v-else>
            <el-alert
              title="登录后参与评论"
              type="info"
              :closable="false"
              show-icon
            />
          </div>
          
          <!-- 评论列表 -->
          <div class="comment-list" v-if="comments && comments.length">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="comment-header">
                <div class="comment-author-info" @click="goToUserProfile(comment.user?.id)" :class="{ clickable: comment.user?.id }">
                  <el-avatar :size="36" :src="comment.user?.avatar || defaultAvatar" />
                  <div class="comment-author-details">
                    <span class="comment-author">
                      {{ comment.user?.nickname || comment.user?.username || comment.author || '匿名用户' }}
                    </span>
                    <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
                  </div>
                </div>
                <el-button 
                  v-if="canDeleteComment(comment)" 
                  @click="deleteCommentFunc(comment.id)" 
                  type="danger" 
                  size="small" 
                  text
                >
                  删除
                </el-button>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
            </div>
          </div>
          
          <el-empty v-else description="暂无评论，快来抢沙发吧~" />
        </div>
      </article>

      <el-empty v-if="!loading && !article" description="文章不存在" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, Calendar, View, Folder, Star, StarFilled, Plus, Check } from '@element-plus/icons-vue'
import { getArticle, likeArticle, getLikeStatus, followUser, getFollowStatus, createComment, getComments, deleteComment } from '@/api/auth'
import TopNavBar from '@/components/TopNavBar.vue'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const article = ref(null)
const comments = ref([])
const newComment = ref('')
const submitting = ref(false)

// 点赞状态
const isLiked = ref(false)
const likeCount = ref(0)

// 关注状态
const isFollowing = ref(false)
const followerCount = ref(0)
const showFollowButton = ref(false)

// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'

// 是否已登录
const isLoggedIn = computed(() => userStore.token && !userStore.isTokenExpired())

// 当前用户 ID
const currentUserId = computed(() => userStore.user?.id)

// 是否可以删除评论（评论作者或管理员）
const canDeleteComment = (comment) => {
  if (!isLoggedIn.value) return false
  const isAdmin = userStore.user?.role === 0
  const isAuthor = comment.user?.id === currentUserId.value
  return isAdmin || isAuthor
}

// 使用简单的 Markdown 转 HTML 方法
const renderedContent = computed(() => {
  try {
    if (!article.value || !article.value.content) {
      return ''
    }
    
    const markdown = article.value.content
    
    // 用于存储代码块的数组
    const codeBlocks = []
    const inlineCodes = []
    
    // 第一步：提取代码块，用特殊占位符（不会被段落处理影响）
    let html = markdown
      .replace(/```([\s\S]*?)```/g, (match, code) => {
        codeBlocks.push(code)
        return `\n___CODE_BLOCK_${codeBlocks.length - 1}___\n`
      })
      .replace(/`(.*?)`/g, (match, code) => {
        inlineCodes.push(code)
        return `___INLINE_CODE_${inlineCodes.length - 1}___`
      })
    
    // 第二步：按行处理，识别不同的块级元素
    const lines = html.split('\n')
    const processedLines = []
    
    for (let line of lines) {
      // 跳过代码块占位符（已经单独处理）
      if (line.trim().match(/^___CODE_BLOCK_\d+___$/)) {
        processedLines.push(line.trim())
        continue
      }
      
      // 处理行内语法（标题、粗体等）
      line = line
        .replace(/^# (.*$)/, '<h1>$1</h1>')
        .replace(/^## (.*$)/, '<h2>$1</h2>')
        .replace(/^### (.*$)/, '<h3>$1</h3>')
        .replace(/^#### (.*$)/, '<h4>$1</h4>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/^> (.*$)/, '<blockquote>$1</blockquote>')
      
      // 处理列表
      if (line.match(/^\- (.*$)/)) {
        line = line.replace(/^\- (.*$)/, '<li>$1</li>')
      } else if (line.match(/^\d+\. (.*$)/)) {
        line = line.replace(/^\d+\. (.*$)/, '<li>$1</li>')
      }
      
      // 处理链接和图片
      line = line
        .replace(/!\[(.*?)\]\((.*?)\)/g, '<img src="$2" alt="$1" />')
        .replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" target="_blank">$1</a>')
      
      // 空行跳过
      if (line.trim() === '') continue
      
      // 普通文本包裹为段落
      if (!line.match(/^<(h[1-6]|blockquote|li|pre|ul|ol)/)) {
        line = `<p>${line}</p>`
      }
      
      processedLines.push(line)
    }
    
    // 第三步：合并所有行
    html = processedLines.join('')
    
    // 第四步：恢复代码块（HTML 转义）
    html = html
      .replace(/___CODE_BLOCK_(\d+)___/g, (match, index) => {
        const code = codeBlocks[parseInt(index)]
        const escaped = code
          .replace(/&/g, '&amp;')
          .replace(/</g, '&lt;')
          .replace(/>/g, '&gt;')
          .replace(/"/g, '&quot;')
          .replace(/'/g, '&#39;')
        return `<pre><code>${escaped}</code></pre>`
      })
      .replace(/___INLINE_CODE_(\d+)___/g, (match, index) => {
        const code = inlineCodes[parseInt(index)]
        const escaped = code
          .replace(/&/g, '&amp;')
          .replace(/</g, '&lt;')
          .replace(/>/g, '&gt;')
          .replace(/"/g, '&quot;')
          .replace(/'/g, '&#39;')
        return `<code>${escaped}</code>`
      })
    
    // 第五步：清理多余标签
    html = html
      .replace(/<p>\s*<\/p>/g, '')
      .replace(/<\/p>\s*<h([1-6])/g, '</p><h$1')
      .replace(/<\/blockquote>\s*<p>/g, '</blockquote>')
      .replace(/<\/li>\s*<li>/g, '</li><li>')
    
    return html
  } catch (error) {
    console.error('Markdown 渲染错误:', error)
    return article.value?.content || ''
  }
})

// 加载文章详情
const loadArticle = async () => {
  loading.value = true
  try {
    const res = await getArticle(route.params.id)
    article.value = res.data.article
    
    // 加载评论
    if (res.data.comments) {
      comments.value = res.data.comments
    } else {
      // 如果后端没有返回评论，单独获取
      try {
        const commentsRes = await getComments(route.params.id)
        comments.value = commentsRes.data.comments || []
      } catch (commentError) {
        // 静默处理，使用空列表
        comments.value = []
      }
    }
    
    // 加载点赞状态
    await loadLikeStatus()
    
    // 加载关注状态（如果不是自己的文章）
    if (article.value.author?.id && article.value.author.id !== currentUserId.value) {
      showFollowButton.value = true
      await loadFollowStatus()
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

// 加载点赞状态
const loadLikeStatus = async () => {
  try {
    const res = await getLikeStatus(route.params.id)
    isLiked.value = res.data.liked
    likeCount.value = res.data.count
  } catch (error) {
    console.error('加载点赞状态失败:', error)
  }
}

// 加载关注状态
const loadFollowStatus = async () => {
  try {
    const res = await getFollowStatus(article.value.author.id)
    isFollowing.value = res.data.following
    followerCount.value = res.data.count
  } catch (error) {
    console.error('加载关注状态失败:', error)
  }
}

// 切换点赞状态
const toggleLike = async () => {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    const res = await likeArticle(route.params.id)
    isLiked.value = res.data.liked
    likeCount.value = res.data.count
    ElMessage.success(res.data.message)
  } catch (error) {
    console.error('点赞失败:', error)
    ElMessage.error(error.response?.data?.message || '点赞失败')
  }
}

// 切换关注状态
const toggleFollow = async () => {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    const res = await followUser(article.value.author.id)
    isFollowing.value = res.data.following
    followerCount.value = res.data.count
    ElMessage.success(res.data.message)
  } catch (error) {
    console.error('关注失败:', error)
    ElMessage.error(error.response?.data?.message || '关注失败')
  }
}

// 发表评论
const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  
  submitting.value = true
  try {
    const res = await createComment(route.params.id, {
      content: newComment.value.trim()
    })
    
    // 添加到评论列表
    comments.value.unshift(res.data.comment)
    newComment.value = ''
    ElMessage.success('评论成功')
  } catch (error) {
    console.error('评论失败:', error)
    ElMessage.error(error.response?.data?.message || '评论失败')
  } finally {
    submitting.value = false
  }
}

// 删除评论
const deleteCommentFunc = async (commentId) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteComment(commentId)
    comments.value = comments.value.filter(c => c.id !== commentId)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除评论失败:', error)
      ElMessage.error(error.response?.data?.message || '删除评论失败')
    }
  }
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

// 返回上一页
const goBack = () => {
  router.back()
}

// 跳转到用户个人主页
const goToUserProfile = (userId) => {
  if (userId) {
    router.push(`/user/${userId}`)
  }
}

onMounted(() => {
  loadArticle()
})
</script>

<style lang="scss" scoped>
.article-detail-page {
  min-height: 100vh;
  padding: 40px 20px;
  background: linear-gradient(135deg, #1e2a4a 0%, #2d1f3d 50%, #1a3a5c 100%);
  
  .article-back {
    max-width: 900px;
    margin: 0 auto 20px;
    
    :deep(.el-button) {
      color: rgba(255, 255, 255, 0.8);
      font-size: 14px;
      
      &:hover {
        color: #fff;
      }
    }
  }
  
  .article-container {
    max-width: 900px;
    margin: 0 auto;
    
    .article-content {
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(20px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 12px;
      padding: 40px;
      animation: fadeIn 0.5s ease-out;
      
      @keyframes fadeIn {
        from {
          opacity: 0;
          transform: translateY(20px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }
      
      .article-header {
        margin-bottom: 30px;
        
        .article-title {
          color: #fff;
          font-size: 32px;
          margin: 0 0 20px;
          font-weight: 400;
          line-height: 1.4;
          letter-spacing: 1px;
        }
        
        // 作者信息区域
        .author-section {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 20px;
          background: rgba(255, 255, 255, 0.08);
          border-radius: 12px;
          margin-bottom: 20px;
          border: 1px solid rgba(255, 255, 255, 0.15);
          
          .author-info {
            display: flex;
            align-items: center;
            gap: 15px;
            
            .author-details {
              .author-name {
                color: #fff;
                font-size: 16px;
                font-weight: 500;
                margin-bottom: 5px;
              }
              
              .author-meta {
                display: flex;
                align-items: center;
                gap: 8px;
                color: rgba(255, 255, 255, 0.5);
                font-size: 13px;
                
                .meta-divider {
                  color: rgba(255, 255, 255, 0.3);
                }
              }
            }
          }
          
          .author-actions {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 8px;
            
            .follower-count {
              color: rgba(255, 255, 255, 0.6);
              font-size: 12px;
            }
          }
        }
        
        .article-tags {
          display: flex;
          gap: 8px;
          flex-wrap: wrap;
          
          .el-tag {
            background: rgba(255, 255, 255, 0.1);
            border-color: rgba(255, 255, 255, 0.2);
            color: rgba(255, 255, 255, 0.8);
          }
        }
      }
      
      .article-cover {
        margin: 30px 0;
        border-radius: 12px;
        overflow: hidden;
        max-height: 500px;
        
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }
      
      .article-body {
        color: rgba(255, 255, 255, 0.9);
        font-size: 16px;
        line-height: 1.8;
        
        :deep(p) {
          margin: 1.5em 0;
        }
        
        :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
          color: #fff;
          margin: 1.5em 0 1em;
          font-weight: 400;
        }
        
        :deep(h1) { font-size: 28px; }
        :deep(h2) { font-size: 24px; }
        :deep(h3) { font-size: 20px; }
        :deep(h4) { font-size: 18px; }
        
        :deep(img) {
          max-width: 100%;
          height: auto;
          border-radius: 8px;
          margin: 1.5em 0;
        }
        
        :deep(a) {
          color: #409eff;
          text-decoration: none;
          
          &:hover {
            text-decoration: underline;
          }
        }
        
        :deep(blockquote) {
          margin: 1.5em 0;
          padding: 15px 20px;
          background: rgba(255, 255, 255, 0.05);
          border-left: 4px solid #409eff;
          border-radius: 4px;
          color: rgba(255, 255, 255, 0.7);
        }
        
        :deep(pre) {
          margin: 1.5em 0;
          padding: 20px;
          background: rgba(0, 0, 0, 0.3);
          border-radius: 8px;
          overflow-x: auto;
          
          code {
            font-family: 'Fira Code', 'Consolas', monospace;
            font-size: 14px;
          }
        }
        
        :deep(ul), :deep(ol) {
          padding-left: 2em;
          margin: 1.5em 0;
        }
        
        :deep(li) {
          margin: 0.5em 0;
        }
        
        :deep(table) {
          width: 100%;
          border-collapse: collapse;
          margin: 1.5em 0;
          
          th, td {
            border: 1px solid rgba(255, 255, 255, 0.2);
            padding: 12px;
            text-align: left;
          }
          
          th {
            background: rgba(255, 255, 255, 0.1);
            font-weight: 400;
          }
        }
      }
      
      // 互动区域
      .article-interactions {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 30px;
        padding: 30px 0;
        margin-top: 30px;
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        
        .interaction-item {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 8px;
          
          .like-button {
            width: 60px;
            height: 60px;
            font-size: 24px;
            transition: all 0.3s ease;
            
            &:hover {
              transform: scale(1.1);
            }
          }
          
          .interaction-count {
            color: rgba(255, 255, 255, 0.8);
            font-size: 16px;
            font-weight: 500;
          }
          
          .interaction-label {
            color: rgba(255, 255, 255, 0.5);
            font-size: 13px;
          }
        }
      }
      
      .article-comments {
        margin-top: 40px;
        padding-top: 30px;
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        
        .comments-title {
          color: #fff;
          font-size: 20px;
          margin: 0 0 25px;
          font-weight: 400;
        }
        
        // 评论表单
        .comment-form {
          margin-bottom: 30px;
          
          .comment-input {
            margin-bottom: 15px;
            
            :deep(.el-textarea__inner) {
              background: rgba(255, 255, 255, 0.08);
              border-color: rgba(255, 255, 255, 0.2);
              color: rgba(255, 255, 255, 0.9);
              
              &::placeholder {
                color: rgba(255, 255, 255, 0.4);
              }
            }
          }
          
          .submit-comment-btn {
            width: 100%;
          }
        }
        
        .no-login-tip {
          margin-bottom: 30px;
          
          :deep(.el-alert) {
            background: rgba(255, 255, 255, 0.08);
            border-color: rgba(255, 255, 255, 0.2);
            
            .el-alert__title {
              color: rgba(255, 255, 255, 0.8);
            }
          }
        }
        
        .comment-list {
          display: flex;
          flex-direction: column;
          gap: 20px;
          
          .comment-item {
            padding: 20px;
            background: rgba(255, 255, 255, 0.05);
            border-radius: 8px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            
            .comment-header {
              display: flex;
              justify-content: space-between;
              align-items: center;
              margin-bottom: 12px;
              
              .comment-author-info {
                display: flex;
                align-items: center;
                gap: 12px;
                cursor: pointer;
                transition: opacity 0.2s;
                
                &:hover {
                  opacity: 0.8;
                }
                
                .comment-avatar {
                  cursor: pointer;
                }
                
                .comment-author-details {
                  display: flex;
                  flex-direction: column;
                  gap: 4px;
                  
                  .comment-author {
                    font-weight: 500;
                    color: rgba(255, 255, 255, 0.9);
                  }
                  
                  .comment-date {
                    font-size: 12px;
                    color: rgba(255, 255, 255, 0.5);
                  }
                }
              }
            }
            
            .comment-content {
              color: rgba(255, 255, 255, 0.8);
              line-height: 1.6;
              font-size: 14px;
              white-space: pre-wrap;
            }
          }
        }
      }
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .article-detail-page {
    padding: 20px 15px;
    
    .article-container {
      .article-content {
        padding: 25px 20px;
        
        .article-header {
          .article-title {
            font-size: 24px;
          }
          
          .author-section {
            flex-direction: column;
            gap: 15px;
            
            .author-info {
              width: 100%;
            }
            
            .author-actions {
              width: 100%;
              flex-direction: row;
              justify-content: center;
            }
          }
        }
        
        .article-body {
          font-size: 15px;
        }
        
        .article-interactions {
          gap: 20px;
          
          .interaction-item {
            .like-button {
              width: 50px;
              height: 50px;
              font-size: 20px;
            }
          }
        }
      }
    }
  }
}
</style>
