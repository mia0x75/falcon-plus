package funcs

import (
	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/nux"
)

func KernelMetrics() (L []*cmodel.MetricValue) {
	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Println(err)
		return
	}
	L = append(L, GaugeValue("kernel.maxfiles", maxFiles))

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Println(err)
		return
	}
	L = append(L, GaugeValue("kernel.maxproc", maxProc))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Println(err)
		return
	}
	L = append(L, GaugeValue("kernel.files.allocated", allocateFiles))
	L = append(L, GaugeValue("kernel.files.left", maxFiles-allocateFiles))
	return
}
