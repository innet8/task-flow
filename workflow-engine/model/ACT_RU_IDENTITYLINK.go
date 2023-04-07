package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Identitylink 用户组同任务的关系
type Identitylink struct {
	Model
	Group      string `gorm:"comment:'集团'" json:"group,omitempty"`
	Type       string `gorm:"comment:'类型'" json:"type,omitempty"`
	UserID     string `gorm:"comment:'用户ID'" json:"userid,omitempty"`
	UserName   string `gorm:"comment:'用户名'" json:"username,omitempty"`
	TaskID     int    `gorm:"comment:'任务ID'" json:"taskID,omitempty"`
	Step       int    `gorm:"comment:'一步'" json:"step"`
	ProcInstID int    `gorm:"comment:'Proc机构ID'" json:"procInstID,omitempty"`
	Company    string `gorm:"comment:'公司'" json:"company,omitempty"`
	Comment    string `gorm:"comment:'评论'" json:"comment,omitempty"`
	State      int    `gorm:"default:0;comment:'状态:0待处理,1通过,2拒绝,3撤回'" json:"state"`
}

// IdentityType 类型
type IdentityType int

const (
	// CANDIDATE 候选
	CANDIDATE IdentityType = iota
	// PARTICIPANT 参与人
	PARTICIPANT
	// MANAGER 上级领导
	MANAGER
	// NOTIFIER 抄送人
	NOTIFIER
)

// IdentityTypes 类型
var IdentityTypes = [...]string{CANDIDATE: "candidate", PARTICIPANT: "participant", MANAGER: "主管", NOTIFIER: "notifier"}

// SaveTx
func (i *Identitylink) SaveTx(tx *gorm.DB) error {
	// if len(i.Company) == 0 {
	// 	return errors.New("Identitylink表的company字段不能为空！！")
	// }
	err := tx.Create(i).Error
	return err
}

// DelCandidateByProcInstID 删除历史候选人
func DelCandidateByProcInstID(procInstID int, tx *gorm.DB) error {
	return tx.Where("proc_inst_id=? and type=?", procInstID, IdentityTypes[CANDIDATE]).Delete(&Identitylink{}).Error
}

// ExistsNotifierByProcInstIDAndGroup 抄送人是否已经存在
func ExistsNotifierByProcInstIDAndGroup(procInstID int, group string) (bool, error) {
	var count int
	err := db.Model(&Identitylink{}).Where(fmt.Sprintf("%sidentitylink.proc_inst_id=? and %sidentitylink.group=? and %sidentitylink.type=?", conf.DbPrefix, conf.DbPrefix, conf.DbPrefix), procInstID, group, IdentityTypes[NOTIFIER]).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IfParticipantByTaskID 针对指定任务判断用户是否已经审批过了
func IfParticipantByTaskID(userID, company string, taskID int) (bool, error) {
	var count int
	err := db.Model(&Identitylink{}).Where("user_id=? and company=? and task_id=?", userID, company, taskID).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// FindParticipantByProcInstID 查询参与审批的人
func FindParticipantByProcInstID(procInstID int) ([]*Identitylink, error) {
	var datas []*Identitylink
	err := db.Select("id,user_id,user_name,step,comment,state").Where("proc_inst_id=? and type=?", procInstID, IdentityTypes[PARTICIPANT]).Order("id asc").Find(&datas).Error
	return datas, err
}

// GetIdentitylinkInfoByTaskID 根据taskID获取信息
func GetIdentitylinkInfoByTaskID(taskID int) (*Identitylink, error) {
	var datas Identitylink
	err := db.Model(&Users{}).Where("task_id=?", taskID).Find(&datas).Error
	return &datas, err
}
