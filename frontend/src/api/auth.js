import request from '@/utils/request'

// 获取验证码
export function getCaptcha(captchaId = '') {
  return request({
    url: '/captcha',
    method: 'get',
    params: { captcha_id: captchaId },
  })
}

// 获取数学验证码
export function getMathCaptcha() {
  return request({
    url: '/captcha/math',
    method: 'get',
  })
}

// 用户登录
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data,
  })
}

// 用户注册
export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data,
  })
}

// 获取当前用户信息
export function getCurrentUser() {
  return request({
    url: '/user/me',
    method: 'get',
  })
}

// 获取用户统计数据（当前用户）
export function getUserStats() {
  return request({
    url: '/user/stats',
    method: 'get',
  })
}

// 获取用户的文章列表（创作中心）
export function getMyArticles(params = {}) {
  return request({
    url: '/my-articles',
    method: 'get',
    params,
  })
}

// 创建文章
export function createArticle(data) {
  return request({
    url: '/articles/create',
    method: 'post',
    data,
  })
}

// 更新文章
export function updateArticle(id, data) {
  return request({
    url: `/articles/${id}`,
    method: 'put',
    data,
  })
}

// 删除文章
export function deleteArticle(id) {
  return request({
    url: `/articles/${id}`,
    method: 'delete',
  })
}

// 获取文章列表
export function getArticles(params = {}) {
  return request({
    url: '/articles',
    method: 'get',
    params,
  })
}

// 获取文章详情
export function getArticle(id) {
  return request({
    url: `/article/${id}`,
    method: 'get',
  })
}

// 按分类获取文章
export function getArticlesByCategory(slug) {
  return request({
    url: `/articles/category/${slug}`,
    method: 'get',
  })
}

// 按标签获取文章
export function getArticlesByTag(slug) {
  return request({
    url: `/articles/tag/${slug}`,
    method: 'get',
  })
}

// 上传封面图片
export function uploadCoverImage(file) {
  const formData = new FormData()
  formData.append('cover', file)
  return request({
    url: '/upload/cover',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

// 获取分类列表
export function getCategories() {
  return request({
    url: '/categories',
    method: 'get',
  })
}

// 获取标签列表
export function getTags() {
  return request({
    url: '/tags',
    method: 'get',
  })
}

// 获取我的分类列表（用户自定义分类）
export function getMyCategories() {
  return request({
    url: '/my-categories',
    method: 'get',
  })
}

// 创建用户分类
export function createUserCategory(data) {
  return request({
    url: '/user-categories',
    method: 'post',
    data,
  })
}

// 更新用户分类
export function updateUserCategory(id, data) {
  return request({
    url: `/user-categories/${id}`,
    method: 'put',
    data,
  })
}

// 删除用户分类
export function deleteUserCategory(id) {
  return request({
    url: `/user-categories/${id}`,
    method: 'delete',
  })
}

// 点赞文章
export function likeArticle(id) {
  return request({
    url: `/article/${id}/like`,
    method: 'post',
  })
}

// 获取点赞状态
export function getLikeStatus(id) {
  return request({
    url: `/article/${id}/like-status`,
    method: 'get',
  })
}

// 关注用户
export function followUser(userId) {
  return request({
    url: `/user/${userId}/follow`,
    method: 'post',
  })
}

// 获取关注状态
export function getFollowStatus(userId) {
  return request({
    url: `/user/${userId}/follow-status`,
    method: 'get',
  })
}

// 发表评论
export function createComment(articleId, data) {
  return request({
    url: `/article/${articleId}/comments`,
    method: 'post',
    data,
  })
}

// 获取评论列表
export function getComments(articleId) {
  return request({
    url: `/article/${articleId}/comments`,
    method: 'get',
  })
}

// 删除评论
export function deleteComment(id) {
  return request({
    url: `/comment/${id}`,
    method: 'delete',
  })
}

// 获取用户资料
export function getUserProfile(userId) {
  return request({
    url: `/user/${userId}/profile`,
    method: 'get',
  })
}

// 获取用户统计信息（指定用户）
export function getUserStatsById(userId) {
  return request({
    url: `/user/${userId}/stats`,
    method: 'get',
  })
}

// 获取用户文章列表
export function getUserArticles(userId, params) {
  return request({
    url: `/user/${userId}/articles`,
    method: 'get',
    params,
  })
}

// 获取用户评论列表
export function getUserComments(userId, params) {
  return request({
    url: `/user/${userId}/comments`,
    method: 'get',
    params,
  })
}
