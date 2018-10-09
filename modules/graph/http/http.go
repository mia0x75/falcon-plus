package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var router *gin.Engine

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func Start() {
	if !g.Config().Http.Enabled {
		log.Println("http.Start warning, not enabled")
		return
	}
	if !g.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()
	configCommonRoutes()
	configProcRoutes()
	configIndexRoutes()

	router.GET("/api/v2/counter/migrate", func(c *gin.Context) {
		counter := rrdtool.GetCounterV2()
		log.Debug("migrating counter v2:", fmt.Sprintf("%+v", counter))
		c.JSON(200, counter)
	})

	//compatible with open-falcon v0.1
	router.GET("/counter/migrate", func(c *gin.Context) {
		cnt := rrdtool.GetCounter()
		log.Debug("migrating counter:", cnt)
		c.JSON(200, gin.H{"msg": "ok", "counter": cnt})
	})

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}
	go router.Run(addr)
}
