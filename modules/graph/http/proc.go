package http

import (
	"github.com/gin-gonic/gin"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
)

func SetupProcRoutes() {
	// statistics
	routes.GET("/statistics/all", func(c *gin.Context) {
		cutils.RenderDataJson(c.Writer, proc.GetAll())
	})
}
