package model

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// ProcInstHistory
type ProcInstHistory struct {
	ProcInst
}

// StartHistoryByMyself 查询我发起的流程
func StartHistoryByMyself(userID, company string, pageIndex, pageSize int) ([]*ProcInstHistory, int, error) {
	maps := map[string]interface{}{
		"start_user_id": userID,
		"company":       company,
	}
	return findProcInstsHistory(maps, pageIndex, pageSize)
}
func findProcInstsHistory(maps map[string]interface{}, pageIndex, pageSize int) ([]*ProcInstHistory, int, error) {
	var datas []*ProcInstHistory
	var count int
	selectDatas := func(in chan<- error, wg *sync.WaitGroup) {
		go func() {
			err := db.Where(maps).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("start_time desc").Find(&datas).Error
			in <- err
			wg.Done()
		}()
	}
	selectCount := func(in chan<- error, wg *sync.WaitGroup) {
		err := db.Model(&ProcInstHistory{}).Where(maps).Count(&count).Error
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

// FindProcHistory 查询历史纪录
func FindProcHistory(userID, procName, company string, sort string, pageIndex, pageSize int) ([]*ProcInstHistory, int, error) {
	var datas []*ProcInstHistory
	var count int
	var err1 error
	var order string
	// 判段排序
	if sort == "asc" {
		order = "start_time asc"
	} else {
		order = "start_time desc"
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	errStream := make(chan error, 2)

	go func() {
		defer wg.Done()

		dbQuery := db.Where(fmt.Sprintf("id in (select distinct proc_inst_id from %sidentitylink_history where company=? and user_id=?)", conf.DbPrefix), company, userID)
		if procName != "" {
			dbQuery = dbQuery.Where("proc_def_name = ?", procName)
		}
		err := dbQuery.Offset((pageIndex - 1) * pageSize).Limit(pageSize).
			Order(order).Find(&datas).Error
		errStream <- err
	}()

	go func() {
		defer wg.Done()

		dbQuery := db.Model(&ProcInstHistory{})
		dbQuery = dbQuery.Where(fmt.Sprintf("id in (select distinct proc_inst_id from %sidentitylink_history where company=? and user_id=?)", conf.DbPrefix), company, userID)
		if procName != "" {
			dbQuery = dbQuery.Where("proc_def_name = ?", procName)
		}
		err := dbQuery.Count(&count).Error
		errStream <- err
	}()

	wg.Wait()
	close(errStream)

	for err := range errStream {
		if err != nil {
			err1 = err
		}
	}
	return datas, count, err1
}

// SaveProcInstHistory
func SaveProcInstHistory(p *ProcInst) error {
	return db.Table(conf.DbPrefix + "proc_inst_history").Create(p).Error
}

// DelProcInstHistoryByID
func DelProcInstHistoryByID(id int) error {
	return db.Where("id=?", id).Delete(&ProcInstHistory{}).Error
}

// SaveProcInstHistoryTx
func SaveProcInstHistoryTx(p *ProcInst, tx *gorm.DB) error {
	return tx.Table(conf.DbPrefix + "proc_inst_history").Create(p).Error
}

// FindProcHistoryNotify 查询抄送我的历史纪录
func FindProcHistoryNotify(userID, procName, company string, groups []string, sort string, pageIndex, pageSize int) ([]*ProcInstHistory, int, error) {
	var datas []*ProcInstHistory
	var count int
	var sql string
	var values []interface{}
	var order string
	// 判断排序
	if sort == "asc" {
		order = "start_time asc"
	} else {
		order = "start_time desc"
	}
	if len(groups) != 0 {
		var placeholders []string
		for range groups {
			placeholders = append(placeholders, "?")
		}
		sql = "SELECT proc_inst_id FROM %sidentitylink_history WHERE type='notifier' AND company=? AND (FIND_IN_SET(?, user_id) OR `group` IN (" + strings.Join(placeholders, ",") + "))"
		values = append(values, company, userID)
		for _, group := range groups {
			values = append(values, group)
		}
	} else {
		sql = "SELECT proc_inst_id FROM %sidentitylink_history WHERE type='notifier' AND company=? AND FIND_IN_SET(?, user_id)"
		values = append(values, company, userID)
	}
	if procName != "" {
		sql += " AND proc_def_name = ?"
		values = append(values, procName)
	}
	sql = fmt.Sprintf(sql, conf.DbPrefix)
	err := db.Where("id in ("+sql+")", values...).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order(order).Find(&datas).Error
	if err != nil {
		return datas, count, err
	}
	err = db.Model(&ProcInstHistory{}).Where("id in ("+sql+")", values...).Count(&count).Error
	if err != nil {
		return nil, count, err
	}
	return datas, count, err
}
