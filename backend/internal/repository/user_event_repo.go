package repository

import (
	"gorm.io/gorm"

	"menu-recommend/internal/model"
)

type UserEventRepo struct {
	db *gorm.DB
}

func NewUserEventRepo(db *gorm.DB) *UserEventRepo {
	return &UserEventRepo{db: db}
}

func (r *UserEventRepo) Create(event *model.UserEvent) error {
	return r.db.Create(event).Error
}
