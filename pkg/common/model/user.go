package model

type User struct {
	ID     int64  `json:"id" `
	Name   string `json:"name"`
	Cnname string `json:"cnname"`
	Passwd string `json:"-"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	IM     string `json:"im" gorm:"column:im"`
	QQ     string `json:"qq" gorm:"column:qq"`
	Role   int    `json:"role"`
}
