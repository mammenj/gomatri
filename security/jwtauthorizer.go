package security

import (
	"fmt"
	"gomatri/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/casbin/casbin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// NewJwtAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewJwtAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &JwtAuthorizer{enforcer: e}

	return func(c *gin.Context) {
		if !a.CheckPermission(c) {
			a.RequirePermission(c)
		}
	}
}

// JwtAuthorizer stores the casbin handler
type JwtAuthorizer struct {
	enforcer *casbin.Enforcer
}

func NewJwtAuth(e *casbin.Enforcer) *JwtAuthorizer {
	return &JwtAuthorizer{e}

}

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *JwtAuthorizer) GetUserName(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *JwtAuthorizer) CheckPermission(c *gin.Context) bool {
	role := a.getRoles(c)
	method := c.Request.Method
	path := c.Request.URL.Path
	return a.enforcer.Enforce(role, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *JwtAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}

func (a *JwtAuthorizer) GetLoggedInUser(c *gin.Context) *models.User {
	session := sessions.Default(c)

	tokenSession := session.Get("jwt")
	if tokenSession == nil {
		return nil
	}
	tokenString := fmt.Sprintf("%v", tokenSession)
	if tokenString == "" {
		log.Println("tokenString is empty ")
		return nil
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET"))

		return hmacSampleSecret, nil
	})

	if token == nil {
		log.Println("IN GetLoggedInUser getting claims Token is nil")

		return nil
	}

	var user *models.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("IN GetLoggedInUser getting claims")

		id, _ := claims["userid"].(string)
		role, _ := claims["role"].(string)
		name, _ := claims["name"].(string)
		uId, _ := strconv.Atoi(id)
		uIntId := uint(uId)
		log.Printf(" Claims for logged in user:: ID:%v, Role:%v, Name:%v \n", id, role, name)
		user = &models.User{Role: role, Name: name}
		user.ID = uIntId

		return user
	} else {
		return nil
	}
}

func (a *JwtAuthorizer) getRoles(c *gin.Context) string {
	//
	session := sessions.Default(c)
	tokenSession := session.Get("jwt")
	//log.Println("old jwt token:: ", tokenString)
	tokenString := fmt.Sprintf("%v", tokenSession)

	///
	//tokenString := r.Header.Get("Authorization")
	log.Println("IN Get Roles ", tokenString)
	if tokenString == "" {
		log.Println("tokenString is empty ")
		return "anonymous"
	}

	//splitToken := strings.Split(tokenString, "Bearer")
	log.Println("JwtAuthorizer Token String: ", tokenString)
	//tokenString = strings.TrimSpace(splitToken[1])

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		//hmacSampleSecret := []byte(config.Jwtsecret)
		// os.Getenv("API_KEY")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET"))

		return hmacSampleSecret, nil
	})
	role := "anonymous"
	if token == nil {
		return role
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, _ = claims["role"].(string)

		log.Println("Role is ", role)
		log.Println("All claims::: ", claims)
	}
	return role
}
