package model

type Expression struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"function"`   // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
	ActionId   int               `json:"actionId"`
}

type ExpressionResponse struct {
	Expressions []*Expression `json:"expressions"`
}
