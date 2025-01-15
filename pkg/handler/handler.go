package handler

import (
	"github.com/aaanger/todo/pkg/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) PathHandler() *gin.Engine {
	r := gin.New()

	auth := r.Group("/auth")
	auth.POST("/signup", h.signUp)
	auth.POST("/signin", h.signIn)

	lists := r.Group("/lists", h.UserIdentity)
	lists.POST("/create", h.createList)
	lists.GET("/all", h.getAllLists)
	lists.GET("/:id", h.getListByID)
	lists.PUT("/:id", h.updateList)
	lists.DELETE("/:id", h.deleteList)

	items := lists.Group("/:id/items")
	items.POST("/newitem", h.createItem)
	items.GET("/", h.getAllItems)

	items = r.Group("/items", h.UserIdentity)
	items.GET("/:id")
	items.PUT("/:id")
	items.DELETE("/:id")

	return r
}
