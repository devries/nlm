package main

import (
	"embed"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/devries/nlm"
)

//go:embed templates/*.html
var templateFiles embed.FS

const size = 5

type Name struct {
	FirstName string
	LastName  string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()

	templates := template.Must(template.New("web").ParseFS(templateFiles, "templates/*"))

	articleBuilder, err := nlm.NewArticleBuilder(size)
	if err != nil {
		log.Fatalf("Error creating article builder: %s", err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		indexTemplate := templates.Lookup("index.html")
		err := indexTemplate.Execute(w, nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/getarticle", func(w http.ResponseWriter, r *http.Request) {
		responseTemplate := templates.Lookup("article.html")
		errorTemplate := templates.Lookup("error.html")

		// data := nlm.GetArticle()
		articleNames, err := nlm.GetArticleList()
		if err != nil {
			log.Print(err)
			errorTemplate.Execute(w, "Unable to render article")
			return
		}

		selectedName := articleNames[rand.Int()%len(articleNames)]
		data, err := nlm.GetNamedArticle(selectedName)
		if err != nil {
			log.Print(err)
			errorTemplate.Execute(w, "Unable to render article")
			return
		}

		err = responseTemplate.Execute(w, data)
		if err != nil {
			log.Print(err)
			errorTemplate.Execute(w, "Unable to render article")
			return
		}
	})

	mux.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		responseTemplate := templates.Lookup("article.html")
		errorTemplate := templates.Lookup("error.html")

		article := articleBuilder.GenerateArticle(120, 5000)

		err = responseTemplate.Execute(w, article)
		if err != nil {
			log.Printf("error writing article template: %s", err)
			errorTemplate.Execute(w, "Unable to render article")
		}
	})

	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingHandler(mux),
	}

	log.Print("Server starting on port 8080")
	log.Fatal(server.ListenAndServe())
}
