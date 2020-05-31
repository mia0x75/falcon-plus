package http

import (
	"net/http"

	"github.com/toolkits/nux"
	"github.com/toolkits/sys"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// SetupKernelRoutes TODO:
func SetupKernelRoutes() {
	http.HandleFunc("/proc/kernel/hostname", func(w http.ResponseWriter, r *http.Request) {
		data, err := g.Hostname()
		cu.AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxproc", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxProc()
		cu.AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxfiles", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxFiles()
		cu.AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/version", func(w http.ResponseWriter, r *http.Request) {
		data, err := sys.CmdOutNoLn("uname", "-r")
		cu.AutoRender(w, data, err)
	})

}
