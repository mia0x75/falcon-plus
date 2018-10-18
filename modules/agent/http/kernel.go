package http

import (
	"net/http"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
	"github.com/toolkits/sys"
)

func SetupKernelRoutes() {
	http.HandleFunc("/proc/kernel/hostname", func(w http.ResponseWriter, r *http.Request) {
		data, err := g.Hostname()
		cutils.AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxproc", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxProc()
		cutils.AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxfiles", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxFiles()
		cutils.AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/version", func(w http.ResponseWriter, r *http.Request) {
		data, err := sys.CmdOutNoLn("uname", "-r")
		cutils.AutoRender(w, data, err)
	})

}
