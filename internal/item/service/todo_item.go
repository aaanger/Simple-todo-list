package service

import (
	"fmt"
	"github.com/aaanger/todo/internal/item/model"
	"github.com/aaanger/todo/internal/item/repository"
	listRepository "github.com/aaanger/todo/internal/list/repository"
)

type ITodoItemService interface {
	CreateItem(userID, listID int, item model.Item) (int, error)
	GetAllItems(userID, listID int) ([]model.Item, error)
	GetItemByID(userID, itemID int) (model.Item, error)
	UpdateItem(userID, itemID int, input model.UpdateItem) error
	DeleteItem(userID, itemID int) error
}

type TodoItemService struct {
	repo     *repository.TodoItemRepository
	listRepo *listRepository.TodoListRepository
}

func NewTodoItemService(repo *repository.TodoItemRepository, listRepo *listRepository.TodoListRepository) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) CreateItem(userID, listID int, item model.Item) (int, error) {
	_, err := s.listRepo.GetListByID(userID, listID)
	if err != nil {
		return 0, fmt.Errorf("service create item: %w", err)
	}

	return s.repo.CreateItem(listID, item)
}

func (s *TodoItemService) GetAllItems(userID, listID int) ([]model.Item, error) {
	return s.repo.GetAllItems(userID, listID)
}

func (s *TodoItemService) GetItemByID(userID, itemID int) (model.Item, error) {
	return s.repo.GetItemByID(userID, itemID)
}

func (s *TodoItemService) UpdateItem(userID, itemID int, input model.UpdateItem) error {
	return s.repo.UpdateItem(userID, itemID, input)
}

func (s *TodoItemService) DeleteItem(userID, itemID int) error {
	return s.repo.DeleteItem(userID, itemID)
}
