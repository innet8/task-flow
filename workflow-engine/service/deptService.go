package service

import (
	"fmt"
	"workflow/workflow-engine/model"

	"workflow/util"
)

// GetAllDepartments 获取所有部门列表
func GetAllDepartments(parentID int) (string, error) {
	departments, err := model.GetAllDept(parentID)
	if err != nil {
		return "", fmt.Errorf("failed to get all departments: %v", err)
	}

	jsonStr, err := util.ToJSONStr(departments)
	if err != nil {
		return "", fmt.Errorf("failed to convert departments to JSON string: %v", err)
	}

	return jsonStr, nil
}

// GetUserAndChildDepts 根据部门ID获取子部门列表和所有子部门下的用户
func GetUserAndChildDepts(deptID int) (string, error) {
	// 获取子部门列表
	childDepts, err := model.GetDeptByParentID(deptID)
	if err != nil {
		return "", fmt.Errorf("failed to get child departments: %v", err)
	}

	// 获取父级部门
	parentDept, err := model.GetDeptByID(deptID)
	if err != nil {
		return "", fmt.Errorf("failed to get parent department: %v", err)
	}

	// 获取部门父级下的用户列表
	users, err := model.GetUsersByDept(parentDept.Name)
	if err != nil {
		return "", fmt.Errorf("failed to get users by department: %v", err)
	}

	dept := map[string]interface{}{
		"childDepartments": childDepts,
		"employees":        users,
		"titleDepartments": parentDept.Name,
	}

	jsonStr, err := util.ToJSONStr(dept)
	if err != nil {
		return "", fmt.Errorf("failed to convert result to JSON string: %v", err)
	}

	return jsonStr, nil
}

// GetUsersByDeptAllTree 根据部门ID获取子部门列表和所有子部门下的用户
func GetUsersByDeptAllTree(parentID int) (string, error) {
	// 获取子部门列表
	childDepts, err := model.GetDeptByParentID(parentID)
	if err != nil {
		return "", fmt.Errorf("failed to get child departments: %v", err)
	}

	// 获取子部门列表中的 ID 和名称，用于查询子部门下的用户
	var deptIDs []int
	var deptNames []string
	for _, dept := range childDepts {
		deptIDs = append(deptIDs, dept.Id)
		deptNames = append(deptNames, dept.Name)
	}

	// 获取子部门下的用户列表
	users, err := model.GetUsersByDeptId(parentID)
	if err != nil {
		return "", fmt.Errorf("failed to get users by department ID: %v", err)
	}

	dept := map[string]interface{}{
		"childDepartments": childDepts,
		"employees":        users,
		"titleDepartments": deptNames,
	}

	jsonStr, err := util.ToJSONStr(dept)
	if err != nil {
		return "", fmt.Errorf("failed to convert result to JSON string: %v", err)
	}

	return jsonStr, nil
}

// GetUserByName 根据员工姓名分页查询员工信息
func GetUserByName(name string, pageNum, pageSize int) (string, error) {
	employees, err := model.GetUserByName(name, pageNum, pageSize)
	if err != nil {
		return "", fmt.Errorf("failed to get employees by name: %v", err)
	}

	// 获取员工总数
	total, err := model.GetUserByNameCount(name)
	if err != nil {
		return "", fmt.Errorf("failed to get employee count by name: %v", err)
	}

	result := map[string]interface{}{
		"pageNum":  pageNum,
		"pageSize": pageSize,
		"total":    total,
		"list":     employees,
	}

	jsonStr, err := util.ToJSONStr(result)
	if err != nil {
		return "", fmt.Errorf("failed to convert result to JSON string: %v", err)
	}

	return jsonStr, nil
}
