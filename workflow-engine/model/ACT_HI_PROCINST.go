package model

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// ProcInst 流程实例
type ProcInst struct {
	Model
	ProcDefID     int    `gorm:"not null;default:0;comment:'流程定义ID'" json:"procDefId"`
	ProcDefName   string `gorm:"comment:'流程定义名'" json:"procDefName"`
	Title         string `gorm:"comment:'标题'" json:"title,omitempty"`
	DepartmentId  int    `gorm:"comment:'用户部门ID'" json:"departmentId"`
	Department    string `gorm:"comment:'用户部门'" json:"department"`
	Company       string `gorm:"comment:'用户公司'" json:"company"`
	NodeID        string `gorm:"comment:'当前节点'" json:"nodeID"`
	Candidate     string `gorm:"comment:'审批人'" json:"candidate"`
	TaskID        int    `gorm:"comment:'当前任务'" json:"taskID"`
	StartTime     string `gorm:"comment:'开始时间'" json:"startTime"`
	EndTime       string `gorm:"comment:'结束时间'" json:"endTime"`
	Duration      int64  `gorm:"comment:'持续时间'" json:"duration"`
	StartUserID   string `gorm:"comment:'开始用户ID'" json:"startUserId"`
	StartUserName string `gorm:"comment:'开始用户名'" json:"startUserName"`
	IsFinished    bool   `gorm:"default:false;comment:'是否完成'" json:"isFinished"`
	Var           string `gorm:"size:4000;comment:'执行流程的附加参数'" json:"var"`
	State         int    `gorm:"not null;default:0;comment:'当前状态: 0待审批，1审批中，2通过，3拒绝，4撤回'" json:"state"`
	LatestComment string `gorm:"size:500;comment:'最新评论'" json:"latestComment"`
}

// GroupsNotNull 候选组
func GroupsNotNull(groups []string, sql string) func(db *gorm.DB) *gorm.DB {
	if len(groups) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Or("candidate in (?) and "+sql, groups)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// DepartmentsNotNull 分管部门
func DepartmentsNotNull(departments []string, sql string) func(db *gorm.DB) *gorm.DB {
	if len(departments) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Or("department in (?) and candidate=? and "+sql, departments, IdentityTypes[MANAGER])
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// StartByMyself 我发起的流程
func StartByMyself(userID, company string, pageIndex, pageSize int) ([]*ProcInst, int, error) {
	maps := map[string]interface{}{
		"start_user_id": userID,
		"company":       company,
	}
	return findProcInsts(maps, pageIndex, pageSize)
}

// FindProcInstByID 根据ID查询流程实例
func FindProcInstByID(id int) (*ProcInst, error) {
	var data = ProcInst{}
	err := db.Where("id=?", id).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// FindProcNotify 查询抄送我的流程
func FindProcNotify(userID, company string, groups []string, pageIndex, pageSize int) ([]*ProcInst, int, error) {
	var datas []*ProcInst
	var count int
	var sql string
	if len(groups) != 0 {
		var s []string
		for _, val := range groups {
			s = append(s, "\""+val+"\"")
		}
		sql = "select proc_inst_id from %sidentitylink i where i.type='notifier' and i.company='" + company + "' and (find_in_set('" + userID + "',i.user_id) or i.group in (" + strings.Join(s, ",") + "))"
	} else {
		sql = "select proc_inst_id from %sidentitylink i where i.type='notifier' and i.company='" + company + "' and find_in_set('" + userID + "',i.user_id)"
	}
	sql = fmt.Sprintf(sql, conf.DbPrefix)
	err := db.Where("id in (" + sql + ")").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("start_time desc").Find(&datas).Error
	if err != nil {
		return datas, count, err
	}
	err = db.Model(&ProcInst{}).Where("id in (" + sql + ")").Count(&count).Error
	if err != nil {
		return nil, count, err
	}
	return datas, count, err
}
func findProcInsts(maps map[string]interface{}, pageIndex, pageSize int) ([]*ProcInst, int, error) {
	var datas []*ProcInst
	var count int
	selectDatas := func(in chan<- error, wg *sync.WaitGroup) {
		go func() {
			err := db.Where(maps).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("start_time desc").Find(&datas).Error
			in <- err
			wg.Done()
		}()
	}
	selectCount := func(in chan<- error, wg *sync.WaitGroup) {
		err := db.Model(&ProcInst{}).Where(maps).Count(&count).Error
		in <- err
		wg.Done()
	}
	var err1 error
	var wg sync.WaitGroup
	numberOfRoutine := 2
	wg.Add(numberOfRoutine)
	errStream := make(chan error, numberOfRoutine)
	// defer fmt.Println("close channel")
	selectDatas(errStream, &wg)
	selectCount(errStream, &wg)
	wg.Wait()
	defer close(errStream) // 关闭通道
	for i := 0; i < numberOfRoutine; i++ {
		// log.Printf("send: %v", <-errStream)
		if err := <-errStream; err != nil {
			err1 = err
		}
	}
	// fmt.Println("结束")
	return datas, count, err1
}

// FindProcInsts 分页查询
func FindProcInsts(userID, procName, company string, groups, departments []string, pageIndex, pageSize int) ([]*ProcInst, int, error) {
	var datas []*ProcInst
	var count int
	var sql = " company='" + company + "' and is_finished=0 "
	if len(procName) > 0 {
		sql += "and proc_def_name='" + procName + "'"
	}
	// fmt.Println(sql)
	selectDatas := func(in chan<- error, wg *sync.WaitGroup) {
		go func() {
			err := db.Scopes(GroupsNotNull(groups, sql), DepartmentsNotNull(departments, sql)).
				Or("candidate=? and "+sql, userID).
				Offset((pageIndex - 1) * pageSize).Limit(pageSize).
				Order("start_time desc").
				Find(&datas).Error
			in <- err
			wg.Done()
		}()
	}
	selectCount := func(in chan<- error, wg *sync.WaitGroup) {
		go func() {
			err := db.Scopes(GroupsNotNull(groups, sql), DepartmentsNotNull(departments, sql)).Model(&ProcInst{}).Or("candidate=? and "+sql, userID).Count(&count).Error
			in <- err
			wg.Done()
		}()
	}
	var err1 error
	var wg sync.WaitGroup
	numberOfRoutine := 2
	wg.Add(numberOfRoutine)
	errStream := make(chan error, numberOfRoutine)
	// defer fmt.Println("close channel")
	selectDatas(errStream, &wg)
	selectCount(errStream, &wg)
	wg.Wait()
	defer close(errStream) // 关闭通道

	for i := 0; i < numberOfRoutine; i++ {
		// log.Printf("send: %v", <-errStream)
		if err := <-errStream; err != nil {
			err1 = err
		}
	}
	// fmt.Println("结束")
	return datas, count, err1
}

// Save
func (p *ProcInst) Save() (int, error) {
	err := db.Create(p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

// SaveTx 事务保存
func (p *ProcInst) SaveTx(tx *gorm.DB) (int, error) {
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return p.ID, nil
}

// DelProcInstByID 删除流程实例
func DelProcInstByID(id int) error {
	return db.Where("id=?", id).Delete(&ProcInst{}).Error
}

// DelProcInstByIDTx 事务
func DelProcInstByIDTx(id int, tx *gorm.DB) error {
	return tx.Where("id=?", id).Delete(&ProcInst{}).Error
}

// UpdateTx 更新
func (p *ProcInst) UpdateTx(tx *gorm.DB) error {
	return tx.Model(&ProcInst{}).Updates(p).Error
}

// FindFinishedProc 查询已完成的流程
func FindFinishedProc() ([]*ProcInst, error) {
	var datas []*ProcInst
	err := db.Where("is_finished=1").Find(&datas).Error
	return datas, err
}
