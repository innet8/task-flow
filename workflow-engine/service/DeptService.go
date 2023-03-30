package service

import (
	"workflow/workflow-engine/model"

	"github.com/mumushuiding/util"
)

// GetAllDept 所有部门列表
func GetAllDept() (string, error) {
	datas, err := model.GetAllDept()
	if err != nil {
		return "", err
	}
	str, err := util.ToJSONStr(datas)
	if err != nil {
		return "", err
	}
	return str, nil
}
