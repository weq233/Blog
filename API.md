# 📡 博客系统 API 接口速查手册

> 所有接口都需要在请求头中添加 JWT Token（除特别说明外）
> 
> **基础 URL**: `http://localhost:8080/api`

---

## 🔐 认证接口

### 1. 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "string (必填，3-20 字符)",
  "email": "string (必填，邮箱格式)",
  "password": "string (必填，6-20 字符)",
  "captcha_id": "string (必填，从获取验证码接口获取)",
  "captcha_ans": "string (必填，验证码答案)"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "nickname": "John",
      "avatar": null,
      "role": 1
    }
  }
}
```

---

### 2. 用户登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "string (必填)",
  "password": "string (必填)",
  "captcha_id": "string (必填)",
  "captcha_ans": "string (必填)"
}
```

**响应示例**: 同注册接口

---

### 3. 获取验证码

```http
GET /api/auth/captcha
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "captcha_id": "xxx-xxx-xxx",
    "captcha_image": "data:image/png;base64,iVBORw0KG..."
  }
}
```

---

### 4. 退出登录

```http
POST /api/auth/logout
Authorization: Bearer {token}
```

---

### 5. 获取当前用户信息

```http
GET /api/auth/me
Authorization: Bearer {token}
```

---

### 6. 更新用户资料

```http
PUT /api/auth/profile
Content-Type: application/json
Authorization: Bearer {token}

{
  "nickname": "新昵称",
  "avatar": "/uploads/avatars/xxx.jpg"
}
```

---

### 7. 修改密码

```http
PUT /api/auth/password
Content-Type: application/json
Authorization: Bearer {token}

{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

---

### 8. 找回密码（邮件）

```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

---

### 9. 重置密码（带令牌）

```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "重置令牌（从邮件中获取）",
  "new_password": "新密码"
}
```

---

## 📝 文章接口

### 1. 获取文章列表

```http
GET /api/articles?page=1&page_size=10&category_id=1&tag_id=2&keyword=test&status=1
Authorization: Bearer {token}
```

**查询参数**:
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10，最大 100）
- `category_id`: 分类 ID（可选）
- `tag_id`: 标签 ID（可选）
- `keyword`: 搜索关键词（可选）
- `status`: 状态筛选（0: 草稿，1: 已发布）

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "total": 100,
    "articles": [
      {
        "id": 1,
        "title": "文章标题",
        "summary": "文章摘要",
        "cover_image": "/uploads/covers/xxx.jpg",
        "view_count": 1234,
        "created_at": "2026-04-01T10:00:00Z",
        "author": {
          "id": 1,
          "username": "author",
          "nickname": "作者",
          "avatar": "/uploads/avatars/xxx.jpg"
        },
        "category": {
          "id": 1,
          "name": "分类名"
        },
        "tags": [
          {"id": 1, "name": "标签 1"}
        ]
      }
    ]
  }
}
```

---

### 2. 获取文章详情

```http
GET /api/articles/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "article": {
      "id": 1,
      "title": "文章标题",
      "content": "Markdown 内容",
      "summary": "摘要",
      "cover_image": "/uploads/covers/xxx.jpg",
      "view_count": 1235,
      "created_at": "2026-04-01T10:00:00Z",
      "updated_at": "2026-04-01T12:00:00Z",
      "author": {...},
      "category": {...},
      "tags": [...],
      "comments": [...]
    }
  }
}
```

---

### 3. 创建文章

```http
POST /api/articles
Content-Type: application/json
Authorization: Bearer {token}

{
  "title": "文章标题 (必填)",
  "content": "Markdown 内容 (必填)",
  "summary": "摘要 (可选，为空自动截取)",
  "cover_image": "/uploads/covers/xxx.jpg (可选)",
  "category_id": 1 (可选，系统分类 ID),
  "user_category_id": 2 (可选，个人分类 ID),
  "tags": [1, 2, 3] (标签 ID 数组),
  "status": 1 (0: 草稿，1: 发布)
}
```

---

### 4. 更新文章

```http
PUT /api/articles/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "title": "新标题",
  "content": "新内容",
  ...
}
```

---

### 5. 删除文章

```http
DELETE /api/articles/:id
Authorization: Bearer {token}
```

---

### 6. 上传封面图

```http
POST /api/articles/upload-cover
Content-Type: multipart/form-data
Authorization: Bearer {token}

FormData:
  cover: File (图片文件，最大 5MB，支持 jpg/png/gif/webp)
```

**响应示例**:
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "url": "/uploads/covers/1234567890_xxx.jpg"
  }
}
```

---

## 🏷️ 分类接口

### 1. 获取所有分类（系统标签）

```http
GET /api/categories
```

---

### 2. 获取我的分类（个人合集）

```http
GET /api/my-categories
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "技术笔记",
        "color": "#409EFF",
        "article_count": 10,
        "created_at": "2026-04-01T10:00:00Z"
      }
    ]
  }
}
```

---

### 3. 创建个人分类

```http
POST /api/user-categories
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "分类名称 (必填)",
  "color": "#409EFF (可选，默认蓝色)"
}
```

---

### 4. 更新个人分类

```http
PUT /api/user-categories/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "新分类名",
  "color": "#67C23A"
}
```

---

### 5. 删除个人分类

```http
DELETE /api/user-categories/:id
Authorization: Bearer {token}
```

**注意**: 有文章使用该分类时无法删除

---

## 🏷️ 标签接口

### 1. 获取所有标签

```http
GET /api/tags
```

---

### 2. 获取热门标签

```http
GET /api/tags/hot?limit=10
```

---

## 💬 评论接口

### 1. 获取文章评论

```http
GET /api/comments?article_id=1&page=1&page_size=20
Authorization: Bearer {token}
```

---

### 2. 发表评论

```http
POST /api/comments
Content-Type: application/json
Authorization: Bearer {token}

{
  "article_id": 1 (必填),
  "content": "评论内容 (必填)",
  "parent_id": null (可选，回复评论时填写父评论 ID)
}
```

---

### 3. 删除评论

```http
DELETE /api/comments/:id
Authorization: Bearer {token}
```

**权限**: 仅评论作者或管理员可删除

---

### 4. 点赞评论

```http
POST /api/comments/:id/like
Authorization: Bearer {token}
```

---

## 👍 互动接口

### 1. 点赞文章

```http
POST /api/articles/:id/like
Authorization: Bearer {token}
```

---

### 2. 收藏文章

```http
POST /api/articles/:id/favorite
Authorization: Bearer {token}
```

---

### 3. 取消收藏

```http
DELETE /api/articles/:id/favorite
Authorization: Bearer {token}
```

---

### 4. 获取我的收藏

```http
GET /api/favorites?page=1&page_size=10
Authorization: Bearer {token}
```

---

## 👤 用户接口

### 1. 获取用户列表（管理后台）

```http
GET /api/admin/users?page=1&page_size=20&keyword=admin&role=1
Authorization: Bearer {token}
```

**权限**: 需要管理员权限

---

### 2. 更新用户角色

```http
PUT /api/admin/users/:id/role
Content-Type: application/json
Authorization: Bearer {token}

{
  "role": 0 (0: 管理员，1: 普通用户)
}
```

---

### 3. 锁定/解锁用户

```http
PUT /api/admin/users/:id/status
Content-Type: application/json
Authorization: Bearer {token}

{
  "status": 0 (0: 锁定，1: 正常),
  "lock_duration": 30 (可选，锁定时长/分钟)
}
```

---

### 4. 删除用户

```http
DELETE /api/admin/users/:id
Authorization: Bearer {token}
```

**权限**: 需要管理员权限

---

## 📊 统计接口

### 1. 获取数据概览（管理后台）

```http
GET /api/admin/stats
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "users": 1234,
    "articles": 567,
    "comments": 890,
    "views": 123456
  }
}
```

---

### 2. 获取我的统计

```http
GET /api/user/stats
Authorization: Bearer {token}
```

---

## 🔍 搜索接口

### 1. 搜索文章

```http
GET /api/search/articles?keyword=Vue&page=1&page_size=20
Authorization: Bearer {token}
```

---

### 2. 搜索用户

```http
GET /api/search/users?keyword=john&page=1&page_size=20
Authorization: Bearer {token}
```

---

## 📋 响应格式规范

### 成功响应

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {...}
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "请求参数错误",
  "data": null
}
```

### 未授权响应

```json
{
  "code": 401,
  "message": "未登录或 Token 已过期",
  "data": null
}
```

### 禁止访问响应

```json
{
  "code": 403,
  "message": "权限不足",
  "data": null
}
```

### 服务器错误响应

```json
{
  "code": 500,
  "message": "服务器内部错误",
  "data": null
}
```

---

## 🔑 状态码说明

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 客户端错误（参数错误、验证失败等） |
| 401 | 未授权（未登录、Token 无效或过期） |
| 403 | 禁止访问（权限不足） |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 📝 使用示例

### JavaScript (Axios)

```javascript
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

// 添加请求拦截器
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 获取文章列表
async function getArticles() {
  const response = await api.get('/articles', {
    params: { page: 1, page_size: 10 }
  })
  return response.data
}

// 创建文章
async function createArticle(article) {
  const response = await api.post('/articles', article)
  return response.data
}

// 使用示例
const articles = await getArticles()
console.log(articles)
```

### Vue 3 Composition API

```vue
<script setup>
import { ref, onMounted } from 'vue'
import { getArticles } from '@/api/article'

const articles = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await getArticles({ page: 1, page_size: 10 })
    articles.value = res.data.articles
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
  }
})
</script>
```

---

## 🌐 Postman 集合

导入以下 JSON 到 Postman 快速测试 API：

```json
{
  "info": {
    "name": "博客系统 API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "认证",
      "item": [
        {
          "name": "用户登录",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"username\": \"admin\",\n  \"password\": \"admin123\",\n  \"captcha_id\": \"xxx\",\n  \"captcha_ans\": \"1234\"\n}"
            },
            "url": {
              "raw": "{{baseUrl}}/auth/login",
              "host": ["{{baseUrl}}"],
              "path": ["auth", "login"]
            }
          }
        }
      ]
    }
  ],
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api"
    }
  ]
}
```

---

**最后更新**: 2026-04-01
