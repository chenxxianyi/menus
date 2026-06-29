package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type MenuRepo struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{db: db}
}

func (r *MenuRepo) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *MenuRepo) CountByUserID(userID uint) int64 {
	var count int64
	r.db.Model(&model.Menu{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (r *MenuRepo) FindByUserID(userID uint, page, pageSize int) ([]model.Menu, int64, error) {
	var menus []model.Menu
	var total int64

	query := r.db.Model(&model.Menu{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&menus).Error
	return menus, total, err
}
