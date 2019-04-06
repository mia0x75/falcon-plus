package http

import (
	"net/http"
	"os"
	"time"

	"github.com/toolkits/file"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

func SetupAdminRoutes() {
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		if hbs.IsTrustable(r.RemoteAddr) {
			w.Write([]byte("exiting..."))
			go func() {
				time.Sleep(time.Second)
				os.Exit(0)
			}()
		} else {
			w.Write([]byte("no privilege"))
		}
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if hbs.IsTrustable(r.RemoteAddr) {
			g.ParseConfig(g.ConfigFile)
			cutils.RenderDataJson(w, g.Config())
		} else {
			w.Write([]byte("no privilege"))
		}
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, file.SelfDir())
	})

	http.HandleFunc("/ips", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, hbs.TrustableIps())
	})
}
