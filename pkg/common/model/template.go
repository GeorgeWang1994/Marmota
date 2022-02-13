package model

type Template struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentId"`
	ActionId int    `json:"actionId"`
	Creator  string `json:"creator"`
}
