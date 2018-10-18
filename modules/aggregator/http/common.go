package http

import (
	"net/http"
	"strings"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	"github.com/toolkits/file"
)

func SetupCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.VERSION))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, file.SelfDir())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			cutils.RenderJson(w, g.Config())
		} else {
			w.Write([]byte("no privilege"))
		}
	})
}
