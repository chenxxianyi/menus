package model

import "time"

type Feedback struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint     `json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Contact   string    `gorm:"type:varchar(100)" json:"contact"`
	Status    int8      `gorm:"default:0;index" json:"status"`
	Reply     string    `gorm:"type:text" json:"reply"`
	CreatedAt time.Time `json:"created_at"`
}

func (Feedback) TableName() string { return "feedback" }
