package http

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/file"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

// SetupCommonRoutes 设置路由
func SetupCommonRoutes() {
	routes.GET("/health", func(c *gin.Context) {
		c.Writer.Write([]byte("ok\n"))
	})

	routes.GET("/version", func(c *gin.Context) {
		c.Writer.Write([]byte(fmt.Sprintf("%s\n", g.Version)))
	})

	routes.GET("/workdir", func(c *gin.Context) {
		c.Writer.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	routes.GET("/config", func(c *gin.Context) {
		cu.RenderDataJSON(c.Writer, g.Config())
	})

	routes.GET("/config/reload", func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			cu.RenderDataJSON(c.Writer, "ok")
		} else {
			cu.RenderDataJSON(c.Writer, "no privilege")
		}
	})
}
