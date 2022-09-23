package model

import "time"

type Detail struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	DepartHour  int64     `json:"depart_hour" gorm:"not null"`
	DepartPlace string    `json:"depart_place" gorm:"not null; size:256"`
	ArrivePlace string    `json:"arrive_place" gorm:"not null; size:256"`
	IsTax       bool      `json:"is_tax" gorm:"not null"`
	Sales       int64     `json:"sales" gorm:"not null"`
	MethodFlg   string    `json:"method_flg" gorm:"not null; size:256"`
	Description string    `json:"description" gorm:"not null"`
	RecordId    uint      `json:"record_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
