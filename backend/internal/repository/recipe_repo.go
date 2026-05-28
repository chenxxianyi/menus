package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type RecipeRepo struct {
	db *gorm.DB
}

func NewRecipeRepo(db *gorm.DB) *RecipeRepo {
	return &RecipeRepo{db: db}
}

func (r *RecipeRepo) FindByID(id uint) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.db.First(&recipe, id).Error
	return &recipe, err
}

func (r *RecipeRepo) List(keyword string, categoryID uint, taste, cookTime, difficulty, healthTags string, page, pageSize int) ([]model.Recipe, int64, error) {
	var recipes []model.Recipe
	var total int64

	query := r.db.Model(&model.Recipe{}).Where("status = 1")
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if taste != "" {
		query = query.Where("taste = ?", taste)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&recipes).Error
	return recipes, total, err
}

func (r *RecipeRepo) FindByIDs(ids []uint) ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := r.db.Where("id IN ?", ids).Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepo) FindRandom(limit int) ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := r.db.Where("status = 1").Order("RAND()").Limit(limit).Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepo) FindHot(limit int) ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := r.db.Where("status = 1").Order("favorite_count DESC, view_count DESC").Limit(limit).Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepo) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Recipe{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *RecipeRepo) IncrementFavoriteCount(id uint) error {
	return r.db.Model(&model.Recipe{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

func (r *RecipeRepo) DecrementFavoriteCount(id uint) error {
	return r.db.Model(&model.Recipe{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - 1, 0)")).Error
}

func (r *RecipeRepo) Create(recipe *model.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *RecipeRepo) Update(recipe *model.Recipe) error {
	return r.db.Save(recipe).Error
}

func (r *RecipeRepo) Delete(id uint) error {
	return r.db.Delete(&model.Recipe{}, id).Error
}
