package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.RecipeCategory, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) Create(cat *model.RecipeCategory) error {
	return s.repo.Create(cat)
}

func (s *CategoryService) Update(cat *model.RecipeCategory) error {
	return s.repo.Update(cat)
}

func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}
