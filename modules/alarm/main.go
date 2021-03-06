package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/cron"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/http"
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
	g.InitRedisConnPool()
	g.InitDB()
	cron.InitSenderWorker()

	http.Start()
	cron.ReadHighEvent()
	cron.ReadLowEvent()
	cron.CombineSMS()
	cron.CombineMail()
	cron.CombineIM()
	cron.ConsumeIM()
	cron.ConsumeSMS()
	cron.ConsumeMail()
	cron.CleanExpiredEvent()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
