package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type FavoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepo(db *gorm.DB) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (r *FavoriteRepo) Add(userID, recipeID uint) error {
	fav := &model.Favorite{UserID: userID, RecipeID: recipeID}
	return r.db.Create(fav).Error
}

func (r *FavoriteRepo) Remove(userID, recipeID uint) error {
	return r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).Delete(&model.Favorite{}).Error
}

func (r *FavoriteRepo) Exists(userID, recipeID uint) bool {
	var count int64
	r.db.Model(&model.Favorite{}).Where("user_id = ? AND recipe_id = ?", userID, recipeID).Count(&count)
	return count > 0
}

func (r *FavoriteRepo) CountByUserID(userID uint) int64 {
	var count int64
	r.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (r *FavoriteRepo) FindByUserID(userID uint, page, pageSize int) ([]model.Favorite, int64, error) {
	var favs []model.Favorite
	var total int64

	query := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&favs).Error
	return favs, total, err
}
