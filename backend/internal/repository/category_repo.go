package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) FindAll() ([]model.RecipeCategory, error) {
	var categories []model.RecipeCategory
	err := r.db.Where("status = 1").Order("sort ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepo) FindByID(id uint) (*model.RecipeCategory, error) {
	var cat model.RecipeCategory
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *CategoryRepo) Create(cat *model.RecipeCategory) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepo) Update(cat *model.RecipeCategory) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.RecipeCategory{}, id).Error
}
