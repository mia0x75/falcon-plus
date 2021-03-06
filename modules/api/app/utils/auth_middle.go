package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

func AuthSessionMidd(c *gin.Context) {
	auth, err := h.SessionChecking(c)
	if !g.Config().SkipAuth {
		if err != nil || auth != true {
			log.Errorf("[E] %v, auth: %v", err, auth)
			c.Set("auth", auth)
			h.JSONR(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}
	}
	c.Set("auth", auth)
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", context.Request.Header.Get("Origin"))
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, X-Falcon-Token")
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
