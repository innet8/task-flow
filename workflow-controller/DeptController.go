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
