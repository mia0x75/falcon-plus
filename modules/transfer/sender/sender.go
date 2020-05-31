package sender

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	rings "github.com/toolkits/consistent/rings"
	nlist "github.com/toolkits/container/list"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cp "github.com/open-falcon/falcon-plus/common/pool"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
)

// 常量定义
const (
	DefaultSendQueueMaxSize = 102400 // 10.24w
)

// 默认参数
var (
	MinStep int // 最小上报周期,单位sec
)

// 服务节点的一致性哈希环
// pk -> node
var (
	JudgeNodeRing *rings.ConsistentHashNodeRing
	GraphNodeRing *rings.ConsistentHashNodeRing
)

// 发送缓存队列
// node -> queue_of_data
var (
	TSDBQueue     *nlist.SafeListLimited
	JudgeQueues   = make(map[string]*nlist.SafeListLimited)
	GraphQueues   = make(map[string]*nlist.SafeListLimited)
	TransferQueue *nlist.SafeListLimited
)

// transfer的主机列表，以及主机名和地址的映射关系
// 用于随机遍历transfer
var (
	TransferMap       = make(map[string]string, 0)
	TransferHostnames = make([]string, 0)
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools     *cp.SafeRPCConnPools
	TSDBConnPoolHelper *cp.TSDBConnPoolHelper
	GraphConnPools     *cp.SafeRPCConnPools
	TransferConnPools  *cp.SafeRPCConnPools
)

// Start 初始化数据发送服务, 在main函数中调用
func Start() {
	go start()
}

func start() {
	// 初始化默认参数
	MinStep = g.Config().MinStep
	if MinStep < 1 {
		MinStep = 30 // 默认30s
	}
	//
	initConnPools()
	initSendQueues()
	initNodeRings()
	// SendTasks依赖基础组件的初始化,要最后启动
	startSendTasks()
	startSenderCron()
	log.Info("[I] send.Start, ok")
}

// Push2JudgeSendQueue 将数据 打入 某个Judge的发送缓存队列, 具体是哪一个Judge 由一致性哈希 决定
func Push2JudgeSendQueue(items []*cm.MetaData) {
	for _, item := range items {
		pk := item.PK()
		node, err := JudgeNodeRing.GetNode(pk)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		// align ts
		step := int(item.Step)
		if step < MinStep {
			step = MinStep
		}
		ts := alignTs(item.Timestamp, int64(step))

		judgeItem := &cm.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}
		Q := JudgeQueues[node]
		isSuccess := Q.PushFront(judgeItem)

		// statistics
		if !isSuccess {
			proc.SendToJudgeDropCnt.Incr()
		}
	}
}

// Push2GraphSendQueue 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
func Push2GraphSendQueue(items []*cm.MetaData) {
	cfg := g.Config().Graph

	for _, item := range items {
		if isMetricIgnored(item) {
			continue
		}
		graphItem, err := convert2GraphItem(item)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}
		pk := item.PK()

		// statistics. 为了效率,放到了这里,因此只有graph是enbale时才能trace
		proc.RecvDataTrace.Trace(pk, item)
		proc.RecvDataFilter.Filter(pk, item.Value, item)

		node, err := GraphNodeRing.GetNode(pk)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		cnode := cfg.ClusterList[node]
		errCnt := 0
		for _, addr := range cnode.Addrs {
			Q := GraphQueues[node+addr]
			if !Q.PushFront(graphItem) {
				errCnt++
			}
		}

		// statistics
		if errCnt > 0 {
			proc.SendToGraphDropCnt.Incr()
		}
	}
}

// convert2GraphItem 打到Graph的数据,要根据rrdtool的特定 来限制 step、counterType、timestamp
func convert2GraphItem(d *cm.MetaData) (*cm.GraphItem, error) {
	item := &cm.GraphItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < MinStep {
		item.Step = MinStep
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == g.GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == g.COUNTER {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == g.DERIVE {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = alignTs(item.Timestamp, int64(item.Step)) // item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}

// Push2TSDBSendQueue 将原始数据入到tsdb发送缓存队列
func Push2TSDBSendQueue(items []*cm.MetaData) {
	for _, item := range items {
		if isMetricIgnored(item) {
			continue
		}
		tsdbItem := convert2TSDBItem(item)
		isSuccess := TSDBQueue.PushFront(tsdbItem)

		if !isSuccess {
			proc.SendToTSDBDropCnt.Incr()
		}
	}
}

// convert2TSDBItem 转化为TSDB格式
func convert2TSDBItem(d *cm.MetaData) *cm.TSDBItem {
	t := cm.TSDBItem{Tags: make(map[string]string)}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Metric = d.Metric
	t.Timestamp = d.Timestamp
	t.Value = d.Value
	return &t
}

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}

func isMetricIgnored(item *cm.MetaData) bool {
	ignoreMetrics := g.Config().IgnoreMetrics
	if b, ok := ignoreMetrics[item.Metric]; ok && b {
		return true
	}
	return false
}

// Push2TransferSendQueue 数据发送到下一个Transfer 当前角色为网关
func Push2TransferSendQueue(items []*cm.MetaData) {
	for _, item := range items {
		isSuccess := TransferQueue.PushFront(item)

		if !isSuccess {
			proc.SendToTransferDropCnt.Incr()
		}
	}
}
