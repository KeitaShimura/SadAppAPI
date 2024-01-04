package models

import "time"

type EventLike struct {
	EventId   uint      `json:"event_id" gorm:"primaryKey;autoIncrement:false"`
	UserId    uint      `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
}
