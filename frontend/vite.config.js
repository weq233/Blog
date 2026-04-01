import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        // 使用新的现代 API
        api: 'modern-compiler',
      },
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 管理后台 API 代理
      '/admin': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // uploads 静态文件代理
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
