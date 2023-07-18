package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Film struct {
	Title    string
	Director string
}

//go:embed templates
var templateFS embed.FS

//go:embed static
var staticFiles embed.FS

func main() {
	r := gin.Default()

	r.StaticFS("/static", http.FS(staticFiles))

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

	r.Run()
}
