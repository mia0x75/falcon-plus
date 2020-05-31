package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/toolkits/file"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

// SetupCommonRoutes 设置路由
func SetupCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", g.Version)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, g.Config())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			cu.RenderDataJSON(w, "ok")
		} else {
			cu.RenderDataJSON(w, "no privilege")
		}
	})
}
