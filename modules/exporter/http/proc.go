package http

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

func SetupProcHttpRoutes() {
	// counter
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		ret := make(map[string]interface{})
		ret["msg"] = "success"
		ret["data"] = proc.GetAll()
		RenderDataJson(w, ret)
	})
}
