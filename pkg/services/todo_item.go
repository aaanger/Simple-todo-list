package services

import (
	"fmt"
	"github.com/aaanger/todo/pkg/models"
	"github.com/aaanger/todo/pkg/repository"
)

type TodoItemService struct {
	repo *repository.Repository
}

func NewTodoItemService(repo *repository.Repository) *TodoItemService {
	return &TodoItemService{
		repo: repo,
	}
}

func (s *TodoItemService) CreateItem(userID, listID int, item models.Item) (int, error) {
	_, err := s.repo.GetListByID(userID, listID)
	if err != nil {
		return 0, fmt.Errorf("service create item: %w", err)
	}

	return s.repo.CreateItem(listID, item)
}

func (s *TodoItemService) GetAllItems(userID, listID int) ([]models.Item, error) {
	return s.repo.GetAllItems(userID, listID)
}

func (s *TodoItemService) GetItemByID(userID, itemID int) (models.Item, error) {
	return s.repo.GetItemByID(userID, itemID)
}

func (s *TodoItemService) UpdateItem(userID, itemID int, input models.UpdateItem) error {
	return s.repo.UpdateItem(userID, itemID, input)
}

func (s *TodoItemService) DeleteItem(userID, itemID int) error {
	return s.repo.DeleteItem(userID, itemID)
}
