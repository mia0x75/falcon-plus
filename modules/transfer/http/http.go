package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

func SetupRouters() {
	SetupCommonRoutes()
	SetupProcHttpRoutes()
	SetupDebugHttpRoutes()
	SetupApiRoutes()
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

	SetupRouters()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.startHttpServer ok, listening", addr)
	log.Fatalln(s.ListenAndServe())
}
