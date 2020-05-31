package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
	"github.com/open-falcon/falcon-plus/modules/api/g"
	"github.com/open-falcon/falcon-plus/modules/api/graph"
	"github.com/open-falcon/falcon-plus/modules/api/http"
	"github.com/open-falcon/falcon-plus/modules/api/rpc"
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
	if err := g.InitDB(); err != nil {
		log.Fatalf("[F] open db fail: %v", err)
		os.Exit(0)
	}
	cache.Init()
	rpc.Start()
	graph.Start()
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
