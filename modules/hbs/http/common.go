package http

import (
	"net/http"
	"strings"

	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"github.com/toolkits/file"

	cu "github.com/open-falcon/falcon-plus/common/utils"
)

// SetupCommonRoutes 设置路由
func SetupCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.Version))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, file.SelfDir())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			cu.RenderDataJSON(w, g.Config())
		} else {
			w.Write([]byte("no privilege"))
		}
	})
}
