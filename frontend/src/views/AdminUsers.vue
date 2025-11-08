<template>
  <div class="admin-users">
    <el-card>
      <div class="toolbar" style="margin-bottom: 12px; display: flex; gap: 8px; align-items: center;">
        <el-input v-model="query.username" placeholder="按用户名搜索" clearable style="width: 220px;" />
        <el-select v-model="query.role_id" placeholder="角色" clearable style="width: 160px;">
          <el-option v-for="r in roles" :key="r.id" :label="r.role_name" :value="r.id" />
        </el-select>
        <el-select v-model="query.is_active" placeholder="状态" clearable style="width: 140px;">
          <el-option label="启用" :value="true" />
          <el-option label="禁用" :value="false" />
        </el-select>
        <el-button type="primary" @click="loadUsers">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
        <el-button type="success" @click="openCreate">新增用户</el-button>
      </div>

      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role.role_name" label="角色" width="120" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-switch :model-value="row.is_active" @change="(val)=>toggleActive(row, val)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button type="primary" link @click="openEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该用户？" @confirm="onDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div style="margin-top: 12px; text-align: right;">
        <el-pagination
          background
          layout="prev, pager, next, jumper, total"
          :total="total"
          :current-page="query.page"
          :page-size="query.page_size"
          @current-change="(p)=>{query.page=p; loadUsers()}"
        />
      </div>
    </el-card>

    <el-dialog :title="isEdit ? '编辑用户' : '新增用户'" v-model="visible" width="480px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="用户名" v-if="!isEdit">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="留空则不修改" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_id" style="width: 100%;">
            <el-option v-for="r in roles" :key="r.id" :label="r.role_name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible=false">取消</el-button>
        <el-button type="primary" @click="onSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
  </template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchUsers, createUser, updateUser, deleteUser, fetchRoles } from '@/api/admin'

const list = ref([])
const total = ref(0)
const roles = ref([])

const query = reactive({
  page: 1,
  page_size: 10,
  username: '',
  role_id: undefined,
  is_active: undefined,
})

const visible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const form = reactive({ username: '', password: '', role_id: undefined, is_active: true })

function resetQuery() {
  query.page = 1
  query.username = ''
  query.role_id = undefined
  query.is_active = undefined
  loadUsers()
}

async function loadRoles() {
  const res = await fetchRoles({ page: 1, page_size: 100 })
  roles.value = (res.data.list || []).map(r => ({ ...r, id: r.ID }))
}

async function loadUsers() {
  const params = { ...query }
  const res = await fetchUsers(params)
  list.value = (res.data.list || []).map(u => ({ ...u, id: u.ID }))
  total.value = res.data.total
}

function openCreate() {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { username: '', password: '', role_id: roles.value[0]?.id, is_active: true })
  visible.value = true
}

function openEdit(row) {
  isEdit.value = true
  currentId.value = row.id
  Object.assign(form, { username: row.username, password: '', role_id: row.role_id, is_active: row.is_active })
  visible.value = true
}

async function onSubmit() {
  try {
    if (isEdit.value) {
      const payload = { }
      if (form.password) payload.password = form.password
      if (form.role_id != null) payload.role_id = form.role_id
      if (form.is_active != null) payload.is_active = form.is_active
      await updateUser(currentId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createUser({ username: form.username, password: form.password, role_id: form.role_id, is_active: form.is_active })
      ElMessage.success('创建成功')
    }
    visible.value = false
    loadUsers()
  } catch {}
}

async function onDelete(row) {
  try {
    await deleteUser(row.id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch {}
}

async function toggleActive(row, val) {
  try {
    await updateUser(row.id, { is_active: val })
    row.is_active = val
    ElMessage.success('状态已更新')
  } catch {}
}

onMounted(async () => {
  await loadRoles()
  await loadUsers()
})
</script>

<style scoped>
</style>


