package model

type Area struct {
	Prefecture string `json:"prefecture" gorm:"not null; size:256"`
	Area       string `json:"area" gorm:"not null; size:256"`
}
