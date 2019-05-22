package http

import (
	"fmt"
	"net/http"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
)

// SetupProcRoutes 设置路由
func SetupProcRoutes() {
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, cache.ExpressionCache.Get())
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		cutils.RenderDataJson(w, cache.Agents.Keys())
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*cmodel.Host, len(cache.MonitoredHosts.Get()))
		for k, v := range cache.MonitoredHosts.Get() {
			data[fmt.Sprint(k)] = v
		}
		cutils.RenderDataJson(w, data)
	})

	http.HandleFunc("/strategies", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*cmodel.Strategy, len(cache.Strategies.GetMap()))
		for k, v := range cache.Strategies.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		cutils.RenderDataJson(w, data)
	})

	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*cmodel.Template, len(cache.TemplateCache.GetMap()))
		for k, v := range cache.TemplateCache.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		cutils.RenderDataJson(w, data)
	})

	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		cutils.RenderDataJson(w, cache.GetPlugins(hostname))
	})
}
