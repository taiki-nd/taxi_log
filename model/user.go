package model

import "time"

type User struct {
	Id              uint      `json:"id" gorm:"primaryKey"`
	Uuid            string    `json:"uuid" gorm:"not null; size:256"`
	Nickname        string    `json:"nickname" gorm:"not null; size:30"`
	ProfileImageUrl string    `json:"profile_image_url" gorm:"not null"`
	Prefecture      string    `json:"prefecture" gorm:"not null; size:256"`
	Company         string    `json:"company" gorm:"not null; size:256"`
	StyleFlg        string    `json:"style_flg" gorm:"not null; size:256"`
	CloseDay        int64     `json:"close_day" gorm:"not null"`
	DailyTarget     int64     `json:"daily_target" gorm:"not null"`
	MonthlyTarget   int64     `json:"monthly_target" gorm:"not null"`
	IsTax           bool      `json:"is_tax" gorm:"not null; default:false"`
	OpenFlg         string    `json:"open_flg" gorm:"not null; size:256; default:open"`
	IsAdmin         bool      `json:"is_admin" gorm:"not null; default:false"`
	Records         []Record  `json:"records" gorm:"foreignKey:UserId"`
	Followings      []*User   `json:"followings" gorm:"many2many:user_followings"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
