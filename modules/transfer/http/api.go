package http

import (
	"encoding/json"
	"net/http"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
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
		var metrics []*cm.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		reply := &cm.TransferResponse{}
		prpc.RecvMetricValues(metrics, reply, "http")

		cu.RenderDataJSON(w, reply)
	})
}
