package http

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	yaag_gin "github.com/betacraft/yaag/gin"
	"github.com/betacraft/yaag/yaag"
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/alarm"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/dashboard_graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/dashboard_screen"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/expression"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/host"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/mockcfg"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/strategy"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/template"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var routes *gin.Engine

func SetupRoutes() {
	routes.Use(utils.CORS())
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	graph.Routes(routes)
	uic.Routes(routes)
	template.Routes(routes)
	strategy.Routes(routes)
	host.Routes(routes)
	expression.Routes(routes)
	mockcfg.Routes(routes)
	dashboard_graph.Routes(routes)
	dashboard_screen.Routes(routes)
	alarm.Routes(routes)
	SetupCommonRoutes()
}

func Start() {
	go startHttpServer()
}

func startHttpServer() {
	if !g.IsDebug() {
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
	//start gin server
	log.Debugf("will start with port:%v", g.Config().Listen)

	SetupRoutes()

	go routes.Run(g.Config().Listen)
}
