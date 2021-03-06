package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cu "github.com/open-falcon/falcon-plus/common/utils"
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
	flag.Parse()

	fmt.Printf(g.Banner, g.Module)
	fmt.Println()
	fmt.Println()
	fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
		"Version", g.Version,
		"Git commit", g.Git,
		"Compile", g.Compile,
		"Distro", g.Distro,
		"Kernel", g.Kernel,
		"Branch", g.Branch,
	)

	if *version {
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)
	cu.InitLog(g.Config().Log.Level)
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
