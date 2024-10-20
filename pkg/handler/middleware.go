package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) UserIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	userID, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid token")
	}

	c.Set("userID", userID)
}

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "user id not found")
		return 0, fmt.Errorf("user id not found")
	}

	userID, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "invalid type of user id")
		return 0, fmt.Errorf("invalid type of user id")
	}

	return userID, nil
}
