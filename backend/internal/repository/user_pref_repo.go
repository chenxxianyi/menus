package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type UserPrefRepo struct {
	db *gorm.DB
}

func NewUserPrefRepo(db *gorm.DB) *UserPrefRepo {
	return &UserPrefRepo{db: db}
}

func (r *UserPrefRepo) FindByUserID(userID uint) (*model.UserPreference, error) {
	var pref model.UserPreference
	err := r.db.Where("user_id = ?", userID).First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &pref, err
}

func (r *UserPrefRepo) Save(pref *model.UserPreference) error {
	var existing model.UserPreference
	err := r.db.Where("user_id = ?", pref.UserID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(pref).Error
	}
	pref.ID = existing.ID
	return r.db.Save(pref).Error
}
