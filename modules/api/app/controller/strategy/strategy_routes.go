package strategy

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
	s := r.Group("/api/v1/strategy")
	s.Use(utils.AuthSessionMidd)
	s.GET("", GetStrategys)
	s.GET("/:sid", GetStrategy)
	s.POST("", CreateStrategy)
	s.PUT("", UpdateStrategy)
	s.DELETE("/:sid", DeleteStrategy)

	m := r.Group("/api/v1/metric")
	m.Use(utils.AuthSessionMidd)
	m.GET("default_list", MetricQuery)
}
