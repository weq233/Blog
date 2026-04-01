# 全局导航栏组件使用说明

## 组件位置
`@/components/TopNavBar.vue`

## 功能特性
- ✅ 自动检测用户登录状态
- ✅ 显示用户头像和用户名（已登录）
- ✅ 显示登录/注册按钮（未登录）
- ✅ 点击 Logo 返回首页
- ✅ 退出登录功能
- ✅ 响应式布局（移动端适配）

## 使用方式

### 1. 在页面中引入组件

```vue
<template>
  <div class="page-container">
    <TopNavBar />
    <!-- 页面内容 -->
  </div>
</template>

<script setup>
import TopNavBar from '@/components/TopNavBar.vue'
</script>
```

### 2. 完整示例

```vue
<template>
  <div class="my-page">
    <TopNavBar />
    
    <div class="content">
      <h1>我的页面</h1>
      <p>页面内容...</p>
    </div>
  </div>
</template>

<script setup>
import TopNavBar from '@/components/TopNavBar.vue'
</script>

<style lang="scss" scoped>
.my-page {
  min-height: 100vh;
  // 页面样式...
}
</style>
```

## 注意事项

1. **无需传递用户信息**：组件内部会自动从 Pinia store 获取用户信息
2. **样式自定义**：组件样式是 scoped 的，不会影响外部样式
3. **响应式支持**：组件已内置移动端响应式样式
4. **路由依赖**：确保项目已配置 Vue Router

## 已应用的页面

- ✅ 首页 (`Home.vue`)
- ✅ 文章列表页 (`Articles.vue`)
- ✅ 文章详情页 (`ArticleDetail.vue`)

## 创建新页面时

只需在页面模板开头添加 `<TopNavBar />` 即可，无需重复编写导航栏代码！
