package models

import "time"

// Post 構造体は、ユーザー情報を表します。
type PostComment struct {
    Id        uint      `json:"id"`
    PostId    uint      `json:"post_id" gorm:"column:post_id"`
    UserId    uint      `json:"user_id" gorm:"column:user_id"`
    User      User      `json:"user"`
    Content string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at" gorm:"nullable"`
}