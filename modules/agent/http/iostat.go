package http

import (
	"net/http"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
)

// SetupIoStatRoutes TODO:
func SetupIoStatRoutes() {
	http.HandleFunc("/page/diskio", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, funcs.IOStatsForPage())
	})
}
