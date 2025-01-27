package handler

import (
	"github.com/aaanger/todo/internal/item/model"
	"github.com/aaanger/todo/internal/item/service"
	"github.com/aaanger/todo/pkg/middleware"
	"github.com/aaanger/todo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TodoItemHandler struct {
	service service.ITodoItemService
}

func NewTodoItemHandler(service service.ITodoItemService) *TodoItemHandler {
	return &TodoItemHandler{
		service: service,
	}
}

func (h *TodoItemHandler) CreateItem(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.Item

	err = c.BindJSON(&input)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemID, err := h.service.CreateItem(userID, listID, input)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"item id": itemID,
	})
}

func (h *TodoItemHandler) GetAllItems(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	items, err := h.service.GetAllItems(userID, listID)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string][]model.Item{
		"items": items,
	})
}

func (h *TodoItemHandler) GetItemByID(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	item, err := h.service.GetItemByID(userID, itemID)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"items": item,
	})
}

func (h *TodoItemHandler) UpdateItem(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.UpdateItem
	err = c.BindJSON(&input)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdateItem(userID, itemID, input)

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.service.GetItemByID(userID, itemID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"updated item": item,
	})
}

func (h *TodoItemHandler) DeleteItem(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.service.DeleteItem(userID, itemID)

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
