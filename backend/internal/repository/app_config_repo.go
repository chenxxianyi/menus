package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type AppConfigRepo struct {
	db *gorm.DB
}

func NewAppConfigRepo(db *gorm.DB) *AppConfigRepo {
	return &AppConfigRepo{db: db}
}

func (r *AppConfigRepo) FindActiveByKeys(keys []string) ([]model.AppConfig, error) {
	var configs []model.AppConfig
	err := r.db.Where("status = 1 AND config_key IN ?", keys).Find(&configs).Error
	return configs, err
}
