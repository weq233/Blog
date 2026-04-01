<template>
  <div class="article-editor-page">
    <TopNavBar />
    
    <div class="editor-container">
      <!-- 页面头部 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon><Edit /></el-icon>
          {{ isEditMode ? '编辑文章' : '创建文章' }}
        </h1>
        <div class="header-actions">
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="saving">
            {{ isEditMode ? '保存修改' : '发布文章' }}
          </el-button>
        </div>
      </div>

      <!-- 编辑表单和预览并排布局 -->
      <div class="editor-form-wrapper">
        <!-- 左侧编辑器 -->
        <div class="editor-panel">
          <el-form :model="form" label-width="100px" label-position="top">
            <el-form-item label="文章标题" required>
              <el-input
                v-model="form.title"
                placeholder="请输入文章标题"
                maxlength="200"
                size="large"
                class="dark-input"
              />
            </el-form-item>

            <el-form-item label="文章摘要" required>
              <el-input
                v-model="form.summary"
                type="textarea"
                :rows="3"
                placeholder="请输入文章摘要（200 字以内）"
                maxlength="200"
                show-word-limit
                class="dark-input"
              />
            </el-form-item>

            <el-form-item label="封面图片">
              <div class="cover-upload-area">
                <el-upload
                  ref="coverUploadRef"
                  class="cover-uploader"
                  :auto-upload="false"
                  :limit="1"
                  :show-file-list="false"
                  :on-change="handleCoverChange"
                  :on-remove="handleCoverRemove"
                  accept="image/*"
                  drag
                >
                  <div class="cover-uploader-content" v-if="!coverPreviewUrl">
                    <el-icon class="cover-uploader-icon"><Upload /></el-icon>
                    <div class="cover-uploader-text">
                      <p class="cover-uploader-title">点击或拖拽上传封面图片</p>
                      <p class="cover-uploader-tip">支持 jpg/png 格式，图片大小不超过 5MB</p>
                    </div>
                  </div>
                  <div v-else class="cover-preview-wrapper">
                    <img :src="coverPreviewUrl" alt="封面预览" @error="handleImageError" class="cover-preview-image" />
                    <div class="cover-overlay">
                      <el-button type="danger" size="small" @click.stop="handleCoverRemove">
                        <el-icon><Close /></el-icon>
                        删除封面
                      </el-button>
                    </div>
                  </div>
                </el-upload>
              </div>
            </el-form-item>

            <el-form-item label="文章内容" required>
              <div class="content-editor-wrapper">
                <!-- 左侧编辑器 -->
                <div class="vditor-wrapper">
                  <div ref="vditorContainer" class="vditor-container"></div>
                </div>
                
                <!-- 右侧预览面板 - 使用 v-show 控制显示/隐藏 -->
                <div class="inline-preview-panel" v-show="showPreview">
                  <div class="preview-header">
                    <h3 class="preview-title">
                      <el-icon><View /></el-icon>
                      文章预览
                    </h3>
                    <el-button text :icon="Close" @click="showPreview = false" class="close-btn"></el-button>
                  </div>
                  <div class="preview-content">
                    <article class="preview-article">
                      <header class="article-header">
                        <h1 class="article-title">{{ form.title || '文章标题' }}</h1>
                        <div class="article-meta">
                          <span class="meta-item">
                            <el-icon><Clock /></el-icon>
                            {{ currentDate }}
                          </span>
                          <span class="meta-item" v-if="coverPreviewUrl">
                            <el-icon><Picture /></el-icon>
                            有封面图
                          </span>
                        </div>
                      </header>

                      <div class="article-cover" v-if="coverPreviewUrl">
                        <img 
                          :src="coverPreviewUrl" 
                          alt="封面图" 
                          @error="handleImageError"
                          @load="handleImageLoad"
                        />
                      </div>

                      <div class="article-summary" v-if="form.summary">
                        <h3>摘要</h3>
                        <p>{{ form.summary }}</p>
                      </div>

                      <!-- 正文内容 -->
                      <div v-if="previewContent" class="article-body vditor-reset" v-html="previewContent"></div>
                      <div v-else class="article-body" style="color: rgba(255,255,255,0.5); text-align: center; padding: 20px;">
                        <p>暂无正文内容</p>
                      </div>
                    </article>
                  </div>
                </div>
                
                <!-- 预览按钮放在编辑器右上角 -->
                <el-button 
                  class="preview-toggle-btn"
                  :type="showPreview ? 'warning' : 'success'" 
                  @click="showPreview = !showPreview" 
                  :icon="showPreview ? 'Hide' : View"
                  size="small"
                  circle>
                </el-button>
              </div>
            </el-form-item>

            <el-form-item label="文章分类">
              <el-tabs type="border-card" class="category-tabs">
                <el-tab-pane label="标签">
                  <el-select
                    v-model="form.category_id"
                    placeholder="请选择标签"
                    style="width: 100%"
                    class="dark-select"
                    clearable
                    @change="handleSystemCategoryChange"
                  >
                    <el-option
                      v-for="category in categories"
                      :key="category.id"
                      :label="category.name"
                      :value="category.id"
                    />
                  </el-select>
                  <p class="category-tip">标签由管理员统一管理</p>
                </el-tab-pane>
                
                <el-tab-pane label="合集">
                  <div class="my-category-wrapper">
                    <el-select
                      v-model="form.user_category_id"
                      placeholder="请选择合集"
                      style="width: 100%"
                      class="dark-select"
                      clearable
                      @change="handleUserCategoryChange"
                    >
                      <el-option
                        v-for="cat in userCategories"
                        :key="cat.id"
                        :label="cat.name"
                        :value="cat.id"
                      >
                        <span>{{ cat.name }}</span>
                        <span 
                          class="category-color-dot" 
                          :style="{ backgroundColor: cat.color }"
                        ></span>
                      </el-option>
                    </el-select>
                    
                    <div class="category-actions">
                      <el-button 
                        type="primary" 
                        size="small" 
                        @click="showCreateCategory"
                        :icon="Plus"
                      >
                        新建合集
                      </el-button>
                      <el-button 
                        v-if="selectedUserCategoryId"
                        type="warning" 
                        size="small" 
                        @click="editCurrentCategory"
                        :icon="Edit"
                      >
                        编辑分类
                      </el-button>
                      <el-button 
                        v-if="selectedUserCategoryId"
                        type="danger" 
                        size="small" 
                        @click="deleteCurrentCategory"
                        :icon="Delete"
                      >
                        删除分类
                      </el-button>
                    </div>
                  </div>
                  <p class="category-tip">我的分类仅自己可见，可自由管理</p>
                </el-tab-pane>
              </el-tabs>
            </el-form-item>

            <!-- 创建/编辑分类对话框 -->
            <el-dialog
              v-model="showCreateCategoryDialog"
              :title="isEditingCategory ? '编辑分类' : '创建新分类'"
              width="400px"
              class="dark-dialog"
              :append-to-body="true"
              @close="resetCategoryForm"
            >
              <el-form :model="categoryForm" label-width="80px">
                <el-form-item label="分类名称" required>
                  <el-input
                    v-model="categoryForm.name"
                    placeholder="请输入分类名称"
                    maxlength="50"
                    show-word-limit
                    class="dark-input"
                  />
                </el-form-item>
                <el-form-item label="颜色标识">
                  <el-color-picker v-model="categoryForm.color" />
                </el-form-item>
              </el-form>
              <template #footer>
                <el-button @click="showCreateCategoryDialog = false">取消</el-button>
                <el-button type="primary" @click="submitCategoryForm" :loading="savingCategory">
                  {{ isEditingCategory ? '保存' : '创建' }}
                </el-button>
              </template>
            </el-dialog>

            <el-form-item label="文章状态">
              <el-radio-group v-model="form.status">
                <el-radio :value="1">已发布</el-radio>
                <el-radio :value="0">草稿</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-form>
        </div>

      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, onBeforeUnmount, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Edit, Upload, View, Clock, Picture, Close, Hide, Plus, Delete } from '@element-plus/icons-vue'
import { 
  getArticle, 
  createArticle, 
  updateArticle, 
  uploadCoverImage, 
  getCategories, 
  getMyCategories,
  createUserCategory, 
  updateUserCategory, 
  deleteUserCategory 
} from '@/api/auth'
import TopNavBar from '@/components/TopNavBar.vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const router = useRouter()
const route = useRoute()

// 判断是否是编辑模式
const isEditMode = computed(() => !!route.params.id)

// 预览控制 - 默认关闭预览面板
const showPreview = ref(false)

// 预览内容状态
const previewContent = ref('')

// Vditor 初始化状态标记
const vditorInitialized = ref(false)

// 当前日期
const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

// 自定义 Markdown 渲染函数（零依赖、高稳定性）
const renderMarkdown = (markdown) => {
  if (!markdown || markdown.trim() === '') {
    return '<div class="article-body">暂无正文内容</div>'
  }
  
  try {
    // 1. 提取代码块和行内代码，用占位符替换
    const codeBlocks = []
    const inlineCodes = []
    
    // 提取多行代码块 ```code```
    let processed = markdown.replace(/```(\w*)\n([\s\S]*?)```/g, (match, lang, code) => {
      const index = codeBlocks.length
      codeBlocks.push({ lang, code: code.trim() })
      return `\n___CODE_BLOCK_${index}___\n`
    })
    
    // 提取单行代码块 `code`
    processed = processed.replace(/`([^`]+)`/g, (match, code) => {
      const index = inlineCodes.length
      inlineCodes.push(code)
      return `___INLINE_CODE_${index}___`
    })
    
    // 2. 分段处理普通文本
    const lines = processed.split('\n')
    const htmlLines = []
    let inParagraph = false
    let inList = false
    let listType = null
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i]
      const trimmedLine = line.trim()
      
      // 跳过空行
      if (trimmedLine === '') {
        if (inParagraph) {
          htmlLines.push('</p>')
          inParagraph = false
        }
        if (inList) {
          htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
          inList = false
          listType = null
        }
        continue
      }
      
      // 检查是否是特殊行（代码块占位符、标题、列表等）
      if (trimmedLine.startsWith('___CODE_BLOCK_')) {
        if (inParagraph) {
          htmlLines.push('</p>')
          inParagraph = false
        }
        if (inList) {
          htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
          inList = false
          listType = null
        }
        const match = trimmedLine.match(/___CODE_BLOCK_(\d+)___/)
        if (match) {
          const index = parseInt(match[1])
          const block = codeBlocks[index]
          // 对代码内容进行 HTML 转义
          const escapedCode = escapeHtml(block.code)
          htmlLines.push(`<pre class="highlight-chroma"><code class="language-${block.lang || 'text'}">${escapedCode}</code></pre>`)
        }
        continue
      }
      
      // 处理标题
      if (trimmedLine.startsWith('#')) {
        if (inParagraph) {
          htmlLines.push('</p>')
          inParagraph = false
        }
        if (inList) {
          htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
          inList = false
          listType = null
        }
        const level = trimmedLine.match(/^#+/)[0].length
        const content = trimmedLine.replace(/^#+\s*/, '')
        htmlLines.push(`<h${level}>${content}</h${level}>`)
        continue
      }
      
      // 处理无序列表
      if (trimmedLine.startsWith('- ') || trimmedLine.startsWith('* ')) {
        if (inParagraph) {
          htmlLines.push('</p>')
          inParagraph = false
        }
        if (!inList || listType !== 'ul') {
          if (inList) {
            htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
          }
          htmlLines.push('<ul>')
          inList = true
          listType = 'ul'
        }
        const content = trimmedLine.replace(/^[-*]\s*/, '')
        htmlLines.push(`<li>${processInlineStyles(content)}</li>`)
        continue
      }
      
      // 处理有序列表
      if (/^\d+\.\s/.test(trimmedLine)) {
        if (inParagraph) {
          htmlLines.push('</p>')
          inParagraph = false
        }
        if (!inList || listType !== 'ol') {
          if (inList) {
            htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
          }
          htmlLines.push('<ol>')
          inList = true
          listType = 'ol'
        }
        const content = trimmedLine.replace(/^\d+\.\s*/, '')
        htmlLines.push(`<li>${processInlineStyles(content)}</li>`)
        continue
      }
      
      // 处理引用
      if (trimmedLine.startsWith('>')) {
        if (inParagraph) {
          htmlLines.push('</p>')
          inParagraph = false
        }
        if (inList) {
          htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
          inList = false
          listType = null
        }
        const content = trimmedLine.replace(/^>\s*/, '')
        htmlLines.push(`<blockquote>${processInlineStyles(content)}</blockquote>`)
        continue
      }
      
      // 关闭列表状态
      if (inList) {
        htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
        inList = false
        listType = null
      }
      
      // 普通文本行
      if (!inParagraph) {
        htmlLines.push('<p>')
        inParagraph = true
      }
      
      // 处理行内格式
      const formattedLine = processInlineStyles(trimmedLine)
      htmlLines.push(formattedLine)
    }
    
    // 关闭未结束的标签
    if (inParagraph) {
      htmlLines.push('</p>')
    }
    if (inList) {
      htmlLines.push(listType === 'ul' ? '</ul>' : '</ol>')
    }
    
    // 3. 合并所有 HTML
    const html = htmlLines.join('\n')
    
    // 如果没有生成任何内容，返回原始包装
    if (!html || html.trim() === '') {
      return `<div class="article-body">${escapeHtml(markdown).replace(/\n/g, '<br>')}</div>`
    }
    
    return html
  } catch (error) {
    console.error('Markdown 渲染错误:', error)
    // 降级方案：简单包装
    return `<div class="article-body">${escapeHtml(markdown).replace(/\n/g, '<br>')}</div>`
  }
}

// 处理行内样式（粗体、斜体、链接、代码等）
const processInlineStyles = (text) => {
  let formatted = escapeHtml(text)
  
  // 粗体 **text**
  formatted = formatted.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  
  // 斜体 *text*
  formatted = formatted.replace(/\*([^*]+)\*/g, '<em>$1</em>')
  
  // 删除线 ~~text~~
  formatted = formatted.replace(/~~([^~]+)~~/g, '<del>$1</del>')
  
  // 链接 [text](url)
  formatted = formatted.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
  
  // 恢复行内代码占位符（如果有）
  formatted = formatted.replace(/___INLINE_CODE_(\d+)___/g, (match, index) => {
    // 这里需要从外部传入 inlineCodes，简化处理暂不支持
    return `<code class="highlight-chroma">${match}</code>`
  })
  
  return formatted
}

// HTML 转义函数（防止 XSS 和标签误解析）
const escapeHtml = (str) => {
  if (!str) return ''
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

// 更新预览内容的函数
const updatePreview = () => {
  if (!vditor) {
    return
  }
  
  try {
    // 获取 Markdown 内容 - 在 wysiwyg 模式下使用 form.content
    let markdown = ''
    
    // 优先使用 form.content（从 change 回调或 input 事件同步）
    if (form.content && form.content.trim() !== '') {
      markdown = form.content
    } 
    // 降级方案：尝试使用 getMarkdown 方法
    else if (typeof vditor.getMarkdown === 'function') {
      markdown = vditor.getMarkdown() || ''
    }
    
    if (!markdown || markdown.trim() === '') {
      previewContent.value = ''
      return
    }
    
    // 使用自定义的 Markdown 转 HTML 方法（更稳定可靠）
    const html = renderMarkdown(markdown)
    
    previewContent.value = html
  } catch (error) {
    console.error('预览渲染失败:', error)
    previewContent.value = ''
  }
}

// 表单数据
const form = reactive({
  title: '',
  summary: '',
  cover_image: '',
  content: '',
  status: 1,
  category_id: undefined, // 添加分类 ID 字段
  user_category_id: undefined // 添加用户分类 ID 字段
})

const saving = ref(false)
const articleId = ref(null)
const coverFile = ref(null)
const coverPreviewUrl = ref('')
const coverUploadRef = ref(null)

// 分类列表
const categories = ref([])

// 用户分类列表
const userCategories = ref([])

// 分类相关状态
const showCreateCategoryDialog = ref(false)
const isEditingCategory = ref(false)
const selectedUserCategoryId = ref(null)

// 分类表单数据
const categoryForm = reactive({
  name: '',
  color: '#409EFF'
})

// 分类保存状态
const savingCategory = ref(false)

// 重置分类表单
const resetCategoryForm = () => {
  categoryForm.name = ''
  categoryForm.color = '#409EFF'
  isEditingCategory.value = false
}

// 显示创建分类对话框
const showCreateCategory = () => {
  isEditingCategory.value = false // 确保是创建模式
  showCreateCategoryDialog.value = true
}

// 编辑器相关
const vditorContainer = ref(null)
let vditor = null

// 初始化 Vditor 编辑器
const initVditor = () => {
  if (!vditorContainer.value) {
    return
  }
  
  try {
    vditor = new Vditor(vditorContainer.value, {
      // 必需配置：CDN 路径，指向 Vditor 静态资源
      cdn: 'https://cdn.jsdelivr.net/npm/vditor@3.10.9',
      height: '100%',
      placeholder: '请输入文章内容（支持 markdown 语法 和可视化编辑）',
      value: form.content,
      theme: 'dark',
      icon: 'material',
      mode: 'wysiwyg', // 所见即所得模式，更稳定
      cache: {
        enable: false
      },
      preview: {
        delay: 1000
      },
      toolbar: [
        'emoji', 'headings', 'bold', 'italic', 'strike', 'link', '|',
        'list', 'ordered-list', 'check', 'outdent', 'indent', '|',
        'quote', 'line', 'code', 'inline-code', 'insert-before', 'insert-after', '|',
        'upload', 'table', '|',
        'undo', 'redo', '|',
        'edit-mode', 'fullscreen', 'download'
      ],
      // 添加空的 customWysiwygToolbar 函数避免内部调用错误
      customWysiwygToolbar: () => {},
      upload: {
        accept: 'image/*,.mp3,.wav,.rar,.gif,.pdf',
        url: '/api/upload/image',
        max: 5 * 1024 * 1024,
        filename(name) {
          return name
            .replace(/[^(a-zA-Z0-9\u4e00-\u9fa5\.)]/g, '')
            .replace(/[\?\\/:|<>\*\[\]\(\)\$%\{\}@~]/g, '')
            .replace(/\s/g, '')
        }
      },
      // 添加错误处理
      error: (msg) => {
        console.error('Vditor 初始化错误:', msg)
        ElMessage.error('编辑器初始化失败：' + msg)
      },
      // 初始化成功回调
      after: () => {
        vditorInitialized.value = true
        
        // 监听输入事件来实时获取编辑器内容（使用防抖避免频繁触发）
        if (vditor.vditor?.wysiwyg?.element) {
          let debounceTimer = null
          const updateContent = () => {
            // 清除上一个定时器
            if (debounceTimer) {
              clearTimeout(debounceTimer)
            }
            
            // 防抖 300ms
            debounceTimer = setTimeout(() => {
              form.content = vditor.getValue()
              // 同步后更新预览
              updatePreview()
            }, 300)
          }
          
          // 监听 input 事件（键盘输入）
          vditor.vditor.wysiwyg.element.addEventListener('input', updateContent)
        }
        
        // 如果是编辑模式且有待加载的内容，在初始化完成后设置
        if (isEditMode.value && form.content) {
          nextTick(() => {
            if (vditor) {
              vditor.setValue(form.content)
              // 等待一下再更新预览，避免访问未初始化的内部属性
              setTimeout(() => {
                updatePreview()
              }, 100)
            }
          })
        } else {
          // 创建模式：初始化完成后立即更新预览
          if (form.content && form.content.trim() !== '') {
            updatePreview()
          }
        }
      },
      // 内容变化回调
      change: (value) => {
        form.content = value
        // 自适应高度：当内容变化时，让编辑器容器自动调整高度
        if (vditorContainer.value) {
          vditorContainer.value.style.height = 'auto'
        }
        // 更新预览内容
        updatePreview()
      }
    })
    
  } catch (error) {
    console.error('Vditor 初始化异常:', error)
    ElMessage.error('编辑器初始化失败，请刷新页面重试')
  }
}

// 设置编辑器内容
const setEditorContent = (content) => {
  if (!vditor) {
    return
  }
  
  if (!vditorInitialized.value) {
    // 等待初始化完成后设置
    const checkInit = setInterval(() => {
      if (vditorInitialized.value) {
        clearInterval(checkInit)
        vditor.setValue(content)
      }
    }, 100)
    
    // 5 秒超时
    setTimeout(() => {
      clearInterval(checkInit)
    }, 5000)
    return
  }
  
  vditor.setValue(content)
}

// 获取编辑器内容
const getEditorContent = () => {
  if (vditor) {
    return vditor.getValue()
  }
  return ''
}

// 加载文章详情（编辑模式）
const loadArticle = async () => {
  if (!isEditMode.value) return
  
  try {
    const res = await getArticle(route.params.id)
    
    if (res.code === 0 && res.data.article) {
      const article = res.data.article
      
      form.title = article.title
      form.summary = article.summary || ''
      form.cover_image = article.cover_image || ''
      form.content = article.content
      form.status = article.status
      form.category_id = article.category?.id || undefined // 加载文章分类
      articleId.value = article.id
      
      // 如果有封面图片，设置预览
      if (article.cover_image) {
        coverPreviewUrl.value = article.cover_image
      }
      
      // 如果 Vditor 已初始化，直接设置内容
      // 否则在 initVditor 的 after 回调中处理
      if (vditorInitialized.value && vditor) {
        nextTick(() => {
          vditor.setValue(article.content)
          updatePreview()
        })
      }
    }
  } catch (error) {
    ElMessage.error('加载文章失败')
  }
}

// 处理封面图片选择
const handleCoverChange = (file) => {
  const rawFile = file.raw
  
  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
  if (!allowedTypes.includes(rawFile.type)) {
    ElMessage.error('只支持 jpg/png/gif/webp 格式的图片')
    coverFile.value = null
    return
  }
  
  // 验证文件大小（最大 5MB）
  const maxSize = 5 * 1024 * 1024
  if (rawFile.size > maxSize) {
    ElMessage.error('图片大小不能超过 5MB')
    coverFile.value = null
    return
  }
  
  coverFile.value = rawFile
  
  // 创建预览 URL
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      coverPreviewUrl.value = e.target.result
    } catch (error) {
      ElMessage.error('图片预览生成失败')
      coverPreviewUrl.value = ''
      coverFile.value = null
    }
  }
  reader.onerror = () => {
    ElMessage.error('图片读取失败，请重试')
    coverPreviewUrl.value = ''
    coverFile.value = null
  }
  reader.readAsDataURL(rawFile)
}

// 处理封面图片移除
const handleCoverRemove = () => {
  coverFile.value = null
  coverPreviewUrl.value = ''
  form.cover_image = ''
  
  // 清除 el-upload 组件内部的文件列表状态
  if (coverUploadRef.value) {
    coverUploadRef.value.clearFiles()
  }
}

// 处理图片加载错误
const handleImageError = (e) => {
  const img = e.target
  
  // 隐藏错误的图片
  img.style.display = 'none'
  
  // 判断是本地预览还是网络图片
  const isLocalPreview = img.src.startsWith('data:image')
  
  // 显示更详细的错误提示
  if (isLocalPreview) {
    ElMessage.warning({
      message: '图片预览失败，请检查：1) 图片格式是否正确（支持 jpg/png/gif/webp）2) 图片是否损坏',
      duration: 5000
    })
  } else {
    ElMessage.warning({
      message: '封面图片加载失败，可能原因：1) 后端服务未启动 2) 图片文件不存在 3) 图片链接已失效',
      duration: 5000
    })
  }
}

// 处理图片加载成功
const handleImageLoad = (e) => {
  // 图片加载成功，无需额外处理
}

// 上传图片到服务器（返回 URL）
const uploadToServer = async (file) => {
  try {
    const response = await uploadCoverImage(file)
    if (response.code === 0 && response.data?.url) {
      return response.data.url
    } else {
      throw new Error(response.message || '上传失败')
    }
  } catch (error) {
    ElMessage.error('封面图片上传失败：' + (error.message || '未知错误'))
    throw error
  }
}

// 提交表单
const handleSubmit = async () => {
  // 验证必填字段
  if (!form.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return
  }
  
  if (!form.summary.trim()) {
    ElMessage.warning('请输入文章摘要')
    return
  }
  
  if (!form.content.trim()) {
    ElMessage.warning('请输入文章内容')
    return
  }
  
  saving.value = true
  
  try {
    // 如果有上传封面图片，先上传图片
    if (coverFile.value) {
      const imageUrl = await uploadToServer(coverFile.value)
      form.cover_image = imageUrl
    }
    
    const submitData = {
      title: form.title,
      summary: form.summary,
      cover_image: form.cover_image,
      content: form.content,
      status: form.status,
      category_id: form.category_id || null, // 系统分类 ID
      user_category_id: form.user_category_id || null, // 用户自定义分类 ID
      tag_ids: [] // 标签，后续可以添加标签选择
    }
    
    if (isEditMode.value) {
      // 更新文章
      await updateArticle(articleId.value, submitData)
      ElMessage.success('文章已更新')
    } else {
      // 创建文章
      const res = await createArticle(submitData)
      ElMessage.success('文章创建成功')
      articleId.value = res.data?.id
    }
    
    // 跳转到创作中心
    setTimeout(() => {
      router.push('/creator')
    }, 500)
  } catch (error) {
    // 获取具体错误信息
    let errorMessage = '保存文章失败'
    
    if (error.response) {
      // 后端返回了错误响应
      errorMessage = error.response.data?.message || 
                    error.response.data?.msg || 
                    error.response.data?.error ||
                    '请求失败'
    } else if (error.request) {
      // 请求已发送但没有收到响应
      errorMessage = '无法连接到后端服务，请检查后端是否启动'
    } else {
      // 其他错误
      errorMessage = error.message || '保存文章失败'
    }
    
    ElMessage.error(errorMessage)
  } finally {
    saving.value = false
  }
}

// 取消操作
const handleCancel = () => {
  router.back()
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await getCategories()
    if (res.code === 0 && res.data.categories) {
      categories.value = res.data.categories
    }
  } catch (error) {
    // 不显示错误提示，避免影响用户体验
  }
}

// 加载用户分类列表
const loadUserCategories = async () => {
  try {
    const res = await getMyCategories()
    if (res.code === 0 && res.data.categories) {
      userCategories.value = res.data.categories
    }
  } catch (error) {
    // 不显示错误提示，避免影响用户体验
  }
}

// 处理系统分类选择
const handleSystemCategoryChange = (value) => {
  form.user_category_id = undefined // 清除用户分类选择
}

// 处理用户分类选择
const handleUserCategoryChange = (value) => {
  selectedUserCategoryId.value = value
  form.category_id = undefined // 清除系统分类选择
}

// 创建/编辑分类表单提交
const submitCategoryForm = async () => {
  if (!categoryForm.name.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }

  savingCategory.value = true
  
  try {
    if (isEditingCategory.value && selectedUserCategoryId.value) {
      // 更新分类
      await updateUserCategory(selectedUserCategoryId.value, categoryForm)
      ElMessage.success('分类已更新')
    } else {
      // 创建分类
      await createUserCategory(categoryForm)
      ElMessage.success('分类创建成功')
    }
    
    showCreateCategoryDialog.value = false
    await loadUserCategories() // 重新加载分类列表
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    savingCategory.value = false
  }
}

// 编辑当前选中的分类
const editCurrentCategory = async () => {
  if (!selectedUserCategoryId.value) {
    ElMessage.warning('请先选择一个分类')
    return
  }
  
  const category = userCategories.value.find(c => c.id === selectedUserCategoryId.value)
  if (category) {
    categoryForm.name = category.name
    categoryForm.color = category.color || '#409EFF'
    isEditingCategory.value = true
    showCreateCategoryDialog.value = true
  }
}

// 删除当前选中的分类
const deleteCurrentCategory = async () => {
  if (!selectedUserCategoryId.value) {
    ElMessage.warning('请先选择一个分类')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      '确定要删除这个分类吗？使用该分类的文章将被取消分类关联。',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await deleteUserCategory(selectedUserCategoryId.value)
    ElMessage.success('分类已删除')
    
    // 清空选择
    selectedUserCategoryId.value = null
    form.user_category_id = undefined
    
    // 重新加载分类列表
    await loadUserCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

// 初始化
onMounted(async () => {
  try {
    // 加载系统分类列表
    await loadCategories()
    
    // 加载用户分类列表
    await loadUserCategories()
    

    // 如果是编辑模式，加载文章详情
    if (isEditMode.value) {
      await loadArticle()
    }
    
    // 初始化编辑器
    initVditor()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '加载分类失败')
  }
})

// 销毁编辑器
onBeforeUnmount(() => {
  if (vditor && typeof vditor.destroy === 'function') {
    vditor.destroy()
    vditor = null
  }
})

</script>



<style lang="scss" scoped>
.article-editor-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e2a4a 0%, #2d1f3d 50%, #1a3a5c 100%);
  padding: 24px;
  
  .editor-container {
    max-width: 1000px;
    margin: 0 auto;
    
    // 页面头部
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;
      
      .page-title {
        font-size: 28px;
        font-weight: 700;
        color: rgba(255, 255, 255, 0.95);
        margin: 0;
        display: flex;
        align-items: center;
        gap: 12px;
        
        .el-icon {
          color: rgba(102, 126, 234, 0.8);
        }
      }
      
      .header-actions {
        display: flex;
        gap: 12px;
        
        .el-button {
          height: 40px;
          padding: 0 20px;
        }
      }
    }
    
    // 编辑器表单包裹器
    .editor-form-wrapper {
      position: relative;
      display: flex;
      gap: 24px;
      align-items: stretch;
      margin-bottom: 24px;
      
      // 编辑器面板
      .editor-panel {
        flex: 1;
        min-width: 0;
        background: rgba(255, 255, 255, 0.05);
        backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        padding: 30px;
        
        :deep(.el-form-item__label) {
          color: rgba(255, 255, 255, 0.9) !important;
          font-weight: 500;
          font-size: 15px;
        }
        
        // 深色输入框样式 - 强制覆盖 Element Plus 默认样式
        :deep(.el-input__wrapper),
        :deep(.el-textarea__wrapper) {
          background: rgba(255, 255, 255, 0.05) !important;
          background-color: rgba(255, 255, 255, 0.05) !important;
          border: 1px solid rgba(255, 255, 255, 0.2) !important;
          box-shadow: none !important;
          padding: 0 !important;
          
          // 覆盖 CSS 变量
          --el-input-bg-color: rgba(255, 255, 255, 0.05) !important;
          --el-fill-color-blank: rgba(255, 255, 255, 0.05) !important;
        }
        
        // 输入框内层统一样式
        :deep(.el-input__inner) {
          background: rgba(255, 255, 255, 0.05) !important;
          background-color: rgba(255, 255, 255, 0.05) !important;
          color: rgba(255, 255, 255, 0.95) !important;
          border: none !important;
          padding: 8px 12px !important;
        }
        
        // 额外强制覆盖 textarea 样式
        :deep(.el-textarea__inner) {
          background: rgba(255, 255, 255, 0.05) !important;
          background-color: rgba(255, 255, 255, 0.05) !important;
          color: rgba(255, 255, 255, 0.95) !important;
          border: none !important;
          padding: 8px 12px !important;
          --el-input-bg-color: rgba(255, 255, 255, 0.05) !important;
          --el-fill-color-blank: rgba(255, 255, 255, 0.05) !important;
        }
        
        // 输入框聚焦时的边框样式
        :deep(.el-input__wrapper.is-focus),
        :deep(.el-textarea__wrapper.is-focus) {
          border-color: rgba(255, 255, 255, 0.4) !important;
          box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.4) inset !important;
        }
        
        // 深色选择器样式
        :deep(.el-select .el-input__wrapper) {
          background: rgba(255, 255, 255, 0.05) !important;
          background-color: rgba(255, 255, 255, 0.05) !important;
          border: 1px solid rgba(255, 255, 255, 0.2) !important;
          box-shadow: none !important;
          --el-input-bg-color: rgba(255, 255, 255, 0.05) !important;
          --el-fill-color-blank: rgba(255, 255, 255, 0.05) !important;
        }
        
        :deep(.el-select .el-input__inner) {
          color: rgba(255, 255, 255, 0.95) !important;
        }
        
        :deep(.el-select .el-input__icon) {
          color: rgba(255, 255, 255, 0.7) !important;
        }
        
        // 选择器下拉菜单样式
        :deep(.el-select-dropdown) {
          background: rgba(30, 42, 74, 0.98) !important;
          border: 1px solid rgba(255, 255, 255, 0.1) !important;
          
          .el-select-dropdown__item {
            color: rgba(255, 255, 255, 0.9) !important;
            
            &:hover {
              background: rgba(102, 126, 234, 0.2) !important;
            }
            
            &.selected {
              color: rgba(102, 126, 234, 0.9) !important;
              background: rgba(102, 126, 234, 0.15) !important;
            }
          }
        }
        
        // 字数统计样式
        :deep(.el-input__count) {
          background: transparent !important;
          background-color: transparent !important;
          border: none !important;
          
          .el-input__count-inner {
            color: rgba(255, 255, 255, 0.7) !important;
            background: transparent !important;
          }
        }
        
        // 封面上传样式
        .cover-upload-area {
          width: 100%;
          overflow: hidden;
        }
        
        .cover-uploader {
          width: 100%;
          margin-bottom: 12px;
          
          :deep(.el-upload-dragger) {
            background: rgba(255, 255, 255, 0.05);
            border: 2px dashed rgba(255, 255, 255, 0.2);
            border-radius: 8px;
            padding: 40px 20px;
            transition: all 0.3s;
            
            &:hover {
              border-color: rgba(64, 158, 255, 0.5);
              background: rgba(64, 158, 255, 0.05);
            }
          }
          
          .el-upload__tip {
            color: rgba(255, 255, 255, 0.6);
            font-size: 13px;
            margin-top: 8px;
          }
          
          .cover-uploader-content {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 16px;
            
            .cover-uploader-icon {
              font-size: 48px;
              color: rgba(64, 158, 255, 0.8);
              margin-bottom: 8px;
            }
            
            .cover-uploader-text {
              text-align: center;
              
              .cover-uploader-title {
                color: rgba(255, 255, 255, 0.9);
                font-size: 16px;
                font-weight: 500;
                margin: 0 0 8px 0;
              }
              
              .cover-uploader-tip {
                color: rgba(255, 255, 255, 0.5);
                font-size: 13px;
                margin: 0;
              }
            }
          }
          
          .cover-preview-wrapper {
            position: relative;
            width: 100%;
            max-width: 100%;
            overflow: hidden;
            
            .cover-preview-image {
              width: 100%;
              max-height: 400px;
              object-fit: contain;
              border-radius: 8px;
              box-shadow: none;
            }
            
            .cover-overlay {
              position: absolute;
              top: 0;
              left: 0;
              right: 0;
              bottom: 0;
              background: rgba(0, 0, 0, 0.6);
              display: flex;
              align-items: center;
              justify-content: center;
              opacity: 0;
              transition: opacity 0.3s;
              border-radius: 8px;
              
              button {
                transform: translateY(10px);
                transition: transform 0.3s;
              }
            }
            
            &:hover .cover-overlay {
              opacity: 1;
              
              button {
                transform: translateY(0);
              }
            }
          }
        }
        
        // 内容编辑器包装器
        .content-editor-wrapper {
          position: relative;
          display: flex;
          width: 100%;
          gap: 0;
          
          // 编辑器区域
          .vditor-wrapper {
            flex: 1;
            min-width: 0;
            border: 1px solid rgba(255, 255, 255, 0.2);
            border-radius: 4px;
            overflow: hidden;
            background: rgba(255, 255, 255, 0.05);
            
            .vditor-container {
              min-height: 600px;
              width: 100%;
              
              // Vditor 编辑区域
              :deep(.vditor) {
                border: none !important;
                background: transparent !important;
                
                // 工具栏
                .vditor-toolbar {
                  background: rgba(255, 255, 255, 0.05) !important;
                  border-bottom: 1px solid rgba(255, 255, 255, 0.2) !important;
                }
                
                // 即时渲染编辑区域
                .vditor-ir {
                  min-height: 600px !important;
                  background: rgba(255, 255, 255, 0.05) !important;
                  color: rgba(255, 255, 255, 0.95) !important;
                }
                
                // 内容区域
                .vditor-reset {
                  min-height: 600px !important;
                  color: rgba(255, 255, 255, 0.95) !important;
                }
              }
            }
          }
          
          // 预览面板 - 使用 flex 布局并排显示
          .inline-preview-panel {
            flex: 0 0 380px;
            min-width: 380px;
            min-height: 600px;
            // 让预览面板高度与编辑器一致
            align-self: stretch;
            background: rgba(255, 255, 255, 0.08);
            backdrop-filter: blur(10px);
            border-left: 2px solid rgba(102, 126, 234, 0.5);
            border-radius: 0;
            overflow: hidden;
            display: flex;
            flex-direction: column;
            z-index: 5;
            box-shadow: -4px 0 12px rgba(0, 0, 0, 0.2);
            transition: all 0.3s ease;
            
            .preview-header {
              flex-shrink: 0;
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding: 12px 16px;
              border-bottom: 1px solid rgba(255, 255, 255, 0.1);
              background: rgba(255, 255, 255, 0.03);
                
                .preview-title {
                  margin: 0;
                  font-size: 16px;
                  font-weight: 600;
                  color: rgba(255, 255, 255, 0.95);
                  display: flex;
                  align-items: center;
                  gap: 8px;
                  
                  .el-icon {
                    color: rgba(102, 126, 234, 0.8);
                  }
                }
                
                .close-btn {
                  color: rgba(255, 255, 255, 0.6);
                  
                  &:hover {
                    color: rgba(255, 255, 255, 0.95);
                  }
                }
              }
              
              .preview-content {
                flex: 1;
                overflow-y: auto;
                overflow-x: hidden;
                padding: 20px;
                background: rgba(255, 255, 255, 0.02);
                scrollbar-width: thin;
                scrollbar-color: rgba(102, 126, 234, 0.5) rgba(255, 255, 255, 0.05);
                min-height: 0;
                scroll-behavior: smooth;
                -webkit-overflow-scrolling: touch;
                
                &::-webkit-scrollbar {
                  width: 8px;
                }
                
                &::-webkit-scrollbar-track {
                  background: rgba(255, 255, 255, 0.05);
                  border-radius: 4px;
                }
                
                &::-webkit-scrollbar-thumb {
                  background: linear-gradient(180deg, rgba(102, 126, 234, 0.6), rgba(102, 126, 234, 0.4));
                  border-radius: 4px;
                  transition: background 0.3s ease;
                  
                  &:hover {
                    background: linear-gradient(180deg, rgba(102, 126, 234, 0.8), rgba(102, 126, 234, 0.6));
                  }
                }
                
                .preview-article {
                  background: rgba(255, 255, 255, 0.05);
                  border-radius: 12px;
                  padding: 24px;
                  border: 1px solid rgba(255, 255, 255, 0.1);
                  margin-bottom: 20px;
                  
                  .article-header {
                    margin-bottom: 24px;
                    padding-bottom: 16px;
                    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
                    
                    .article-title {
                      font-size: 22px;
                      font-weight: 700;
                      color: rgba(255, 255, 255, 0.95);
                      margin: 0 0 12px 0;
                      line-height: 1.4;
                    }
                    
                    .article-meta {
                      display: flex;
                      gap: 20px;
                      color: rgba(255, 255, 255, 0.6);
                      font-size: 13px;
                      margin-top: 12px;
                      
                      .meta-item {
                        display: flex;
                        align-items: center;
                        gap: 6px;
                        
                        .el-icon {
                          font-size: 15px;
                        }
                      }
                    }
                  }
                  
                  .article-cover {
                    margin: 20px 0;
                    border-radius: 8px;
                    overflow: hidden;
                    background: transparent;
                    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
                    
                    img {
                      width: 100%;
                      max-height: 300px;
                      object-fit: cover;
                      border-radius: 8px;
                      display: block;
                      transition: transform 0.3s ease;
                      
                      &:hover {
                        transform: scale(1.02);
                      }
                    }
                  }
                  
                  .article-summary {
                    margin: 20px 0;
                    padding: 16px 20px;
                    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(102, 126, 234, 0.05));
                    border-left: 3px solid rgba(102, 126, 234, 0.8);
                    border-radius: 6px;
                    
                    h3 {
                      margin: 0 0 10px 0;
                      font-size: 16px;
                      color: rgba(255, 255, 255, 0.9);
                      font-weight: 600;
                    }
                    
                    p {
                      margin: 0;
                      color: rgba(255, 255, 255, 0.75);
                      font-size: 14px;
                      line-height: 1.8;
                    }
                  }
                  
                  .article-body {
                    color: rgba(255, 255, 255, 0.9);
                    line-height: 1.8;
                    font-size: 15px;
                    padding: 16px 20px;
                    
                    // Vditor wysiwyg 模式输出的样式支持 - 完整复制编辑器样式
                    :deep(.vditor-reset), :deep(.vditor-wysiwyg) {
                      color: rgba(255, 255, 255, 0.95) !important;
                      line-height: 1.8 !important;
                      font-size: 15px !important;
                      background: transparent !important;
                      
                      // 标题样式
                      h1, h2, h3, h4, h5, h6 {
                        color: rgba(255, 255, 255, 0.95) !important;
                        margin-top: 32px !important;
                        margin-bottom: 16px !important;
                        font-weight: 600 !important;
                        line-height: 1.4 !important;
                      }
                      
                      h1 { 
                        font-size: 26px !important;
                        padding-bottom: 12px !important;
                        border-bottom: 2px solid rgba(102, 126, 234, 0.3) !important;
                      }
                      h2 { 
                        font-size: 22px !important;
                        padding-bottom: 8px !important;
                        border-bottom: 1px solid rgba(102, 126, 234, 0.2) !important;
                      }
                      h3 { font-size: 19px !important; }
                      h4 { font-size: 17px !important; }
                      
                      // 段落样式
                      p {
                        margin-bottom: 20px !important;
                        line-height: 1.9 !important;
                        text-align: justify !important;
                        color: rgba(255, 255, 255, 0.95) !important;
                      }
                      
                      // 列表样式
                      ul, ol {
                        margin-bottom: 20px !important;
                        padding-left: 24px !important;
                        
                        li {
                          margin-bottom: 10px !important;
                          line-height: 1.8 !important;
                          color: rgba(255, 255, 255, 0.95) !important;
                        }
                      }
                      
                      // 引用样式
                      blockquote {
                        margin: 24px 0 !important;
                        padding: 16px 20px !important;
                        border-left: 4px solid rgba(102, 126, 234, 0.8) !important;
                        background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(102, 126, 234, 0.05)) !important;
                        border-radius: 6px !important;
                        color: rgba(255, 255, 255, 0.85) !important;
                        font-style: italic !important;
                      }
                      
                      // 代码块样式
                      pre {
                        margin: 24px 0 !important;
                        padding: 20px !important;
                        background: rgba(0, 0, 0, 0.4) !important;
                        border-radius: 8px !important;
                        overflow-x: auto !important;
                        border: 1px solid rgba(255, 255, 255, 0.1) !important;
                        
                        code {
                          padding: 0 !important;
                          background: transparent !important;
                          color: rgba(255, 255, 255, 0.95) !important;
                          font-family: 'Consolas', 'Monaco', 'Courier New', monospace !important;
                          font-size: 14px !important;
                          line-height: 1.6 !important;
                        }
                      }
                      
                      // 行内代码样式
                      code {
                        padding: 4px 10px !important;
                        background: rgba(102, 126, 234, 0.2) !important;
                        border-radius: 4px !important;
                        color: rgba(255, 255, 255, 0.95) !important;
                        font-family: 'Consolas', 'Monaco', 'Courier New', monospace !important;
                        font-size: 14px !important;
                      }
                      
                      // 链接样式
                      a {
                        color: #66b1ff !important;
                        text-decoration: none !important;
                        transition: all 0.3s !important;
                        border-bottom: 1px dashed rgba(102, 177, 255, 0.5) !important;
                        
                        &:hover {
                          color: #90caff !important;
                          border-bottom-style: solid !important;
                        }
                      }
                      
                      // 图片样式
                      img {
                        max-width: 100% !important;
                        border-radius: 8px !important;
                        margin: 20px 0 !important;
                        display: block !important;
                        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2) !important;
                        transition: transform 0.3s ease !important;
                        
                        &:hover {
                          transform: scale(1.01) !important;
                        }
                      }
                      
                      // 表格样式
                      table {
                        width: 100% !important;
                        border-collapse: collapse !important;
                        margin: 20px 0 !important;
                        border-radius: 8px !important;
                        overflow: hidden !important;
                        
                        th, td {
                          border: 1px solid rgba(255, 255, 255, 0.15) !important;
                          padding: 12px 16px !important;
                          text-align: left !important;
                        }
                        
                        th {
                          background: linear-gradient(135deg, rgba(102, 126, 234, 0.3), rgba(102, 126, 234, 0.1)) !important;
                          color: rgba(255, 255, 255, 0.95) !important;
                          font-weight: 600 !important;
                        }
                        
                        td {
                          color: rgba(255, 255, 255, 0.85) !important;
                        }
                        
                        tr:nth-child(even) {
                          background: rgba(255, 255, 255, 0.03) !important;
                        }
                        
                        tr:hover {
                          background: rgba(102, 126, 234, 0.1) !important;
                        }
                      }
                      
                      // 粗体和斜体
                      strong, b {
                        font-weight: 700 !important;
                        color: rgba(255, 255, 255, 0.95) !important;
                      }
                      
                      em, i {
                        font-style: italic !important;
                        color: rgba(255, 255, 255, 0.95) !important;
                      }
                    }
                  }
                }
              }
          }
          
          // 预览切换按钮 - 固定在编辑器右上角
          .preview-toggle-btn {
            position: absolute;
            top: 12px;
            right: 12px;
            z-index: 100;
            width: 36px;
            height: 36px;
            border: none;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
            transition: all 0.3s ease;
            background: rgba(102, 126, 234, 0.8) !important;
            color: rgba(255, 255, 255, 0.95) !important;
            
            &:hover {
              transform: scale(1.1);
              box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
              background: rgba(102, 126, 234, 1) !important;
            }
            
            .el-icon {
              font-size: 18px;
            }
          }
          
          // 内容编辑器包装器整体样式
          &.content-editor-wrapper {
            position: relative;
          }
          
          // 当预览面板显示时，调整按钮位置到预览面板右侧
          &:has(.inline-preview-panel[style*="display: flex"]) .preview-toggle-btn,
          .preview-toggle-btn {
            top: 12px;
            right: 12px;
          }
          
          // 预览面板隐藏时，按钮移回编辑器右上角
          .inline-preview-panel[style*="display: none"] + .preview-toggle-btn {
            right: 12px;
          }
          
          // Vditor 深色主题覆盖
          :deep(.vditor) {
            border: none !important;
            background: rgba(255, 255, 255, 0.05) !important;
            width: 100% !important;
            height: 100% !important;
            display: flex !important;
            flex-direction: column !important;
            
            // 即时渲染编辑区域占满整个容器
            .vditor-ir {
              width: 100% !important;
              height: 100% !important;
              min-height: 700px !important;
              background: rgba(255, 255, 255, 0.05) !important;
              color: rgba(255, 255, 255, 0.95) !important;
              flex: 1 !important;
            }
            
            // 工具栏
            .vditor-toolbar {
              background: rgba(255, 255, 255, 0.05) !important;
              border-bottom: 1px solid rgba(255, 255, 255, 0.2) !important;
              width: 100% !important;
              
              .vditor-toolbar__item {
                span {
                  color: rgba(200, 200, 200, 0.8) !important;
                  fill: rgba(200, 200, 200, 0.8) !important;
                }
                
                &:hover span {
                  color: rgba(255, 255, 255, 0.95) !important;
                  fill: rgba(255, 255, 255, 0.95) !important;
                }
                
                &.vditor-toolbar__item--current span {
                  color: rgba(255, 255, 255, 0.95) !important;
                  fill: rgba(255, 255, 255, 0.95) !important;
                }
              }
            }
            
            // 编辑区域内容
            .vditor-reset {
              color: rgba(255, 255, 255, 0.95) !important;
              min-height: 600px !important;
              width: 100% !important;
              
              h1, h2, h3, h4, h5, h6 {
                color: rgba(255, 255, 255, 0.95) !important;
              }
              
              p, li, td, th, div, span {
                color: rgba(255, 255, 255, 0.95) !important;
              }
              
              strong, b {
                color: rgba(255, 255, 255, 0.95) !important;
              }
              
              em, i {
                color: rgba(255, 255, 255, 0.95) !important;
              }
              
              code, pre {
                background: rgba(255, 255, 255, 0.1) !important;
                color: rgba(255, 255, 255, 0.95) !important;
              }
              
              blockquote {
                background: rgba(255, 255, 255, 0.05) !important;
                color: rgba(255, 255, 255, 0.9) !important;
                border-left-color: rgba(255, 255, 255, 0.3) !important;
              }
              
              a {
                color: #409EFF !important;
              }
            }
            
            // 对话框深色主题样式
            .dark-dialog {
              // 确保对话框可见
              z-index: 9999 !important;
              
              .el-dialog__header {
                background: rgba(255, 255, 255, 0.05) !important;
                border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
              }
              
              .el-dialog__title {
                color: rgba(255, 255, 255, 0.9) !important;
              }
              
              .el-dialog__body {
                background: rgba(255, 255, 255, 0.05) !important;
                color: rgba(255, 255, 255, 0.9) !important;
              }
              
              .el-dialog__footer {
                background: rgba(255, 255, 255, 0.05) !important;
                border-top: 1px solid rgba(255, 255, 255, 0.1) !important;
              }
            }
            
            // 滚动条
            ::-webkit-scrollbar {
              background: rgba(255, 255, 255, 0.05) !important;
            }
            
            ::-webkit-scrollbar-thumb {
              background: rgba(255, 255, 255, 0.2) !important;
              
              &:hover {
                background: rgba(255, 255, 255, 0.3) !important;
              }
            }
          }
        }
      }
      
      // 分类选择器样式
      .dark-select {
        .el-input__wrapper {
          background-color: rgba(255, 255, 255, 0.05) !important;
          box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.2) inset !important;
        }
        
        .el-input__inner {
          color: rgba(255, 255, 255, 0.95) !important;
        }
        
        .el-select__selected-item {
          color: rgba(255, 255, 255, 0.95) !important;
        }
        
        .el-select__caret {
          color: rgba(255, 255, 255, 0.6) !important;
        }
        
        &:hover {
          .el-input__wrapper {
            box-shadow: 0 0 0 1px rgba(102, 126, 234, 0.5) inset !important;
          }
        }
        
        &.is-focus {
          .el-input__wrapper {
            box-shadow: 0 0 0 1px rgba(102, 126, 234, 1) inset !important;
          }
        }
      }
      
      :deep(.el-select-dropdown) {
        background-color: rgba(30, 30, 40, 0.98) !important;
        border: 1px solid rgba(255, 255, 255, 0.2) !important;
        
        .el-select-dropdown__item {
          color: rgba(255, 255, 255, 0.95) !important;
          
          &.is-hovering {
            background-color: rgba(102, 126, 234, 0.3) !important;
          }
          
          &.is-selected {
            color: rgba(102, 126, 234, 1) !important;
            font-weight: 600;
          }
        }
      }
      
      // 分类标签样式
      .category-tabs {
        // 覆盖 border-card 类型的白色背景
        background: transparent !important;
        border: none !important;
        
        :deep(.el-tabs__header) {
          background: rgba(255, 255, 255, 0.05) !important;
          border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
        }
        
        :deep(.el-tabs__nav) {
          border: none !important;
        }
        
        :deep(.el-tabs__nav-wrap) {
          padding: 0 12px !important;
        }
        
        :deep(.el-tabs__item) {
          color: rgba(255, 255, 255, 0.9) !important;
          padding: 12px 24px !important;
          border: none !important;
          border-bottom: 2px solid transparent !important;
          transition: all 0.3s !important;
          
          &:hover {
            color: rgba(102, 126, 234, 1) !important;
          }
          
          &.is-active {
            color: rgba(102, 126, 234, 1) !important;
            border-bottom-color: rgba(102, 126, 234, 1) !important;
          }
        }
        
        :deep(.el-tabs__content) {
          background: rgba(255, 255, 255, 0.05) !important;
          border: 1px solid rgba(255, 255, 255, 0.1) !important;
          border-top: none !important;
          padding: 24px !important;
          border-radius: 8px !important;
        }
      }
      
      .category-tip {
        color: rgba(255, 255, 255, 0.6) !important;
        font-size: 13px !important;
        margin-top: 12px !important;
      }
      
      .my-category-wrapper {
        display: flex;
        flex-direction: column;
        gap: 12px;
      }
      
      .category-color-dot {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        margin-left: 8px;
      }
      
      .category-actions {
        display: flex;
        gap: 12px;
        margin-top: 12px;
      }
    }
  }
}
</style>
