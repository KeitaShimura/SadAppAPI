package models

import "time"

// User 構造体は、ユーザー情報を表します。
type Follow struct {
	Id          uint      `json:"id"`
	FollowingId uint      `json:"following_id"`
	FollowerId  uint      `json:"follower_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"nullable"`
}