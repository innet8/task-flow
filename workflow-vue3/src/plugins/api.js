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
    return http.get(`/workflow/dootask/getAllDept`, { params: data })
}

/**
 * 获取职员
 * @param {*} data 
 * @returns 
 */
export function getEmployees(data) {
    return http.get(`employees.json`, { params: data })
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
    return http.get(`/workflow/procdef/findById`, { params: data })
}

/**
 * 设置审批数据
 * @param {*} data 
 * @returns 
 */
export function setWorkFlowData(data) {
    return http.post(`/workflow/procdef/save`, data)
}