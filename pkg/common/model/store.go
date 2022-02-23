package model

import (
	"bytes"
	"fmt"
	"marmota/pkg/utils/bufferPool"
	"marmota/pkg/utils/md5"
	"marmota/pkg/utils/tag"
	"math"
	"strconv"
)

type StoreItem struct {
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

func UUID(endpoint, metric string, tags map[string]string, dstype string, step int) string {
	ret := bufferPool.BufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.BufferPool.Put(ret)

	if tags == nil || len(tags) == 0 {
		ret.WriteString(endpoint)
		ret.WriteString("/")
		ret.WriteString(metric)
		ret.WriteString("/")
		ret.WriteString(dstype)
		ret.WriteString("/")
		ret.WriteString(strconv.Itoa(step))

		return ret.String()
	}
	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(tag.SortedTags(tags))
	ret.WriteString("/")
	ret.WriteString(dstype)
	ret.WriteString("/")
	ret.WriteString(strconv.Itoa(step))

	return ret.String()
}

func Checksum(endpoint string, metric string, tags map[string]string) string {
	pk := PK(endpoint, metric, tags)
	return md5.Md5(pk)
}

func PK(endpoint, metric string, tags map[string]string) string {
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

func (s *StoreItem) PrimaryKey() string {
	return PK(s.Endpoint, s.Metric, s.Tags)
}

func (s *StoreItem) Checksum() string {
	return Checksum(s.Endpoint, s.Metric, s.Tags)
}

func (s *StoreItem) UUID() string {
	return UUID(s.Endpoint, s.Metric, s.Tags, s.DsType, s.Step)
}

type JsonFloat float64

func (v JsonFloat) MarshalJSON() ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return []byte("null"), nil
	} else {
		return []byte(fmt.Sprintf("%f", f)), nil
	}
}

type RRDData struct {
	Timestamp int64     `json:"timestamp"`
	Value     JsonFloat `json:"value"`
}

func NewRRDData(ts int64, val float64) *RRDData {
	return &RRDData{Timestamp: ts, Value: JsonFloat(val)}
}
