package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/falcon-plus/common/sdk/sender"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
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

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Printf(g.Banner, "Aggregator")
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

	g.ParseConfig(*cfg)
	cutils.InitLog(g.Config().Log.Level)
	if err := db.InitDB(); err != nil {
		os.Exit(0)
	}

	http.Start()
	cron.UpdateItems()

	// sdk configuration
	sender.Debug = cutils.IsDebug()
	sender.PostPushUrl = g.Config().API.Agent
	// sender
	sender.StartSender()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
