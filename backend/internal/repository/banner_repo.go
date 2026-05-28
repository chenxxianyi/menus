package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type BannerRepo struct {
	db *gorm.DB
}

func NewBannerRepo(db *gorm.DB) *BannerRepo {
	return &BannerRepo{db: db}
}

func (r *BannerRepo) FindActive() ([]model.Banner, error) {
	var banners []model.Banner
	err := r.db.Where("status = 1").Order("sort ASC").Find(&banners).Error
	return banners, err
}

func (r *BannerRepo) FindAll() ([]model.Banner, error) {
	var banners []model.Banner
	err := r.db.Order("sort ASC").Find(&banners).Error
	return banners, err
}

func (r *BannerRepo) Create(banner *model.Banner) error {
	return r.db.Create(banner).Error
}

func (r *BannerRepo) Update(banner *model.Banner) error {
	return r.db.Save(banner).Error
}

func (r *BannerRepo) Delete(id uint) error {
	return r.db.Delete(&model.Banner{}, id).Error
}
