package expression

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	g := r.Group("/api/v1/expression")
	g.Use(utils.AuthSessionMidd)
	g.GET("", GetExpressionList)
	g.GET("/:eid", GetExpression)
	g.POST("", CreateExrpession)
	g.PUT("", UpdateExrpession)
	g.DELETE("/:eid", DeleteExpression)
}
