package services

import (
	"github.com/aaanger/p1/pkg/models"
	"github.com/aaanger/p1/pkg/repository"
)

type TodoListService struct {
	repo *repository.Repository
}

func NewTodoListService(repo *repository.Repository) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (ts *TodoListService) CreateList(userID int, list models.TodoList) (int, error) {
	return ts.repo.CreateList(userID, list)
}

func (ts *TodoListService) GetAllLists(userID int) ([]models.TodoList, error) {
	return ts.repo.GetAllLists(userID)
}

func (ts *TodoListService) GetListByID(userID, listID int) (models.TodoList, error) {
	return ts.repo.GetListByID(userID, listID)
}

func (ts *TodoListService) UpdateList(userID, listID int, input models.UpdateTodoList) error {
	return ts.repo.UpdateList(userID, listID, input)
}

func (ts *TodoListService) DeleteList(userID, listID int) error {
	return ts.repo.DeleteList(userID, listID)
}
