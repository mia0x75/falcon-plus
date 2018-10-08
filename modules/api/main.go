package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	"github.com/spf13/viper"
)

func initGraph() {
	graph.Start(viper.GetStringMapString("graphs.cluster"))
}

func main() {
	cfgTmp := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()
	cfg := *cfgTmp
	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./api/config")
	cfg = strings.Replace(cfg, ".json", "", 1)
	viper.SetConfigName(cfg)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = g.InitLog(viper.GetString("log.level"))
	if err != nil {
		log.Fatal(err)
	}
	err = g.InitDB(viper.GetBool("db.db_bug"), viper.GetViper())
	if err != nil {
		log.Fatalf("db conn failed with error %s", err.Error())
	}

	//TODO:
	if viper.GetString("log.level") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	routes := gin.Default()
	if false {
		routes.Use(statsd.New(statsd.Options{Port: 8089}))
	}
	if viper.GetBool("gen_doc") {
		yaag.Init(&yaag.Config{
			On:       true,
			DocTitle: "Gin",
			DocPath:  viper.GetString("gen_doc_path"),
			BaseUrls: map[string]string{"Production": "/api/v1", "Staging": "/api/v1"},
		})
		routes.Use(yaag_gin.Document())
	}
	initGraph()
	//start gin server
	log.Debugf("will start with port:%v", viper.GetString("listen"))
	go controller.StartGin(viper.GetString("listen"), routes)
	if viper.GetBool("rpc.enabled") {
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
