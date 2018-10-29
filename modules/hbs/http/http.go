package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcRoutes()
}

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
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	SetupRoutes()

	log.Printf("http listening %s", addr)
	log.Fatalln(s.ListenAndServe())
}
