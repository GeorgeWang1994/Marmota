package stat

import (
	nproc "github.com/toolkits/proc"
)

var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 3)
)

var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt       = nproc.NewSCounterQps("RecvCnt")
	RpcRecvCnt    = nproc.NewSCounterQps("RpcRecvCnt")
	HttpRecvCnt   = nproc.NewSCounterQps("HttpRecvCnt")
	SocketRecvCnt = nproc.NewSCounterQps("SocketRecvCnt")

	SendToJudgeCnt    = nproc.NewSCounterQps("SendToJudgeCnt")
	SendToTsdbCnt     = nproc.NewSCounterQps("SendToTsdbCnt")
	SendToGraphCnt    = nproc.NewSCounterQps("SendToGraphCnt")
	SendToTransferCnt = nproc.NewSCounterQps("SendToTransferCnt")
	SendToInfluxdbCnt = nproc.NewSCounterQps("SendToInfluxdbCnt")

	SendToJudgeDropCnt    = nproc.NewSCounterQps("SendToJudgeDropCnt")
	SendToTsdbDropCnt     = nproc.NewSCounterQps("SendToTsdbDropCnt")
	SendToGraphDropCnt    = nproc.NewSCounterQps("SendToGraphDropCnt")
	SendToTransferDropCnt = nproc.NewSCounterQps("SendToTransferDropCnt")
	SendToInfluxdbDropCnt = nproc.NewSCounterQps("SendToTsdbDropCnt")

	SendToJudgeFailCnt    = nproc.NewSCounterQps("SendToJudgeFailCnt")
	SendToTsdbFailCnt     = nproc.NewSCounterQps("SendToTsdbFailCnt")
	SendToGraphFailCnt    = nproc.NewSCounterQps("SendToGraphFailCnt")
	SendToTransferFailCnt = nproc.NewSCounterQps("SendToTransferFailCnt")
	SendToInfluxdbFailCnt = nproc.NewSCounterQps("SendToInfluxdbFailCnt")

	// 发送缓存大小
	JudgeQueuesCnt    = nproc.NewSCounterBase("JudgeSendCacheCnt")
	TsdbQueuesCnt     = nproc.NewSCounterBase("TsdbSendCacheCnt")
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
