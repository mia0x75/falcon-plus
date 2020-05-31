package sender

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/list"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
)

// 常量定义
const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second    // ProcCron的周期,默认1s
	DefaultLogCronPeriod  = time.Duration(3600) * time.Second // LogCron的周期,默认300s
)

// startSenderCron cron程序入口
func startSenderCron() {
	go startProcCron()
	go startLogCron()
}

func startProcCron() {
	for range time.Tick(DefaultProcCronPeriod) {
		refreshSendingCacheSize()
	}
}

func startLogCron() {
	for range time.Tick(DefaultLogCronPeriod) {
		logConnPoolsProc()
	}
}

func refreshSendingCacheSize() {
	cfg := g.Config()

	if cfg.Judge.Enabled {
		proc.JudgeQueuesCnt.SetCnt(calcSendCacheSize(JudgeQueues))
	}

	if cfg.Graph.Enabled {
		proc.GraphQueuesCnt.SetCnt(calcSendCacheSize(GraphQueues))
	}

	if cfg.TSDB.Enabled {
		proc.TSDBQueuesCnt.SetCnt(int64(TSDBQueue.Len()))
	}

	if cfg.Transfer.Enabled {
		proc.TransferQueuesCnt.SetCnt(int64(TransferQueue.Len()))
	}
}

func calcSendCacheSize(mapList map[string]*list.SafeListLimited) int64 {
	var cnt int64
	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}

func logConnPoolsProc() {
	cfg := g.Config()
	if cfg.Judge.Enabled {
		log.Infof("[I] judge connPools proc: \n%v", strings.Join(JudgeConnPools.Proc(), "\n"))
	}
	if cfg.Graph.Enabled {
		log.Infof("[I] graph connPools proc: \n%v", strings.Join(GraphConnPools.Proc(), "\n"))
	}
	if cfg.Transfer.Enabled {
		log.Infof("[I] transfer connPools proc: \n%v", strings.Join(TransferConnPools.Proc(), "\n"))
	}
}
