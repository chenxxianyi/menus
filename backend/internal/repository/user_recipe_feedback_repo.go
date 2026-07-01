package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"menu-recommend/internal/model"
)

type UserRecipeFeedbackRepo struct {
	db *gorm.DB
}

func NewUserRecipeFeedbackRepo(db *gorm.DB) *UserRecipeFeedbackRepo {
	return &UserRecipeFeedbackRepo{db: db}
}

func (r *UserRecipeFeedbackRepo) Upsert(feedback *model.UserRecipeFeedback) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "recipe_id"},
			{Name: "feedback_type"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"source", "updated_at"}),
	}).Create(feedback).Error
}

func (r *UserRecipeFeedbackRepo) Delete(userID, recipeID uint, feedbackType string) error {
	return r.db.
		Where("user_id = ? AND recipe_id = ? AND feedback_type = ?", userID, recipeID, feedbackType).
		Delete(&model.UserRecipeFeedback{}).Error
}

func (r *UserRecipeFeedbackRepo) FindByUserAndRecipe(userID, recipeID uint) ([]model.UserRecipeFeedback, error) {
	var items []model.UserRecipeFeedback
	err := r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).Find(&items).Error
	return items, err
}

func (r *UserRecipeFeedbackRepo) FindByUser(userID uint) ([]model.UserRecipeFeedback, error) {
	var items []model.UserRecipeFeedback
	err := r.db.Where("user_id = ?", userID).Find(&items).Error
	return items, err
}
