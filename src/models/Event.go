package models

import "time"

// Event 構造体は、ユーザー情報を表します。
type Event struct {
	Id          uint      `json:"id"`
	UserId      uint      `json:"user_id"`
	User        User      `json:"user"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Event_URL   string    `json:"event_url"`
	EventDate   time.Time `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"nullable"`
}
