package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Film struct {
	Title    string
	Director string
}

//go:embed templates
var templateFS embed.FS

func main1() {
	fmt.Println("Go app...")

	// handler function #1 - returns the index.html template, with film data
	h1 := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html"))

		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		tmpl.Execute(w, films)
	}

	// handler function #2 - returns the template block with the newly added film, as an HTMX response
	h2 := func(w http.ResponseWriter, r *http.Request) {
		//time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		//tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
	}

	// define handlers
	http.HandleFunc("/", h1)
	http.HandleFunc("/add-film/", h2)
	http.Handle("/static/", http.FileServer(http.Dir("/static/")))

	log.Fatal(http.ListenAndServe(":8000", nil))

}

func main() {
	r := gin.Default()

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

	r.POST("/add-film/", func(c *gin.Context) {

		title := c.PostForm("title")
		director := c.PostForm("director")
		tmpl := template.Must(template.ParseFS(templateFS,
			"templates/index.html"))
		tmpl.ExecuteTemplate(c.Writer, "film-list-element", Film{Title: title, Director: director})
	})

	r.Run()
}
