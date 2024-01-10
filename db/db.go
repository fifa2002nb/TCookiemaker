package db

import (
	"TCookiemaker/config"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	MysqlDB *gorm.DB
)

func InitDB(options *config.Options) error {
	if nil == options {
		return errors.New("nil == options")
	}
	var err error
	if nil != MysqlDB {
		err = MysqlDB.Close()
	}
	if nil != err {
		return err
	}

	MysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", options.MysqlUser, options.MysqlPasswd, options.Mysqls[0], options.MysqlDB))
	if nil != err {
		return err
	}
	MysqlDB.DB().SetMaxIdleConns(10)
	MysqlDB.DB().SetMaxOpenConns(100)
	MysqlDB.LogMode(options.MysqlLog)
	return err
}
