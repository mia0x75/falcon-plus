package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/gateway/g"
	"github.com/open-falcon/falcon-plus/modules/gateway/http"
	"github.com/open-falcon/falcon-plus/modules/gateway/receiver"
	"github.com/open-falcon/falcon-plus/modules/gateway/sender"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)
	g.InitLog(g.Config().Log.Level)

	// receiver
	receiver.Start()
	// sender
	sender.Start()
	// http
	http.Start()

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
