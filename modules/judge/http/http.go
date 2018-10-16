package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func SetupRoutes() {
	SetupCommonRoutes()
	SetupInfoRoutes()
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
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	SetupRoutes()

	log.Println("http listening", addr)
	log.Fatalln(s.ListenAndServe())
}
