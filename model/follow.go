package model

type Follow struct {
	UserId      uint `json:"user_id" gorm:"not null"`
	FollowingId uint `json:"following_id" gorm:"not null"`
}
