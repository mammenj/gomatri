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

var rootTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/index.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var contactTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/contact.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var groomTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/grooms.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var brideTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/brides.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var placeAdsTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/placead.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var loginTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/loginregister.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var tncTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/tnc.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var logoutTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/user.html"))

type pageData struct {
	User  models.User
	AdMap map[string][]models.Ad
}

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
		c.Header("User-Agent", "Unreal-Minna_Minny")
	})

	r.StaticFS("/static", http.Dir("static"))

	e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv", true)

	auth := security.NewJwtAuth(e)

	r.Use(security.NewJwtAuthorizer(e))

	r.GET("/", func(c *gin.Context) {
		var page pageData

		user := auth.GetLoggedInUser(c)
		if user != nil {

			page = pageData{*user, nil}
		} else {
			page = pageData{models.User{Name: ""}, nil}
		}
		rootTemplate.Execute(c.Writer, page)
	})

	r.GET("/contact.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		contactTemplate.Execute(c.Writer, page)
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
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			log.Println("In groom user !=null ")
			page = pageData{*user, admap}
		} else {
			page = pageData{models.User{}, admap}
		}

		groomTemplate.Execute(c.Writer, page)
	})

	r.GET("/brides.html", func(c *gin.Context) {
		var page pageData
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
		//brideTemplate.Execute(c.Writer, admap)
		user := auth.GetLoggedInUser(c)
		if user != nil {
			log.Println("In bride user !=null ")
			page = pageData{*user, admap}
		} else {
			page = pageData{models.User{}, admap}
		}
		//log.Println("Page in brides ", page)
		brideTemplate.Execute(c.Writer, page)
	})

	r.GET("/ads.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		placeAdsTemplate.Execute(c.Writer, page)
	})

	r.GET("/login.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}

		log.Println("Logged in Page is :::::", page)
		loginTemplate.Execute(c.Writer, page)
	})

	r.GET("/tnc.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		tncTemplate.Execute(c.Writer, page)
	})

	/// TEST CODE FOR EMBED END

	userHandler := handlers.CreateNewUserHandler()
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.PATCH("/users", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/login", security.Login)
	r.POST("/logout", func(c *gin.Context) {

		session := sessions.Default(c)
		mypage := pageData{models.User{}, nil}
		log.Println(mypage)
		session.Delete("jwt")
		session.Clear()
		session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0
		session.Save()
		//c.Redirect(http.StatusTemporaryRedirect, "/")
		//c.Redirect(http.StatusFound, "/")
		logoutTemplate.Execute(c.Writer, "Logged Out")

	})

	adHandler := handlers.CreateNewAdHandler()
	r.POST("/ads", adHandler.CreateAd)
	r.GET("/ads", adHandler.GetAds)
	r.PATCH("/ads", adHandler.UpdateAd)
	r.DELETE("/ads/:id", adHandler.DeleteAd)
	r.GET("/ads/:id", adHandler.GetAd)

	r.Run()
}
