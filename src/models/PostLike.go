package models

import "time"

type PostLike struct {
	PostId    uint      `json:"post_id" gorm:"primaryKey;autoIncrement:false"`
	UserId    uint      `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
}
