import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import App from './App.vue'
import router from './router'
import { useUserStore } from '@/stores/user'

// 导入全局样式
import './styles/index.scss'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

// 在路由初始化前尝试自动登录
const userStore = useUserStore()
userStore.autoLogin().finally(() => {
  // 确保用户信息恢复后再挂载应用
  app.use(router)
  app.use(ElementPlus, {
    locale: zhCn,
  })
  app.mount('#app')
})
