package consume

type ImDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Content  string `json:"content"`
	IM       string `json:"im"`
	Status   string `json:"status"`
}
