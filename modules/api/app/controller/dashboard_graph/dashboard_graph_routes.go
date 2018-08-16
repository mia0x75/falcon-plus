package dashboard_graph

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed
const TMP_GRAPH_FILED_DELIMITER = "|"

func Routes(r *gin.Engine) {
	db = config.Con()
	g := r.Group("/api/v1/dashboard")
	g.Use(utils.AuthSessionMidd)
	g.POST("/tmpgraph", DashboardTmpGraphCreate)
	g.GET("/tmpgraph/:id", DashboardTmpGraphQuery)
	g.POST("/graph", DashboardGraphCreate)
	g.PUT("/graph/:id", DashboardGraphUpdate)
	g.GET("/graph/:id", DashboardGraphGet)
	g.DELETE("/graph/:id", DashboardGraphDelete)
	g.GET("/graphs/screen/:screen_id", DashboardGraphGetsByScreenID)
}
