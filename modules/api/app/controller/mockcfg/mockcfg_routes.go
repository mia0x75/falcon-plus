package mockcfg

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
	g := r.Group("/api/v1/nodata")
	g.Use(utils.AuthSessionMidd)
	g.GET("", GetNoDataList)
	g.GET("/:nid", GetNoData)
	g.POST("/", CreateNoData)
	g.PUT("/", UpdateNoData)
	g.DELETE("/:nid", DeleteNoData)
}
