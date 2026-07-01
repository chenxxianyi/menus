package service

import (
	"encoding/json"
	"strings"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type UserEventService struct {
	repo *repository.UserEventRepo
}

func NewUserEventService(repo *repository.UserEventRepo) *UserEventService {
	return &UserEventService{repo: repo}
}

func (s *UserEventService) Track(userID uint, eventName, entityType string, entityID uint, payload map[string]interface{}) error {
	eventName = strings.TrimSpace(eventName)
	if eventName == "" {
		eventName = "unknown"
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		payloadJSON = []byte(`{}`)
	}
	return s.repo.Create(&model.UserEvent{
		UserID:      userID,
		EventName:   eventName,
		EntityType:  strings.TrimSpace(entityType),
		EntityID:    entityID,
		PayloadJSON: model.JSON(payloadJSON),
	})
}
