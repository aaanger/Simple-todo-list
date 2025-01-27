package handler

import (
	"github.com/aaanger/todo/internal/list/model"
	"github.com/aaanger/todo/internal/list/service"
	"github.com/aaanger/todo/pkg/middleware"
	"github.com/aaanger/todo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TodoListHandler struct {
	service service.ITodoListService
}

func NewTodoListHandler(service service.ITodoListService) *TodoListHandler {
	return &TodoListHandler{
		service: service,
	}
}

func (h *TodoListHandler) CreateList(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var list model.TodoList

	err = c.BindJSON(&list)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listID, err := h.service.CreateList(userID, list)

	c.JSON(http.StatusOK, map[string]interface{}{
		"listID": listID,
	})
}

func (h *TodoListHandler) GetAllLists(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	lists, err := h.service.GetAllLists(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string][]model.TodoList{
		"lists": lists,
	})
}

func (h *TodoListHandler) GetListByID(c *gin.Context) {
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

	list, err := h.service.GetListByID(userID, listID)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
}

func (h *TodoListHandler) UpdateList(c *gin.Context) {
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

	var input model.UpdateTodoList
	err = c.BindJSON(&input)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdateList(userID, listID, input)

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.service.GetListByID(userID, listID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"updated list": list,
	})
}

func (h *TodoListHandler) DeleteList(c *gin.Context) {
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

	err = h.service.DeleteList(userID, listID)

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
