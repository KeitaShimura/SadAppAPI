package models

import "time"

// User 構造体は、ユーザー情報を表します。
type Follow struct {
	Id          uint      `json:"id"`
	FollowingId uint      `json:"following_id"`
	Following   User      `json:"following" gorm:"foreignKey:FollowingId;references:Id"`
	FollowerId  uint      `json:"follower_id"`
	Follower    User      `json:"follower" gorm:"foreignKey:FollowerId;references:Id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
