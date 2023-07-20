package handlers

import (
	"fmt"
	"gomatri/models"
	"gomatri/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store *storage.UserSqlliteStore
}

func CreateNewUserHandler() *UserHandler {
	return &UserHandler{
		store: storage.NewSqliteUserStore(),
	}
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	log.Println("IN GET handler")

	users, err := uh.store.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER List", users)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	log.Println("IN PATCH  handler")
	var input models.User
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("IN PATCH  handler ", &input)
	ID, err := uh.store.Update(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER Updated ID: ", ID)

}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	log.Println("IN Delete handler")

	id := c.Param("id")

	ID, err := uh.store.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER deleted ID: ", ID)
}
func (uh *UserHandler) GetUser(c *gin.Context) {
	log.Println("IN GET one handler")
	id := c.Param("id")

	user, err := uh.store.GetOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("...... Get user: ", user)
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	log.Println("IN Create handler")
	var input models.User
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//userStore := storage.NewSqliteUserStore()
	ID, err := uh.store.Create(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER CREATED ID: ", ID)
}
