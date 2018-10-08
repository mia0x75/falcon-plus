package dashboard_screen

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var db g.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = g.Con()
	g := r.Group("/api/v1/dashboard")
	g.Use(utils.AuthSessionMidd)
	g.POST("/screen", ScreenCreate)
	g.GET("/screen/:screen_id", ScreenGet)
	g.GET("/screens/pid/:pid", ScreenGetsByPid)
	g.GET("/screens", ScreenGetsAll)
	g.DELETE("/screen/:screen_id", ScreenDelete)
	g.PUT("/screen/:screen_id", ScreenUpdate)
}
