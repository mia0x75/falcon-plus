package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/http"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	"github.com/open-falcon/falcon-plus/modules/transfer/receiver"
	"github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	fmt.Printf(g.Banner, "Transfer")
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
	proc.Start()
	sender.Start()
	receiver.Start()
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
