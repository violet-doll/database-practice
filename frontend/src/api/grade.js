import request from '@/utils/request'

// 全部成绩列表（分页，可按 student_id、course_id 过滤）
export const getGrades = (params) => {
    return request({
        url: '/grades',
        method: 'get',
        params,
    })
}

// 录入成绩
export const createGrade = (data) => {
    return request({
        url: '/grades',
        method: 'post',
        data,
    })
}

// 按学生查询成绩（返回选课+成绩明细）
export const getGradesByStudent = (studentId, params) => {
    return request({
        url: `/grades/student/${studentId}`,
        method: 'get',
        params,
    })
}

// 按课程查询成绩（返回选课+成绩明细）
export const getGradesByCourse = (courseId, params) => {
    return request({
        url: `/grades/course/${courseId}`,
        method: 'get',
        params,
    })
}


