package model

import "time"

type UserRecipeFeedback struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null;uniqueIndex:idx_user_recipe_feedback;index" json:"user_id"`
	RecipeID     uint      `gorm:"not null;uniqueIndex:idx_user_recipe_feedback;index" json:"recipe_id"`
	FeedbackType string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_user_recipe_feedback" json:"feedback_type"`
	Source       string    `gorm:"type:varchar(30)" json:"source"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (UserRecipeFeedback) TableName() string { return "user_recipe_feedbacks" }
