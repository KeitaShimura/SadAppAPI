package models

import "time"

// Post 構造体は、ユーザー情報を表します。
type Post struct {
	Id          uint      `json:"id"`
	UserId      uint      `json:"user_id" gorm:"column:user_id"`
	User        User      `json:"user"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"nullable"`
}
