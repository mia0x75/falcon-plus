package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func SetupRoutes() {
	SetupAdminRoutes()
	SetupCpuRoutes()
	SetupDfRoutes()
	SetupHealthRoutes()
	SetupIoStatRoutes()
	SetupKernelRoutes()
	SetupMemoryRoutes()
	SetupPageRoutes()
	SetupPluginRoutes()
	SetupPushRoutes()
	SetupRunRoutes()
	SetupSystemRoutes()
}

func Start() {
	go startHttpServer()
}

func startHttpServer() {
	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	SetupRoutes()

	log.Println("listening", addr)
	log.Fatalln(s.ListenAndServe())
}
