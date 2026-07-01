package repository

import (
	"menu-recommend/internal/model"

	"gorm.io/gorm"
)

type AIGenerationLogRepo struct {
	db *gorm.DB
}

func NewAIGenerationLogRepo(db *gorm.DB) *AIGenerationLogRepo {
	return &AIGenerationLogRepo{db: db}
}

func (r *AIGenerationLogRepo) Create(log *model.AIGenerationLog) error {
	return r.db.Create(log).Error
}
