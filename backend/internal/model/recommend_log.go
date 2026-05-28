package model

import "time"

type RecommendLog struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        *uint     `gorm:"index" json:"user_id"`
	RecommendType string    `gorm:"type:varchar(50);index" json:"recommend_type"`
	InputJSON     JSON      `gorm:"type:json" json:"input_json"`
	ResultJSON    JSON      `gorm:"type:json" json:"result_json"`
	CreatedAt     time.Time `json:"created_at"`
}

func (RecommendLog) TableName() string { return "recommend_logs" }
