package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func Start() {
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
	routes := gin.Default()
	routes.GET("/version", Version)
	routes.GET("/health", Health)
	routes.GET("/workdir", Workdir)
	log.Println("http listening", addr)
	go routes.Run(addr)
}
