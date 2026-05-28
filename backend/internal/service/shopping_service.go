package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type ShoppingService struct {
	repo *repository.ShoppingRepo
}

func NewShoppingService(repo *repository.ShoppingRepo) *ShoppingService {
	return &ShoppingService{repo: repo}
}

func (s *ShoppingService) GetLists(userID uint) ([]model.ShoppingList, error) {
	return s.repo.FindByUserID(userID)
}

func (s *ShoppingService) CreateList(list *model.ShoppingList) error {
	return s.repo.Create(list)
}

func (s *ShoppingService) UpdateList(list *model.ShoppingList) error {
	return s.repo.Update(list)
}

func (s *ShoppingService) DeleteList(id uint) error {
	return s.repo.Delete(id)
}

func (s *ShoppingService) GetListByID(id uint) (*model.ShoppingList, error) {
	return s.repo.FindByID(id)
}
