package http

import (
	"encoding/json"
	"net/http"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	prpc "github.com/open-falcon/falcon-plus/modules/transfer/receiver/rpc"
)

// SetupAPIRoutes 设置路由
func SetupAPIRoutes() {
	http.HandleFunc("/api/push", func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var metrics []*cmodel.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		reply := &cmodel.TransferResponse{}
		prpc.RecvMetricValues(metrics, reply, "http")

		cutils.RenderDataJson(w, reply)
	})
}
