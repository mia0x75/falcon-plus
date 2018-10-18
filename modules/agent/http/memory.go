package http

import (
	"net/http"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/toolkits/nux"
)

func SetupMemoryRoutes() {
	http.HandleFunc("/page/memory", func(w http.ResponseWriter, r *http.Request) {
		mem, err := nux.MemInfo()
		if err != nil {
			cutils.RenderMsgJson(w, err.Error())
			return
		}

		memFree := mem.MemFree + mem.Buffers + mem.Cached
		//if mem.MemAvailable > 0 {
		//	memFree = mem.MemAvailable
		//}
		memUsed := mem.MemTotal - memFree
		var t uint64 = 1024 * 1024
		cutils.RenderDataJson(w, []interface{}{mem.MemTotal / t, memUsed / t, memFree / t})
	})

	http.HandleFunc("/proc/memory", func(w http.ResponseWriter, r *http.Request) {
		mem, err := nux.MemInfo()
		if err != nil {
			cutils.RenderMsgJson(w, err.Error())
			return
		}

		memFree := mem.MemFree + mem.Buffers + mem.Cached
		//if mem.MemAvailable > 0 {
		//	memFree = mem.MemAvailable
		//}
		memUsed := mem.MemTotal - memFree

		cutils.RenderDataJson(w, map[string]interface{}{
			"total": mem.MemTotal,
			"free":  memFree,
			"used":  memUsed,
		})
	})
}
