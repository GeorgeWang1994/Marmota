package model

import (
	"bytes"
	"marmota/pkg/utils/bufferPool"
	"marmota/pkg/utils/tag"
)

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

type MetaData struct {
	Metric      string            `json:"metric"`
	Endpoint    string            `json:"endpoint"`
	Timestamp   int64             `json:"timestamp"`
	Step        int64             `json:"step"`
	Value       float64           `json:"value"`
	CounterType string            `json:"counterType"`
	Tags        map[string]string `json:"tags"`
}

func pk(endpoint, metric string, tags map[string]string) string {
	ret := bufferPool.BufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.BufferPool.Put(ret)

	if tags == nil || len(tags) == 0 {
		ret.WriteString(endpoint)
		ret.WriteString("/")
		ret.WriteString(metric)

		return ret.String()
	}
	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(tag.SortedTags(tags))
	return ret.String()
}

func (t *MetaData) PK() string {
	return pk(t.Endpoint, t.Metric, t.Tags)
}
