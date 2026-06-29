package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type ShoppingRepo struct {
	db *gorm.DB
}

func NewShoppingRepo(db *gorm.DB) *ShoppingRepo {
	return &ShoppingRepo{db: db}
}

func (r *ShoppingRepo) FindByUserID(userID uint) ([]model.ShoppingList, error) {
	var lists []model.ShoppingList
	err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&lists).Error
	return lists, err
}

func (r *ShoppingRepo) CountByUserID(userID uint) int64 {
	var count int64
	r.db.Model(&model.ShoppingList{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (r *ShoppingRepo) FindByID(id uint) (*model.ShoppingList, error) {
	var list model.ShoppingList
	err := r.db.First(&list, id).Error
	return &list, err
}

func (r *ShoppingRepo) Create(list *model.ShoppingList) error {
	return r.db.Create(list).Error
}

func (r *ShoppingRepo) Update(list *model.ShoppingList) error {
	return r.db.Save(list).Error
}

func (r *ShoppingRepo) Delete(id uint) error {
	return r.db.Delete(&model.ShoppingList{}, id).Error
}
