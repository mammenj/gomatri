package main

import (
	"embed"
	"gomatri/handlers"
	"gomatri/models"
	"gomatri/storage"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Film struct {
	Title    string
	Director string
}

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFiles embed.FS

func main() {
	r := gin.Default()

	static, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	r.StaticFS("/static", http.FS(static))

	/// TEST CODE FOR EMBED

	r.GET("/", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)

	})

	r.GET("/matri.html", func(c *gin.Context) {

		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/matri.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)

	})

	r.GET("/contact.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/contact.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)

	})

	r.GET("/grooms.html", func(c *gin.Context) {

		adStore := storage.NewSqliteAdsStore()
		ads, err := adStore.Get()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		admap := map[string][]models.Ad{"Ads": ads}
		//log.Println("Ad Map", admap)
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/grooms.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, admap)

	})

	r.GET("/brides.html", func(c *gin.Context) {
		adStore := storage.NewSqliteAdsStore()
		ads, err := adStore.Get()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		admap := map[string][]models.Ad{"Ads": ads}
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/brides.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, admap)

	})

	r.GET("/ads.html", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/placead.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(c.Writer, nil)

	})

	r.POST("/films", func(c *gin.Context) {

		title := c.PostForm("title")
		director := c.PostForm("director")
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html", "templates/header.html", "templates/footer.html"))
		tmpl.ExecuteTemplate(c.Writer, "matri-list-element", Film{Title: title, Director: director})
	})

	/// TEST CODE FOR EMBED END

	userHandler := handlers.CreateNewUserHandler()
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.PATCH("/users", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.GET("/users/:id", userHandler.GetUser)

	adHandler := handlers.CreateNewAdHandler()
	r.POST("/ads", adHandler.CreateAd)
	r.GET("/ads", adHandler.GetAds)
	r.PATCH("/ads", adHandler.UpdateAd)
	r.DELETE("/ads/:id", adHandler.DeleteAd)
	r.GET("/ads/:id", adHandler.GetAd)

	r.Run()
}
