package http

import (
	"net/http"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/index"
)

// SetupIndexRoutes 设置路由
func SetupIndexRoutes() {
	http.HandleFunc("/index/delete", func(w http.ResponseWriter, r *http.Request) {
		index.DeleteIndex()
		cu.RenderDataJSON(w, "ok")
	})

	http.HandleFunc("/index/updateAll", func(w http.ResponseWriter, r *http.Request) {
		index.UpdateAllIndex()
		cu.RenderDataJSON(w, "ok")
	})
}
