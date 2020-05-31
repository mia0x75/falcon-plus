package helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/g"
)

const (
	HTTPOK                  = http.StatusOK
	HTTPBadRequest          = http.StatusBadRequest
	HTTPUnauthorized        = http.StatusUnauthorized
	HTTPExpectationFailed   = http.StatusExpectationFailed
	HTTPInternalServerError = http.StatusInternalServerError
	StandardErrorMessage    = "An error occurred while %s. More information about this error may be available in the server error log."
	// Please contact the server administrator and inform them of the time the error occurred and anything you might have done that may have caused the error.
)

// InternelError
func InternelError(c *gin.Context, op string, err error) {
	log.Errorf("[E] An error occurred while %s, error: %v", op, err)
	JSONR(c, HTTPInternalServerError, fmt.Sprintf(StandardErrorMessage, op))
}

// cu.RespJSON ?
type RespJSON struct {
	Error string `json:"error,omitempty"`
	Msg   string `json:"message,omitempty"`
}

func JSONR(c *gin.Context, arg ...interface{}) (werror error) {
	var (
		code int
		msg  interface{}
	)
	if len(arg) == 1 {
		code = http.StatusOK
		msg = arg[0]
	} else {
		code = arg[0].(int)
		msg = arg[1]
	}
	need_doc := g.Config().GenDoc
	var body interface{}
	defer func() {
		if need_doc {
			ds, _ := json.Marshal(body)
			bodys := string(ds)
			log.Debugf("[D] body: %v, bodys: %v ", body, bodys)
			c.Set("body_doc", bodys)
		}
	}()
	if code == 200 {
		switch msg.(type) {
		case string:
			body = RespJSON{
				Msg: msg.(string),
			}
			c.JSON(http.StatusOK, body)
		default:
			c.JSON(http.StatusOK, msg)
			body = msg
		}
	} else {
		switch msg.(type) {
		case string:
			body = RespJSON{
				Error: msg.(string),
			}
		case error:
			body = RespJSON{
				Error: msg.(error).Error(),
			}
		default:
			body = RespJSON{
				Error: "system type error. please ask admin for help",
			}
		}
		c.JSON(code, body)
	}
	return
}
