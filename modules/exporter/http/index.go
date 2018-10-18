package http

import (
	"net/http"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/index"
)

func SetupIndexHttpRoutes() {
	http.HandleFunc("/index/delete", func(w http.ResponseWriter, r *http.Request) {
		index.DeleteIndex()
		cutils.RenderDataJson(w, "ok")
	})

	http.HandleFunc("/index/updateAll", func(w http.ResponseWriter, r *http.Request) {
		index.UpdateAllIndex()
		cutils.RenderDataJson(w, "ok")
	})
}
