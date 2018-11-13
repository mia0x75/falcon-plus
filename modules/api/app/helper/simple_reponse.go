package helper

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

type RespJson struct {
	Error string `json:"error,omitempty"`
	Msg   string `json:"message,omitempty"`
}

func JSONR(c *gin.Context, arg ...interface{}) (werror error) {
	var (
		wcode int
		msg   interface{}
	)
	if len(arg) == 1 {
		wcode = http.StatusOK
		msg = arg[0]
	} else {
		wcode = arg[0].(int)
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
	if wcode == 200 {
		switch msg.(type) {
		case string:
			body = RespJson{Msg: msg.(string)}
			c.JSON(http.StatusOK, body)
		default:
			c.JSON(http.StatusOK, msg)
			body = msg
		}
	} else {
		switch msg.(type) {
		case string:
			body = RespJson{Error: msg.(string)}
			c.JSON(wcode, body)
		case error:
			body = RespJson{Error: msg.(error).Error()}
			c.JSON(wcode, body)
		default:
			body = RespJson{Error: "system type error. please ask admin for help"}
			c.JSON(wcode, body)
		}
	}
	return
}
