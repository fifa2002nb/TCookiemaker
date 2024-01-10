package config

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/widuu/goini"
	"strconv"
	"strings"
)

// 配置项
type Options struct {
	// for common
	ProxyPort int
	Workers   int
	StoreFile string
	// for mysql query
	Mysqls      []string //ip1:port1...
	MysqlDB     string
	MysqlUser   string
	MysqlPasswd string
	MysqlLog    bool
}

// 解析配置文件
func ParseConf(c *cli.Context) (*Options, error) {
	if c.IsSet("configure") || c.IsSet("C") {
		options := &Options{}
		var conf *goini.Config
		var err error
		if c.IsSet("configure") {
			conf = goini.SetConfig(c.String("configure"))
		} else {
			conf = goini.SetConfig(c.String("C"))
		}
		// mysql configures
		ms := conf.GetValue("mysql", "mysqls")
		options.Mysqls = strings.Split(ms, ",")
		if 0 == len(options.Mysqls) {
			return nil, errors.New(fmt.Sprintf("mysql is required  to start a mysql job. See '%s start --help'.", c.App.Name))
		}
		options.MysqlDB = conf.GetValue("mysql", "mysqldb")
		if "" == options.MysqlDB {
			options.MysqlDB = "uaq_speed"
		}
		options.MysqlUser = conf.GetValue("mysql", "mysqluser")
		options.MysqlPasswd = conf.GetValue("mysql", "mysqlpasswd")
		mysqllog := conf.GetValue("mysql", "mysqlLog")
		if "true" == mysqllog {
			options.MysqlLog = true
		} else {
			options.MysqlLog = false
		}
		// main configure
		wkers := conf.GetValue("main", "workers")
		if options.Workers, err = strconv.Atoi(wkers); nil != err {
			log.Errorf("%v", err)
			options.Workers = 10
		}
		pPort := conf.GetValue("main", "proxyPort")
		if options.ProxyPort, err = strconv.Atoi(pPort); nil != err {
			log.Errorf("%v", err)
			options.ProxyPort = 8899
		}
		options.StoreFile = conf.GetValue("main", "storeFile")
		if "" == options.StoreFile {
			options.StoreFile = "./db.file"
		}
		return options, nil
	} else {
		return nil, errors.New(fmt.Sprintf("configure is required to run a job. See '%s start --help'.", c.App.Name))
	}
}
