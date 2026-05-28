package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type IngredientRepo struct {
	db *gorm.DB
}

func NewIngredientRepo(db *gorm.DB) *IngredientRepo {
	return &IngredientRepo{db: db}
}

func (r *IngredientRepo) FindAll() ([]model.Ingredient, error) {
	var items []model.Ingredient
	err := r.db.Find(&items).Error
	return items, err
}

func (r *IngredientRepo) FindByCategory(category string) ([]model.Ingredient, error) {
	var items []model.Ingredient
	err := r.db.Where("category = ?", category).Find(&items).Error
	return items, err
}

func (r *IngredientRepo) Create(item *model.Ingredient) error {
	return r.db.Create(item).Error
}

func (r *IngredientRepo) Update(item *model.Ingredient) error {
	return r.db.Save(item).Error
}

func (r *IngredientRepo) Delete(id uint) error {
	return r.db.Delete(&model.Ingredient{}, id).Error
}

func (r *IngredientRepo) List(keyword, category string, page, pageSize int) ([]model.Ingredient, int64, error) {
	var items []model.Ingredient
	var total int64

	query := r.db.Model(&model.Ingredient{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id ASC").Find(&items).Error
	return items, total, err
}
