package model

import "fmt"

// AgentReportRequest 上报探针基础信息
type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
}

// AgentUpdateInfo 更新探针信息
type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

// AgentHeartbeatRequest 探针心跳信息
type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

type AgentMetric struct {
	Metric string
	Tags   string
}

func (a *AgentMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		a.Metric,
		a.Tags,
	)
}

type AgentMetricResponse struct {
	Metrics   []*AgentMetric
	Checksum  string
	Timestamp int64
}

type AgentMetricSlice []*AgentMetric

func (a AgentMetricSlice) Len() int {
	return len(a)
}
func (a AgentMetricSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AgentMetricSlice) Less(i, j int) bool {
	return a[i].String() < a[j].String()
}

type AgentPluginsResponse struct {
	Plugins   []string
	Timestamp int64
}
