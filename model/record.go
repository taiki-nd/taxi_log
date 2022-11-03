package model

import "time"

type Record struct {
	Id            uint      `json:"id" gorm:"primaryKey"`
	Date          time.Time `json:"date" gorm:"not null"`
	DayOfWeek     string    `json:"day_of_week" gorm:"not null; size:10"`
	StyleFlg      string    `json:"style_flg" gorm:"not null; size:256"`
	StartHour     int64     `json:"start_hour" gorm:"not null"`
	RunningTime   int64     `json:"running_time" gorm:"not null"`
	RunningKm     int64     `json:"running_km" gorm:"not null"`
	OccupancyRate float64   `json:"occupancy_rate" gorm:"not null"`
	NumberOfTime  int64     `json:"number_of_time" gorm:"not null"`
	IsTax         bool      `json:"is_tax" gorm:"not null"`
	DailySales    int64     `json:"daily_sales" gorm:"not null"`
	UserId        uint      `json:"user_id"`
	Details       []Detail  `json:"details" gorm:"foreignKey:RecordId"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
