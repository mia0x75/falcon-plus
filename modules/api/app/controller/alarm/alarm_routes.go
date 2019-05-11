package alarm

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var db g.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

// Routes 路由表
func Routes(r *gin.Engine) {
	db = g.Con()
	g := r.Group("/api/v1/alarm")
	g.Use(utils.AuthSessionMidd)
	g.POST("/eventcases", Lists)
	g.GET("/eventcases", Lists)
	g.POST("/events", EventsGet)
	g.GET("/events", EventsGet)
	g.POST("/event_note", AddNotesToAlarm)
	g.GET("/event_note", GetNotesOfAlarm)
}
