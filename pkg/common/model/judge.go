package model

type JudgeItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	JudgeType string            `json:"judgeType"`
	Tags      map[string]string `json:"tags"`
}
