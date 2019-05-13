package http

import (
	"fmt"
	"net/http"

	"github.com/toolkits/core"
	"github.com/toolkits/nux"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
)

// SetupDfRoutes TODO:
func SetupDfRoutes() {
	http.HandleFunc("/page/df", func(w http.ResponseWriter, r *http.Request) {
		mountPoints, err := nux.ListMountPoint()
		if err != nil {
			cutils.RenderMsgJson(w, err.Error())
			return
		}

		ret := make([][]interface{}, 0)
		for idx := range mountPoints {
			var du *nux.DeviceUsage
			du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
			if err == nil {
				ret = append(ret,
					[]interface{}{
						du.FsSpec,
						core.ReadableSize(float64(du.BlocksAll)),
						core.ReadableSize(float64(du.BlocksUsed)),
						core.ReadableSize(float64(du.BlocksFree)),
						fmt.Sprintf("%.1f%%", du.BlocksUsedPercent),
						du.FsFile,
						core.ReadableSize(float64(du.InodesAll)),
						core.ReadableSize(float64(du.InodesUsed)),
						core.ReadableSize(float64(du.InodesFree)),
						fmt.Sprintf("%.1f%%", du.InodesUsedPercent),
						du.FsVfstype,
					})
			}
		}

		cutils.RenderDataJson(w, ret)
	})
}
