package http

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/toolkits/file"
)

func SetupCommonRoutes() {
	routes.GET("/health", func(c *gin.Context) {
		c.Writer.Write([]byte("ok\n"))
	})

	routes.GET("/version", func(c *gin.Context) {
		c.Writer.Write([]byte(fmt.Sprintf("%s\n", g.VERSION)))
	})

	routes.GET("/workdir", func(c *gin.Context) {
		c.Writer.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	routes.GET("/config", func(c *gin.Context) {
		cutils.RenderDataJson(c.Writer, g.Config())
	})

	routes.GET("/config/reload", func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			cutils.RenderDataJson(c.Writer, "ok")
		} else {
			cutils.RenderDataJson(c.Writer, "no privilege")
		}
	})
}
