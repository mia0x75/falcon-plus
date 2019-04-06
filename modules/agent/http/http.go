package http

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"

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
	go start()
}

func start() {
	addr := g.Config().Http.Listen
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
