package model

type Team struct {
	ID      int64  `json:"id,"`
	Name    string `json:"name"`
	Resume  string `json:"resume"`
	Creator int64  `json:"creator"`
}
