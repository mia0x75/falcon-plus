package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// KernelMetrics TODO:
func KernelMetrics() (L []*cm.MetricValue) {
	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}
	L = append(L, GaugeValue("kernel.maxfiles", maxFiles))

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}
	L = append(L, GaugeValue("kernel.maxproc", maxProc))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}
	L = append(L, GaugeValue("kernel.files.allocated", allocateFiles))
	L = append(L, GaugeValue("kernel.files.left", maxFiles-allocateFiles))
	return
}
