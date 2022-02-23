package model

type Sms struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type IM struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}
