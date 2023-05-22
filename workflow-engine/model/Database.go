package model

import (
	"fmt"
	"log"
	"strconv"

	config "workflow/workflow-config"

	"github.com/jinzhu/gorm"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

// Model 其它数据结构的公共部分
type Model struct {
	ID int `gorm:"primary_key" json:"id,omitempty"`
}

// 配置
var conf = *config.Config

// Setup 初始化一个db连接
func Setup() {
	var err error
	log.Println("数据库初始化！！")
	db, err = gorm.Open(conf.DbType, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName))
	if err != nil {
		log.Fatalf("数据库连接失败 err: %v", err)
	}
	// 启用Logger，显示详细日志
	mode, _ := strconv.ParseBool(conf.DbLogMode)

	db.LogMode(mode)

	db.SingularTable(true) //全局设置表名不可以为复数形式
	// db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	idle, err := strconv.Atoi(conf.DbMaxIdleConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(idle)
	open, err := strconv.Atoi(conf.DbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(open)

	//添加表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		// 判断defaultTableName等于Department时返回defaultTableName
		if defaultTableName == "user_departments" || defaultTableName == "users" {
			return "pre_" + defaultTableName
		}
		return conf.DbPrefix + defaultTableName
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(&Procdef{}).
		AutoMigrate(&Execution{}).
		AutoMigrate(&Task{}).
		AutoMigrate(&ProcInst{}).
		AutoMigrate(&ProcMsgs{}).
		AutoMigrate(&Identitylink{}).
		AutoMigrate(&ExecutionHistory{}).
		AutoMigrate(&IdentitylinkHistory{}).
		AutoMigrate(&ProcInstHistory{}).
		AutoMigrate(&TaskHistory{}).
		AutoMigrate(&ProcdefHistory{})

	// 添加索引
	db.Model(&Procdef{}).AddIndex("idx_id", "id")
	db.Model(&ProcInst{}).AddIndex("idx_id", "id")
	db.Model(&Execution{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&Identitylink{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&Task{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	//---------------------历史纪录------------------------------
	db.Model(&ProcInstHistory{}).AddIndex("idx_id", "id")
	db.Model(&ExecutionHistory{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst_history(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&IdentitylinkHistory{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst_history(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&TaskHistory{}).AddIndex("idx_id", "id")

}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

// GetDB getdb
func GetDB() *gorm.DB {
	return db
}

// GetTx GetTx
func GetTx() *gorm.DB {
	return db.Begin()
}
