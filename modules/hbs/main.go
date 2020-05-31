package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"github.com/open-falcon/falcon-plus/modules/hbs/http"
	"github.com/open-falcon/falcon-plus/modules/hbs/rpc"
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

	g.ParseConfig(*cfg)
	cu.InitLog(g.Config().Log.Level)
	if err := db.InitDB(); err != nil {
		os.Exit(0)
	}
	cache.Init()

	cache.DeleteStaleAgents()

	http.Start()
	rpc.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		db.DB.Close()
		os.Exit(0)
	}()

	select {}
}
