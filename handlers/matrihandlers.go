package handlers

import (
	"fmt"
	"gomatri/models"
	"gomatri/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	log.Println("IN GET handler")
	userStore := storage.NewSqliteUserStore()
	users, err := userStore.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER List", users)
}

func UpdateUser(c *gin.Context) {
	log.Println("IN PATCH  handler")
	var input models.User
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("IN PATCH  handler ", &input)
	userStore := storage.NewSqliteUserStore()
	ID, err := userStore.Update(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER Updated ID: ", ID)

}

func DeleteUser(c *gin.Context) {
	log.Println("IN Delete handler")

	id := c.Param("id")

	userStore := storage.NewSqliteUserStore()
	ID, err := userStore.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER deleted ID: ", ID)
}
func GetUser(c *gin.Context) {
	log.Println("IN GET one handler")
	id := c.Param("id")

	userStore := storage.NewSqliteUserStore()
	user, err := userStore.GetOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("...... Get user: ", user)
}

func CreateUser(c *gin.Context) {
	log.Println("IN Create handler")
	var input models.User
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStore := storage.NewSqliteUserStore()
	ID, err := userStore.Create(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER CREATED ID: ", ID)
}
