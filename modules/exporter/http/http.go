package http

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcHttpRoutes()
	SetupIndexHttpRoutes()
}

// start http server
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

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http:startHttpServer, ok, listening ", addr)
	log.Fatalln(s.ListenAndServe())
}
