package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/nodata/collector"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
	"github.com/open-falcon/falcon-plus/modules/nodata/http"
	"github.com/open-falcon/falcon-plus/modules/nodata/judge"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	fmt.Printf(g.Banner, "Nodata")
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
	cutils.InitLog(g.Config().Log.Level)
	// proc
	g.StartProc()

	// config
	config.Start()
	// collector
	collector.Start()
	// judge
	judge.Start()
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
