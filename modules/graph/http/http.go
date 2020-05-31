package http

import (
	"net"
	// Blank import
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

var routes *gin.Engine

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept TODO:
func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// SetupRoutes 设置路由
func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcRoutes()
	SetupIndexRoutes()
}

// Start 启动服务
func Start() {
	go start()
}

func start() {
	if !g.Config().HTTP.Enabled {
		log.Info("[I] http.Start warning, not enabled")
		return
	}
	if !cu.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()

	SetupRoutes()

	addr := g.Config().HTTP.Listen
	if addr == "" {
		return
	}
	log.Infof("[I] http listening %s", addr)
	go routes.Run(addr)
}
