package rpc

import (
	"marmota/pkg/common/model"
	"marmota/pkg/utils/tag"
	"marmota/transfer/cc"
	"marmota/transfer/gg"
	"marmota/transfer/sender/queue"
	"marmota/transfer/stat"
	"strconv"
	"time"
)

type Transfer int

func (t *Transfer) Ping(req model.NullRpcRequest, resp *model.RpcResponse) error {
	return nil
}

func (t *Transfer) Update(args []*model.MetricValue, reply *model.TransferResponse) error {
	return RecvMetricValues(args, reply, "connpool")
}

// RecvMetricValues process new metric values
func RecvMetricValues(args []*model.MetricValue, reply *model.TransferResponse, from string) error {
	start := time.Now()
	reply.Invalid = 0

	var items []*model.MetaData
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		if v.Metric == "" || v.Endpoint == "" {
			reply.Invalid += 1
			continue
		}

		if v.Type != gg.COUNTER && v.Type != gg.GAUGE && v.Type != gg.DERIVE {
			reply.Invalid += 1
			continue
		}

		if v.Value == "" {
			reply.Invalid += 1
			continue
		}

		if v.Step <= 0 {
			reply.Invalid += 1
			continue
		}

		if len(v.Metric)+len(v.Tags) > 510 {
			reply.Invalid += 1
			continue
		}

		// TODO 呵呵,这里需要再优雅一点
		now := start.Unix()
		if v.Timestamp <= 0 || v.Timestamp > now*2 {
			v.Timestamp = now
		}

		fv := &model.MetaData{
			Metric:      v.Metric,
			Endpoint:    v.Endpoint,
			Timestamp:   v.Timestamp,
			Step:        v.Step,
			CounterType: v.Type,
			Tags:        tag.TagsDict(v.Tags), //TODO tags键值对的个数,要做一下限制
		}

		valid := true
		var vv float64
		var err error

		switch cv := v.Value.(type) {
		case string:
			vv, err = strconv.ParseFloat(cv, 64)
			if err != nil {
				valid = false
			}
		case float64:
			vv = cv
		case int64:
			vv = float64(cv)
		default:
			valid = false
		}

		if !valid {
			reply.Invalid += 1
			continue
		}

		fv.Value = vv
		items = append(items, fv)
	}

	// statistics 统计数据
	cnt := int64(len(items))
	stat.RecvCnt.IncrBy(cnt)
	if from == "connpool" {
		stat.RpcRecvCnt.IncrBy(cnt)
	} else if from == "http" {
		stat.HttpRecvCnt.IncrBy(cnt)
	}

	// 发送数据对应的队列中
	cfg := cc.Config()

	if cfg.Graph.Enabled {
		queue.Push2GraphSendQueue(items)
	}

	if cfg.Judge.Enabled {
		queue.Push2JudgeSendQueue(items)
	}

	if cfg.Tsdb.Enabled {
		queue.Push2TsdbSendQueue(items)
	}

	if cfg.Transfer.Enabled {
		queue.Push2TransferSendQueue(items)
	}

	if cfg.Influxdb.Enabled {
		queue.Push2InfluxdbSendQueue(items)
	}

	reply.Message = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
