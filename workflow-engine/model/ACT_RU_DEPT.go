package model

// 部门结构体 user_departments 包括【id, name, dialog_id, parent_id, owner_userid, created_at, updated_at】
type UserDepartments struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DialogId    string `json:"dialog_id"`
	ParentId    int    `json:"parent_id"`
	OwnerUserid string `json:"owner_userid"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 用户结构体 users 包括【userid, identity, department, nickname, profession, userimg, created_at, updated_at】
type Users struct {
	Userid     string `json:"userid"`
	Identity   string `json:"identity"`
	Department string `json:"department"`
	Nickname   string `json:"nickname"`
	Profession string `json:"profession"`
	Userimg    string `json:"userimg"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// GetAllDept 所有部门列表
func GetAllDept(parentId int) ([]*UserDepartments, error) {
	var datas []*UserDepartments
	dbModel := db
	if parentId != -1 {
		dbModel = dbModel.Where("parent_id=?", parentId)
	}
	err := dbModel.Find(&datas).Error
	return datas, err
}

// GetDeptByParentID 根据父级部门ID获取子部门列表
func GetDeptByParentID(parentID int) ([]*UserDepartments, error) {
	var datas []*UserDepartments
	err := db.Where("parent_id=?", parentID).Find(&datas).Error
	return datas, err
}

// GetDeptByID 根据部门ID获取部门信息
func GetDeptByID(deptID int) (*UserDepartments, error) {
	var data UserDepartments
	err := db.Where("id=?", deptID).Find(&data).Error
	return &data, err
}

// GetUsersByDept 根据部门名称获取用户列表
func GetUsersByDept(deptName string) ([]*Users, error) {
	var datas []*Users
	err := db.Where("department=?", deptName).Find(&datas).Error
	return datas, err
}
