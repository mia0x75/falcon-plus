package http

import (
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

var routes *gin.Engine

func SetupRoutes() {
	SetupCommonRoutes()
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

	if !g.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()

	SetupRoutes()

	go routes.Run(addr)
}
