package model

import (
	"workflow/util"

	"github.com/jinzhu/gorm"
)

// Task 流程任务表
type Task struct {
	Model
	NodeID        string `gorm:"comment:'当前执行流所在的节点ID'" json:"nodeId"`
	Step          int    `gorm:"comment:'一步'" json:"step"`
	ProcInstID    int    `gorm:"comment:'流程实例id'" json:"procInstID"`
	Assignee      string `gorm:"comment:'处理人'" json:"assignee"`
	CreateTime    string `gorm:"autoCreateTime;not null;comment:'创建时间'" json:"createTime"`
	ClaimTime     string `gorm:"not null;comment:'处理时间'" json:"claimTime"`
	MemberCount   int8   `gorm:"default:1;comment:'还未审批的用户数，等于0代表会签已经全部审批结束，默认值为1'" json:"memberCount"`
	UnCompleteNum int8   `gorm:"default:1;comment:'完整的数字'" json:"unCompleteNum"`
	AgreeNum      int8   `gorm:"comment:'审批通过数'" json:"agreeNum"`
	ActType       string `gorm:"default:'or';comment:'and 为会签，or为或签，默认为or'" json:"actType"`
	IsFinished    bool   `gorm:"default:false;comment:'是否完成'" json:"isFinished"`
}

// CreateTask 新建任务
func (t *Task) CreateTask() (int, error) {
	err := db.Create(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}

// UpdateTask 更新任务
func (t *Task) UpdateTaskTx(tx *gorm.DB) error {
	err := tx.Model(&Task{}).Updates(t).Error
	return err
}

// GetTaskByID 根据任务ID获取任务
func GetTaskByID(id int) (*Task, error) {
	var t Task
	err := db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTasksByProcInstID 根据流程实例ID获取任务列表
func GetTasksByProcInstID(procInstID int) ([]*Task, error) {
	var tasks []*Task
	err := db.Where("proc_inst_id=?", procInstID).Order("claim_time desc").Find(&tasks).Error
	if err != nil || len(tasks) == 0 {
		var taskHistories []*TaskHistory
		err = db.Where("proc_inst_id=?", procInstID).Order("claim_time desc").Find(&taskHistories).Error
		if err != nil {
			return nil, err
		}
		strJSON, _ := util.ToJSONStr(&taskHistories)
		util.Str2Struct(strJSON, &tasks)
	}
	return tasks, nil
}

// GetLastFinishedTaskByProcInstID 根据流程实例ID获取上一个已完成的任务
func GetLastFinishedTaskByProcInstID(procInstID int) (*Task, error) {
	var t Task
	err := db.Where("proc_inst_id=? and is_finished=1", procInstID).Order("claim_time desc").First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// CreateTaskTx 开启事务并新建任务
func (t *Task) CreateTaskTx(tx *gorm.DB) (int, error) {
	err := tx.Create(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}

// DeleteTask 根据任务ID删除任务
func DeleteTask(id int) error {
	err := db.Where("id=?", id).Delete(&Task{}).Error
	return err
}
