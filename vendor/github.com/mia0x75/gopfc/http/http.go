package http

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/mia0x75/gopfc/g"
)

func Start() {
	configCommonRoutes()
	configProcRoutes()
	addr := g.Config().Http.Listen
	debug := g.Config().Debug
	if len(addr) >= 9 {
		s := &http.Server{
			Addr:           addr,
			MaxHeaderBytes: 1 << 30,
		}
		go func() {
			if debug {
				log.Println("[perfcounter] http server start, listening on", addr)
			}
			s.ListenAndServe()
			if debug {
				log.Println("[perfcounter] http server stop,", addr)
			}
		}()
	}
}

func isLocalReq(raddr string) bool {
	if strings.HasPrefix(raddr, "127.0.0.1") {
		return true
	}
	return false
}

// render
func RenderJson(w http.ResponseWriter, data interface{}) {
	renderJson(w, Dto{Msg: "success", Data: data})
}

func RenderString(w http.ResponseWriter, msg string) {
	renderJson(w, map[string]string{"msg": msg})
}

func renderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

// common http return
type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
