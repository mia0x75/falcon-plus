package http

import (
	"net/http"

	"github.com/toolkits/nux"

	cu "github.com/open-falcon/falcon-plus/common/utils"
)

// SetupMemoryRoutes TODO:
func SetupMemoryRoutes() {
	http.HandleFunc("/page/memory", func(w http.ResponseWriter, r *http.Request) {
		mem, err := nux.MemInfo()
		if err != nil {
			cu.RenderMsgJSON(w, err.Error())
			return
		}

		memFree := mem.MemFree + mem.Buffers + mem.Cached
		// if mem.MemAvailable > 0 {
		// 	memFree = mem.MemAvailable
		// }
		memUsed := mem.MemTotal - memFree
		var t uint64 = 1024 * 1024
		cu.RenderDataJSON(w, []interface{}{mem.MemTotal / t, memUsed / t, memFree / t})
	})

	http.HandleFunc("/proc/memory", func(w http.ResponseWriter, r *http.Request) {
		mem, err := nux.MemInfo()
		if err != nil {
			cu.RenderMsgJSON(w, err.Error())
			return
		}

		memFree := mem.MemFree + mem.Buffers + mem.Cached
		// if mem.MemAvailable > 0 {
		// 	memFree = mem.MemAvailable
		// }
		memUsed := mem.MemTotal - memFree

		cu.RenderDataJSON(w, map[string]interface{}{
			"total": mem.MemTotal,
			"free":  memFree,
			"used":  memUsed,
		})
	})
}
