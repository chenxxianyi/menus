package model

import "time"

type UserEvent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	EventName   string    `gorm:"type:varchar(80);not null;index" json:"event_name"`
	EntityType  string    `gorm:"type:varchar(50);index" json:"entity_type"`
	EntityID    uint      `gorm:"index" json:"entity_id"`
	PayloadJSON JSON      `gorm:"type:json" json:"payload_json"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

func (UserEvent) TableName() string { return "user_events" }
