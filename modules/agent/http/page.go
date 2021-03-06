package http

import (
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// SetupPageRoutes TODO:
func SetupPageRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			f := filepath.Join(g.Config().HTTP.Root, "/public", r.URL.Path, "index.html")
			log.Debugf("[D] %s", f)
			if !file.IsExist(f) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Config().HTTP.Root, "/public"))).ServeHTTP(w, r)
		f := filepath.Join(g.Config().HTTP.Root, "/public")
		log.Debugf("[D] %s", f)
		http.FileServer(http.Dir(f)).ServeHTTP(w, r)
	})
}
