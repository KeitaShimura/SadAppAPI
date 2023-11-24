package models

import "time"

// Post 構造体は、ユーザー情報を表します。
type Post struct {
	Id          uint          `json:"id"`
	UserId      uint          `json:"user_id" gorm:"column:user_id"`
	User        User          `json:"user"`
	Content     string        `json:"content"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"nullable"`
	PostComment []PostComment `json:"comments" gorm:"foreignKey:PostId"`
	PostLikes   []User        `json:"post_likes" gorm:"many2many:post_likes"`
}
