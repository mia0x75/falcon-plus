package funcs

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

// ProcMetrics TODO:
func ProcMetrics() (L []*cm.MetricValue) {
	ps, err := nux.AllProcs()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	pslen := len(ps)
	L = append(L, GaugeValue(g.PROC_NUM, pslen))

	procs := hbs.ReportProcs()
	sz := len(procs)
	if sz == 0 {
		return
	}

	for tags, m := range procs {
		cnt := 0
		for i := 0; i < pslen; i++ {
			if isA(ps[i], m) {
				cnt++
			}
		}

		L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
	}

	return
}

func isA(p *nux.Proc, m map[int]string) bool {
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
