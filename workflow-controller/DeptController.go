package controller

import (
	"net/http"
	"workflow/workflow-engine/service"

	"github.com/mumushuiding/util"
)

// GetAllDept 所有部门列表
func GetAllDept(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		util.ResponseErr(writer, "只支持get方法！！")
		return
	}
	result, err := service.GetAllDept()
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseData(writer, result)
}
