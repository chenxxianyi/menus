package model

import "time"

type AIGenerationLog struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	GenerationType string    `gorm:"type:varchar(50);not null;index" json:"generation_type"`
	Scene          string    `gorm:"type:varchar(80);index" json:"scene"`
	Model          string    `gorm:"type:varchar(120)" json:"model"`
	InputJSON      JSON      `gorm:"type:json" json:"input_json"`
	OutputJSON     JSON      `gorm:"type:json" json:"output_json"`
	Status         string    `gorm:"type:varchar(20);not null;index" json:"status"`
	ErrorMessage   string    `gorm:"type:text" json:"error_message"`
	DurationMS     int64     `gorm:"index" json:"duration_ms"`
	RecipeIDsJSON  JSON      `gorm:"type:json" json:"recipe_ids_json"`
	CreatedAt      time.Time `gorm:"index" json:"created_at"`
}

func (AIGenerationLog) TableName() string { return "ai_generation_logs" }
