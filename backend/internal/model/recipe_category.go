package model

import "time"

type RecipeCategory struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Icon      string    `gorm:"type:varchar(500)" json:"icon"`
	Sort      int       `gorm:"default:0;index" json:"sort"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (RecipeCategory) TableName() string { return "recipe_categories" }
