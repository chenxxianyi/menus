package model

import "time"

type BrowseHistory struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_recipe_history;index" json:"user_id"`
	RecipeID  uint      `gorm:"not null;uniqueIndex:idx_user_recipe_history;index" json:"recipe_id"`
	ViewedAt  time.Time `gorm:"not null;index" json:"viewed_at"`
	Recipe    Recipe    `gorm:"foreignKey:RecipeID" json:"recipe"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BrowseHistory) TableName() string { return "browse_histories" }
