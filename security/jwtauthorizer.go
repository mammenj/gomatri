package security

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	role := a.getRoles(c.Request)
	method := c.Request.Method
	path := c.Request.URL.Path
	return a.enforcer.Enforce(role, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *JwtAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}

func (a *JwtAuthorizer) getRoles(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, "Bearer")
	tokenString = strings.TrimSpace(splitToken[1])
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
		hmacSampleSecret := []byte(os.Getenv("SECRET"))
		return hmacSampleSecret, nil
	})
	var role string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, _ = claims["role"].(string)
	} else {
		log.Println("Error getting claims:: ", err)
	}
	return role
}
