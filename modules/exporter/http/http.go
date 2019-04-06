package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcHttpRoutes()
	SetupIndexHttpRoutes()
}

// start http server
func Start() {
	go start()
}

func start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	SetupRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Infof("[I] http listening %s", addr)
	log.Fatalf("[F] %v", s.ListenAndServe())
}
