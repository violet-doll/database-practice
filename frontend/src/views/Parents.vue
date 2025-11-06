<template>
  <div class="parents">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>家长联系方式</span>
          <div>
            <el-input v-model="query.student_no" placeholder="按学生学号筛选" clearable style="width: 200px; margin-right: 8px" />
            <el-button @click="loadList">查询</el-button>
            <el-button type="primary" @click="openCreate">新增</el-button>
          </div>
        </div>
      </template>
      <el-table :data="list" stripe>
        <el-table-column prop="student_id" label="学生ID" width="100" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="relation" label="关系" width="120" />
        <el-table-column prop="phone" label="电话" width="160" />
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="openEdit(scope.row)">编辑</el-button>
            <el-popconfirm title="确认删除?" @confirm="remove(scope.row.ID)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
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

    <el-dialog v-model="visible" :title="editing ? '编辑家长' : '新增家长'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="学生学号">
          <el-input v-model="form.student_no" placeholder="学生学号" :disabled="editing" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.name" placeholder="姓名" />
        </el-form-item>
        <el-form-item label="关系">
          <el-select v-model="form.relation" placeholder="选择关系">
            <el-option label="父亲" value="父亲" />
            <el-option label="母亲" value="母亲" />
            <el-option label="监护人" value="监护人" />
          </el-select>
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.phone" placeholder="电话" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible=false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { fetchParents, createParent, updateParent, deleteParent } from '@/api/parent'

const query = ref({ page: 1, page_size: 10, student_no: '' })
const list = ref([])
const total = ref(0)

const visible = ref(false)
const editing = ref(false)
const currentId = ref(null)
const saving = ref(false)

const emptyForm = () => ({ student_no: '', name: '', relation: '', phone: '' })
const form = ref(emptyForm())

async function loadList() {
  const res = await fetchParents(query.value)
  list.value = res.data.list || []
  total.value = res.data.total || 0
}

function openCreate() {
  editing.value = false
  currentId.value = null
  form.value = emptyForm()
  visible.value = true
}

function openEdit(row) {
  editing.value = true
  currentId.value = row.ID
  // 编辑时不允许改关联学生，为避免后端误判，学号置空并禁用输入
  form.value = { student_no: '', name: row.name, relation: row.relation, phone: row.phone }
  visible.value = true
}

async function submit() {
  saving.value = true
  try {
    if (editing.value) {
      await updateParent(currentId.value, { name: form.value.name, relation: form.value.relation, phone: form.value.phone })
    } else {
      await createParent({ student_no: String(form.value.student_no).trim(), name: form.value.name, relation: form.value.relation, phone: form.value.phone })
    }
    visible.value = false
    await loadList()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteParent(id)
  await loadList()
}

onMounted(loadList)
</script>

<style scoped>
.parents {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.table-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
</style>


