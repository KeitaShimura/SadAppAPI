package models

import "time"

type EventComment struct {
	Id        uint      `json:"id"`
	EventId   uint      `json:"event_id" gorm:"column:event_id"`
	UserId    uint      `json:"user_id" gorm:"column:user_id"`
	User      User      `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"nullable"`
}
