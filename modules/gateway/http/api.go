package http

import (
	"encoding/json"
	"net/http"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	trpc "github.com/open-falcon/falcon-plus/modules/gateway/receiver/rpc"
)

// SetupAPIRoutes 设置路由
func SetupAPIRoutes() {
	http.HandleFunc("/api/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*cmodel.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		reply := &cmodel.TransferResponse{}
		trpc.RecvMetricValues(metrics, reply, "http")

		cutils.RenderDataJson(w, reply)
	})
}
