package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	yaag_gin "github.com/mia0x75/yaag/gin"
	"github.com/mia0x75/yaag/yaag"
	log "github.com/sirupsen/logrus"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var routes *gin.Engine

func SetupRoutes() {
	routes.Use(utils.CORS())
	routes.Use(utils.AuthSessionMidd)
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	controller.Routes(routes)
	SetupCommonRoutes()
}

func Start() {
	go start()
}

func start() {
	if !cu.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()
	if g.Config().GenDoc {
		yaag.Init(&yaag.Config{
			On:       true,
			DocTitle: "Gin",
			DocPath:  g.Config().GenDocPath,
			BaseUrls: map[string]string{"Production": "/api/v1", "Staging": "/api/v1"},
		})
		routes.Use(yaag_gin.Document())
	}
	// Start gin server
	addr := g.Config().Listen
	log.Infof("[I] http listening %s", addr)

	SetupRoutes()

	go routes.Run(addr)
}
