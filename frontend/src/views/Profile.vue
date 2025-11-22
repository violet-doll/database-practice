<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
          <el-button type="primary" @click="showPasswordDialog">修改密码</el-button>
        </div>
      </template>
      
      <div class="profile-content" v-loading="loading">
        <div class="basic-info">
          <el-avatar :size="100" :icon="UserFilled" class="avatar" />
          <h2 class="username">{{ userInfo.username }}</h2>
          <el-tag :type="roleTagType" class="role-tag">{{ roleName }}</el-tag>
        </div>

        <el-divider />

        <el-descriptions :column="1" border class="info-list">
          <el-descriptions-item label="用户ID">{{ userInfo.ID }}</el-descriptions-item>
          <el-descriptions-item label="角色">{{ roleName }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(userInfo.CreatedAt) }}</el-descriptions-item>
          
          <!-- 学生特有信息 -->
          <template v-if="userType === 'student' && detailedInfo.student">
            <el-descriptions-item label="学号">{{ detailedInfo.student.student_id }}</el-descriptions-item>
            <el-descriptions-item label="姓名">{{ detailedInfo.student.name }}</el-descriptions-item>
            <el-descriptions-item label="性别">{{ detailedInfo.student.gender }}</el-descriptions-item>
            <el-descriptions-item label="班级" v-if="detailedInfo.student.class">
              {{ detailedInfo.student.class.class_name }}
            </el-descriptions-item>
            <el-descriptions-item label="联系电话">{{ detailedInfo.student.phone }}</el-descriptions-item>
            <el-descriptions-item label="家庭住址">{{ detailedInfo.student.address }}</el-descriptions-item>
          </template>

          <!-- 教师特有信息 -->
          <template v-if="userType === 'teacher' && detailedInfo.teacher">
            <el-descriptions-item label="工号">{{ detailedInfo.teacher.teacher_id }}</el-descriptions-item>
            <el-descriptions-item label="姓名">{{ detailedInfo.teacher.name }}</el-descriptions-item>
            <el-descriptions-item label="联系电话">{{ detailedInfo.teacher.phone }}</el-descriptions-item>
            <el-descriptions-item label="电子邮箱">{{ detailedInfo.teacher.email }}</el-descriptions-item>
          </template>

          <!-- 家长特有信息 -->
          <template v-if="userType === 'parent' && detailedInfo.parent">
            <el-descriptions-item label="姓名">{{ detailedInfo.parent.name }}</el-descriptions-item>
            <el-descriptions-item label="联系电话">{{ detailedInfo.parent.phone }}</el-descriptions-item>
            <el-descriptions-item label="与学生关系">{{ detailedInfo.parent.relation }}</el-descriptions-item>
            <el-descriptions-item label="关联学生" v-if="detailedInfo.student">
              {{ detailedInfo.student.name }} ({{ detailedInfo.student.student_id }})
            </el-descriptions-item>
          </template>
        </el-descriptions>
      </div>
    </el-card>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="400px">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="旧密码" prop="old_password">
          <el-input v-model="passwordForm.old_password" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input v-model="passwordForm.confirm_password" type="password" show-password></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPasswordChange" :loading="passwordLoading">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, reactive } from 'vue'
import { UserFilled } from '@element-plus/icons-vue'
import { getCurrentUser, updatePassword } from '@/api/auth'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const userInfo = ref({})
const detailedInfo = ref({})
const userType = ref('')

const roleName = computed(() => {
  if (userInfo.value.role) {
    const roleMap = {
      'admin': '管理员',
      'student': '学生',
      'teacher': '教师',
      'parent': '家长'
    }
    return roleMap[userInfo.value.role.role_name] || userInfo.value.role.role_name
  }
  return '未知角色'
})

const roleTagType = computed(() => {
  const role = userInfo.value.role?.role_name
  if (role === 'admin') return 'danger'
  if (role === 'teacher') return 'warning'
  if (role === 'student') return 'success'
  return 'info'
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

const fetchProfile = async () => {
  loading.value = true
  try {
    const res = await getCurrentUser()
    if (res.code === 200) {
      userInfo.value = res.data.user
      // The API returns detailed info in 'student', 'teacher', or 'parent' keys directly in data
      // based on the structure seen in auth.go GetCurrentUser
      if (res.data.student) {
        detailedInfo.value.student = res.data.student
        userType.value = 'student'
      } else if (res.data.teacher) {
        detailedInfo.value.teacher = res.data.teacher
        userType.value = 'teacher'
      } else if (res.data.parent) {
        detailedInfo.value.parent = res.data.parent
        detailedInfo.value.student = res.data.student // Parent also has student info
        userType.value = 'parent'
      } else {
        userType.value = 'admin' // Or other types without specific details
      }
    }
  } catch (error) {
    console.error('获取个人信息失败:', error)
    ElMessage.error('获取个人信息失败')
  } finally {
    loading.value = false
  }
}

// 修改密码相关逻辑
const passwordDialogVisible = ref(false)
const passwordLoading = ref(false)
const passwordFormRef = ref(null)
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const passwordRules = {
  old_password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const showPasswordDialog = () => {
  passwordForm.old_password = ''
  passwordForm.new_password = ''
  passwordForm.confirm_password = ''
  passwordDialogVisible.value = true
}

const submitPasswordChange = async () => {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      passwordLoading.value = true
      try {
        const res = await updatePassword({
            old_password: passwordForm.old_password,
            new_password: passwordForm.new_password
        })
        if (res.code === 200) {
          ElMessage.success('密码修改成功')
          passwordDialogVisible.value = false
        } else {
           // 如果后端返回非200，通常request拦截器会处理，但这里做个保险
           ElMessage.error(res.message || '密码修改失败')
        }
      } catch (error) {
        console.error(error)
      } finally {
        passwordLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>
.profile-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profile-card {
  min-height: 500px;
}

.basic-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.avatar {
  margin-bottom: 15px;
  background-color: #409eff;
}

.username {
  margin: 0 0 10px 0;
  font-size: 24px;
  color: #303133;
}

.role-tag {
  margin-bottom: 10px;
}

.info-list {
  margin-top: 20px;
}

:deep(.el-descriptions__label) {
  width: 120px;
  justify-content: flex-end;
}
</style>
