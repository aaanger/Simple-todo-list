package repository

import (
	"database/sql"
	"fmt"
	"github.com/aaanger/todo/internal/list/model"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListRepository struct {
	DB *sql.DB
}

func NewTodoListRepository(db *sql.DB) *TodoListRepository {
	return &TodoListRepository{
		DB: db,
	}
}

func (r *TodoListRepository) CreateList(userID int, list model.TodoList) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("repository create list: %w", err)
	}

	row := tx.QueryRow(`INSERT INTO lists (title, description) VALUES($1, $2) RETURNING id;`, list.Title, list.Description)
	err = row.Scan(&list.ID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repository create list: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO user_lists (user_id, list_id) VALUES($1, $2);`, userID, list.ID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repository create list: %w", err)
	}

	return list.ID, tx.Commit()
}

func (r *TodoListRepository) GetAllLists(userID int) ([]model.TodoList, error) {
	var lists []model.TodoList

	rows, err := r.DB.Query(`SELECT title, description FROM lists l INNER JOIN user_lists ul ON l.id = ul.list_id WHERE ul.user_id=$1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("repository get all lists: %w", err)
	}

	for rows.Next() {
		var list model.TodoList

		err = rows.Scan(&list.Title, &list.Description)
		if err != nil {
			return nil, fmt.Errorf("repository get all lists: %w", err)
		}
		lists = append(lists, list)
	}

	return lists, rows.Err()
}

func (r *TodoListRepository) GetListByID(userID, listID int) (model.TodoList, error) {
	list := model.TodoList{
		ID: listID,
	}
	row := r.DB.QueryRow(`SELECT title, description FROM lists l INNER JOIN user_lists ul ON l.id = ul.list_id WHERE ul.list_id=$1 AND ul.user_id=$2;`, listID, userID)

	err := row.Scan(&list.Title, &list.Description)
	if err != nil {
		return model.TodoList{}, fmt.Errorf("repository get list by id: %w", err)
	}

	return list, nil
}

func (r *TodoListRepository) UpdateList(userID, listID int, input model.UpdateTodoList) error {
	keys := make([]string, 0)
	values := make([]interface{}, 0)

	arg := 1

	if input.Title != nil {
		keys = append(keys, fmt.Sprintf(`title=$%d`, arg))
		values = append(values, *input.Title)
		arg++
	}

	if input.Description != nil {
		keys = append(keys, fmt.Sprintf(`description=$%d`, arg))
		values = append(values, *input.Description)
		arg++
	}

	joinQuery := strings.Join(keys, ", ")

	query := fmt.Sprintf(`UPDATE lists SET %s WHERE id=$%d AND user_id=$%d;`, joinQuery, arg, arg+1)
	values = append(values, listID, userID)

	logrus.Debugf("update query: %s", query)
	logrus.Debugf("values: %s", values)

	_, err := r.DB.Exec(query, values...)

	return err
}

func (r *TodoListRepository) DeleteList(userID, listID int) error {
	_, err := r.DB.Exec(`DELETE FROM lists l USING user_lists ul WHERE l.id = ul.list_id AND l.id=$1 AND ul.user_id=$2`, listID, userID)
	return err
}
