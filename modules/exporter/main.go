package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/falcon-plus/modules/exporter/collector"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
	"github.com/open-falcon/falcon-plus/modules/exporter/http"
	"github.com/open-falcon/falcon-plus/modules/exporter/index"
	"github.com/open-falcon/falcon-plus/modules/exporter/monitor"
	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if *versionGit {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)
	g.InitLog(g.Config().Log.Level)
	// proc
	proc.Start()

	// graph index
	index.Start()
	// collector
	collector.Start()
	// monitor
	monitor.Start()
	// http
	http.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
