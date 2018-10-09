package http

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	statsd "github.com/amalfra/gin-statsd/middleware"
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

func Start() {
	if !g.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Printf("%+v", g.Config().Statsd)
	routes := gin.Default()
	if g.Config().Statsd.Enabled {
		log.Println("start gin-statsd ...")
		routes.Use(statsd.New(statsd.Options{Port: g.Config().Statsd.Port}))
	}
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
	go routes.Run(g.Config().Listen)
}
