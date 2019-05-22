package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pfc "github.com/mia0x75/gopfc"
	pfcg "github.com/mia0x75/gopfc/g"
	log "github.com/sirupsen/logrus"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/api"
	"github.com/open-falcon/falcon-plus/modules/graph/cron"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/http"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
)

func start_signal(pid int, cfg *g.GlobalConfig) {
	sigs := make(chan os.Signal, 1)
	log.Infof("[I] %d register signal notify", pid)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Infof("[I] recv: %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Info("[I] graceful shut down")
			if cfg.RPC.Enabled {
				api.Close_chan <- 1
				<-api.Close_done_chan
			}
			log.Info("[I] rpc stop ok")

			rrdtool.Out_done_chan <- 1
			rrdtool.FlushAll(true)
			log.Info("[I] rrdtool stop ok")

			g.DB.Close()

			log.Infof("[I] %d exit", pid)
			os.Exit(0)
		}
	}
}

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
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

	// global config
	g.ParseConfig(*cfg)
	cutils.InitLog(g.Config().Log.Level)
	// init db
	if err := g.InitDB(); err != nil {
		os.Exit(0)
	}
	if g.Config().PerfCounter != nil {
		log.Debugf("[D] pfc config: %v", g.Config().PerfCounter)
		pfcg.PFCWithConfig(g.Config().PerfCounter)
		pfc.Start()
	}

	// rrdtool init
	rrdtool.InitChannel()
	// rrdtool before api for disable loopback connection
	rrdtool.Start()
	// start api
	api.Start()
	// start indexing
	index.Start()
	// start http server
	http.Start()
	cron.CleanCache()

	start_signal(os.Getpid(), g.Config())
}
