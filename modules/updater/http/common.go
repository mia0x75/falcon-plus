package http

import (
	"net/http"
	"strings"

	"github.com/toolkits/file"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/updater/g"
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
		cutils.RenderDataJson(w, file.SelfDir())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			err := g.ParseConfig(g.ConfigFile)
			cutils.AutoRender(w, g.Config(), err)
		} else {
			w.Write([]byte("no privilege"))
		}
	})
}
