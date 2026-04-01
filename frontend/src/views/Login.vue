<template>
  <div class="login">
    <el-card class="login-card">
      <h2>用户登录</h2>
      <el-form :model="form" :rules="rules" ref="loginFormRef" label-width="80px">
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
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item label="验证码" prop="captcha_ans">
          <div class="captcha-row">
            <el-input 
              v-model="form.captcha_ans" 
              placeholder="请输入答案"
              style="width: 150px; margin-right: 10px;"
              @keyup.enter="handleLogin"
            />
            <div class="captcha-question">
              {{ captchaQuestion }}
            </div>
            <el-button 
              type="primary" 
              link 
              @click="refreshCaptcha"
              style="margin-left: 10px;"
            >
              换一张
            </el-button>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleLogin" 
            :loading="loading"
            style="width: 100%;"
          >
            登录
          </el-button>
          <div class="extra-links">
            <el-link type="primary" underline="never" @click="$router.push('/register')">
              还没有账号？立即注册
            </el-link>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login, getMathCaptcha } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const loginFormRef = ref(null)
const loading = ref(false)
const captchaId = ref('')
const captchaQuestion = ref('')

// 表单数据
const form = reactive({
  username: '',
  password: '',
  captcha_id: '',
  captcha_ans: '',
})

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
  captcha_ans: [
    { required: true, message: '请输入验证码答案', trigger: 'blur' },
  ],
}

// 获取验证码
const refreshCaptcha = async () => {
  try {
    const res = await getMathCaptcha()
    captchaId.value = res.data.captcha_id
    captchaQuestion.value = res.data.question
    form.captcha_id = res.data.captcha_id
  } catch (error) {
    console.error('获取验证码失败:', error)
    ElMessage.error('获取验证码失败，请刷新重试')
  }
}

// 登录处理
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      
      try {
        const res = await login({
          username: form.username,
          password: form.password,
          captcha_id: form.captcha_id,
          captcha_ans: form.captcha_ans,
        })
        
        // 登录成功，保存 token 和用户信息
        userStore.setToken(res.data.token, 24 * 60 * 60 * 1000) // 24 小时
        userStore.setUserInfo(res.data.user)
        
        ElMessage.success('登录成功')
        
        // 跳转到首页或之前访问的页面
        router.push('/')
      } catch (error) {
        console.error('登录失败:', error)
        // 刷新验证码
        refreshCaptcha()
        form.captcha_ans = ''
      } finally {
        loading.value = false
      }
    } else {
      return false
    }
  })
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style lang="scss" scoped>
.login {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
  // 使用深色山脉背景图
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.9) 0%, rgba(50, 30, 60, 0.85) 100%),
              url('https://images.unsplash.com/photo-1519681393784-d120267933ba?w=1920&q=80') center/cover no-repeat;
  position: relative;
  
  .login-card {
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
      margin-top: 0;
      letter-spacing: 3px;
      text-transform: uppercase;
    }
    
    :deep(.el-form-item) {
      margin-bottom: 30px;
      
      .el-form-item__label {
        display: none;
      }
      
      .el-form-item__content {
        .el-input__wrapper {
          background: transparent;
          box-shadow: none;
          border: none;
          border-bottom: 1px solid rgba(255, 255, 255, 0.3);
          border-radius: 0;
          padding: 12px 0;
          transition: all 0.3s ease;
          
          &:hover {
            border-bottom-color: rgba(255, 255, 255, 0.6);
          }
          
          &.is-focus {
            box-shadow: none;
            border-bottom-color: #fff;
          }
          
          .el-input__inner {
            color: #fff;
            font-size: 15px;
            letter-spacing: 0.5px;
            
            &::placeholder {
              color: rgba(255, 255, 255, 0.5);
              font-weight: 300;
            }
          }
        }
        
        .el-input__suffix {
          .el-icon {
            color: rgba(255, 255, 255, 0.5);
          }
        }
      }
    }
    
    .captcha-row {
      display: flex;
      align-items: center;
      gap: 15px;
      width: 100%;
      margin-top: 10px;
      
      .el-input {
        flex: 1;
        
        .el-input__wrapper {
          background: rgba(255, 255, 255, 0.05);
          border-radius: 8px;
          border: 1px solid rgba(255, 255, 255, 0.2);
          padding: 10px 15px;
          
          .el-input__inner {
            &::placeholder {
              color: rgba(255, 255, 255, 0.4);
            }
          }
        }
      }
      
      .captcha-question {
        font-size: 18px;
        font-weight: 400;
        color: #fff;
        padding: 10px 20px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        text-align: center;
        white-space: nowrap;
        flex-shrink: 0;
        border: 1px solid rgba(255, 255, 255, 0.2);
        letter-spacing: 1px;
      }
      
      .el-button {
        font-size: 13px;
        padding: 8px 15px;
        flex-shrink: 0;
        color: rgba(255, 255, 255, 0.8);
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.3);
        transition: all 0.3s ease;
        
        &:hover {
          color: #fff;
          border-color: #fff;
          background: rgba(255, 255, 255, 0.1);
        }
      }
    }
    
    .el-button[type="primary"] {
      height: 50px;
      font-size: 15px;
      font-weight: 400;
      border-radius: 8px;
      background: #fff;
      color: #2d3748;
      border: none;
      margin-top: 20px;
      transition: all 0.3s ease;
      letter-spacing: 2px;
      text-transform: uppercase;
      
      &:hover {
        background: rgba(255, 255, 255, 0.9);
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(255, 255, 255, 0.3);
      }
      
      &:active {
        transform: translateY(0);
      }
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

// 滚动条美化
.login-card::-webkit-scrollbar {
  width: 6px;
}

.login-card::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.login-card::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
}
</style>
