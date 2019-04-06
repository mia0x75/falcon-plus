package http

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
)

func SetupCommonRoutes() {
	// compatible with anteye
	routes.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	//compatible with open-falcon v0.1
	routes.GET("/counter/migrate", func(c *gin.Context) {
		cnt := rrdtool.GetCounter()
		log.Debugf("[D] migrating counter: %s", cnt)
		c.JSON(200, gin.H{"msg": "ok", "counter": cnt})
	})
}
