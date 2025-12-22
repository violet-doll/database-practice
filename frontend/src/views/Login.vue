<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>学生管理系统</h2>
          <p>Student Management System</p>
        </div>
      </template>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="80px"
        size="large"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="User"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-tips">
        <p>默认账号: admin</p>
        <p>默认密码: admin123</p>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { login as loginApi } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  try {
    await loginFormRef.value.validate()
    
    loading.value = true
    const response = await loginApi(loginForm)
    
    if (response.token) {
      // 保存登录信息
      userStore.login({
        token: response.token,
        user: response.user,
        permissions: response.permissions || []
      })
      
      ElMessage.success('登录成功')
      router.push('/')
    } else {
      ElMessage.error('登录失败，请检查用户名和密码')
    }
  } catch (error) {
    console.error('登录错误:', error)
    if (error.response) {
      ElMessage.error(error.response.data?.error || '登录失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 450px;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.card-header p {
  margin: 8px 0 0 0;
  font-size: 14px;
  color: #999;
}

.login-tips {
  margin-top: 20px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 13px;
  color: #666;
  text-align: center;
}

.login-tips p {
  margin: 4px 0;
}
</style>
