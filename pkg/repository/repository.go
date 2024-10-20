package repository

import (
	"database/sql"
	"fmt"
	"github.com/aaanger/p1/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) CreateUser(user models.User) (int, error) {
	row := r.DB.QueryRow(`INSERT INTO users (username, password_hash) VALUES($1, $2) RETURNING id;`, user.Username, user.Password)
	err := row.Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("repository create user: %w", err)
	}
	return user.ID, nil
}

func (r *Repository) AuthUser(username, password string) (*models.User, error) {
	user := models.User{
		Username: username,
	}
	row := r.DB.QueryRow(`SELECT id, password_hash FROM users WHERE username=$1`, username)
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("repository authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("repository authenticate: %w", err)
	}
	return &user, nil
}
