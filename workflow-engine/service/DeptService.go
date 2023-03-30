package service

import (
	"workflow/workflow-engine/model"

	"workflow/util"
)

// GetAllDept 所有部门列表
func GetAllDept(parentId int) (string, error) {
	datas, err := model.GetAllDept(parentId)
	if err != nil {
		return "", err
	}
	str, err := util.ToJSONStr(datas)
	if err != nil {
		return "", err
	}
	return str, nil
}
