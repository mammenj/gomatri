package security

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// config, err := config.GetConfiguration("config.json")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//hmacSampleSecret := []byte(config.Jwtsecret)
		// os.Getenv("API_KEY")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET"))

		return hmacSampleSecret, nil
	})
	var role string
	if token == nil {
		return "anonymous"
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, _ = claims["role"].(string)
		log.Println("Role is ", role)
	} else {
		log.Println("Error getting claims:: ", err)
	}
	return role
}
