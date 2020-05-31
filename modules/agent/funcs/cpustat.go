package funcs

import (
	"math"
	"runtime"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

// CPUUsages TODO:
type CPUUsages struct {
	User    float64
	Nice    float64
	System  float64
	Idle    float64
	Busy    float64
	Iowait  float64
	Irq     float64
	SoftIrq float64
	Steal   float64
	Guest   float64
}

// UpdateCPUStats TODO:
func UpdateCPUStats() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func cpunumTotal() uint64 {
	num := runtime.NumCPU()
	ss := strconv.Itoa(num)
	b, err := strconv.ParseUint(ss, 10, 64)
	if err != nil {
		log.Errorf("[E] %v", err)
	}
	return b
}

// CPUUsagesSummary TODO:
func CPUUsagesSummary() (cpuUsages *CPUUsages, switches uint64, prepared bool) {
	psLock.RLock()
	defer psLock.RUnlock()

	// cpuUsages = make([]float64, 0, 10)
	switches = 0

	prepared = true
	if procStatHistory[1] == nil {
		prepared = false
	}

	// procStatHistory[1] equals to nil
	if !prepared {
		return
	}

	// procStatHistory[1] alway does not equal nil here
	dt := procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total

	if dt == 0 {
		cpuUsages = &CPUUsages{
			User:    0.0,
			Nice:    0.0,
			System:  0.0,
			Idle:    0.0,
			Busy:    100.0,
			Iowait:  0.0,
			Irq:     0.0,
			SoftIrq: 0.0,
			Steal:   0.0,
			Guest:   0.0,
		}

		switches = procStatHistory[0].Ctxt
	} else {
		invQuotient := 100.00 / float64(dt)

		user := math.Round(float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User)*invQuotient*100) / 100
		nice := math.Round(float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice)*invQuotient*100) / 100
		system := math.Round(float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System)*invQuotient*100) / 100
		iowait := math.Round(float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait)*invQuotient*100) / 100
		irq := math.Round(float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq)*invQuotient*100) / 100
		softirq := math.Round(float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq)*invQuotient*100) / 100
		steal := math.Round(float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal)*invQuotient*100) / 100
		guest := math.Round(float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest)*invQuotient*100) / 100
		busy := user + nice + system + iowait + irq + softirq + steal + guest
		idle := 100.0 - busy

		cpuUsages = &CPUUsages{
			User:    user,
			Nice:    nice,
			System:  system,
			Idle:    idle,
			Busy:    busy,
			Iowait:  iowait,
			Irq:     irq,
			SoftIrq: softirq,
			Steal:   steal,
			Guest:   guest,
		}
		switches = procStatHistory[0].Ctxt
	}

	return
}

// CPUMetrics TODO:
func CPUMetrics() []*cm.MetricValue {
	cpuUsages, currentCPUSwitches, prepared := CPUUsagesSummary()

	if !prepared {
		return []*cm.MetricValue{}
	}

	cpunum := GaugeValue("cpu.num", cpunumTotal())
	idle := GaugeValue("cpu.idle", cpuUsages.Idle)
	busy := GaugeValue("cpu.busy", cpuUsages.Busy)
	user := GaugeValue("cpu.user", cpuUsages.User)
	nice := GaugeValue("cpu.nice", cpuUsages.Nice)
	system := GaugeValue("cpu.system", cpuUsages.System)
	iowait := GaugeValue("cpu.iowait", cpuUsages.Iowait)
	irq := GaugeValue("cpu.irq", cpuUsages.Irq)
	softirq := GaugeValue("cpu.softirq", cpuUsages.SoftIrq)
	steal := GaugeValue("cpu.steal", cpuUsages.Steal)
	guest := GaugeValue("cpu.guest", cpuUsages.Guest)
	switches := CounterValue("cpu.switches", currentCPUSwitches)
	return []*cm.MetricValue{cpunum, idle, busy, user, nice, system, iowait, irq, softirq, steal, guest, switches}
}
