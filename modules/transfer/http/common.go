package http

import (
	"fmt"
	"net/http"
	"strings"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/toolkits/file"
)

func SetupCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", g.VERSION)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, g.Config())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			cutils.RenderDataJson(w, "ok")
		} else {
			cutils.RenderDataJson(w, "no privilege")
		}
	})
}
