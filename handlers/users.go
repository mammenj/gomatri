package handlers

import (
	"gomatri/models"
	"gomatri/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	store storage.UserSqlliteStore
}

func CreateUsersHandler(store *storage.UserSqlliteStore) *UsersHandler {
	return &UsersHandler{store: *store}
}

func (u UsersHandler) Create(c *gin.Context) {
	// Valiate input

	var input models.User
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID, err := u.store.Create(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "A New User Created with ID: " + ID})
}
