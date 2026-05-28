package model

import "time"

type Banner struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(100)" json:"title"`
	Image      string    `gorm:"type:varchar(500);not null" json:"image"`
	LinkType   string    `gorm:"type:varchar(20)" json:"link_type"`
	LinkValue  string    `gorm:"type:varchar(500)" json:"link_value"`
	Sort       int       `gorm:"default:0" json:"sort"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Banner) TableName() string { return "banners" }
