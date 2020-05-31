package http

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

var routes *gin.Engine

// SetupRoutes 设置路由
func SetupRoutes() {
	SetupCommonRoutes()
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

	if !cu.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()
	// Start gin server
	log.Infof("[I] http listening %s", addr)

	SetupRoutes()

	go routes.Run(addr)
}
