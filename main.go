package main

import (
	//"embed"
	"gomatri/handlers"
	"gomatri/models"
	"gomatri/security"
	"gomatri/storage"
	"html/template"

	//"io/fs"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/casbin/casbin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

/*go:embed templates*/
//var templateFS embed.FS

/*go:embed static*/
//var staticFiles embed.FS

// var templateFS fs.FS
//var staticFiles fs.FS

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cookie_secret := []byte(os.Getenv("COOKIE_SECRET"))
	r := gin.Default()

	store := cookie.NewStore([]byte(cookie_secret))
	r.Use(sessions.Sessions("jwt-session", store))

	r.Use(func(c *gin.Context) {
		c.Header("User-Agent", "Unreal-monay")
	})

	//static, err := fs.Sub(staticFiles, "static")
	//if err != nil {
	//	panic(err)
	//}

	r.StaticFS("/static", http.Dir("static"))
	//r.StaticFS("/static", http.FS(static))

	/// TEST CODE FOR EMBED

	e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv", true)
	r.Use(security.NewJwtAuthorizer(e))

	r.GET("/", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/index.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)
	})

	r.GET("/matri.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/matri.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)
	})

	r.GET("/contact.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/contact.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)
	})

	r.GET("/grooms.html", func(c *gin.Context) {
		adStore := storage.NewSqliteAdsStore()
		offset := c.Query("offset")
		offsetInt, _ := strconv.Atoi(offset)
		log.Println(" Groom offset ", offset)
		ads, err := adStore.GetSection("Groom Wanted", offsetInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		admap := map[string][]models.Ad{"Ads": ads}
		//log.Println("Ad Map", admap)
		tmpl := template.Must(template.ParseFiles(
			"templates/grooms.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, admap)
	})

	r.GET("/brides.html", func(c *gin.Context) {
		adStore := storage.NewSqliteAdsStore()
		offset := c.Query("offset")
		offsetInt, _ := strconv.Atoi(offset)
		log.Println(" Bride offset ", offset)
		ads, err := adStore.GetSection("Bride Wanted", offsetInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		admap := map[string][]models.Ad{"Ads": ads}
		tmpl := template.Must(template.ParseFiles(
			"templates/brides.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, admap)
	})

	r.GET("/ads.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/placead.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)
	})

	r.GET("/login.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/loginregister.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)
	})

	r.GET("/tnc.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/tnc.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)
	})

	/// TEST CODE FOR EMBED END

	userHandler := handlers.CreateNewUserHandler()
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.PATCH("/users", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/login", security.Login)

	adHandler := handlers.CreateNewAdHandler()
	r.POST("/ads", adHandler.CreateAd)
	r.GET("/ads", adHandler.GetAds)
	r.PATCH("/ads", adHandler.UpdateAd)
	r.DELETE("/ads/:id", adHandler.DeleteAd)
	r.GET("/ads/:id", adHandler.GetAd)

	r.Run()
}
