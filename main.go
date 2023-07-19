package main

import (
	"embed"
	"fmt"
	"gomatri/models"
	"gomatri/storage"
	"html/template"
	"io/fs"
	"log"
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

	r.GET("/", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html"))

		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		tmpl.Execute(c.Writer, films)

	})

	r.POST("/films", func(c *gin.Context) {

		title := c.PostForm("title")
		director := c.PostForm("director")
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html"))
		tmpl.ExecuteTemplate(c.Writer, "matri-list-element", Film{Title: title, Director: director})
	})

	r.POST("/users", func(c *gin.Context) {
		log.Println("IN Create handler")
		var input models.User
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userStore := storage.NewSqliteUserStore()
		ID, createErr := userStore.Create(&input)
		if createErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("USER CREATED ID: ", ID)
	})

	r.GET("/users", func(c *gin.Context) {
		log.Println("IN GET handler")
		userStore := storage.NewSqliteUserStore()
		users, createErr := userStore.Get()
		if createErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("USER List", users)
	})

	r.PATCH("/users", func(c *gin.Context) {
		log.Println("IN PATCH  handler")
		var input models.User
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println("IN PATCH  handler ", &input)
		userStore := storage.NewSqliteUserStore()
		ID, updatedErr := userStore.Update(&input)
		if updatedErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("USER Updated ID: ", ID)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {

		log.Println("IN Delete handler")

		id := c.Param("id")

		userStore := storage.NewSqliteUserStore()
		ID, deleteErr := userStore.Delete(id)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("USER deleted ID: ", ID)
	})

	r.GET("/users/:id", func(c *gin.Context) {

		log.Println("IN GET one handler")
		id := c.Param("id")

		userStore := storage.NewSqliteUserStore()
		user, err := userStore.GetOne(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println("...... Get user: ", user)
	})

	r.Run()
}
