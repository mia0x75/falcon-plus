package http

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
)

func configIoStatRoutes() {
	http.HandleFunc("/page/diskio", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.IOStatsForPage())
	})
}
