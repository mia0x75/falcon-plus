package http

import (
	"github.com/gin-gonic/gin"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
)

// SetupProcRoutes 设置路由
func SetupProcRoutes() {
	// statistics
	routes.GET("/statistics/all", func(c *gin.Context) {
		cu.RenderDataJSON(c.Writer, proc.GetAll())
	})
}
