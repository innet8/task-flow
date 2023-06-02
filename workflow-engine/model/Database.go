package model

import (
	"fmt"
	"log"
	"strconv"

	config "workflow/workflow-config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Model 是其他数据结构的公共部分
type Model struct {
	ID int `gorm:"primary_key" json:"id,omitempty"`
}

// Configuration
var conf = *config.Config

// Setup 初始化一个db连接
func Setup() {
	var err error
	log.Println("正在初始化数据库...")
	db, err = gorm.Open(conf.DbType, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName))
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	// 启用Logger，显示详细日志
	mode, _ := strconv.ParseBool(conf.DbLogMode)
	db.LogMode(mode)

	// 全局设置表名不可以为复数形式
	db.SingularTable(true)

	// 设置最大空闲连接数
	maxIdleConns, err := strconv.Atoi(conf.DbMaxIdleConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(maxIdleConns)

	// 设置最大打开连接数
	maxOpenConns, err := strconv.Atoi(conf.DbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(maxOpenConns)

	// 添加表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == "user_departments" || defaultTableName == "users" {
			return "pre_" + defaultTableName
		}
		return conf.DbPrefix + defaultTableName
	}

	// 自动迁移表
	db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(&Procdef{}, &Execution{}, &Task{}, &ProcInst{}, &ProcMsgs{}, &Identitylink{}, &ExecutionHistory{}, &IdentitylinkHistory{}, &ProcInstHistory{}, &TaskHistory{}, &ProcdefHistory{})

	// 添加索引
	db.Model(&Procdef{}).AddIndex("idx_id", "id")
	db.Model(&ProcInst{}).AddIndex("idx_id", "id")
	db.Model(&Execution{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&Identitylink{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&Task{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")

	// 添加历史记录索引
	db.Model(&ProcInstHistory{}).AddIndex("idx_id", "id")
	db.Model(&ExecutionHistory{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst_history(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&IdentitylinkHistory{}).AddForeignKey("proc_inst_id", conf.DbPrefix+"proc_inst_history(id)", "CASCADE", "RESTRICT").AddIndex("idx_id", "id")
	db.Model(&TaskHistory{}).AddIndex("idx_id", "id")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	defer db.Close()
}

// GetDB 返回数据库连接
func GetDB() *gorm.DB {
	return db
}

// GetTx 返回一个事务
func GetTx() *gorm.DB {
	return db.Begin()
}
