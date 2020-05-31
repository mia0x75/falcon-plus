package http

import (
	"net/http"
	"strconv"
	"strings"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	"github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

// SetupProcRoutes 设置路由
func SetupProcRoutes() {
	// counter
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, proc.GetAll())
	})

	// step
	http.HandleFunc("/proc/step", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, map[string]interface{}{"min_step": sender.MinStep})
	})

	// trace
	http.HandleFunc("/trace/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/trace/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		tags := make(map[string]string)
		if argsLen > 2 {
			tagVals := strings.Split(args[2], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		proc.RecvDataTrace.SetPK(cu.PK(endpoint, metric, tags))
		cu.RenderDataJSON(w, proc.RecvDataTrace.GetAllTraced())
	})

	// filter
	http.HandleFunc("/filter/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/filter/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		opt := args[2]

		threadholdStr := args[3]
		threadhold, err := strconv.ParseFloat(threadholdStr, 64)
		if err != nil {
			cu.RenderDataJSON(w, "bad threadhold")
			return
		}

		tags := make(map[string]string)
		if argsLen > 4 {
			tagVals := strings.Split(args[4], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}

		err = proc.RecvDataFilter.SetFilter(cu.PK(endpoint, metric, tags), opt, threadhold)
		if err != nil {
			cu.RenderDataJSON(w, err.Error())
			return
		}

		cu.RenderDataJSON(w, proc.RecvDataFilter.GetAllFiltered())
	})
}
