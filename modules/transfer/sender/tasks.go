package sender

import (
	"bytes"
	"math/rand"
	"time"

	pfc "github.com/mia0x75/gopfc/metric"
	log "github.com/sirupsen/logrus"
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
)

// 常量定义
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 // 默认睡眠间隔为50ms
)

// startSendTasks 添加对发送任务的控制,比如stop等
func startSendTasks() {
	cfg := g.Config()

	// init semaphore
	judgeConcurrent := cfg.Judge.MaxConnections
	graphConcurrent := cfg.Graph.MaxConnections
	tsdbConcurrent := cfg.TSDB.MaxConnections
	transferConcurrent := cfg.Transfer.MaxConns

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}

	if graphConcurrent < 1 {
		graphConcurrent = 1
	}

	if transferConcurrent < 1 {
		transferConcurrent = 1
	}
	// init send go-routines
	if cfg.Judge.Enabled {
		for node := range cfg.Judge.Cluster {
			queue := JudgeQueues[node]
			go forward2JudgeTask(queue, node, judgeConcurrent)
		}
	}

	if cfg.Graph.Enabled {
		for node, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				queue := GraphQueues[node+addr]
				go forward2GraphTask(queue, node, addr, graphConcurrent)
			}
		}
	}

	if cfg.TSDB.Enabled {
		go forward2TSDBTask(tsdbConcurrent)
	}

	if cfg.Transfer.Enabled {
		concurrent := transferConcurrent * len(cfg.Transfer.Cluster)
		go forward2TransferTask(TransferQueue, concurrent)
	}
}

// forward2JudgeTask Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeListLimited, node string, concurrent int) {
	batch := g.Config().Judge.Batch // 一次发送,最多batch条数据
	addr := g.Config().Judge.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems := make([]*cm.JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*cm.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, judgeItems []*cm.JudgeItem, count int) {
			defer sema.Release()

			resp := &cm.SimpleRPCResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { // 最多重试3次
				err = JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Errorf("[E] send judge %s: %s fail: %v", node, addr, err)
				proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

// forward2GraphTask Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *list.SafeListLimited, node string, addr string, concurrent int) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*cm.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*cm.GraphItem)
		}

		sema.Acquire()
		go func(addr string, graphItems []*cm.GraphItem, count int) {
			defer sema.Release()

			resp := &cm.SimpleRPCResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { // 最多重试3次
				err = GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Errorf("[E] send to graph %s: %s fail: %v", node, addr, err)
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToGraphCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}

// forward2TSDBTask TSDB定时任务, 将数据通过api发送到tsdb
func forward2TSDBTask(concurrent int) {
	batch := g.Config().TSDB.Batch // 一次发送,最多batch条数据
	retry := g.Config().TSDB.MaxRetry
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := TSDBQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()

			var tsdbBuffer bytes.Buffer
			for i := 0; i < len(itemList); i++ {
				tsdbItem := itemList[i].(*cm.TSDBItem)
				tsdbBuffer.WriteString(tsdbItem.TsdbString())
				tsdbBuffer.WriteString("\n")
			}

			var err error
			for i := 0; i < retry; i++ {
				err = TSDBConnPoolHelper.Send(tsdbBuffer.Bytes())
				if err == nil {
					proc.SendToTSDBCnt.IncrBy(int64(len(itemList)))
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				proc.SendToTSDBFailCnt.IncrBy(int64(len(itemList)))
				log.Errorf("[E] %v", err)
				return
			}
		}(items)
	}
}

// forward2TransferTask Transfer定时任务, 将Transfer发送缓存中的数据 通过rpc连接池 发送到Transfer(此时transfer仅仅起到转发数据的功能)
func forward2TransferTask(Q *list.SafeListLimited, concurrent int) {
	cfg := g.Config()
	batch := cfg.Transfer.Batch // 一次发送,最多batch条数据
	maxConns := int64(cfg.Transfer.MaxConns)
	retry := cfg.Transfer.MaxRetry //最多尝试发送retry次
	if retry < 1 {
		retry = 1
	}

	sema := nsema.NewSemaphore(concurrent)
	transNum := len(TransferHostnames)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		transferItems := make([]*cm.MetricValue, count)
		for i := 0; i < count; i++ {
			transferItems[i] = convert(items[i].(*cm.MetaData))
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(transferItems []*cm.MetricValue, count int) {
			defer sema.Release()

			// 随机遍历transfer列表，直到数据发送成功 或者 遍历完;随机遍历，可以缓解慢transfer
			resp := &cm.TransferResponse{}
			var err error
			sendOk := false

			for j := 0; j < retry && !sendOk; j++ {
				rint := rand.Int()
				for i := 0; i < transNum && !sendOk; i++ {
					idx := (i + rint) % transNum
					host := TransferHostnames[idx]
					addr := TransferMap[host]

					// 过滤掉建连缓慢的host, 否则会严重影响发送速率
					cc := pfc.GetCounterCount(host)
					if cc >= maxConns {
						continue
					}

					pfc.Counter(host, 1)
					err = TransferConnPools.Call(addr, "Transfer.Update", transferItems, resp)
					pfc.Counter(host, -1)

					// statistics
					if err == nil {
						sendOk = true
						proc.SendToTransferCnt.IncrBy(int64(count))
					} else {
						log.Printf("transfer update fail, transfer hostname: %s, transfer instance: %s, items size:%d, error:%v, resp:%v", host, addr, len(transferItems), err, resp)
						proc.SendToTransferFailCnt.IncrBy(int64(count))
					}
				}
			}
		}(transferItems, count)
	}
}

// convert 把模型 cm.MetaData 适配为 cm.MetricValue
func convert(v *cm.MetaData) *cm.MetricValue {
	return &cm.MetricValue{
		Metric:    v.Metric,
		Endpoint:  v.Endpoint,
		Timestamp: v.Timestamp,
		Step:      v.Step,
		Type:      v.CounterType,
		Tags:      cu.SortedTags(v.Tags),
		Value:     v.Value,
	}
}
