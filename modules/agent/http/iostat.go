package http

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
)

func SetupIoStatRoutes() {
	http.HandleFunc("/page/diskio", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.IOStatsForPage())
	})
}
