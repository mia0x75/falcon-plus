package http

import (
	"net"
	_ "net/http/pprof"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

var routes *gin.Engine

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func SetupRoutes() {
	SetupCommonRoutes()
	SetupProcRoutes()
	SetupIndexRoutes()
}

func Start() {
	go startHttpServer()
}

func startHttpServer() {
	if !g.Config().Http.Enabled {
		log.Println("http.Start warning, not enabled")
		return
	}
	if !cutils.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()

	SetupRoutes()

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}
	log.Printf("http listening %s", addr)
	go routes.Run(addr)
}
