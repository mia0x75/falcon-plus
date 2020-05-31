package http

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/toolkits/nux"

	cu "github.com/open-falcon/falcon-plus/common/utils"
)

// SetupSystemRoutes TODO:
func SetupSystemRoutes() {
	http.HandleFunc("/system/date", func(w http.ResponseWriter, req *http.Request) {
		cu.RenderDataJSON(w, time.Now().Format("2006-01-02 15:04:05"))
	})

	http.HandleFunc("/page/system/uptime", func(w http.ResponseWriter, req *http.Request) {
		days, hours, mins, err := nux.SystemUptime()
		cu.AutoRender(w, fmt.Sprintf("%d days %d hours %d minutes", days, hours, mins), err)
	})

	http.HandleFunc("/proc/system/uptime", func(w http.ResponseWriter, req *http.Request) {
		days, hours, mins, err := nux.SystemUptime()
		if err != nil {
			cu.RenderMsgJSON(w, err.Error())
			return
		}

		cu.RenderDataJSON(w, map[string]interface{}{
			"days":  days,
			"hours": hours,
			"mins":  mins,
		})
	})

	http.HandleFunc("/page/system/loadavg", func(w http.ResponseWriter, req *http.Request) {
		cpuNum := runtime.NumCPU()
		load, err := nux.LoadAvg()
		if err != nil {
			cu.RenderMsgJSON(w, err.Error())
			return
		}

		ret := [3][2]interface{}{
			{load.Avg1min, int64(load.Avg1min * 100.0 / float64(cpuNum))},
			{load.Avg5min, int64(load.Avg5min * 100.0 / float64(cpuNum))},
			{load.Avg15min, int64(load.Avg15min * 100.0 / float64(cpuNum))},
		}
		cu.RenderDataJSON(w, ret)
	})

	http.HandleFunc("/proc/system/loadavg", func(w http.ResponseWriter, req *http.Request) {
		data, err := nux.LoadAvg()
		cu.AutoRender(w, data, err)
	})
}
