package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	yaag_gin "github.com/mia0x75/yaag/gin"
	"github.com/mia0x75/yaag/yaag"
	log "github.com/sirupsen/logrus"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
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
	go start()
}

func start() {
	if !cutils.IsDebug() {
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
	addr := g.Config().Listen
	log.Infof("[I] http listening %s", addr)

	SetupRoutes()

	go routes.Run(addr)
}
