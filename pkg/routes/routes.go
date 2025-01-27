package routes

import (
	"database/sql"
	itemHandler "github.com/aaanger/todo/internal/item/handler"
	itemRepository "github.com/aaanger/todo/internal/item/repository"
	itemService "github.com/aaanger/todo/internal/item/service"
	listHandler "github.com/aaanger/todo/internal/list/handler"
	listRepository "github.com/aaanger/todo/internal/list/repository"
	listService "github.com/aaanger/todo/internal/list/service"
	"github.com/aaanger/todo/internal/users/handler"
	"github.com/aaanger/todo/internal/users/repository"
	"github.com/aaanger/todo/internal/users/service"
	"github.com/aaanger/todo/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func PathHandler(db *sql.DB) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	listRepo := listRepository.NewTodoListRepository(db)
	listSvc := listService.NewTodoListService(listRepo)
	listHandler := listHandler.NewTodoListHandler(listSvc)

	itemRepo := itemRepository.NewTodoItemRepository(db)
	itemSvc := itemService.NewTodoItemService(itemRepo, listRepo)
	itemHandler := itemHandler.NewTodoItemHandler(itemSvc)

	r := gin.New()

	auth := r.Group("/auth")
	auth.POST("/signup", userHandler.SignUp)
	auth.POST("/signin", userHandler.SignIn)

	lists := r.Group("/lists", middleware.UserIdentity)
	lists.POST("/create", listHandler.CreateList)
	lists.GET("/all", listHandler.GetAllLists)
	lists.GET("/:id", listHandler.GetListByID)
	lists.PUT("/:id", listHandler.UpdateList)
	lists.DELETE("/:id", listHandler.DeleteList)

	items := lists.Group("/:id/items")
	items.POST("/newitem", itemHandler.CreateItem)
	items.GET("/", itemHandler.GetAllItems)

	items = r.Group("/items", middleware.UserIdentity)
	items.GET("/:id", itemHandler.GetItemByID)
	items.PUT("/:id", itemHandler.UpdateItem)
	items.DELETE("/:id", itemHandler.DeleteItem)

	return r
}
