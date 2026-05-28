package service

import (
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

func (s *UserService) ListUsers(keyword string, page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(keyword, page, pageSize)
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}
