package rpc

import (
	"bytes"
	"marmota/pivas/cache"
	"marmota/pkg/common/model"
	"marmota/pkg/utils/md5"
	"sort"
	"time"
)

type Agent int

func (t *Agent) MinePlugins(args model.AgentHeartbeatRequest, reply *model.AgentPluginsResponse) error {
	if args.Hostname == "" {
		return nil
	}

	reply.Plugins = cache.GetPlugins(args.Hostname)
	reply.Timestamp = time.Now().Unix()

	return nil
}

func (t *Agent) ReportStatus(args *model.AgentReportRequest, reply *model.RpcResponse) error {
	if args.Hostname == "" {
		reply.Code = 1
		return nil
	}

	cache.Agents.Put(args)

	return nil
}

func DigestBuiltinMetrics(items []*model.AgentMetric) string {
	sort.Sort(model.AgentMetricSlice(items))

	var buf bytes.Buffer
	for _, m := range items {
		buf.WriteString(m.String())
	}

	return md5.Md5(buf.String())
}

// AgentMetrics agent按照server端的配置，按需采集的metric，比如net.port.listen port=22 或者 proc.num name=zabbix_agentd
func (t *Agent) AgentMetrics(args *model.AgentHeartbeatRequest, reply *model.AgentMetricResponse) error {
	if args.Hostname == "" {
		return nil
	}

	metrics, err := cache.GetBuiltinMetrics(args.Hostname)
	if err != nil {
		return nil
	}

	checksum := ""
	if len(metrics) > 0 {
		checksum = DigestBuiltinMetrics(metrics)
	}

	if args.Checksum == checksum {
		reply.Metrics = []*model.AgentMetric{}
	} else {
		reply.Metrics = metrics
	}
	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return nil
}

