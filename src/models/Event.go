package models

import "time"

// Event 構造体は、ユーザー情報を表します。
type Event struct {
	Id           uint           `json:"id"`
	UserId       uint           `json:"user_id"`
	User         User           `json:"user"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Event_URL    string         `json:"event_url"`
	EventDate    string         `json:"event_date"`
	Image        string         `json:"image" gorm:"type:text;nullable"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"nullable"`
	EventComment []EventComment `json:"comments" gorm:"foreignKey:EventId"`
	Participants []User         `json:"participants" gorm:"many2many:event_participants"`
	EventLikes   []User         `json:"event_likes" gorm:"many2many:event_likes"`
}
