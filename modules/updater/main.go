package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/sys"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/updater/cron"
	"github.com/open-falcon/falcon-plus/modules/updater/g"
	"github.com/open-falcon/falcon-plus/modules/updater/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
			"Version", g.Version,
			"Git commit", g.Git,
			"Compile", g.Compile,
			"Distro", g.Distro,
			"Kernel", g.Kernel,
			"Branch", g.Branch,
		)
		os.Exit(0)
	}

	if err := g.ParseConfig(*cfg); err != nil {
		log.Fatalf("[F] %v", err)
	}

	cutils.InitLog(g.Config().Log.Level)
	g.InitGlobalVariables()

	CheckDependency()

	http.Start()
	cron.Heartbeat()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}

func CheckDependency() {
	_, err := sys.CmdOut("wget", "--help")
	if err != nil {
		log.Fatal("[F] dependency wget not found")
	}

	_, err = sys.CmdOut("md5sum", "--help")
	if err != nil {
		log.Fatal("[F] dependency md5sum not found")
	}

	_, err = sys.CmdOut("tar", "--help")
	if err != nil {
		log.Fatal("[F] dependency tar not found")
	}
}
