package base

import (
	_ "TCookiemaker/cert"
	"TCookiemaker/config"
	"TCookiemaker/db"
	"TCookiemaker/spider"
	"TCookiemaker/utils/store"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
	"os/signal"
)

var (
	insertChan chan string
)

func init() {
	insertChan = make(chan string, 10000)
}

func Start(c *cli.Context) {
	var (
		err     error
		options *config.Options
	)
	options, err = config.ParseConf(c)
	if nil != err {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	if err = store.Load(options.StoreFile); nil != err {
		log.Fatal(err.Error())
		os.Exit(2)
	}
	go store.Run()

	if err = db.InitDB(options); nil != err {
		log.Fatal(err.Error())
		os.Exit(3)
	}
	spider.InitConfig(&spider.Config{
		Verbose:  true, // Open to see detail logs
		Compress: false,
	})
	spider.Regist(&CustomProcessor{})
	go func() {
		spider.Run(fmt.Sprintf("%d", options.ProxyPort))
	}()
	go DBInserter()
	waitingForExit()
}

//Just to implement Output Method of interface{} Processor
type CustomProcessor struct {
	spider.BaseProcessor
}

func (c *CustomProcessor) Output() {
	log.Infof("[output] %s => %s", c.DetailResult().Url, c.DetailResult().Cookie)
}

func DBInserter() {
	for cookie := range insertChan {
		log.Info("insert => %v", cookie)
	}
}

func waitingForExit() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	killing := false
	for range sc {
		if killing {
			log.Info("Second interrupt: exiting")
			os.Exit(1)
		}
		killing = true
		go func() {
			log.Info("Interrupt: closing down...")
			log.Info("done")
			os.Exit(1)
		}()
	}
}
