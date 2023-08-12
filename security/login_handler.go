package security

import (
	"gomatri/models"
	"gomatri/storage"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwt_lib "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"errors"
	//"fmt"
	"log"
	//"net/http"
	"time"
)

// Login for JWT
func Login(c *gin.Context) {
	// get the user from db
	//var user *models.User
	user, err := loginUser(c)

	if err != nil {
		log.Println("Error when  trying to log in :", err)
		//c.JSON(200, gin.H{"message": "Invalid User or Password"})
		//errMessage := "Invalid user/password for " + user.Email
		//c.HTML(http.StatusOK, "templates/user.html", err)
		c.String(http.StatusOK, err.Error(), nil)
		return
	}
	// Create the token
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	// Set some claims
	// role := "user"
	// if user.Username == "admin" {
	// 	role = "admin"
	// }
	token.Claims = jwt_lib.MapClaims{
		"username": user.Email,
		"expiry":   time.Now().Add(time.Hour * 1).Unix(),
		"role":     user.Role,
	}
	// Sign and get the complete encoded token as a string
	//config, err := config.GetConfiguration("config.json")
	// if err != nil {
	// 	log.Fatal(err)
	// 	c.JSON(500, gin.H{"message": "Error getting configuration"})
	// 	return
	// }
	//tokenString, err := token.SignedString([]byte(config.Jwtsecret))

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Could not generate token")
		c.JSON(500, gin.H{"message": "Could not generate token"})
		return
	}
	c.Header("token", tokenString)
	//c.JSON(200, gin.H{"token": tokenString})
	c.String(http.StatusOK, tokenString, nil)
}

// Login get the User by username/password
func loginUser(c *gin.Context) (*models.User, error) {
	log.Println("loginUser username :")
	var user models.User
	if err := c.Bind(&user); err != nil {
		// c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		// 	"error": err.Error(),
		// })
		log.Println("Failed to bind user from context", user)
		return nil, err
	}
	log.Println("loginUser username :", user)
	username := user.Email
	password := user.Password

	store := storage.NewSqliteUserStore()
	if err := store.DB.Where("email = ?", username).First(&user).Error; err != nil {
		log.Println("Failed to GetUser in db")
		//c.AbortWithStatus(http.StatusNotFound)
		return nil, err
	}

	// db := db.GetDB()
	// if err := db.Where("username = ?", username).First(&user).Error; err != nil {
	// 	log.Println("Failed to GetUser in db")
	// 	c.AbortWithStatus(http.StatusNotFound)
	// 	return user, err
	// }
	success := comparePasswords(user.Password, []byte(password))
	if !success {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}
