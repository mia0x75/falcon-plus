package funcs

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
	"github.com/toolkits/nux"
)

func ProcMetrics() (L []*cmodel.MetricValue) {
	ps, err := nux.AllProcs()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	pslen := len(ps)

	procs := hbs.ReportProcs()
	sz := len(procs)
	if sz == 0 {
		return
	}

	for tags, m := range procs {
		cnt := 0
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], m) {
				cnt++
			}
		}

		L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
	}

	return
}

func is_a(p *nux.Proc, m map[int]string) bool {
	// only one kv pair
	for key, val := range m {
		if key == 1 {
			// name
			if val != p.Name {
				return false
			}
		} else if key == 2 {
			// cmdline
			if !strings.Contains(p.Cmdline, val) {
				return false
			}
		}
	}
	return true
}
