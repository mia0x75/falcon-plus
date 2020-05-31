package http

import (
	"fmt"
	"net/http"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
)

// SetupProcRoutes 设置路由
func SetupProcRoutes() {
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, cache.ExpressionCache.Get())
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		cu.RenderDataJSON(w, cache.Agents.Keys())
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*cm.Host, len(cache.MonitoredHosts.Get()))
		for k, v := range cache.MonitoredHosts.Get() {
			data[fmt.Sprint(k)] = v
		}
		cu.RenderDataJSON(w, data)
	})

	http.HandleFunc("/strategies", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*cm.Strategy, len(cache.Strategies.GetMap()))
		for k, v := range cache.Strategies.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		cu.RenderDataJSON(w, data)
	})

	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*cm.Template, len(cache.TemplateCache.GetMap()))
		for k, v := range cache.TemplateCache.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		cu.RenderDataJSON(w, data)
	})

	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		cu.RenderDataJSON(w, cache.GetPlugins(hostname))
	})
}
