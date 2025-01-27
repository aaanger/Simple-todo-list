package service

import (
	"github.com/aaanger/todo/internal/list/model"
	"github.com/aaanger/todo/internal/list/repository"
)

type ITodoListService interface {
	CreateList(userID int, list model.TodoList) (int, error)
	GetAllLists(userID int) ([]model.TodoList, error)
	GetListByID(userID, listID int) (model.TodoList, error)
	UpdateList(userID, listID int, input model.UpdateTodoList) error
	DeleteList(userID, listID int) error
}

type TodoListService struct {
	repo *repository.TodoListRepository
}

func NewTodoListService(repo *repository.TodoListRepository) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (ts *TodoListService) CreateList(userID int, list model.TodoList) (int, error) {
	return ts.repo.CreateList(userID, list)
}

func (ts *TodoListService) GetAllLists(userID int) ([]model.TodoList, error) {
	return ts.repo.GetAllLists(userID)
}

func (ts *TodoListService) GetListByID(userID, listID int) (model.TodoList, error) {
	return ts.repo.GetListByID(userID, listID)
}

func (ts *TodoListService) UpdateList(userID, listID int, input model.UpdateTodoList) error {
	return ts.repo.UpdateList(userID, listID, input)
}

func (ts *TodoListService) DeleteList(userID, listID int) error {
	return ts.repo.DeleteList(userID, listID)
}
