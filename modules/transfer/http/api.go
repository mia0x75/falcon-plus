package http

import (
	"encoding/json"
	"net/http"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	prpc "github.com/open-falcon/falcon-plus/modules/transfer/receiver/rpc"
)

func api_push_datapoints(rw http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(rw, "blank body", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var metrics []*cmodel.MetricValue
	err := decoder.Decode(&metrics)
	if err != nil {
		http.Error(rw, "decode error", http.StatusBadRequest)
		return
	}

	reply := &cmodel.TransferResponse{}
	prpc.RecvMetricValues(metrics, reply, "http")

	cutils.RenderDataJson(rw, reply)
}

func SetupApiRoutes() {
	http.HandleFunc("/api/push", api_push_datapoints)
}
