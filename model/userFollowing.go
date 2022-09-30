package model

type UserFollowing struct {
	Id          uint `json:"id" gorm:"primaryKey"`
	UserId      uint `json:"user_id" gorm:"not null"`
	FollowingId uint `json:"following_id" gorm:"not null"`
	Permission  bool `json:"permission" gorm:"not null; default:false"`
}
