<template>
  <div class="admin-roles">
    <el-card>
      <div class="toolbar" style="margin-bottom: 12px; display:flex; gap:8px; align-items:center;">
        <el-input v-model="query.role_name" placeholder="按角色名搜索" clearable style="width: 220px;" />
        <el-button type="primary" @click="loadRoles">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
        <el-button type="success" @click="openCreate">新增角色</el-button>
      </div>

      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="role_name" label="角色名" />
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button type="primary" link @click="openEdit(row)">重命名</el-button>
            <el-button type="success" link @click="openPermissionDialog(row)">权限设置</el-button>
            <el-popconfirm title="确定删除该角色？" @confirm="onDelete(row)">
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
          @current-change="(p)=>{query.page=p; loadRoles()}"
        />
      </div>
    </el-card>

    <el-dialog :title="isEdit ? '重命名角色' : '新增角色'" v-model="visible" width="420px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="角色名">
          <el-input v-model="form.role_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible=false">取消</el-button>
        <el-button type="primary" @click="onSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限设置弹窗 -->
    <el-dialog title="权限设置" v-model="permissionVisible" width="600px">
      <div v-loading="permissionLoading">
        <el-tree
          ref="permissionTreeRef"
          :data="permissionTreeData"
          show-checkbox
          node-key="permission"
          :default-checked-keys="checkedPermissions"
          :props="{ children: 'children', label: 'name' }"
        />
      </div>
      <template #footer>
        <el-button @click="permissionVisible=false">取消</el-button>
        <el-button type="primary" @click="onSavePermissions" :loading="savingPermissions">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchRoles, createRole, updateRole, deleteRole, fetchPermissions, fetchRolePermissions, updateRolePermissions } from '@/api/admin'

const list = ref([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 10, role_name: '' })

const visible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const form = reactive({ role_name: '' })

// 权限管理相关
const permissionVisible = ref(false)
const permissionLoading = ref(false)
const savingPermissions = ref(false)
const permissionTreeRef = ref(null)
const allPermissions = ref([])
const permissionTreeData = ref([])
const checkedPermissions = ref([])
const currentRoleId = ref(null)

function resetQuery() {
  query.page = 1
  query.role_name = ''
  loadRoles()
}

async function loadRoles() {
  const res = await fetchRoles({ ...query })
  list.value = (res.data.list || []).map(r => ({ ...r, id: r.ID }))
  total.value = res.data.total
}

function openCreate() {
  isEdit.value = false
  currentId.value = null
  form.role_name = ''
  visible.value = true
}

function openEdit(row) {
  isEdit.value = true
  currentId.value = row.id
  form.role_name = row.role_name
  visible.value = true
}

async function onSubmit() {
  try {
    if (isEdit.value) {
      await updateRole(currentId.value, { role_name: form.role_name })
      ElMessage.success('更新成功')
    } else {
      await createRole({ role_name: form.role_name })
      ElMessage.success('创建成功')
    }
    visible.value = false
    loadRoles()
  } catch {}
}

async function onDelete(row) {
  try {
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    loadRoles()
  } catch {}
}

// 打开权限设置弹窗
async function openPermissionDialog(row) {
  currentRoleId.value = row.id
  permissionVisible.value = true
  permissionLoading.value = true
  checkedPermissions.value = []

  try {
    // 同时获取所有权限和该角色的权限
    const [allPermsRes, rolePermsRes] = await Promise.all([
      fetchPermissions(),
      fetchRolePermissions(row.id)
    ])

    allPermissions.value = allPermsRes.data.list || []
    const rolePerms = rolePermsRes.data.permissions || []

    // 构建树形结构数据（按分组）
    const grouped = {}
    allPermissions.value.forEach(perm => {
      if (!grouped[perm.group]) {
        grouped[perm.group] = []
      }
      grouped[perm.group].push({
        name: perm.name,
        permission: perm.permission,
        group: perm.group
      })
    })

    permissionTreeData.value = Object.keys(grouped).map(group => ({
      name: getGroupName(group),
      permission: group,
      children: grouped[group]
    }))

    // 设置已选中的权限
    checkedPermissions.value = rolePerms
  } catch (err) {
    ElMessage.error('加载权限失败')
    console.error(err)
  } finally {
    permissionLoading.value = false
  }
}

// 获取分组名称（中文显示）
function getGroupName(group) {
  const groupNames = {
    student: '学生管理',
    class: '班级管理',
    course: '课程管理',
    schedule: '排课管理',
    enrollment: '选课管理',
    grade: '成绩管理',
    attendance: '考勤管理',
    reward: '奖惩管理',
    parent: '家长管理',
    notification: '通知管理',
    admin: '管理员'
  }
  return groupNames[group] || group
}

// 保存权限设置
async function onSavePermissions() {
  if (!permissionTreeRef.value) return

  savingPermissions.value = true
  try {
    const checked = permissionTreeRef.value.getCheckedKeys()
    
    // 过滤掉分组节点（只保留权限标识，权限标识包含冒号）
    const permissions = checked.filter(key => {
      // 检查是否是权限标识（包含冒号，例如 "student:read"）
      return key.includes(':')
    })

    await updateRolePermissions(currentRoleId.value, permissions)
    ElMessage.success('权限更新成功')
    permissionVisible.value = false
  } catch (err) {
    ElMessage.error('保存权限失败')
    console.error(err)
  } finally {
    savingPermissions.value = false
  }
}

onMounted(() => {
  loadRoles()
})
</script>

<style scoped>
</style>


