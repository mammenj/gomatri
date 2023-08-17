package handlers

import (
	"fmt"
	"gomatri/models"
	"gomatri/security"
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

	//fmt.Println("USER List", users)
	c.JSON(http.StatusOK, &users)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	log.Println("IN PATCH  handler")
	var user models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = security.HashAndSalt([]byte(user.Password))
	log.Println("IN PATCH  handler ", &user)
	ID, err := uh.store.Update(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("USER Updated ID: ", ID)
	c.JSON(http.StatusOK, &user)
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
	c.JSON(http.StatusOK, "Deleted User ID: "+ID)
}
func (uh *UserHandler) GetUser(c *gin.Context) {
	log.Println("IN GET one handler")
	id := c.Param("id")
	value := c.GetHeader("Authorization")
	log.Println("GetHeader in UserHandler controller :: ", value)
	user, err := uh.store.GetOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("...... Get user: ", user)
	c.JSON(http.StatusOK, &user)
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	log.Println("IN Create handler User")

	log.Println("CreateUser in db")

	var user models.User
	if err := c.Bind(&user); err != nil {
		//c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		c.String(http.StatusOK, err.Error(), nil)
		return
	}
	log.Println("CreateUser bound user", user)
	user.Password = security.HashAndSalt([]byte(user.Password))
	ID, err := uh.store.Create(&user)
	log.Println("CreateUser user hashed password", user.Password)
	if err != nil {
		c.String(http.StatusOK, err.Error(), nil)
		return
	}

	fmt.Println("USER CREATED ID: ", ID)
	c.String(http.StatusOK, "User created ID:: "+ID, "Keep the ID for reference")
}
