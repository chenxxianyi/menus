package repository

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"menu-recommend/internal/model"
)

type BrowseHistoryRepo struct {
	db *gorm.DB
}

func NewBrowseHistoryRepo(db *gorm.DB) *BrowseHistoryRepo {
	return &BrowseHistoryRepo{db: db}
}

func (r *BrowseHistoryRepo) Record(userID, recipeID uint) error {
	if userID == 0 || recipeID == 0 {
		return nil
	}

	history := &model.BrowseHistory{
		UserID:   userID,
		RecipeID: recipeID,
		ViewedAt: time.Now(),
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "recipe_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"viewed_at", "updated_at"}),
	}).Create(history).Error
}

func (r *BrowseHistoryRepo) List(userID uint, page, pageSize int) ([]model.BrowseHistory, int64, error) {
	var histories []model.BrowseHistory
	var total int64

	query := r.db.Model(&model.BrowseHistory{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.
		Preload("Recipe").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("viewed_at DESC, id DESC").
		Find(&histories).Error

	return histories, total, err
}

func (r *BrowseHistoryRepo) Clear(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.BrowseHistory{}).Error
}
