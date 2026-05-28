package model

import "time"

type Favorite struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_recipe" json:"user_id"`
	RecipeID  uint      `gorm:"not null;uniqueIndex:idx_user_recipe" json:"recipe_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Favorite) TableName() string { return "favorites" }
