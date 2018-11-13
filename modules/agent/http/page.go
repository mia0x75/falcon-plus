package http

import (
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
)

func SetupPageRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			f := filepath.Join(g.Config().Http.Root, "/public", r.URL.Path, "index.html")
			log.Debugf("[D] %s", f)
			if !file.IsExist(f) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Config().Http.Root, "/public"))).ServeHTTP(w, r)
		f := filepath.Join(g.Config().Http.Root, "/public")
		log.Debugf("[D] %s", f)
		http.FileServer(http.Dir(f)).ServeHTTP(w, r)
	})
}
