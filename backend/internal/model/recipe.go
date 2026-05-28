package model

import "time"

type Recipe struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title         string    `gorm:"type:varchar(100);not null;index" json:"title"`
	Cover         string    `gorm:"type:varchar(500)" json:"cover"`
	Description   string    `gorm:"type:text" json:"description"`
	CategoryID    uint      `gorm:"index" json:"category_id"`
	CategoryName  string    `gorm:"-" json:"category_name"`
	CookTime      int       `json:"cook_time"`
	Difficulty    string    `gorm:"type:varchar(20)" json:"difficulty"`
	PeopleCount   int       `json:"people_count"`
	Taste         string    `gorm:"type:varchar(50)" json:"taste"`
	HealthTags    JSON      `gorm:"type:json" json:"health_tags"`
	Ingredients   JSON      `gorm:"type:json" json:"ingredients"`
	Seasonings    JSON      `gorm:"type:json" json:"seasonings"`
	Steps         JSON      `gorm:"type:json" json:"steps"`
	Tips          string    `gorm:"type:text" json:"tips"`
	Nutrition     JSON      `gorm:"type:json" json:"nutrition"`
	ViewCount     int       `gorm:"default:0" json:"view_count"`
	FavoriteCount int       `gorm:"default:0" json:"favorite_count"`
	Status        int8      `gorm:"default:1;index" json:"status"`
	IsFavorited   bool      `gorm:"-" json:"is_favorited"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Recipe) TableName() string { return "recipes" }
