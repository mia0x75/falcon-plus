package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	statsd "github.com/amalfra/gin-statsd/middleware"
	yaag_gin "github.com/betacraft/yaag/gin"
	"github.com/betacraft/yaag/yaag"
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller"
	"github.com/open-falcon/falcon-plus/modules/api/g"
	"github.com/open-falcon/falcon-plus/modules/api/graph"
	"github.com/open-falcon/falcon-plus/modules/api/rpc"
)

func initGraph() {
	// TODO:
	graph.Start(g.Config().Graphs.Cluster)
}

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()
	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitLog(g.Config().Log.Level)

	var err error
	err = g.InitLog(g.Config().Log.Level)
	if err != nil {
		log.Fatal(err)
	}
	err = g.InitDB()
	if err != nil {
		log.Fatalf("db conn failed with error %s", err.Error())
	}

	//TODO:
	if !g.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes := gin.Default()
	if false {
		routes.Use(statsd.New(statsd.Options{Port: 8089}))
	}
	if g.Config().GenDoc {
		yaag.Init(&yaag.Config{
			On:       true,
			DocTitle: "Gin",
			DocPath:  g.Config().GenDocPath,
			BaseUrls: map[string]string{"Production": "/api/v1", "Staging": "/api/v1"},
		})
		routes.Use(yaag_gin.Document())
	}
	initGraph()
	//start gin server
	log.Debugf("will start with port:%v", g.Config().Listen)
	go controller.StartGin(g.Config().Listen, routes)
	if g.Config().Rpc.Enabled {
		go rpc.Start()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()
	select {}
}
