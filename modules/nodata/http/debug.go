package http

import (
	"net/http"
	"time"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/nodata/collector"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/sender"
)

func SetupDebugRoutes() {
	http.HandleFunc("/debug/collector/collect", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := collector.CollectDataOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		cu.RenderDataJSON(w, ret)
	})

	http.HandleFunc("/debug/config/sync", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := config.SyncNdConfigOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		cu.RenderDataJSON(w, ret)
	})

	http.HandleFunc("/debug/sender/send", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := sender.SendMockOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		cu.RenderDataJSON(w, ret)
	})
}
