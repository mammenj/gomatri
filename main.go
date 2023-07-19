package main

import (
	"embed"
	"gomatri/handlers"
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

	/// TEST CODE FOR EMBED END

	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.PATCH("/users", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.GET("/users/:id", handlers.GetUser)

	r.Run()
}
