package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/updater/g"
)

func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcRoutes()
}

func Start() {
	go startHttpServer()
}

func startHttpServer() {
	if !g.Config().Http.Enabled {
		return
	}
	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	SetupRoutes()

	log.Println("http listening", addr)
	err := http.ListenAndServe(addr, nil)
	log.Fatalln(err)
}
