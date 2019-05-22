package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

// USES TODO:
var USES = map[string]struct{}{
	"PruneCalled":        {},
	"LockDroppedIcmps":   {},
	"ArpFilter":          {},
	"TW":                 {},
	"DelayedACKLocked":   {},
	"ListenOverflows":    {},
	"ListenDrops":        {},
	"TCPPrequeueDropped": {},
	"TCPTSReorder":       {},
	"TCPDSACKUndo":       {},
	"TCPLoss":            {},
	"TCPLostRetransmit":  {},
	"TCPLossFailures":    {},
	"TCPFastRetrans":     {},
	"TCPTimeouts":        {},
	"TCPSchedulerFailed": {},
	"TCPAbortOnMemory":   {},
	"TCPAbortOnTimeout":  {},
	"TCPAbortFailed":     {},
	"TCPMemoryPressures": {},
	"TCPSpuriousRTOs":    {},
	"TCPBacklogDrop":     {},
	"TCPMinTTLDrop":      {},
}

// NetstatMetrics TODO:
func NetstatMetrics() (L []*cmodel.MetricValue) {
	tcpExts, err := nux.Netstat("TcpExt")

	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	cnt := len(tcpExts)
	if cnt == 0 {
		return
	}

	for key, val := range tcpExts {
		if _, ok := USES[key]; !ok {
			continue
		}
		L = append(L, CounterValue("net.tcp."+key, val))
	}

	return
}
