<template>
  <div class="notifications">
    <!-- 发布通知 -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>发布通知</span>
        </div>
      </template>
      <el-form :model="form" label-width="90px" class="notify-form">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item label="标题">
              <el-input v-model="form.title" placeholder="请输入标题" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="目标">
              <div class="target-inline">
                <el-select v-model="form.target" placeholder="选择目标" style="width: 140px">
                  <el-option label="全部用户" value="all" />
                  <el-option label="学生ID" value="student:" />
                  <el-option label="家长ID" value="parent:" />
                  <el-option label="班级ID" value="class:" />
                </el-select>
                <el-input v-model="targetId" :disabled="form.target==='all'" placeholder="如 12" style="width: 180px" />
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="内容">
              <el-input v-model="form.content" type="textarea" :rows="4" placeholder="请输入内容" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="渠道">
              <el-checkbox-group v-model="form.channels">
                <el-checkbox label="sms">短信</el-checkbox>
                <el-checkbox label="email">邮件</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <div class="form-actions">
              <el-button @click="resetForm">重置</el-button>
              <el-button type="primary" :loading="sending" @click="handleSend">发布</el-button>
            </div>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 通知列表 -->
    <el-card style="margin-top: 16px">
      <template #header>
        <div class="card-header">
          <span>通知列表</span>
          <div>
            <el-input v-model="query.keyword" placeholder="按标题/内容关键字" clearable style="width: 240px; margin-right: 8px" />
            <el-input v-model="query.target" placeholder="按目标筛选 (如 student:12)" clearable style="width: 240px; margin-right: 8px" />
            <el-button type="primary" @click="loadList">查询</el-button>
          </div>
        </div>
      </template>
      <el-table :data="list" stripe>
        <el-table-column prop="title" label="标题" min-width="180" />
        <el-table-column prop="content" label="内容" min-width="300">
          <template #default="scope">
            <div class="content-wrap">{{ scope.row.content }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="target" label="目标" width="200" />
        <el-table-column prop="CreatedAt" label="时间" width="200" />
      </el-table>
      <div class="table-footer">
        <el-pagination
          background
          layout="prev, pager, next, total"
          :total="total"
          :page-size="query.page_size"
          :current-page="query.page"
          @current-change="(p)=>{query.page=p; loadList()}"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { createNotification, fetchNotifications } from '@/api/notification'

const form = ref({
  title: '',
  content: '',
  target: 'all',
  channels: ['sms'],
})
const targetId = ref('')
const sending = ref(false)

const query = ref({ page: 1, page_size: 10, target: '', keyword: '' })
const list = ref([])
const total = ref(0)

function resetForm() {
  form.value = { title: '', content: '', target: 'all', channels: ['sms'] }
  targetId.value = ''
}

async function handleSend() {
  if (!form.value.title || !form.value.content) return
  sending.value = true
  const payload = {
    ...form.value,
    target: form.value.target === 'all' ? 'all' : `${form.value.target}${targetId.value || ''}`,
  }
  try {
    await createNotification(payload)
    resetForm()
    query.value.page = 1
    await loadList()
  } finally {
    sending.value = false
  }
}

async function loadList() {
  const res = await fetchNotifications(query.value)
  list.value = res.data.list || []
  total.value = res.data.total || 0
}

onMounted(loadList)
</script>

<style scoped>
.notifications {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.notify-form {
  padding-top: 4px;
}
.target-inline {
  display: flex;
  gap: 8px;
  align-items: center;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
.table-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
.content-wrap {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
