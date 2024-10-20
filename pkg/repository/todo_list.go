package repository

import (
	"fmt"
	"github.com/aaanger/p1/pkg/models"
	"github.com/sirupsen/logrus"
	"strings"
)

func (r *Repository) CreateList(userID int, list models.TodoList) (int, error) {
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

func (r *Repository) GetAllLists(userID int) ([]models.TodoList, error) {
	var lists []models.TodoList

	rows, err := r.DB.Query(`SELECT title, description FROM lists l INNER JOIN user_lists ul ON l.id = ul.list_id WHERE ul.user_id=$1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("repository get all lists: %w", err)
	}

	for rows.Next() {
		var list models.TodoList

		err = rows.Scan(&list.Title, &list.Description)
		if err != nil {
			return nil, fmt.Errorf("repository get all lists: %w", err)
		}
		lists = append(lists, list)
	}

	return lists, rows.Err()
}

func (r *Repository) GetListByID(userID, listID int) (models.TodoList, error) {
	list := models.TodoList{
		ID: listID,
	}
	row := r.DB.QueryRow(`SELECT title, description FROM lists l INNER JOIN user_lists ul ON l.id = ul.list_id WHERE ul.list_id=$1 AND ul.user_id=$2;`, listID, userID)

	err := row.Scan(&list.Title, &list.Description)
	if err != nil {
		return models.TodoList{}, fmt.Errorf("repository get list by id: %w", err)
	}

	return list, nil
}

func (r *Repository) UpdateList(userID, listID int, input models.UpdateTodoList) error {
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

func (r *Repository) DeleteList(userID, listID int) error {
	_, err := r.DB.Exec(`DELETE FROM lists l USING user_lists ul WHERE l.id = ul.list_id AND l.id=$1 AND ul.user_id=$2`, listID, userID)
	return err
}
