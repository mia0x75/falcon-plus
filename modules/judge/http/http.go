package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

// SetupRoutes TODO:
func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcRoutes()
}

// Start TODO:
func Start() {
	go start()
}

func start() {
	if !g.Config().HTTP.Enabled {
		return
	}

	addr := g.Config().HTTP.Listen
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	SetupRoutes()

	log.Infof("[I] http listening %s", addr)
	log.Fatalf("[F] %v", s.ListenAndServe())
}
