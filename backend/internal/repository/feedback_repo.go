package repository

import (
	"gorm.io/gorm"
	"menu-recommend/internal/model"
)

type FeedbackRepo struct {
	db *gorm.DB
}

func NewFeedbackRepo(db *gorm.DB) *FeedbackRepo {
	return &FeedbackRepo{db: db}
}

func (r *FeedbackRepo) Create(fb *model.Feedback) error {
	return r.db.Create(fb).Error
}

func (r *FeedbackRepo) List(status *int8, page, pageSize int) ([]model.Feedback, int64, error) {
	var items []model.Feedback
	var total int64

	query := r.db.Model(&model.Feedback{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&items).Error
	return items, total, err
}

func (r *FeedbackRepo) UpdateStatus(id uint, status int8, reply string) error {
	return r.db.Model(&model.Feedback{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"reply":  reply,
	}).Error
}
