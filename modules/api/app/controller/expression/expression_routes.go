package expression

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
	g := r.Group("/api/v1/expression")
	g.Use(utils.AuthSessionMidd)
	g.GET("", GetExpressionList)
	g.GET("/:eid", GetExpression)
	g.POST("", CreateExrpession)
	g.PUT("", UpdateExrpession)
	g.DELETE("/:eid", DeleteExpression)
}
