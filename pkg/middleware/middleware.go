package middleware

import (
	"fmt"
	"github.com/aaanger/todo/internal/users/service"
	"github.com/aaanger/todo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		response.ErrorResponse(c, http.StatusUnauthorized, "Empty authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	userID, err := service.ParseToken(headerParts[1])
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
	}

	c.Set("userID", userID)
}

func GetUserID(c *gin.Context) (int, error) {
	id, ok := c.Get("userID")
	if !ok {
		response.ErrorResponse(c, http.StatusBadRequest, "user id not found")
		return 0, fmt.Errorf("user id not found")
	}

	userID, ok := id.(int)
	if !ok {
		response.ErrorResponse(c, http.StatusInternalServerError, "invalid type of user id")
		return 0, fmt.Errorf("invalid type of user id")
	}

	return userID, nil
}
