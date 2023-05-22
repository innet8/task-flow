/*
 * @Date: 2022-08-25 14:06:59
 * @LastEditors: StavinLi 495727881@qq.com
 * @LastEditTime: 2022-09-21 14:36:58
 * @FilePath: /Workflow-Vue3/src/plugins/api.js
 */
import http from '@/plugins/axios'

/**
 * 获取角色
 * @param {*} data 
 * @returns 
 */
export function getRoles(data) {
    return http.get(`roles.json`, { params: data })
}

/**
 * 获取部门
 * @param {*} data 
 * @returns 
 */
export function getDepartments(data) {
    return http.get(`/workflow/dootask/getAllDeptUserByDept`, { params: data })
}

/**
 * 获取职员
 * @param {*} data 
 * @returns 
 */
export function getEmployees(data) {
    return http.get(`/workflow/dootask/getUserByName`, { params: data })
}

/**
 * 获取条件字段
 * @param {*} data 
 * @returns 
 */
export function getConditions(data) {
    return http.get(`conditions.json`, { params: data })
}

/**
 * 获取审批数据
 * @param {*} data 
 * @returns 
 */
export function getWorkFlowData(data) {
    // return http.get(`${baseUrl}data.json`, { params: data })
    return http.get(`/workflow/procdef/findByName`, { params: data })
}

/**
 * 设置审批数据
 * @param {*} data 
 * @returns 
 */
export function setWorkFlowData(data) {
    return http.post(`/workflow/procdef/save`, data)
}

/**
 * 获取待办列表
 * @param {*} data 
 * @returns 
 */
export function getBacklogData(data) {
    return http.post(`/workflow/process/findTask`, data)
}



/**
 * 获取已办列表
 * @param {*} data 
 * @returns 
 */
export function getDoneData(data) {
    return http.post(`/workflow/procHistory/findTask`, data)
}


/**
 * 获取详情
 * @param {*} id 
 * @returns 
 */
export function getProcessData(id) {
    return http.get(`/workflow/process/findById?id=${id}` )
}

/**
 * 审批中心同意或拒绝申请
 * @param {*} data 
 * @returns 
 */
export function agreeOrRefuse(data) {
    return http.post(`/workflow/task/complete`, data)
}

/**
 * 审批中心撤销
 * @param {*} data 
 * @returns 
 */
export function revocationTask(data) {
    return http.post(`/workflow/task/withdraw`, data)
}

/**
 * 获取抄送列表
 * @param {*} data 
 * @returns 
 */
export function getNotifyData(data) {
    return http.post(`/workflow/procHistory/findProcNotify`, data)
}

/**
 * 获取已发起列表
 * @param {*} data 
 * @returns 
 */
export function getInitiatedData(data) {
    return http.post(`/workflow/process/startByMyselfAll`, data)
}