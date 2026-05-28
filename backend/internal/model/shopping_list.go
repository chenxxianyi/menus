package model

import "time"

type ShoppingList struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	ItemsJSON JSON      `gorm:"type:json" json:"items_json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ShoppingList) TableName() string { return "shopping_lists" }
