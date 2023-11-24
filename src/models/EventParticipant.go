package models

import "time"

type EventParticipant struct {
	EventId   uint      `json:"event_id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
