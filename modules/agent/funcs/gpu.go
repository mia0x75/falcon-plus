package funcs

import (
	"github.com/mindprince/gonvml"
	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

// GpuMetrics 需要load libnvidia-ml.so.1库
func GpuMetrics() (L []*cmodel.MetricValue) {
	if err := gonvml.Initialize(); err != nil {
		log.Errorf("[E] Initialize error: %v", err)
		return
	}

	defer gonvml.Shutdown()

	count, err := gonvml.DeviceCount()
	if err != nil {
		log.Errorf("[E] DeviceCount error: %v", err)
		return
	}

	if count == 0 {
		return
	}

	temperature := uint(0)
	totalMemory := uint64(0)
	usedMemory := uint64(0)
	gpuUtilization := uint(0)
	memoryUtilization := uint(0)
	powerUsage := uint(0)
	allUtilization := uint(0)
	allMemoryUtilization := uint(0)

	for i := 0; i < int(count); i++ {
		dev, err := gonvml.DeviceHandleByIndex(uint(i))
		if err != nil {
			log.Errorf("[E] DeviceHandleByIndex error: %v", err)
			continue
		}

		uuid, err := dev.UUID()
		if err != nil {
			log.Errorf("[E] dev.UUID error: %v", err)
		}

		tag := "uuid=" + uuid

		// 不是所有gpu都有风扇
		fanSpeed, err := dev.FanSpeed()
		if err != nil {
			log.Errorf("[E] dev.FanSpeed error: %v", err)
		} else {
			L = append(L, GaugeValue("gpu.fan.speed", fanSpeed, tag))
		}

		temperature, err = dev.Temperature()
		if err != nil {
			log.Errorf("[E] dev.Temperature error: %v", err)
			continue
		}

		totalMemory, usedMemory, err = dev.MemoryInfo()
		if err != nil {
			log.Errorf("[E] dev.MemoryInfo error: %v", err)
			continue
		}

		// 单位换算为兆
		totalBillion := float64(totalMemory / 1024 / 1024)
		usedBillion := float64(usedMemory / 1024 / 1024)

		gpuUtilization, memoryUtilization, err = dev.UtilizationRates()
		if err != nil {
			log.Errorf("[E] dev.UtilizationRates error: %v", err)
			continue
		}

		allUtilization += gpuUtilization
		allMemoryUtilization += memoryUtilization

		powerUsage, err = dev.PowerUsage()
		if err != nil {
			log.Errorf("[E] dev.PowerUsage error: %v", err)
		}

		// 单位换算为瓦特
		powerWatt := float64(powerUsage / 1000)

		L = append(L, GaugeValue("gpu.temperature", temperature, tag))
		L = append(L, GaugeValue("gpu.memory.total", totalBillion, tag))
		L = append(L, GaugeValue("gpu.memory.used", usedBillion, tag))
		L = append(L, GaugeValue("gpu.memory.util", memoryUtilization, tag))
		L = append(L, GaugeValue("gpu.util", gpuUtilization, tag))
		L = append(L, GaugeValue("gpu.power.usage", powerWatt, tag))
	}

	L = append(L, GaugeValue("gpu.count", count))
	L = append(L, GaugeValue("gpu.util.avg", allUtilization/count))
	L = append(L, GaugeValue("gpu.memory.util.avg", allMemoryUtilization/count))
	return L
}
