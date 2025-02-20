package service

import (
	"errors"
	"fmt"
	"github.com/aaanger/todo/internal/users/model"
	"github.com/aaanger/todo/internal/users/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	signingTokenKey = "joASdeDS3i#kjmFDSk3i303904lXSDds"
	tokenExpire     = 12 * time.Hour
)

type IUserService interface {
	CreateUser(user model.User) (int, error)
	AuthUser(username, password string) (string, error)
}

type UserService struct {
	repo *repository.UserRepository
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"id"`
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) CreateUser(user model.User) (int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("service create user: %w", err)
	}
	user.Password = string(passwordHash)
	return us.repo.CreateUser(user)
}

func (us *UserService) AuthUser(username, password string) (string, error) {
	user, err := us.repo.AuthUser(username, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, user.ID,
	})
	return token.SignedString([]byte(signingTokenKey))
}

func ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingTokenKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("parse token: %w", err)
	}
	return claims.UserID, nil
}
