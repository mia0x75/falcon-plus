package http

import (
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
)

func SetupProcRoutes() {
	// statistics
	routes.GET("/statistics/all", func(c *gin.Context) {
		ret := make(map[string]interface{})
		ret["msg"] = "success"
		ret["data"] = proc.GetAll()
		JSONR(c, 200, ret)
	})
}
