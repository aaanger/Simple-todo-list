package handler

import (
	"github.com/aaanger/todo/internal/users/model"
	"github.com/aaanger/todo/internal/users/service"
	"github.com/aaanger/todo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}
	userID, err := h.service.CreateUser(user)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": userID,
	})
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.service.AuthUser(user.Username, user.Password)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
