package http

import (
	"net/http"
	"strings"

	"github.com/mia0x75/gopfc/metric"
)

// routers
func configProcRoutes() {
	http.HandleFunc("/pfc/proc/metrics/json", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, metric.RawMetrics())
	})
	http.HandleFunc("/pfc/proc/metrics/falcon", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, metric.FalconMetrics())
	})
	// url=/pfc/proc/metric/{json,falcon}
	http.HandleFunc("/pfc/proc/metrics/", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		urlParam := r.URL.Path[len("/pfc/proc/metrics/"):]
		args := strings.Split(urlParam, "/")
		argsLen := len(args)
		if argsLen != 2 {
			RenderJson(w, "")
			return
		}

		types := []string{}
		typeslice := strings.Split(args[0], ",")
		for _, t := range typeslice {
			nt := strings.TrimSpace(t)
			if nt != "" {
				types = append(types, nt)
			}
		}

		if args[1] == "json" {
			RenderJson(w, metric.RawMetric(types))
			return
		}
		if args[1] == "falcon" {
			RenderJson(w, metric.FalconMetric(types))
			return
		}
	})

	http.HandleFunc("/pfc/proc/metrics/size", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, metric.RawSizes())
	})
}
