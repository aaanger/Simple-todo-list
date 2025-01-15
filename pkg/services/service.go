package services

import (
	"github.com/aaanger/todo/pkg/models"
	"github.com/aaanger/todo/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) (int, error)
	AuthUser(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	CreateList(userID int, list models.TodoList) (int, error)
	GetAllLists(userID int) ([]models.TodoList, error)
	GetListByID(userID, listID int) (models.TodoList, error)
	UpdateList(userID, listID int, input models.UpdateTodoList) error
	DeleteList(userID, listID int) error
}

type TodoItem interface {
	CreateItem(userID, listID int, item models.Item) (int, error)
	GetAllItems(userID, listID int) ([]models.Item, error)
	GetItemByID(userID, itemID int) (models.Item, error)
	UpdateItem(userID, itemID int, input models.UpdateItem) error
	DeleteItem(userID, itemID int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewUserService(repo),
		TodoList:      NewTodoListService(repo),
		TodoItem:      NewTodoItemService(repo),
	}
}
