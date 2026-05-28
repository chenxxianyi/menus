package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone        string    `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string    `gorm:"type:varchar(50)" json:"nickname"`
	Avatar       string    `gorm:"type:varchar(500)" json:"avatar"`
	Gender       int8      `gorm:"default:0" json:"gender"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
