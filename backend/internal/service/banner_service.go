package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type BannerService struct {
	repo *repository.BannerRepo
}

func NewBannerService(repo *repository.BannerRepo) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) GetActive() ([]model.Banner, error) {
	return s.repo.FindActive()
}

func (s *BannerService) GetAll() ([]model.Banner, error) {
	return s.repo.FindAll()
}

func (s *BannerService) Create(banner *model.Banner) error {
	return s.repo.Create(banner)
}

func (s *BannerService) Update(banner *model.Banner) error {
	return s.repo.Update(banner)
}

func (s *BannerService) Delete(id uint) error {
	return s.repo.Delete(id)
}
