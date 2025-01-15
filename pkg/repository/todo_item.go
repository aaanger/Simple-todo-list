package repository

import (
	"fmt"
	"github.com/aaanger/todo/pkg/models"
	"github.com/sirupsen/logrus"
	"strings"
)

func (r *Repository) CreateItem(listID int, item models.Item) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("repository create item: %w", err)
	}

	row := tx.QueryRow(`INSERT INTO items (title, description) VALUES($1, $2) RETURNING id;`, item.Title, item.Description)
	err = row.Scan(&item.ID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repository create item: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO list_items (list_id, item_id) VALUES($1, $2);`, listID, item.ID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repository create item: %w", err)
	}

	return item.ID, tx.Commit()
}

func (r *Repository) GetAllItems(userID, listID int) ([]models.Item, error) {
	var items []models.Item

	rows, err := r.DB.Query(`SELECT title, description, done FROM items i INNER JOIN list_items li ON i.id = li.item_id 
    INNER JOIN user_lists ul ON ul.list_id = li.list_id WHERE li.list_id=$1 AND ul.user_id=$2`, listID, userID)
	if err != nil {
		return nil, fmt.Errorf("repository get all items: %w", err)
	}

	for rows.Next() {
		var item models.Item

		err = rows.Scan(&item.Title, &item.Description, &item.Done)
		if err != nil {
			return nil, fmt.Errorf("repository get all items: %w", err)
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) GetItemByID(userID, itemID int) (models.Item, error) {
	var item models.Item

	row := r.DB.QueryRow(`SELECT title, description, done FROM items i INNER JOIN list_items li ON i.id = li.item_id 
    INNER JOIN user_lists ul ON li.list_id = ul.list_id WHERE ul.user_id = $1 AND i.id = $2`, userID, itemID)
	err := row.Scan(&item.Title, &item.Description, &item.Done)
	if err != nil {
		return item, fmt.Errorf("repository get all items: %w", err)
	}

	return item, nil
}

func (r *Repository) UpdateItem(userID, itemID int, input models.UpdateItem) error {
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

	if input.Done != nil {
		keys = append(keys, fmt.Sprintf(`done=$%d`, arg))
		values = append(values, *input.Done)
		arg++
	}

	joinQuery := strings.Join(keys, ", ")

	query := fmt.Sprintf(`UPDATE items i SET %s FROM list_items li, user_lists ul WHERE i.id = li.item_id AND li.list_id = ul.list_id AND i.id=$%d AND user_id=$%d;`, joinQuery, arg, arg+1)
	values = append(values, itemID, userID)

	logrus.Debugf("update query: %s", query)
	logrus.Debugf("values: %s", values)

	_, err := r.DB.Exec(query, values...)

	return err
}

func (r *Repository) DeleteItem(userID, itemID int) error {
	_, err := r.DB.Exec(`DELETE FROM items i USING list_items li, user_lists ul WHERE li.item_id = i.id AND li.list_id = ul.list_id AND i.id=$1 AND ul.user_id=$2`, itemID, userID)
	return err
}
