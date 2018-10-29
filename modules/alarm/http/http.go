package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

var routes *gin.Engine

func SetupRoutes() {
	SetupCommonRoutes()
}

func Start() {
	go start()
}

func start() {
	if !g.Config().Http.Enabled {
		return
	}
	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	if !cutils.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()
	//start gin server
	log.Printf("http listening %s", addr)

	SetupRoutes()

	go routes.Run(addr)
}
