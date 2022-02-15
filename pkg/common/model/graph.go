package model

type GraphItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	// DsType 即RRD中的Datasource的类型：GAUGE|COUNTER|DERIVE
	DsType    string `json:"dstype"`
	Step      int    `json:"step"`
	Heartbeat int    `json:"heartbeat"`
	Min       string `json:"min"`
	Max       string `json:"max"`
}
