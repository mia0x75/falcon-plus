package funcs

import (
	"fmt"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

var (
	diskStatsMap = make(map[string][2]*nux.DiskStats)
	dsLock       = new(sync.RWMutex)
)

// UpdateDiskStats TODO:
func UpdateDiskStats() error {
	dsList, err := nux.ListDiskStats()
	if err != nil {
		return err
	}
	dsLock.Lock()
	defer dsLock.Unlock()
	for i := 0; i < len(dsList); i++ {
		device := dsList[i].Device
		diskStatsMap[device] = [2]*nux.DiskStats{dsList[i], diskStatsMap[device][0]}
	}
	return nil
}

// IOReadRequests TODO:
func IOReadRequests(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadRequests - arr[1].ReadRequests
}

// IOReadMerged TODO:
func IOReadMerged(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadMerged - arr[1].ReadMerged
}

// IOReadSectors TODO:
func IOReadSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadSectors - arr[1].ReadSectors
}

// IOMsecRead TODO:
func IOMsecRead(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecRead - arr[1].MsecRead
}

// IOWriteRequests TODO:
func IOWriteRequests(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteRequests - arr[1].WriteRequests
}

// IOWriteMerged TODO:
func IOWriteMerged(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteMerged - arr[1].WriteMerged
}

// IOWriteSectors TODO:
func IOWriteSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteSectors - arr[1].WriteSectors
}

// IOMsecWrite TODO:
func IOMsecWrite(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecWrite - arr[1].MsecWrite
}

// IOMsecTotal TODO:
func IOMsecTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecTotal - arr[1].MsecTotal
}

// IOMsecWeightedTotal TODO:
func IOMsecWeightedTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecWeightedTotal - arr[1].MsecWeightedTotal
}

// TS TODO:
func TS(arr [2]*nux.DiskStats) uint64 {
	return uint64(arr[0].TS.Sub(arr[1].TS).Nanoseconds() / 1000000)
}

// IODelta TODO:
func IODelta(device string, f func([2]*nux.DiskStats) uint64) uint64 {
	val, ok := diskStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

// DiskIOMetrics TODO:
func DiskIOMetrics() (L []*cm.MetricValue) {
	dsList, err := nux.ListDiskStats()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	for _, ds := range dsList {
		if !ShouldHandleDevice(ds.Device) {
			continue
		}

		device := "device=" + ds.Device

		L = append(L, CounterValue("disk.io.read_requests", ds.ReadRequests, device))
		L = append(L, CounterValue("disk.io.read_merged", ds.ReadMerged, device))
		L = append(L, CounterValue("disk.io.read_sectors", ds.ReadSectors, device))
		L = append(L, CounterValue("disk.io.msec_read", ds.MsecRead, device))
		L = append(L, CounterValue("disk.io.write_requests", ds.WriteRequests, device))
		L = append(L, CounterValue("disk.io.write_merged", ds.WriteMerged, device))
		L = append(L, CounterValue("disk.io.write_sectors", ds.WriteSectors, device))
		L = append(L, CounterValue("disk.io.msec_write", ds.MsecWrite, device))
		L = append(L, CounterValue("disk.io.ios_in_progress", ds.IosInProgress, device))
		L = append(L, CounterValue("disk.io.msec_total", ds.MsecTotal, device))
		L = append(L, CounterValue("disk.io.msec_weighted_total", ds.MsecWeightedTotal, device))
	}
	return
}

// IOStatsMetrics TODO:
func IOStatsMetrics() (L []*cm.MetricValue) {
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device := range diskStatsMap {
		if !ShouldHandleDevice(device) {
			continue
		}

		tags := "device=" + device
		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)
		rsec := IODelta(device, IOReadSectors)
		wsec := IODelta(device, IOWriteSectors)
		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		use := IODelta(device, IOMsecTotal)
		io := rio + wio
		avgrqSize := 0.0
		await := 0.0
		svctm := 0.0
		if io != 0 {
			avgrqSize = float64(rsec+wsec) / float64(io)
			await = float64(ruse+wuse) / float64(io)
			svctm = float64(use) / float64(io)
		}

		duration := IODelta(device, TS)

		L = append(L, GaugeValue("disk.io.read_bytes", float64(rsec)*512.0, tags))
		L = append(L, GaugeValue("disk.io.write_bytes", float64(wsec)*512.0, tags))
		L = append(L, GaugeValue("disk.io.avgrq_sz", avgrqSize, tags))
		L = append(L, GaugeValue("disk.io.avgqu_sz", float64(IODelta(device, IOMsecWeightedTotal))/1000.0, tags))
		L = append(L, GaugeValue("disk.io.await", await, tags))
		L = append(L, GaugeValue("disk.io.svctm", svctm, tags))

		var tmp float64
		if duration == 0 {
			tmp = 0
		} else {
			tmp = float64(use) * 100.0 / float64(duration)
			if tmp > 100.0 {
				tmp = 100.0
			}
		}
		L = append(L, GaugeValue("disk.io.util", tmp, tags))
	}

	return
}

// IOStatsForPage TODO:
func IOStatsForPage() (L [][]string) {
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device := range diskStatsMap {
		if !ShouldHandleDevice(device) {
			continue
		}

		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)

		rsec := IODelta(device, IOReadSectors)
		wsec := IODelta(device, IOWriteSectors)

		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		use := IODelta(device, IOMsecTotal)
		io := rio + wio
		avgrqSize := 0.0
		await := 0.0
		svctm := 0.0
		if io != 0 {
			avgrqSize = float64(rsec+wsec) / float64(io)
			await = float64(ruse+wuse) / float64(io)
			svctm = float64(use) / float64(io)
		}

		item := []string{
			device,
			fmt.Sprintf("%d", IODelta(device, IOReadMerged)),
			fmt.Sprintf("%d", IODelta(device, IOWriteMerged)),
			fmt.Sprintf("%d", rio),
			fmt.Sprintf("%d", wio),
			fmt.Sprintf("%.2f", float64(rsec)/2.0),
			fmt.Sprintf("%.2f", float64(wsec)/2.0),
			fmt.Sprintf("%.2f", avgrqSize),                                            // avgrq_sz: delta(rsect+wsect)/delta(rio+wio)
			fmt.Sprintf("%.2f", float64(IODelta(device, IOMsecWeightedTotal))/1000.0), // avgqu_sz: delta(aveq)/s/1000
			fmt.Sprintf("%.2f", await),                                                // await: delta(ruse+wuse)/delta(rio+wio)
			fmt.Sprintf("%.2f", svctm),                                                // svctm: delta(use)/delta(rio+wio)
			fmt.Sprintf("%.2f%%", float64(use)/10.0),                                  // %util: delta(use)/s/1000 * 100%
		}
		L = append(L, item)
	}

	return
}

// ShouldHandleDevice TODO:
func ShouldHandleDevice(device string) bool {
	normal := len(device) == 3 && (strings.HasPrefix(device, "sd") || strings.HasPrefix(device, "vd"))
	aws := len(device) >= 4 && strings.HasPrefix(device, "xvd")
	flash := len(device) >= 4 && (strings.HasPrefix(device, "fio") || strings.HasPrefix(device, "nvme"))
	return normal || aws || flash
}
