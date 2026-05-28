package model

import "time"

type Menu struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	MenuName         string    `gorm:"type:varchar(100)" json:"menu_name"`
	MealType         string    `gorm:"type:varchar(20)" json:"meal_type"`
	PeopleCount      int       `json:"people_count"`
	Taste            string    `gorm:"type:varchar(50)" json:"taste"`
	HealthGoal       string    `gorm:"type:varchar(50)" json:"health_goal"`
	DishesJSON       JSON      `gorm:"type:json" json:"dishes_json"`
	ShoppingListJSON JSON      `gorm:"type:json" json:"shopping_list_json"`
	Reason           string    `gorm:"type:text" json:"reason"`
	CreatedAt        time.Time `json:"created_at"`
}

func (Menu) TableName() string { return "menus" }
