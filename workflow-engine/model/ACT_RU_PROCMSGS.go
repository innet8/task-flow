package model

// ProcMsgs 流程信息关联表
type ProcMsgs struct {
	Model
	ProcInstID int    `gorm:"comment:'流程实例ID'" json:"procInstID"`
	Userid     int    `gorm:"comment:'会员ID'" json:"userid"`
	MsgID      int    `gorm:"comment:'消息ID'" json:"msgId"`
	CreatedAt  string `gorm:"autoCreateTime;not null;comment:'创建时间'" json:"createdAt"`
	UpdatedAt  string `gorm:"autoCreateTime;not null;comment:'更新时间'" json:"updatedAt"`
}

// 插入数据
func InsertProcMsgs(procInstID, userID, msgID int) error {
	var data = ProcMsgs{
		ProcInstID: procInstID,
		Userid:     userID,
		MsgID:      msgID,
	}
	err := db.Create(&data).Error
	return err
}

// 根据流程实例ID和用户ID获取消息ID
func GetMsgIDByProcInstIDAndUserID(procInstID, userID int) (int, error) {
	var data ProcMsgs
	err := db.Where("proc_inst_id=? and userid=?", procInstID, userID).Find(&data).Error
	if err != nil {
		return 0, err
	}
	// 如果查找data为空或record not found，返回0
	if data.ID == 0 {
		return 0, nil
	}

	return data.MsgID, nil
}

// 根据流程实例ID获取列表
func GetProcMsgsByProcInstID(procInstID int) ([]ProcMsgs, error) {
	var data []ProcMsgs
	err := db.Where("proc_inst_id=?", procInstID).Find(&data).Error
	return data, err
}
