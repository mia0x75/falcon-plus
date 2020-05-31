package http

import (
	"net/http"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

// SetupProcRoutes 设置路由
func SetupProcRoutes() {
	// Counter
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, proc.GetAll())
	})
}
