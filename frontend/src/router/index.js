import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
  },
  {
    path: '/articles',
    name: 'Articles',
    component: () => import('@/views/Articles.vue'),
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: () => import('@/views/ArticleDetail.vue'),
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/creator',
    name: 'CreatorCenter',
    component: () => import('@/views/CreatorCenter.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/creator/edit/:id?',
    name: 'ArticleEditor',
    component: () => import('@/views/ArticleEditor.vue'),
    meta: { requiresAuth: true },
  },
  // 用户详情页路由
  {
    path: '/user/:id',
    name: 'UserProfile',
    component: () => import('@/views/UserProfile.vue'),
  },
  // 管理后台路由
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: () => import('@/views/admin/AdminDashboard.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/articles',
    name: 'AdminArticles',
    component: () => import('@/views/admin/AdminArticles.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/article/edit/:id',
    name: 'AdminEditArticle',
    component: () => import('@/views/ArticleEditor.vue'), // 复用文章编辑器组件
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/categories',
    name: 'AdminCategories',
    component: () => import('@/views/admin/AdminCategories.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/tags',
    name: 'AdminTags',
    component: () => import('@/views/admin/AdminTags.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: () => import('@/views/admin/AdminUsers.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 全局前置路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 公开的路由
  const publicRoutes = ['/', '/login', '/register', '/articles', '/article']
  const isPublicRoute = publicRoutes.some(route => to.path.startsWith(route))
  
  // 检查 localStorage 中是否有 token（页面刷新时 store 会重置，但 localStorage 还在）
  const hasToken = localStorage.getItem('token')
  
  // 检查是否需要登录
  if (to.meta.requiresAuth && !hasToken) {
    ElMessage.warning('请先登录')
    next('/login')
    return
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin && hasToken) {
    const user = userStore.user
    if (user && user.role !== 'admin') {
      ElMessage.error('没有权限访问')
      next('/')
      return
    }
  }
  
  // 如果已登录（通过 token 判断），访问登录页则重定向到首页
  if (hasToken && to.path === '/login') {
    next('/')
  } else {
    next()
  }
})

export default router