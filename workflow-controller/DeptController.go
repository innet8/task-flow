package controller

import (
	"net/http"
	"strconv"

	"workflow/util"
	"workflow/workflow-engine/service"
)

// GetAllDept 所有部门列表
func GetAllDept(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		util.ResponseErr(writer, "只支持get方法！！")
		return
	}
	request.ParseForm()

	parentId := -1
	if len(request.Form["parentId"]) > 0 {
		parentId, _ = strconv.Atoi(request.Form["parentId"][0])
	}

	result, err := service.GetAllDept(parentId)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseData(writer, result)
}

// 根据部门单一获取所有用户，包括子部门，树形结构
func GetDeptUserByDept(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		util.ResponseErr(writer, "只支持get方法！！")
		return
	}
	request.ParseForm()
	if len(request.Form["deptID"]) == 0 {
		util.ResponseErr(writer, "部门ID不能为空")
		return
	}
	if len(request.Form["deptName"]) == 0 {
		util.ResponseErr(writer, "部门名称不能为空")
		return
	}
	var deptID int
	deptID, _ = strconv.Atoi(request.Form["deptID"][0])
	deptName := request.Form["deptName"][0]
	result, err := service.GetUsersByDeptTree(deptID, deptName)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseData(writer, result)
}

// 根据部门全部获取所有用户，包括子部门，树形结构
func GetAllDeptUserByDept(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		util.ResponseErr(writer, "只支持get方法！！")
		return
	}
	request.ParseForm()
	if len(request.Form["deptID"]) == 0 {
		util.ResponseErr(writer, "部门ID不能为空")
		return
	}
	if len(request.Form["deptName"]) == 0 {
		util.ResponseErr(writer, "部门名称不能为空")
		return
	}
	var deptID int
	deptID, _ = strconv.Atoi(request.Form["deptID"][0])
	deptName := request.Form["deptName"][0]
	result, err := service.GetUsersByDeptAllTree(deptID, deptName)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseData(writer, result)
}
