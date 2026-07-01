package service

import (
	"encoding/json"
	"strings"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepo
	prefRepo *repository.UserPrefRepo
}

func NewUserService(userRepo *repository.UserRepo, prefRepo *repository.UserPrefRepo) *UserService {
	return &UserService{userRepo: userRepo, prefRepo: prefRepo}
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID uint, nickname, avatar string, gender int8) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender > 0 {
		user.Gender = gender
	}
	return s.userRepo.Update(user)
}

func (s *UserService) GetPreferences(userID uint) (*model.UserPreference, error) {
	pref, err := s.prefRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		return &model.UserPreference{UserID: userID}, nil
	}
	return pref, nil
}

func (s *UserService) UpdatePreferences(pref *model.UserPreference) error {
	return s.prefRepo.Save(pref)
}

type PreferenceStatus struct {
	Completed     bool     `json:"completed"`
	Completeness  int      `json:"completeness"`
	MissingFields []string `json:"missing_fields"`
}

func (s *UserService) GetPreferenceStatus(userID uint) (*PreferenceStatus, error) {
	pref, err := s.prefRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	status := &PreferenceStatus{}
	required := []struct {
		key string
		ok  bool
	}{
		{key: "people_count", ok: pref != nil && pref.PeopleCount > 0},
		{key: "taste_preference", ok: pref != nil && jsonArrayHasItems(pref.TastePreference)},
		{key: "cook_time_preference", ok: pref != nil && strings.TrimSpace(pref.CookTimePreference) != ""},
		{key: "health_goal", ok: pref != nil && strings.TrimSpace(pref.HealthGoal) != ""},
	}

	completedCount := 0
	for _, item := range required {
		if item.ok {
			completedCount++
		} else {
			status.MissingFields = append(status.MissingFields, item.key)
		}
	}
	status.Completed = completedCount == len(required)
	status.Completeness = completedCount * 100 / len(required)
	return status, nil
}

func jsonArrayHasItems(raw model.JSON) bool {
	if len(raw) == 0 || string(raw) == "null" {
		return false
	}
	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return false
	}
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}

func (s *UserService) ListUsers(keyword string, page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(keyword, page, pageSize)
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}
