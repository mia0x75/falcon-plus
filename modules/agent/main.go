package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/cron"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

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

	g.ParseConfig(*cfg)

	cu.InitLog(g.Config().Log.Level)
	g.InitRootDir()
	g.InitLocalIP()
	g.InitRPCClients()

	funcs.BuildMappers()

	cron.InitDataHistory()
	cron.ReportAgentStatus()
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()
	cron.Collect()

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
