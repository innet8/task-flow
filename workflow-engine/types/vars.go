package types

import (
	"errors"
	"strings"
	"workflow/util"
)

// 参数
type Vars struct {
	Type        string `json:"type"`        //请假类型，
	Description string `json:"description"` //请假事由
	StartTime   string `json:"startTime"`   //开始时间
	EndTime     string `json:"endTime"`     //结束时间
}

var VacateTypes = []string{"年假", "事假", "病假", "调休假", "婚假", "产假", "陪产假", "丧假", "哺乳假"}

// 验证请假类型
func CheckVacateVars(vars *Vars) (bool, error) {
	if vars == nil {
		return false, errors.New("var 不存在")
	}
	if len(vars.Type) == 0 {
		return false, errors.New("var.type 不存在")
	}
	if !strings.Contains(strings.Join(VacateTypes, ","), vars.Type) {
		return false, errors.New("var.type错误,只允许：" + strings.Join(VacateTypes, ","))
	}
	if len(vars.StartTime) == 0 {
		return false, errors.New("var.startTime 不存在")
	}
	if len(vars.EndTime) == 0 {
		return false, errors.New("var.endTime 不存在")
	}
	if len(vars.Description) == 0 {
		return false, errors.New("var.description 不存在")
	}
	if !util.ValidateTimeFormat(vars.StartTime) {
		return false, errors.New("startTime 无效的时间格式")
	}
	if !util.ValidateTimeFormat(vars.EndTime) {
		return false, errors.New("endTime 无效的时间格式")
	}
	if vars.StartTime > vars.EndTime {
		return false, errors.New("startTime不能大于endTime")
	}
	return true, nil
}

// 验证加班申请类型
func CheckOvertimeVars(vars *Vars) (bool, error) {
	if vars == nil {
		return false, errors.New("var 不存在")
	}
	if len(vars.StartTime) == 0 {
		return false, errors.New("var.startTime 不存在")
	}
	if len(vars.EndTime) == 0 {
		return false, errors.New("var.endTime 不存在")
	}
	if len(vars.Description) == 0 {
		return false, errors.New("var.description 不存在")
	}
	return true, nil
}
