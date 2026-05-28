package model

import "time"

type AdminUser struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(20);default:admin" json:"role"`
	Status       int8      `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func (AdminUser) TableName() string { return "admin_users" }
