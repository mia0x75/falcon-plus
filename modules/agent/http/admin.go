package http

import (
	"net/http"
	"os"
	"time"

	"github.com/toolkits/file"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

// SetupAdminRoutes TODO:
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
			cu.RenderDataJSON(w, g.Config())
		} else {
			w.Write([]byte("no privilege"))
		}
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, file.SelfDir())
	})

	http.HandleFunc("/ips", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, hbs.TrustableIps())
	})
}
