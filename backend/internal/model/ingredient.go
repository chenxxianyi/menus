package model

import "time"

type Ingredient struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"type:varchar(50);not null;index" json:"name"`
	Category       string    `gorm:"type:varchar(50);index" json:"category"`
	Image          string    `gorm:"type:varchar(500)" json:"image"`
	NutritionInfo  JSON      `gorm:"type:json" json:"nutrition_info"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Ingredient) TableName() string { return "ingredients" }
