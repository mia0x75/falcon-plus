package http

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/toolkits/nux"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
)

// SetupCPURoutes TODO:
func SetupCPURoutes() {
	http.HandleFunc("/proc/cpu/num", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, runtime.NumCPU())
	})

	http.HandleFunc("/proc/cpu/mhz", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.CpuMHz()
		cu.AutoRender(w, data, err)
	})

	http.HandleFunc("/page/cpu/usage", func(w http.ResponseWriter, r *http.Request) {
		cpuUsages, _, prepared := funcs.CPUUsagesSummary()
		if !prepared {
			cu.RenderMsgJSON(w, "not prepared")
			return
		}

		item := [10]string{
			fmt.Sprintf("%.1f%%", cpuUsages.Idle),
			fmt.Sprintf("%.1f%%", cpuUsages.Busy),
			fmt.Sprintf("%.1f%%", cpuUsages.User),
			fmt.Sprintf("%.1f%%", cpuUsages.Nice),
			fmt.Sprintf("%.1f%%", cpuUsages.System),
			fmt.Sprintf("%.1f%%", cpuUsages.Iowait),
			fmt.Sprintf("%.1f%%", cpuUsages.Irq),
			fmt.Sprintf("%.1f%%", cpuUsages.SoftIrq),
			fmt.Sprintf("%.1f%%", cpuUsages.Steal),
			fmt.Sprintf("%.1f%%", cpuUsages.Guest),
		}

		cu.RenderDataJSON(w, [][10]string{item})
	})

	http.HandleFunc("/proc/cpu/usage", func(w http.ResponseWriter, r *http.Request) {
		cpuUsages, _, prepared := funcs.CPUUsagesSummary()
		if !prepared {
			cu.RenderMsgJSON(w, "not prepared")
			return
		}

		cu.RenderDataJSON(w, map[string]interface{}{
			"idle":    fmt.Sprintf("%.1f%%", cpuUsages.Idle),
			"busy":    fmt.Sprintf("%.1f%%", cpuUsages.Busy),
			"user":    fmt.Sprintf("%.1f%%", cpuUsages.User),
			"nice":    fmt.Sprintf("%.1f%%", cpuUsages.Nice),
			"system":  fmt.Sprintf("%.1f%%", cpuUsages.System),
			"iowait":  fmt.Sprintf("%.1f%%", cpuUsages.Iowait),
			"irq":     fmt.Sprintf("%.1f%%", cpuUsages.Irq),
			"softirq": fmt.Sprintf("%.1f%%", cpuUsages.SoftIrq),
			"steal":   fmt.Sprintf("%.1f%%", cpuUsages.Steal),
			"guest":   fmt.Sprintf("%.1f%%", cpuUsages.Guest),
		})
	})
}
