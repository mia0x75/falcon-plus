package http

import (
	"net/http"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

// SetupProcRoutes 设置路由
func SetupProcRoutes() {
	// counter
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, proc.GetAll())
	})
}
