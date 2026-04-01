<template>
  <div class="register">
    <el-card class="register-card">
      <h2>用户注册</h2>
      <el-form :model="form" :rules="rules" ref="registerFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="form.username" 
            placeholder="请输入用户名"
            clearable
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="form.password" 
            type="password" 
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="form.confirmPassword" 
            type="password" 
            placeholder="请再次输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input 
            v-model="form.email" 
            placeholder="请输入邮箱"
            clearable
          />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input 
            v-model="form.nickname" 
            placeholder="请输入昵称（可选）"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleRegister" 
            :loading="loading"
            style="width: 100%;"
          >
            注册
          </el-button>
          <div class="extra-links">
            <el-link type="primary" underline="never" @click="$router.push('/login')">
              已有账号？立即登录
            </el-link>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '@/api/auth'

const router = useRouter()
const registerFormRef = ref(null)
const loading = ref(false)

// 表单数据
const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  nickname: '',
})

// 自定义验证器：确认密码
const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  nickname: [
    { max: 50, message: '昵称长度不能超过 50 个字符', trigger: 'blur' },
  ],
}

// 注册处理
const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      
      try {
        await register({
          username: form.username,
          password: form.password,
          email: form.email,
          nickname: form.nickname,
        })
        
        ElMessage.success('注册成功，请登录')
        router.push('/login')
      } catch (error) {
        console.error('注册失败:', error)
      } finally {
        loading.value = false
      }
    } else {
      return false
    }
  })
}
</script>

<style lang="scss" scoped>
.register {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
  // 使用深色山脉背景图（与登录页面一致）
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.9) 0%, rgba(50, 30, 60, 0.85) 100%),
              url('https://images.unsplash.com/photo-1519681393784-d120267933ba?w=1920&q=80') center/cover no-repeat;
  position: relative;
  
  .register-card {
    width: 420px;
    padding: 60px 50px;
    border-radius: 20px;
    box-sizing: border-box;
    max-height: 90vh;
    overflow-y: auto;
    overflow-x: hidden;
    // 毛玻璃效果
    background: rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    animation: fadeIn 0.6s ease-out;
    
    @keyframes fadeIn {
      from {
        opacity: 0;
        transform: scale(0.95);
      }
      to {
        opacity: 1;
        transform: scale(1);
      }
    }
    
    h2 {
      text-align: center;
      margin-bottom: 50px;
      color: #fff;
      font-size: 32px;
      font-weight: 300;
      letter-spacing: 3px;
      text-transform: uppercase;
    }
    
    .extra-links {
      margin-top: 25px;
      text-align: center;
      
      .el-link {
        font-size: 13px;
        font-weight: 300;
        color: rgba(255, 255, 255, 0.6);
        transition: all 0.3s ease;
        
        &:hover {
          color: #fff;
        }
      }
    }
  }
}
</style>
