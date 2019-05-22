package http

import (
	"net/http"
	_ "net/http/pprof" // TODO:

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/updater/g"
)

// SetupRoutes 设置路由
func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcRoutes()
}

// Start 启动服务
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

	SetupRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Infof("[I] http listening %s", addr)
	log.Fatalf("[F] %v", s.ListenAndServe())
}
