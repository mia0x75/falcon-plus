package http

import (
	"fmt"
	"net/http"

	"github.com/mia0x75/gopfc/g"
)

func configCommonRoutes() {
	http.HandleFunc("/pfc/health", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/pfc/version", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		w.Write([]byte(fmt.Sprintf("%s\n", g.VERSION)))
	})

	http.HandleFunc("/pfc/config", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, g.Config())
	})
}
