package proc

import (
	log "github.com/sirupsen/logrus"
	nproc "github.com/toolkits/proc"
)

// trace
var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 3)
)

// filter
var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt       = nproc.NewSCounterQps("RecvCnt")
	RPCRecvCnt    = nproc.NewSCounterQps("RPCRecvCnt")
	HTTPRecvCnt   = nproc.NewSCounterQps("HTTPRecvCnt")
	SocketRecvCnt = nproc.NewSCounterQps("SocketRecvCnt")

	SendToJudgeCnt    = nproc.NewSCounterQps("SendToJudgeCnt")
	SendToTSDBCnt     = nproc.NewSCounterQps("SendToTSDBCnt")
	SendToGraphCnt    = nproc.NewSCounterQps("SendToGraphCnt")
	SendToTransferCnt = nproc.NewSCounterQps("SendToTransferCnt")

	SendToJudgeDropCnt    = nproc.NewSCounterQps("SendToJudgeDropCnt")
	SendToTSDBDropCnt     = nproc.NewSCounterQps("SendToTSDBDropCnt")
	SendToGraphDropCnt    = nproc.NewSCounterQps("SendToGraphDropCnt")
	SendToTransferDropCnt = nproc.NewSCounterQps("SendToTransferDropCnt")

	SendToJudgeFailCnt    = nproc.NewSCounterQps("SendToJudgeFailCnt")
	SendToTSDBFailCnt     = nproc.NewSCounterQps("SendToTSDBFailCnt")
	SendToGraphFailCnt    = nproc.NewSCounterQps("SendToGraphFailCnt")
	SendToTransferFailCnt = nproc.NewSCounterQps("SendToTransferFailCnt")

	// 发送缓存大小
	JudgeQueuesCnt    = nproc.NewSCounterBase("JudgeSendCacheCnt")
	TSDBQueuesCnt     = nproc.NewSCounterBase("TSDBSendCacheCnt")
	GraphQueuesCnt    = nproc.NewSCounterBase("GraphSendCacheCnt")
	TransferQueuesCnt = nproc.NewSCounterBase("TransferSendCacheCnt")

	// http请求次数
	HistoryRequestCnt = nproc.NewSCounterQps("HistoryRequestCnt")
	InfoRequestCnt    = nproc.NewSCounterQps("InfoRequestCnt")
	LastRequestCnt    = nproc.NewSCounterQps("LastRequestCnt")
	LastRawRequestCnt = nproc.NewSCounterQps("LastRawRequestCnt")

	// http回执的监控数据条数
	HistoryResponseCounterCnt = nproc.NewSCounterQps("HistoryResponseCounterCnt")
	HistoryResponseItemCnt    = nproc.NewSCounterQps("HistoryResponseItemCnt")
	LastRequestItemCnt        = nproc.NewSCounterQps("LastRequestItemCnt")
	LastRawRequestItemCnt     = nproc.NewSCounterQps("LastRawRequestItemCnt")
)

func Start() {
	log.Info("[I] proc.Start, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// recv cnt
	ret = append(ret, RecvCnt.Get())
	ret = append(ret, RPCRecvCnt.Get())
	ret = append(ret, HTTPRecvCnt.Get())
	ret = append(ret, SocketRecvCnt.Get())

	// send cnt
	ret = append(ret, SendToJudgeCnt.Get())
	ret = append(ret, SendToTSDBCnt.Get())
	ret = append(ret, SendToGraphCnt.Get())
	ret = append(ret, SendToTransferCnt.Get())

	// drop cnt
	ret = append(ret, SendToJudgeDropCnt.Get())
	ret = append(ret, SendToTSDBDropCnt.Get())
	ret = append(ret, SendToGraphDropCnt.Get())
	ret = append(ret, SendToTransferDropCnt.Get())

	// send fail cnt
	ret = append(ret, SendToJudgeFailCnt.Get())
	ret = append(ret, SendToTSDBFailCnt.Get())
	ret = append(ret, SendToGraphFailCnt.Get())
	ret = append(ret, SendToTransferFailCnt.Get())

	// cache cnt
	ret = append(ret, JudgeQueuesCnt.Get())
	ret = append(ret, TSDBQueuesCnt.Get())
	ret = append(ret, GraphQueuesCnt.Get())
	ret = append(ret, TransferQueuesCnt.Get())

	// http request
	ret = append(ret, HistoryRequestCnt.Get())
	ret = append(ret, InfoRequestCnt.Get())
	ret = append(ret, LastRequestCnt.Get())
	ret = append(ret, LastRawRequestCnt.Get())

	// http response
	ret = append(ret, HistoryResponseCounterCnt.Get())
	ret = append(ret, HistoryResponseItemCnt.Get())
	ret = append(ret, LastRequestItemCnt.Get())
	ret = append(ret, LastRawRequestItemCnt.Get())

	return ret
}
