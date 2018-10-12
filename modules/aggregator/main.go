package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/sdk/sender"
	"github.com/open-falcon/falcon-plus/modules/aggregator/cron"
	"github.com/open-falcon/falcon-plus/modules/aggregator/db"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	"github.com/open-falcon/falcon-plus/modules/aggregator/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	db.InitDB()
	g.InitLog(g.Config().Log.Level)

	http.Start()
	cron.UpdateItems()

	// sdk configuration
	sender.Debug = g.IsDebug()
	sender.PostPushUrl = g.Config().Api.Agent
	// sender
	sender.StartSender()

	log.Infoln("service ready ...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
