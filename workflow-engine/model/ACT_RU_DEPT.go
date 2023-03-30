package model

// 部门结构体 user_departments 包括【id，name，dialog_id，parent_id，owner_userid，created_at，updated_at】
type UserDepartments struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DialogId    string `json:"dialog_id"`
	ParentId    int    `json:"parent_id"`
	OwnerUserid string `json:"owner_userid"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
