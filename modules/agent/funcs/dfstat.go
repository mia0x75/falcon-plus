package funcs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func DeviceMetrics() (L []*cmodel.MetricValue) {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Errorf("[E] collect device metrics fail: %v", err)
		return
	}

	myMountPoints := make(map[string]bool)

	if len(g.Config().Collector.System.MountPoint) > 0 {
		for _, mp := range g.Config().Collector.System.MountPoint {
			myMountPoints[mp] = true
		}
	}

	var diskTotal uint64 = 0
	var diskUsed uint64 = 0

	for idx := range mountPoints {
		fsSpec, fsFile, fsVfstype := mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2]
		if len(myMountPoints) > 0 {
			if _, ok := myMountPoints[fsFile]; !ok {
				log.Debugf("[D] mount point not matched with config %s ignored.", fsFile)
				continue
			}
		}

		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(fsSpec, fsFile, fsVfstype)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		if du.BlocksAll == 0 {
			continue
		}

		if du.FsVfstype == "iso9660" {
			continue
		}

		diskTotal += du.BlocksAll
		diskUsed += du.BlocksUsed

		tags := fmt.Sprintf("mount=%s,fstype=%s", du.FsFile, du.FsVfstype)
		L = append(L, GaugeValue("df.bytes.total", du.BlocksAll, tags))
		L = append(L, GaugeValue("df.bytes.used", du.BlocksUsed, tags))
		L = append(L, GaugeValue("df.bytes.free", du.BlocksFree, tags))
		L = append(L, GaugeValue("df.bytes.used.percent", du.BlocksUsedPercent, tags))
		L = append(L, GaugeValue("df.bytes.free.percent", du.BlocksFreePercent, tags))

		if du.InodesAll == 0 {
			continue
		}

		L = append(L, GaugeValue("df.inodes.total", du.InodesAll, tags))
		L = append(L, GaugeValue("df.inodes.used", du.InodesUsed, tags))
		L = append(L, GaugeValue("df.inodes.free", du.InodesFree, tags))
		L = append(L, GaugeValue("df.inodes.used.percent", du.InodesUsedPercent, tags))
		L = append(L, GaugeValue("df.inodes.free.percent", du.InodesFreePercent, tags))

	}

	if len(L) > 0 && diskTotal > 0 {
		L = append(L, GaugeValue("df.statistics.total", float64(diskTotal)))
		L = append(L, GaugeValue("df.statistics.used", float64(diskUsed)))
		L = append(L, GaugeValue("df.statistics.used.percent", float64(diskUsed)*100.0/float64(diskTotal)))
	}

	return
}

func DeviceMetricsCheck() bool {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Errorf("[E] collect device metrics fail: %v", err)
		return false
	}

	if len(mountPoints) <= 0 {
		return false
	}

	return true
}
