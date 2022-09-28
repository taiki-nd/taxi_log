package model

type UserFollowing struct {
	UserId      uint `json:"user_id" gorm:"not null"`
	FollowingId uint `json:"following_id" gorm:"not null"`
	Permission  bool `json:"permission" gorm:"not null; default:false"`
}
