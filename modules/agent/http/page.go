package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			if !file.IsExist(filepath.Join(g.Config().Http.Root, "/public", r.URL.Path, "index.html")) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Config().Http.Root, "/public"))).ServeHTTP(w, r)
	})

}
