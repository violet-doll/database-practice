<template>
  <div class="grades">
    <el-card>
      <div class="toolbar">
        <el-button type="primary" @click="dialogVisible = true">录入成绩</el-button>
      </div>

      <el-tabs v-model="activeTab" style="margin-top: 16px">
        <el-tab-pane label="全部成绩" name="all">
          <el-table :data="allGrades" style="width: 100%" v-loading="loadingAll">
            <el-table-column label="学生ID" width="120">
              <template #default="{ row }">{{ row.student_id }}</template>
            </el-table-column>
            <el-table-column label="学生姓名" min-width="160">
              <template #default="{ row }">{{ (row.Student && row.Student.name) }}</template>
            </el-table-column>
            <el-table-column label="课程ID" width="120">
              <template #default="{ row }">{{ row.course_id }}</template>
            </el-table-column>
            <el-table-column label="课程名称" min-width="180">
              <template #default="{ row }">{{ (row.Course && row.Course.course_name) }}</template>
            </el-table-column>
            <el-table-column label="成绩明细">
              <template #default="{ row }">
                <div v-if="(row.Grades && row.Grades.length)">
                  <el-tag v-for="g in row.Grades" :key="g.ID" style="margin-right: 6px; margin-bottom: 6px">
                    {{ g.score_type }}: {{ g.score }}
                  </el-tag>
                </div>
                <span v-else>无</span>
              </template>
            </el-table-column>
          </el-table>
          <div class="pagination">
            <el-pagination
              v-model:current-page="allPage.page"
              v-model:page-size="allPage.page_size"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="allPage.total"
              @size-change="() => fetchAll()"
              @current-change="() => fetchAll()"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="按学生查询" name="byStudent">
          <div class="search-box">
            <el-input v-model="studentQuery.student_id" placeholder="学生ID" style="width: 180px; margin-right: 10px" />
            <el-button type="primary" @click="fetchByStudent">查询</el-button>
          </div>

          <el-table :data="studentGrades" style="width: 100%; margin-top: 16px" v-loading="loadingStudent">
            <el-table-column label="学生ID" width="120">
              <template #default="{ row }">{{ row.student_id }}</template>
            </el-table-column>
            <el-table-column label="课程ID" width="120">
              <template #default="{ row }">{{ row.course_id }}</template>
            </el-table-column>
            <el-table-column label="课程名称" min-width="180">
              <template #default="{ row }">{{ (row.Course && row.Course.course_name) }}</template>
            </el-table-column>
            <el-table-column label="成绩明细">
              <template #default="{ row }">
                <div v-if="(row.Grades && row.Grades.length)">
                  <el-tag v-for="g in row.Grades" :key="g.ID" style="margin-right: 6px; margin-bottom: 6px">
                    {{ g.score_type }}: {{ g.score }}
                  </el-tag>
                </div>
                <span v-else>无</span>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="studentPage.page"
              v-model:page-size="studentPage.page_size"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="studentPage.total"
              @size-change="() => fetchByStudent()"
              @current-change="() => fetchByStudent()"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="按课程查询" name="byCourse">
          <div class="search-box">
            <el-input v-model="courseQuery.course_id" placeholder="课程ID" style="width: 180px; margin-right: 10px" />
            <el-button type="primary" @click="fetchByCourse">查询</el-button>
          </div>

          <el-table :data="courseGrades" style="width: 100%; margin-top: 16px" v-loading="loadingCourse">
            <el-table-column label="课程ID" width="120">
              <template #default="{ row }">{{ row.course_id }}</template>
            </el-table-column>
            <el-table-column label="学生姓名" min-width="160">
              <template #default="{ row }">{{ (row.Student && row.Student.name) }}</template>
            </el-table-column>
            <el-table-column label="成绩明细">
              <template #default="{ row }">
                <div v-if="(row.Grades && row.Grades.length)">
                  <el-tag v-for="g in row.Grades" :key="g.ID" style="margin-right: 6px; margin-bottom: 6px">
                    {{ g.score_type }}: {{ g.score }}
                  </el-tag>
                </div>
                <span v-else>无</span>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="coursePage.page"
              v-model:page-size="coursePage.page_size"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="coursePage.total"
              @size-change="() => fetchByCourse()"
              @current-change="() => fetchByCourse()"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="dialogVisible" title="录入成绩" width="520px">
      <el-form ref="gradeFormRef" :model="gradeForm" :rules="gradeRules" label-width="100px">
        <el-form-item label="学生ID" prop="student_id">
          <el-input v-model="gradeForm.student_id" />
        </el-form-item>
        <el-form-item label="课程ID" prop="course_id">
          <el-input v-model="gradeForm.course_id" />
        </el-form-item>
        <el-form-item label="成绩类型" prop="score_type">
          <el-select v-model="gradeForm.score_type" style="width: 100%">
            <el-option label="平时成绩" value="平时成绩" />
            <el-option label="期末成绩" value="期末成绩" />
            <el-option label="总评" value="总评" />
          </el-select>
        </el-form-item>
        <el-form-item label="分数" prop="score">
          <el-input-number v-model="gradeForm.score" :min="0" :max="100" :step="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGrade">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { createGrade, getGradesByStudent, getGradesByCourse, getGrades } from '@/api/grade'

const activeTab = ref('all')

// 全部成绩
const allGrades = ref([])
const loadingAll = ref(false)
const allPage = reactive({ page: 1, page_size: 10, total: 0 })
const fetchAll = async () => {
  loadingAll.value = true
  try {
    const res = await getGrades({ page: allPage.page, page_size: allPage.page_size })
    allGrades.value = res.data.list || []
    allPage.total = res.data.total || 0
  } finally {
    loadingAll.value = false
  }
}

// 录入成绩
const dialogVisible = ref(false)
const gradeFormRef = ref(null)
const gradeForm = reactive({ student_id: '', course_id: '', score_type: '平时成绩', score: 0 })
const gradeRules = {
  student_id: [{ required: true, message: '请输入学生ID', trigger: 'blur' }],
  course_id: [{ required: true, message: '请输入课程ID', trigger: 'blur' }],
  score_type: [{ required: true, message: '请选择成绩类型', trigger: 'change' }],
  score: [{ required: true, message: '请输入分数', trigger: 'change' }],
}

const submitGrade = async () => {
  if (!gradeFormRef.value) return
  await gradeFormRef.value.validate(async (valid) => {
    if (!valid) return
    await createGrade({
      student_id: Number(gradeForm.student_id),
      course_id: Number(gradeForm.course_id),
      score_type: gradeForm.score_type,
      score: Number(gradeForm.score),
    })
    ElMessage.success('录入成功')
    
    // 修复:只在有查询条件时才刷新对应的列表
    if (activeTab.value === 'byStudent' && studentQuery.student_id) {
      await fetchByStudent()
    }
    if (activeTab.value === 'byCourse' && courseQuery.course_id) {
      await fetchByCourse()
    }
    // 默认列表也刷新
    await fetchAll()
    
    dialogVisible.value = false
  })
}

onMounted(() => {
  fetchAll()
})

// 按学生查询
const studentQuery = reactive({ student_id: '' })
const studentGrades = ref([])
const loadingStudent = ref(false)
const studentPage = reactive({ page: 1, page_size: 10, total: 0 })

const fetchByStudent = async () => {
  if (!studentQuery.student_id) {
    studentGrades.value = []
    studentPage.total = 0
    return
  }
  loadingStudent.value = true
  try {
    const res = await getGradesByStudent(Number(studentQuery.student_id), {
      page: studentPage.page,
      page_size: studentPage.page_size,
    })
    studentGrades.value = res.data.list || []
    studentPage.total = res.data.total || 0
  } finally {
    loadingStudent.value = false
  }
}

// 按课程查询
const courseQuery = reactive({ course_id: '' })
const courseGrades = ref([])
const loadingCourse = ref(false)
const coursePage = reactive({ page: 1, page_size: 10, total: 0 })

const fetchByCourse = async () => {
  if (!courseQuery.course_id) {
    courseGrades.value = []
    coursePage.total = 0
    return
  }
  loadingCourse.value = true
  try {
    const res = await getGradesByCourse(Number(courseQuery.course_id), {
      page: coursePage.page,
      page_size: coursePage.page_size,
    })
    courseGrades.value = res.data.list || []
    coursePage.total = res.data.total || 0
  } finally {
    loadingCourse.value = false
  }
}
</script>

<style scoped>
.grades { padding: 20px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.search-box { display: flex; align-items: center; }
.pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
